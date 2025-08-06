package azure

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Hyodar/tdxs/pkg/api"
	"github.com/Hyodar/tdxs/pkg/logger"
	"github.com/Hyodar/tdxs/pkg/validator"
)

type SimulatorValidator struct {
	validator.Validator

	logger logger.Logger
}

func NewSimulatorValidator(logger logger.Logger) *SimulatorValidator {
	return &SimulatorValidator{
		logger: logger,
	}
}

func (i *SimulatorValidator) Start(_ context.Context) error {
	return nil
}

func (i *SimulatorValidator) Validate(_ context.Context, req *api.ValidateRequest) *api.ValidateResponse {
	if req.Document == nil {
		return &api.ValidateResponse{Error: fmt.Errorf("document is nil")}
	}

	type Document struct {
		UserData string `json:"userData"`
		Nonce    string `json:"nonce"`
	}

	var doc Document
	if err := json.Unmarshal(req.Document, &doc); err != nil {
		return &api.ValidateResponse{Error: err}
	}

	userData, err := hex.DecodeString(doc.UserData)
	if err != nil {
		return &api.ValidateResponse{Error: err}
	}

	nonce, err := hex.DecodeString(doc.Nonce)
	if err != nil {
		return &api.ValidateResponse{Error: err}
	}

	return &api.ValidateResponse{UserData: userData, Valid: bytes.Equal(req.Nonce, nonce)}
}
