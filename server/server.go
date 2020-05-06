package server

import (
  "context"
  "sync"

  "github.com/gofrs/uuid"

  pbExample "github.com/grpc-gateway/proto"
)

// Backend implements the protobuf interface
type Backend struct {
  mu    *sync.RWMutex
  users []*pbExample.User
}

// New initializes a new Backend struct.
func New() *Backend {
  return &Backend{
	mu: &sync.RWMutex{},
  }
}

// AddUser adds a user to the in-memory store.
func (b *Backend) AddUser(ctx context.Context, _ *pbExample.AddUserRequest) (*pbExample.User, error) {
  b.mu.Lock()
  defer b.mu.Unlock()

  user := &pbExample.User{
	Id:   uuid.Must(uuid.NewV4()).String(),
	Name: uuid.Must(uuid.NewV4()).String(),
  }
  b.users = append(b.users, user)
  return user, nil
}

// ListUsers lists all users in the store.
func (b *Backend) ListUsers(_ *pbExample.ListUsersRequest, srv pbExample.UserService_ListUsersServer) error {
  b.mu.RLock()
  defer b.mu.RUnlock()

  for _, user := range b.users {
	err := srv.Send(user)
	if err != nil {
	  return err
	}
  }
  return nil
}

// ListUsers lists all users in the store.
func (b *Backend) GetUser(req *pbExample.ListUserIdRequest, srv pbExample.UserService_GetUserServer) error {
  b.mu.RLock()
  defer b.mu.RUnlock()

  for _, user := range b.users {
	if req.Id == user.Id {
	  println("Data found = " + user.String())
	  err := srv.Send(user)
	  if err != nil {
		return err
	  }
	  return nil
	}
  }
  return nil
}
