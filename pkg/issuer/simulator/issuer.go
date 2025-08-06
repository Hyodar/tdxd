package azure

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/Hyodar/tdxs/pkg/api"
	"github.com/Hyodar/tdxs/pkg/issuer"
	"github.com/Hyodar/tdxs/pkg/logger"
)

type SimulatorIssuer struct {
	issuer.Issuer

	logger logger.Logger
}

func NewSimulatorIssuer(logger logger.Logger) *SimulatorIssuer {
	return &SimulatorIssuer{
		logger: logger,
	}
}

func (i *SimulatorIssuer) Start(_ context.Context) error {
	return nil
}

func (i *SimulatorIssuer) Issue(ctx context.Context, req *api.IssueRequest) *api.IssueResponse {
	type Document struct {
		UserData string `json:"userData"`
		Nonce    string `json:"nonce"`
	}

	doc := Document{
		UserData: hex.EncodeToString(req.UserData),
		Nonce:    hex.EncodeToString(req.Nonce),
	}

	jsonDoc, err := json.Marshal(doc)
	if err != nil {
		return &api.IssueResponse{Error: err}
	}

	return &api.IssueResponse{Document: jsonDoc}
}
