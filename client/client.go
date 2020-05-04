package main

import (
  "context"
  "fmt"
  pbExample "github.com/grpc-gateway/proto"
  "google.golang.org/grpc"
  "log"
  "time"
)

func main() {

  addr := "0.0.0.0:10000"
  dialAddr := fmt.Sprintf("dns:///%s", addr)
  conn, err := grpc.DialContext(
	context.Background(),
	dialAddr,
	grpc.WithInsecure(),
	grpc.WithBlock(),
  )
  if err != nil {
	log.Fatalln("Failed to dial server:", err)
  }

  defer conn.Close()
  c := pbExample.NewUserServiceClient(conn)
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()

  r, err := c.GetUser(ctx, &pbExample.ListUserIdRequest{Id: "03616299-23fd-4b27-90ca-f2e9d24bd11e"})
  if err != nil {
	log.Fatalf("could not greet: %v", err)
  }

  recv, err := r.Recv()
  log.Println("Recieved " + recv.GetId())
}
