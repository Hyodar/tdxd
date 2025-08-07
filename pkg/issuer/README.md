# Issuer Package

The issuer package provides interfaces and implementations for issuing TDX (Trusted Domain Extensions) attestation documents.

## Overview

This package handles the creation of attestation documents that cryptographically bind user-provided data with hardware-based attestation evidence. The attestation documents can be used to prove that specific data was processed within a trusted execution environment.

## Interface

```go
type Issuer interface {
    Start(ctx context.Context) error
    Issue(ctx context.Context, req *api.IssueRequest) *api.IssueResponse
}
```

## Available Implementations

### Azure Issuer
- **Type**: `azure`
- **Description**: Production implementation that interfaces with Azure Confidential Computing's attestation service
- **Config**: No additional configuration required (uses ambient Azure credentials)
- **Use Case**: Production environments running on Azure confidential VMs with TDX support

### Simulator Issuer
- **Type**: `simulator`
- **Description**: Mock implementation for development and testing
- **Config**: No configuration required
- **Use Case**: Local development, testing, and environments without TDX hardware

## Usage Example

```yaml
# In config.yaml
issuer:
  type: azure  # or "simulator" for testing
```
