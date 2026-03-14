# pyzanode/shared

Shared Go types and config for PyzaNode controller and agent.

- `config` — controller config (data dir, secrets file, HTTP addr).
- `types` — structs used over the API and WebSocket (presets, metrics, etc.).

No binary to run. Other repos `go get github.com/pyzanode/shared` (or use a replace in go.mod when developing in-tree).

## License

See the [PyzaNode](https://github.com/PyzaNode/PyzaNode) repo for the project license.
