package validator

import (
	"context"

	"github.com/Hyodar/tdxs/pkg/api"
)

type Validator interface {
	Start(ctx context.Context) error
	Validate(ctx context.Context, req *api.ValidateRequest) *api.ValidateResponse
}

type ValidatorType string

const (
	ValidatorTypeAzure     ValidatorType = "azure"
	ValidatorTypeSimulator ValidatorType = "simulator"
)
