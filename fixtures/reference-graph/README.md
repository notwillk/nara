# reference-graph

This fixture represents a small multi-entity graph where entries reference each other through `$ref` fields.

It is useful for exercising recursive resolution, edge creation, and compilation of several related YAML files into one resolved graph.

Sample commands from the repository root:

```bash
go run . --config fixtures/reference-graph/nara.yaml list entries
go run . --config fixtures/reference-graph/nara.yaml lint
go run . --config fixtures/reference-graph/nara.yaml validate 'fixtures/reference-graph/entries/*.person.yaml'
go run . --config fixtures/reference-graph/nara.yaml compile 'fixtures/reference-graph/entries/*.person.yaml' --format json --out /tmp/reference-graph.json
```
