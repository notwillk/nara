# nara

`nara` is a CLI for validating, resolving, and compiling entity graphs from YAML and JSON files.

It is built around a small workspace model:

- `nara.yaml` defines schema discovery, entry resolution, metadata keys, and path aliases.
- `schemas/*.cue` contains CUE schemas used to validate entities.
- entry files such as `entries/hello.note.yaml` or `entries/release.task.json` hold source entities.
- `$ref` links let one entity reference another, and `nara` resolves those links into a compiled graph.

In practice, `nara` helps with a few concrete jobs:

- validating targeted files before committing a change
- linting an entire workspace without mutating it
- listing discovered schemas and entries
- scaffolding a new project
- compiling resolved entities to YAML, JSON, or SQLite

The [`fixtures/`](./nara/fixtures) directory contains representative sample workspaces you can use to try the CLI.

## Overview

Typical workflow:

1. Define one or more schemas in CUE.
2. Create entity files whose filenames and metadata match the configured pattern.
3. Use `nara validate` or `nara lint` to confirm the workspace is valid.
4. Use `nara compile` to emit resolved output for downstream tools.

A few examples from the repository root:

```bash
go run . --config fixtures/basic-yaml/nara.yaml list schemas
go run . --config fixtures/reference-graph/nara.yaml lint
go run . --config fixtures/mixed-json-yaml/nara.yaml compile 'fixtures/mixed-json-yaml/entries/*.*' --format json --out /tmp/nara.json
```

## Developing

The recommended development environment is the repo devcontainer. It installs the toolchain used by this project, including Go, `just`, `checksy`, `watchexec`, Docker-in-Docker support, the `devcontainer` CLI, and the other utilities expected by CI.

### VS Code

Open the repository in VS Code and choose `Dev Containers: Reopen in Container`. On first start, the container runs `just doctor --fix`, and on subsequent starts it runs `just doctor`.

### Devcontainer CLI

If you prefer the CLI:

```bash
devcontainer up --workspace-folder .
devcontainer exec --workspace-folder . just doctor
devcontainer exec --workspace-folder . just test
```

If you want an interactive shell inside the container:

```bash
devcontainer exec --workspace-folder . bash
```

## Justfile

The main developer workflows are exposed through `just`.

### `just build`

Compile the CLI to `bin/nara`.

```bash
just build
VERSION=v0.1.0 just build
```

### `just build --watch`

Rebuild automatically on file changes.

```bash
just build --watch
```

### `just clean`

Print the normal clean step.

```bash
just clean
```

### `just clean --deep`

Remove generated artifacts, including files under `.schemas/`.

```bash
just clean --deep
```

### `just doctor`

Check the local environment and tool availability.

```bash
just doctor
```

### `just doctor --fix`

Run the same environment checks and apply supported fixes.

```bash
just doctor --fix
```

### `just format`

Run the repository formatting workflow.

```bash
just format
```

### `just format --check`

Check formatting without modifying files.

```bash
just format --check
```

### `just static`

Run static checks.

```bash
just static
```

### `just static --fix`

Run static checks and apply supported fixes.

```bash
just static --fix
```

### `just test`

Run Go tests and the devcontainer feature test.

```bash
just test
```

If you only want Go tests and want to skip the devcontainer feature test:

```bash
SKIP_DEVCONTAINER_FEATURE_TEST=1 just test
```

### `just test --watch`

Re-run tests on file changes.

```bash
just test --watch
```

### `just release [part]`

Create and push a release version using `release.sh`. Supported parts are `patch`, `minor`, and `major`.

```bash
just release
just release minor
```

### `just help`

Print the built-in command summary.

```bash
just help
```

## Useful CLI Commands

Once the project is built, or when using `go run .`, the core CLI commands are:

- `nara init`
- `nara list schemas`
- `nara list entries`
- `nara validate`
- `nara lint`
- `nara compile`
- `nara format`

Example:

```bash
go run . --config fixtures/aliased-paths/nara.yaml validate 'fixtures/aliased-paths/entries/*.product.yaml'
go run . --config fixtures/aliased-paths/nara.yaml compile 'fixtures/aliased-paths/entries/*.product.yaml' --format sqlite --out /tmp/aliased-paths.db
```
