package main

import (
  "context"
  "fmt"
  "io/ioutil"
  "mime"
  "net"
  "net/http"
  "os"
  "strings"

  "github.com/grpc-ecosystem/grpc-gateway/runtime"
  "github.com/rakyll/statik/fs"
  "google.golang.org/grpc"
  "google.golang.org/grpc/grpclog"

  pbExample "github.com/grpc-gateway/proto"
  "github.com/grpc-gateway/server"

  // Static files
  _ "github.com/grpc-gateway/statik"
)

func getOpenAPIHandler() http.Handler {
  mime.AddExtensionType(".svg", "image/svg+xml")

  statikFS, err := fs.New()
  if err != nil {
	panic("creating OpenAPI filesystem: " + err.Error())
  }

  return http.FileServer(statikFS)
}

func main() {

  log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
  grpclog.SetLoggerV2(log)

  addr := "0.0.0.0:10000"
  lis, err := net.Listen("tcp", addr)
  if err != nil {
	log.Fatalln("Failed to listen:", err)
  }

  grpcServer := grpc.NewServer(
  )
  pbExample.RegisterUserServiceServer(grpcServer, server.New())

  // Serve gRPC Server
  log.Info("Serving gRPC on https://", addr)
  go func() {
	log.Fatal(grpcServer.Serve(lis))
  }()

  dialAddr := fmt.Sprintf("dns:///%s", addr)
  grpcConn, err := grpc.DialContext(
	context.Background(),
	dialAddr,
	grpc.WithInsecure(),
	grpc.WithBlock(),
  )
  if err != nil {
	log.Fatalln("Failed to dial server:", err)
  }

  jsonpb := &runtime.JSONPb{
	EmitDefaults: true,
	Indent:       "  ",
	OrigName:     true,
  }
  gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb))
  err = pbExample.RegisterUserServiceHandler(context.Background(), gwmux, grpcConn)
  if err != nil {
	log.Fatalln("Failed to register gateway:", err)
  }

  apiHandler := getOpenAPIHandler()

  port := os.Getenv("PORT")
  if port == "" {
	port = "11000"
  }
  gatewayAddr := "0.0.0.0:" + port
  gatewayServer := &http.Server{
	Addr: gatewayAddr,
	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  if strings.HasPrefix(r.URL.Path, "/api") {
		gwmux.ServeHTTP(w, r)
		return
	  }
	  apiHandler.ServeHTTP(w, r)
	}),
  }
  os.Setenv("SERVE_HTTP", "true")
  if strings.ToLower(os.Getenv("SERVE_HTTP")) == "true" {
	log.Info("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	log.Fatalln(gatewayServer.ListenAndServe())
  }
}
