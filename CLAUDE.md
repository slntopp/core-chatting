- Hard Rules
  - Keep changes scoped to this repository.
  - Do not commit changes unless explicitly requested.
  - Keep `go.mod` and `go.sum` in sync with Go dependency changes.
  - Use existing build files and package scripts before adding new tooling.
  - Treat `.proto` files and `buf*.yaml` as generation authority.
  - Do not hand-edit generated protobuf outputs when a `.proto` change is required.

- Authority & Links
  - Go module: `github.com/slntopp/core-chatting`.
  - `go.mod`.
  - `Dockerfile`.
  - `buf.yaml`.
  - `buf.gen.yaml`.
  - Plugin package scripts: `plugin/package.json`.

- Setup / Test
  - Use `go mod download` to fetch Go module dependencies.
  - Use `go test ./...` for Go package checks.
  - Use `cd plugin && pnpm install` before plugin package scripts.
  - No `test` script is defined in `plugin/package.json`.

- Workflow
  - `go mod download`
  - `go test ./...`
  - `buf generate`
  - `cd plugin && pnpm install`
  - `cd plugin && pnpm run build`

- Stop Conditions
  - Ask before using or changing secrets, credentials, or `.env` files.
  - Ask before broad regeneration or formatting that changes unrelated files.
  - Ask if required external services, credentials, or generator plugins are missing.
  - Refuse destructive git operations unless explicitly requested.
  - Omit uncertain repository rules instead of guessing.
