package transport

import (
	"context"

	"github.com/Hyodar/tdxs/pkg/api"
)

type TransportQueues struct {
	IssueQueue    chan *api.IssueRequestWrapper
	ValidateQueue chan *api.ValidateRequestWrapper
}

type Transport interface {
	Start(ctx context.Context, queues *TransportQueues) error
}

type TransportType string

const (
	TransportTypeSocket TransportType = "socket"
)
