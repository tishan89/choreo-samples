# Repository Guidelines

## Project Structure & Module Organization
`main.go` is the application entrypoint. Runtime code lives under `internal/`: `config/` loads environment variables and shared types, `routes/` registers Gin handlers, `service/` wraps calls to the reading-list API, and `templates/` contains server-rendered `.tmpl` views such as `dashboard.tmpl` and `addNewBook.tmpl`. Root-level files include `go.mod` and `go.sum` for dependencies, `README.md` for deployment and local setup, and `sample.env` as the configuration template.

## Build, Test, and Development Commands
Use standard Go tooling from the repository root.

- `go run main.go`: start the web app directly for local development.
- `go build -o reading-list-app`: build a local binary matching the README flow.
- `./reading-list-app`: run the compiled binary.
- `go test ./...`: run all package tests once dependencies are available.
- `go fmt ./...`: apply canonical Go formatting before opening a PR.

Set `API_URL` and optionally `PORT` in `.env` before running locally. Copy values from `sample.env` and point `API_URL` at the deployed reading-list service.

## Coding Style & Naming Conventions
Follow idiomatic Go: tabs for indentation, exported names in `PascalCase`, unexported helpers in `camelCase`, and package names in short lowercase form (`config`, `routes`, `service`). Keep handlers and service methods narrowly scoped. Prefer `go fmt` output over manual formatting and preserve the existing package-level organization when adding code.

## Testing Guidelines
This repository currently has no checked-in `*_test.go` files. New behavior should include package-level Go tests placed next to the code they exercise, using the standard `testing` package and names like `TestFetchBooks` or `TestHandleDashboard`. Cover route validation, config parsing, and service error handling where practical. Run `go test ./...` before submitting changes.

## Commit & Pull Request Guidelines
Recent history uses short, imperative commit subjects such as `Update go mod files` and `Fix ...`, with merge commits handled by the platform. Prefer concise imperative messages, optionally scoped to the affected area. Pull requests should explain the user-visible change, note any config updates (`API_URL`, `PORT`, authentication behavior), link the relevant issue when applicable, and include screenshots for template or page changes.

## Configuration & Deployment Notes
Managed authentication is part of the app contract. If you add protected routes, update the deployment configuration documented in `README.md` so Choreo protects the same paths the code expects.
