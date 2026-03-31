# basic-yaml

This fixture represents the smallest useful `nara` workspace: one schema and one YAML entry file.

It is useful for checking the happy path for config loading, schema discovery, targeted validation, and simple compilation with no references.

Sample commands from the repository root:

```bash
go run . --config fixtures/basic-yaml/nara.yaml list schemas
go run . --config fixtures/basic-yaml/nara.yaml list entries
go run . --config fixtures/basic-yaml/nara.yaml validate 'fixtures/basic-yaml/entries/*.note.yaml'
go run . --config fixtures/basic-yaml/nara.yaml compile 'fixtures/basic-yaml/entries/*.note.yaml' --format yaml --out /tmp/basic-yaml.yaml
```
