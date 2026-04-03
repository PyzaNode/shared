# pyzanode/shared

Shared Go types and config for PyzaNode controller and agent.

- `config`: controller config (data dir, secrets file, HTTP addr).
- `types`: structs used over the API and WebSocket (presets, metrics, etc.).

No binary to run. Other repos `go get github.com/pyzanode/shared` (or use a replace in go.mod when developing in-tree).

## Release label

**`VERSION`** in this repo is a single line (for example `Beta-0.3.0`). Build scripts (`scripts/build-all.*`), the site (`landing` Vite + `npm run releases`), and embedded controller/agent strings all read **`shared/VERSION`** when you lay repos out as siblings next to `shared/`. That matches split GitHub repos: every component checkout can see the same file without a fake workspace root repo.

## License

See the [project license](https://github.com/PyzaNode/.github/blob/main/profile/LICENSE).
