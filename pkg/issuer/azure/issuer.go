package azure

import (
	"context"

	azuretdx "github.com/Hyodar/tdxs/internal/constellation/attestation/azure/tdx"

	"github.com/Hyodar/tdxs/pkg/api"
	"github.com/Hyodar/tdxs/pkg/issuer"
	"github.com/Hyodar/tdxs/pkg/logger"
)

type AzureIssuer struct {
	issuer.Issuer

	logger  logger.Logger
	backend *azuretdx.Issuer
}

func NewAzureIssuer(logger logger.Logger) *AzureIssuer {
	return &AzureIssuer{
		backend: azuretdx.NewIssuer(logger),
		logger:  logger,
	}
}

func (i *AzureIssuer) Start(_ context.Context) error {
	return nil
}

func (i *AzureIssuer) Issue(ctx context.Context, req *api.IssueRequest) *api.IssueResponse {
	doc, err := i.backend.Issue(ctx, req.UserData, req.Nonce)
	if err != nil {
		return &api.IssueResponse{Error: err}
	}
	return &api.IssueResponse{Document: doc}
}
