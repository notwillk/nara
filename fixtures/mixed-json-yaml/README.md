# mixed-json-yaml

This fixture represents a workspace with mixed input formats: one YAML entry and one JSON entry under the same schema.

It is useful for checking extension handling, mixed-file validation, and compilation when references cross file formats.

Sample commands from the repository root:

```bash
go run . --config fixtures/mixed-json-yaml/nara.yaml list entries
go run . --config fixtures/mixed-json-yaml/nara.yaml lint
go run . --config fixtures/mixed-json-yaml/nara.yaml validate 'fixtures/mixed-json-yaml/entries/*.*'
go run . --config fixtures/mixed-json-yaml/nara.yaml compile 'fixtures/mixed-json-yaml/entries/*.*' --format yaml --out /tmp/mixed-json-yaml.yaml
```
