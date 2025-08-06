package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/Hyodar/tdxs/pkg/api"
	"github.com/Hyodar/tdxs/pkg/logger"
	"github.com/Hyodar/tdxs/pkg/transport"
)

type SocketTransport struct {
	transport.Transport

	cfg      *SocketTransportConfig
	queues   *transport.TransportQueues
	listener net.Listener
	logger   logger.Logger
}

type SocketTransportConfig struct {
	FilePath string      `yaml:"file_path"`
	Owner    string      `yaml:"owner"`
	Group    string      `yaml:"group"`
	Perm     os.FileMode `yaml:"perm"`
}

func NewSocketTransport(cfg *SocketTransportConfig, logger logger.Logger) transport.Transport {
	return &SocketTransport{
		cfg:    cfg,
		logger: logger,
	}
}

func (t *SocketTransport) Start(ctx context.Context, queues *transport.TransportQueues) error {
	t.queues = queues

	if err := os.RemoveAll(t.cfg.FilePath); err != nil {
		return fmt.Errorf("failed to remove existing socket: %w", err)
	}

	listener, err := net.Listen("unix", t.cfg.FilePath)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}
	t.listener = listener

	if err := os.Chmod(t.cfg.FilePath, t.cfg.Perm); err != nil {
		listener.Close()
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}

	go t.acceptConnections(ctx)

	return nil
}

func (t *SocketTransport) acceptConnections(ctx context.Context) {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				// Log error and continue
				continue
			}
		}

		go t.handleConnection(ctx, conn)
	}
}

func (t *SocketTransport) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var req SocketTransportRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			encoder.Encode(NewIssueResponseFromError(fmt.Errorf("failed to decode request: %w", err)))
			continue
		}

		apiRequest, err := req.UnmarshalData()
		if err != nil {
			encoder.Encode(NewIssueResponseFromError(fmt.Errorf("failed to unmarshal request data: %w", err)))
			continue
		}

		switch req.Method {
		case SocketTransportRequestMethodIssue:
			issueReq := apiRequest.(*api.IssueRequest)
			wrapper := &api.IssueRequestWrapper{
				Request:  issueReq,
				Response: make(chan *api.IssueResponse, 1),
			}

			select {
			case t.queues.IssueQueue <- wrapper:
				select {
				case resp := <-wrapper.Response:
					encoder.Encode(NewIssueResponseFromAPI(resp))
				case <-ctx.Done():
					encoder.Encode(NewIssueResponseFromError(fmt.Errorf("context cancelled")))
					return
				}
			case <-ctx.Done():
				encoder.Encode(NewIssueResponseFromError(fmt.Errorf("context cancelled")))
				return
			}

		case SocketTransportRequestMethodValidate:
			validateReq := apiRequest.(*api.ValidateRequest)
			wrapper := &api.ValidateRequestWrapper{
				Request:  validateReq,
				Response: make(chan *api.ValidateResponse, 1),
			}

			select {
			case t.queues.ValidateQueue <- wrapper:
				select {
				case resp := <-wrapper.Response:
					encoder.Encode(NewValidateResponseFromAPI(resp))
				case <-ctx.Done():
					encoder.Encode(NewValidateResponseFromError(fmt.Errorf("context cancelled")))
					return
				}
			case <-ctx.Done():
				encoder.Encode(NewValidateResponseFromError(fmt.Errorf("context cancelled")))
				return
			}

		default:
			encoder.Encode(NewIssueResponseFromError(fmt.Errorf("unknown method: %s", req.Method)))
		}
	}
}
