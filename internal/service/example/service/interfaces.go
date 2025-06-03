package service

import (
	"context"
)

type RepoI interface {
	Send(ctx context.Context, phone string, msg string, sync bool) error
}
