package manager

import (
	"context"
	"fmt"

	"github.com/Hyodar/tdxs/pkg/api"
	"github.com/Hyodar/tdxs/pkg/issuer"
	azureissuer "github.com/Hyodar/tdxs/pkg/issuer/azure"
	simulatorissuer "github.com/Hyodar/tdxs/pkg/issuer/simulator"
	"github.com/Hyodar/tdxs/pkg/logger"
	"github.com/Hyodar/tdxs/pkg/transport"
	sockettransport "github.com/Hyodar/tdxs/pkg/transport/socket"
	"github.com/Hyodar/tdxs/pkg/validator"
	azurevalidator "github.com/Hyodar/tdxs/pkg/validator/azure"
	simulatorvalidator "github.com/Hyodar/tdxs/pkg/validator/simulator"
	"gopkg.in/yaml.v3"
)

type Manager struct {
	transport transport.Transport
	issuer    issuer.Issuer
	validator validator.Validator
	logger    logger.Logger
}

type ManagerConfig struct {
	Transport *TransportConfig `json:"transport" yaml:"transport"`
	Issuer    *IssuerConfig    `json:"issuer" yaml:"issuer"`
	Validator *ValidatorConfig `json:"validator" yaml:"validator"`
}

func NewManager(cfg *ManagerConfig, logger logger.Logger) (*Manager, error) {
	if cfg.Transport == nil {
		return nil, fmt.Errorf("transport config is required")
	}
	if cfg.Issuer == nil && cfg.Validator == nil {
		return nil, fmt.Errorf("issuer or validator config is required")
	}

	transport, err := createTransport(cfg.Transport, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	var issuer issuer.Issuer
	if cfg.Issuer != nil {
		issuer, err = createIssuer(cfg.Issuer, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create issuer: %w", err)
		}
	}

	var validator validator.Validator
	if cfg.Validator != nil {
		validator, err = createValidator(cfg.Validator, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create validator: %w", err)
		}
	}

	return &Manager{
		logger:    logger,
		transport: transport,
		issuer:    issuer,
		validator: validator,
	}, nil
}

type TransportConfig struct {
	Type   transport.TransportType `yaml:"-"`
	Config interface{}             `yaml:"-"`
}

func (t *TransportConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type transportConfigHelper struct {
		Type   transport.TransportType `yaml:"type"`
		Config yaml.Node               `yaml:"config"`
	}
	var tc transportConfigHelper
	if err := unmarshal(&tc); err != nil {
		return err
	}

	t.Type = tc.Type

	switch t.Type {
	case transport.TransportTypeSocket:
		var cfg sockettransport.SocketTransportConfig
		if err := tc.Config.Decode(&cfg); err != nil {
			return err
		}
		t.Config = cfg
	default:
		return fmt.Errorf("invalid transport type: %s", t.Type)
	}

	return nil
}

type IssuerConfig struct {
	Type   issuer.IssuerType `yaml:"-"`
	Config interface{}       `yaml:"-"`
}

func (i *IssuerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type issuerConfigHelper struct {
		Type   issuer.IssuerType `yaml:"type"`
		Config *yaml.Node        `yaml:"config"`
	}
	var ic issuerConfigHelper
	if err := unmarshal(&ic); err != nil {
		return err
	}

	i.Type = ic.Type

	switch i.Type {
	case issuer.IssuerTypeAzure, issuer.IssuerTypeSimulator:
		if ic.Config != nil {
			return fmt.Errorf("issuer config is not supported for type: %s", i.Type)
		}
	default:
		return fmt.Errorf("invalid issuer type: %s", i.Type)
	}

	return nil
}

type ValidatorConfig struct {
	Type   validator.ValidatorType `yaml:"-"`
	Config interface{}             `yaml:"-"`
}

func (v *ValidatorConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type validatorConfigHelper struct {
		Type   validator.ValidatorType `yaml:"type"`
		Config *yaml.Node              `yaml:"config"`
	}
	var vc validatorConfigHelper
	if err := unmarshal(&vc); err != nil {
		return err
	}

	v.Type = vc.Type

	switch v.Type {
	case validator.ValidatorTypeAzure:
		var cfg azurevalidator.AzureValidatorConfig
		if err := vc.Config.Decode(&cfg); err != nil {
			return err
		}
		v.Config = cfg
	case validator.ValidatorTypeSimulator:
		if vc.Config != nil {
			return fmt.Errorf("validator config is not supported for type: %s", v.Type)
		}
	default:
		return fmt.Errorf("invalid validator type: %s", v.Type)
	}

	return nil
}

func createTransport(cfg *TransportConfig, logger logger.Logger) (transport.Transport, error) {
	switch cfg.Type {
	case transport.TransportTypeSocket:
		innerCfg, ok := cfg.Config.(sockettransport.SocketTransportConfig)
		if !ok {
			return nil, fmt.Errorf("invalid transport config type: %T", cfg.Config)
		}
		return sockettransport.NewSocketTransport(&innerCfg, logger), nil
	default:
		return nil, fmt.Errorf("invalid transport type: %s", cfg.Type)
	}
}

func createIssuer(cfg *IssuerConfig, logger logger.Logger) (issuer.Issuer, error) {
	switch cfg.Type {
	case issuer.IssuerTypeAzure:
		return azureissuer.NewAzureIssuer(logger), nil
	case issuer.IssuerTypeSimulator:
		return simulatorissuer.NewSimulatorIssuer(logger), nil
	default:
		return nil, fmt.Errorf("invalid issuer type: %s", cfg.Type)
	}
}

func createValidator(cfg *ValidatorConfig, logger logger.Logger) (validator.Validator, error) {
	switch cfg.Type {
	case validator.ValidatorTypeAzure:
		innerCfg, ok := cfg.Config.(azurevalidator.AzureValidatorConfig)
		if !ok {
			return nil, fmt.Errorf("invalid validator config type: %T", cfg.Config)
		}
		return azurevalidator.NewAzureValidator(&innerCfg, logger), nil
	case validator.ValidatorTypeSimulator:
		return simulatorvalidator.NewSimulatorValidator(logger), nil
	default:
		return nil, fmt.Errorf("invalid validator type: %s", cfg.Type)
	}
}

func (m *Manager) Start(ctx context.Context) error {
	queues := &transport.TransportQueues{
		IssueQueue:    make(chan *api.IssueRequestWrapper, 100),
		ValidateQueue: make(chan *api.ValidateRequestWrapper, 100),
	}

	transportCtx, transportCancel := context.WithCancel(ctx)
	defer transportCancel()

	errChan := make(chan error, 1)
	go func() {
		if err := m.transport.Start(transportCtx, queues); err != nil {
			errChan <- fmt.Errorf("transport error: %w", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Manager shutting down")
			return ctx.Err()
		case err := <-errChan:
			return err
		case req := <-queues.IssueQueue:
			go m.handleIssueRequest(ctx, req)
		case req := <-queues.ValidateQueue:
			go m.handleValidateRequest(ctx, req)
		}
	}
}

func (m *Manager) handleIssueRequest(ctx context.Context, wrapper *api.IssueRequestWrapper) {
	response := m.issuer.Issue(ctx, wrapper.Request)
	select {
	case wrapper.Response <- response:
	case <-ctx.Done():
	}
}

func (m *Manager) handleValidateRequest(ctx context.Context, wrapper *api.ValidateRequestWrapper) {
	response := m.validator.Validate(ctx, wrapper.Request)
	select {
	case wrapper.Response <- response:
	case <-ctx.Done():
	}
}
