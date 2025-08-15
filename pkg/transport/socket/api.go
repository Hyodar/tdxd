package socket

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Hyodar/tdxs/pkg/api"
)

type SocketTransportRequestMethod string

const (
	SocketTransportRequestMethodIssue    SocketTransportRequestMethod = "issue"
	SocketTransportRequestMethodMetadata SocketTransportRequestMethod = "metadata"
	SocketTransportRequestMethodValidate SocketTransportRequestMethod = "validate"
)

type SocketTransportRequest struct {
	Method SocketTransportRequestMethod `json:"type"`
	Data   json.RawMessage              `json:"data"`
}

func (r *SocketTransportRequest) UnmarshalData() (any, error) {
	switch r.Method {
	case SocketTransportRequestMethodIssue:
		var issueRequest SocketTransportIssueRequest
		if err := json.Unmarshal(r.Data, &issueRequest); err != nil {
			return nil, fmt.Errorf("failed to unmarshal issue request: %w", err)
		}
		return issueRequest.ToAPIRequest()
	case SocketTransportRequestMethodValidate:
		var validateRequest SocketTransportValidateRequest
		if err := json.Unmarshal(r.Data, &validateRequest); err != nil {
			return nil, fmt.Errorf("failed to unmarshal validate request: %w", err)
		}
		return validateRequest.ToAPIRequest()
	}
	return nil, fmt.Errorf("invalid method: %s", r.Method)
}

type SocketTransportIssueRequest struct {
	UserData string `json:"userData"`
	Nonce    string `json:"nonce"`
}

func (r *SocketTransportIssueRequest) ToAPIRequest() (*api.IssueRequest, error) {
	userData, err := hex.DecodeString(r.UserData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user data: %w", err)
	}

	nonce, err := hex.DecodeString(r.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	return &api.IssueRequest{UserData: userData, Nonce: nonce}, nil
}

type SocketTransportValidateRequest struct {
	Document string `json:"document"`
	Nonce    string `json:"nonce"`
}

func (r *SocketTransportValidateRequest) ToAPIRequest() (*api.ValidateRequest, error) {
	document, err := hex.DecodeString(r.Document)
	if err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	nonce, err := hex.DecodeString(r.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	return &api.ValidateRequest{Document: document, Nonce: nonce}, nil
}

type SocketTransportIssueResponse struct {
	Document string `json:"document"`
	Error    string `json:"error"`
}

func NewIssueResponseFromError(err error) *SocketTransportIssueResponse {
	return &SocketTransportIssueResponse{
		Error: fmt.Sprintf("transport error: %s", err.Error()),
	}
}

func NewIssueResponseFromAPI(response *api.IssueResponse) *SocketTransportIssueResponse {
	if response.Error != nil {
		return &SocketTransportIssueResponse{
			Error: fmt.Sprintf("validator error: %s", response.Error.Error()),
		}
	}

	return &SocketTransportIssueResponse{
		Document: hex.EncodeToString(response.Document),
	}
}

type SocketTransportMetadataResponse struct {
	IssuerType string `json:"issuerType"`
	UserData   string `json:"userData"`
	Nonce      string `json:"nonce"`
	Metadata   any    `json:"metadata"`
	Error      string `json:"error"`
}

func NewMetadataResponseFromError(err error) *SocketTransportMetadataResponse {
	return &SocketTransportMetadataResponse{
		Error: fmt.Sprintf("transport error: %s", err.Error()),
	}
}

func NewMetadataResponseFromAPI(response *api.MetadataResponse) *SocketTransportMetadataResponse {
	if response.Error != nil {
		return &SocketTransportMetadataResponse{
			Error: fmt.Sprintf("validator error: %s", response.Error.Error()),
		}
	}

	return &SocketTransportMetadataResponse{
		IssuerType: response.IssuerType,
		UserData:   hex.EncodeToString(response.UserData),
		Nonce:      hex.EncodeToString(response.Nonce),
		Metadata:   response.Metadata,
	}
}

type SocketTransportValidateResponse struct {
	UserData string `json:"userData"`
	Valid    bool   `json:"valid"`
	Error    string `json:"error"`
}

func NewValidateResponseFromError(err error) *SocketTransportValidateResponse {
	return &SocketTransportValidateResponse{
		Error: fmt.Sprintf("transport error: %s", err.Error()),
	}
}

func NewValidateResponseFromAPI(response *api.ValidateResponse) *SocketTransportValidateResponse {
	if response.Error != nil {
		return &SocketTransportValidateResponse{
			Error: fmt.Sprintf("validator error: %s", response.Error.Error()),
		}
	}

	return &SocketTransportValidateResponse{
		UserData: hex.EncodeToString(response.UserData),
		Valid:    response.Valid,
	}
}
