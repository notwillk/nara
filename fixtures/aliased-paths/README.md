# aliased-paths

This fixture represents a workspace that uses a configured path alias (`~shared`) to resolve references outside the main `entries/` directory.

It is useful for checking alias-based resolution and compilation across multiple configured roots in one project.

Sample commands from the repository root:

```bash
go run . --config fixtures/aliased-paths/nara.yaml list schemas
go run . --config fixtures/aliased-paths/nara.yaml list entries
go run . --config fixtures/aliased-paths/nara.yaml validate 'fixtures/aliased-paths/entries/*.product.yaml'
go run . --config fixtures/aliased-paths/nara.yaml compile 'fixtures/aliased-paths/entries/*.product.yaml' --format sqlite --out /tmp/aliased-paths.db
```
