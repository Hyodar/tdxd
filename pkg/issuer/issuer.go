package issuer

import (
	"context"

	"github.com/Hyodar/tdxs/pkg/api"
)

type Issuer interface {
	Start(ctx context.Context) error
	Issue(ctx context.Context, req *api.IssueRequest) *api.IssueResponse
	Metadata(ctx context.Context, req *api.MetadataRequest) *api.MetadataResponse
}

type IssuerType string

const (
	IssuerTypeAzure     IssuerType = "azure"
	IssuerTypeSimulator IssuerType = "simulator"
)

const (
	MetadataUserData = "userData"
	MetadataNonce    = "nonce"
)
