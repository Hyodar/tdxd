# TDXS - TDX Quote Service

A high-performance service for issuing and validating Intel TDX (Trust Domain Extensions) quotes across various deployment environments.

## Overview

TDXS provides a robust service for managing TDX attestation workflows, enabling secure quote generation and validation for confidential computing environments. Built on attestation infrastructure from [Constellation](https://github.com/edgelesssys/constellation/) (v2.23.1, pre-license change).

## Requirements

- Intel TDX-enabled hardware
- Linux kernel with TDX support (5.19+)
- Go 1.21 or higher (for building from source)

## Installation

```bash
git clone https://github.com/Hyodar/tdxs.git
cd tdxs
make build
make install
``

Run the daemon:

```bash
tdxs start --config /etc/tdxs/config.toml
```

## License

This project is licensed under the Gnu Affero General Public License 3.0 - see the [LICENSE](LICENSE) file for details.
