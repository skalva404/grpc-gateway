package server

import (
  "context"
  "strconv"
  "sync"

  pbExample "github.com/grpc-gateway/proto"
)

var counter int64 = 0

// Backend implements the protobuf interface
type Backend struct {
  mu      *sync.RWMutex
  deinfos []*pbExample.DeInfo
}

// New initializes a new Backend struct.
func New() *Backend {
  return &Backend{
	mu: &sync.RWMutex{},
  }
}

// AddUser adds a user to the in-memory store.
func (b *Backend) AddDeInfo(ctx context.Context, deinfo *pbExample.DeInfo) (*pbExample.DeInfo, error) {
  b.mu.Lock()
  defer b.mu.Unlock()
  counter++
  deinfo.Id = counter
  b.deinfos = append(b.deinfos, deinfo)
  return deinfo, nil
}

func (b *Backend) ListDeInfos(_ *pbExample.Empty, srv pbExample.DeInfoService_ListDeInfosServer) error {
  b.mu.RLock()
  defer b.mu.RUnlock()

  for _, deinfo := range b.deinfos {
	err := srv.Send(deinfo)
	if err != nil {
	  return err
	}
  }
  return nil
}

// ListUsers lists all users in the store.
func (b *Backend) GetDeInfo(req *pbExample.GetDeInfoRequest, srv pbExample.DeInfoService_GetDeInfoServer) error {
  b.mu.RLock()
  defer b.mu.RUnlock()

  for _, deinfo := range b.deinfos {
	deid, _ := strconv.ParseInt(req.Id, 10, 64)
	var did int64 = deid
	if did == deinfo.Id {
	  println("Data found = " + deinfo.String())
	  err := srv.Send(deinfo)
	  if err != nil {
		return err
	  }
	  return nil
	}
  }
  return nil
}
