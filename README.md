# 🦕 Nessie

(work in progress) CLI for the NS API.
A toy project to experiment with go.

# Requirements

- go
- A valid NS API key, subscribed to the "Ns-App" product.
  See [here](https://apiportal.ns.nl/startersguide).

# Setup

```bash
go build
```

# Run

```bash
./nessie -station Delft
```

# Development guide

## Debugging with delve

```bash
go get github.com/go-delve/delve/cmd/dlv
dlv exec ./nessie -- [command line options]
```
