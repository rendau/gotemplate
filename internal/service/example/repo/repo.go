package repo

import (
	"context"
	"google.golang.org/grpc"
)

type Repo struct {
}

func New(grpcClient *grpc.ClientConn) *Repo {
	result := &Repo{}

	return result
}

func (r *Repo) Send(ctx context.Context, msg string) error {
	return nil
}
