# Transport Package

The transport package provides the communication layer for receiving attestation requests and forwarding them to the appropriate issuer or validator.

## Overview

This package defines the transport interface and implementations for handling client connections and request routing. It acts as the bridge between external clients and the internal attestation services.

## Available Implementations

### Socket Transport
- **Type**: `socket`
- **Description**: Unix domain socket implementation for local IPC with optional systemd socket activation
- **Config Options**:

  **Option 1: Manual socket creation**
  ```yaml
  config:
    file_path: "/path/to/socket.sock"  # Required when systemd is false
    owner: "username"        # Optional: socket file owner
    group: "groupname"       # Optional: socket file group  
    perm: 0600              # Optional: socket file permissions (octal)
  ```

  **Option 2: Systemd socket activation**
  ```yaml
  config:
    systemd: true           # Use systemd socket activation
    # No other fields allowed when systemd is true
  ```

- **Use Case**: Local inter-process communication, containerized environments, systemd-managed services

## API Schema

The socket transport uses JSON-based messaging over Unix domain sockets.

### Request Format

All requests follow this general structure:
```json
{
    "type": "issue|validate",
    "data": {
        // Method-specific payload
    }
}
```

### Issue Method

**Request:**
```json
{
    "type": "issue",
    "data": {
        "userData": "68656c6c6f20776f726c64",  // hex-encoded user data
        "nonce": "0123456789abcdef"             // hex-encoded nonce
    }
}
```

**Response:**
```json
{
    "document": "7b2274797065223a2261747465737461...",  // hex-encoded attestation document
    "error": ""  // Empty if successful, error message if failed
}
```

### Validate Method

**Request:**
```json
{
    "type": "validate",
    "data": {
        "document": "7b2274797065223a2261747465737461...",  // hex-encoded attestation document
        "nonce": "0123456789abcdef"                          // hex-encoded nonce
    }
}
```

**Response:**
```json
{
    "userData": "68656c6c6f20776f726c64",  // hex-encoded extracted user data
    "valid": true,                          // validation result
    "error": ""                             // Empty if successful, error message if failed
}
```

## Usage Example

### Configuration Examples

**Manual socket:**
```yaml
# In config.yaml
transport:
  type: socket
  config:
    file_path: "/var/run/tdxd.sock"
    owner: "root"       # Optional
    group: "tdxd"       # Optional
    perm: 0660          # Optional
```

**Systemd socket activation:**
```yaml
# In config.yaml
transport:
  type: socket
  config:
    systemd: true
```

For systemd socket activation, create a socket unit file:
```ini
# /etc/systemd/system/tdxd.socket
[Unit]
Description=TDXD Socket

[Socket]
ListenStream=/var/run/tdxd.sock
SocketMode=0660
SocketUser=root
SocketGroup=tdxd

[Install]
WantedBy=sockets.target
```

### Client Example (Shell)
```bash
# Issue attestation
echo '{"type":"issue","data":{"userData":"48656c6c6f","nonce":"0123456789"}}' | \
  nc -U /var/run/tdxd.sock

# Validate attestation
echo '{"type":"validate","data":{"document":"...","nonce":"0123456789"}}' | \
  nc -U /var/run/tdxd.sock
```
