# Validator Package

The validator package provides interfaces and implementations for validating TDX attestation documents and extracting user data.

## Overview

This package handles the verification of attestation documents to ensure they were created by genuine TDX hardware and haven't been tampered with. After successful validation, it extracts the original user data that was embedded during issuance.

## Interface

```go
type Validator interface {
    Start(ctx context.Context) error
    Validate(ctx context.Context, req *api.ValidateRequest) *api.ValidateResponse
}
```

## Available Implementations

### Azure Validator
- **Type**: `azure`
- **Description**: Production implementation that validates Azure TDX attestation documents
- **Config**: Embeds the Constellation AzureTDX configuration
  ```yaml
  config:
    measurements:           # TPM measurements map
      15:
        expected: "hex-value"
        validationOpt: "WarnOnly|Enforce"
    qeSVN:                 # QE security version
      value: 2
      isLatest: false      # Optional: use latest version
    pceSVN:                # PCE security version
      value: 13
      isLatest: false
    teeTCBSVN:             # TEE TCB security version (hex)
      value: "02060000000000000000000000000000"
      isLatest: false
    qeVendorID:            # QE vendor ID (hex)
      value: "939a7233f79c4ca9940a0db3957f0607"
      isLatest: false
    mrSeam: "hex-value"    # Optional: MR_SEAM value (48 bytes hex)
    xfam:                  # XFAM field (8 bytes hex)
      value: "e742060000000000"
      isLatest: false
    intelRootKey: "-----BEGIN CERTIFICATE-----\n..."  # Intel root certificate
  ```
- **Use Case**: Production environments that need to verify Azure TDX attestation documents

### Simulator Validator
- **Type**: `simulator`
- **Description**: Mock implementation that validates test attestation documents
- **Config**: No configuration required
- **Use Case**: Local development and testing environments

## Usage Example

```yaml
# In config.yaml
validator:
  type: azure
  config:
    measurements:
      15:
        expected: "0000000000000000000000000000000000000000000000000000000000000000"
        validationOpt: "WarnOnly"
    qeSVN:
      value: 2
    pceSVN:
      value: 13
    teeTCBSVN:
      value: "02060000000000000000000000000000"
    qeVendorID:
      value: "939a7233f79c4ca9940a0db3957f0607"
    xfam:
      value: "e742060000000000"
    intelRootKey: |
      -----BEGIN CERTIFICATE-----
      MIICjjCCAjSgAwIBAgIUImUM1lqdNInzg7SVUr9QGzknBqwwCgYIKoZIzj0EAwIw
      ...
      -----END CERTIFICATE-----

# Or for simulator:
validator:
  type: simulator
  # No config needed
```
