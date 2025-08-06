package azure

import (
	"context"
	"fmt"

	azure "github.com/Hyodar/tdxs/internal/constellation/attestation/azure/tdx"
	"github.com/Hyodar/tdxs/internal/constellation/config"

	"github.com/Hyodar/tdxs/pkg/api"
	"github.com/Hyodar/tdxs/pkg/logger"
	"github.com/Hyodar/tdxs/pkg/validator"
)

type AzureValidator struct {
	validator.Validator

	logger  logger.Logger
	backend *azure.Validator
}

type AzureValidatorConfig struct {
	*config.AzureTDX `yaml:",inline"`
}

func NewAzureValidator(cfg *AzureValidatorConfig, logger logger.Logger) *AzureValidator {
	return &AzureValidator{
		backend: azure.NewValidator(cfg.AzureTDX, logger),
		logger:  logger,
	}
}

func (i *AzureValidator) Start(_ context.Context) error {
	return nil
}

func (i *AzureValidator) Validate(ctx context.Context, req *api.ValidateRequest) *api.ValidateResponse {
	if i.backend == nil {
		return &api.ValidateResponse{Error: fmt.Errorf("backend not initialized")}
	}

	userData, err := i.backend.Validate(ctx, req.Document, req.Nonce)
	if err != nil {
		return &api.ValidateResponse{Error: err}
	}
	return &api.ValidateResponse{UserData: userData, Valid: true}
}
