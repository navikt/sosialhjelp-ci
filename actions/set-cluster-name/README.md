# Set Cluster Name

Pr. nå brukes både `mock.yml` og `dev-gcp.yml` mot samme cluster.
Denne pakker den logikken inn i en action, og returnerer
riktig cluster basert på input.

### Inputs
* `resource-name` Navn på k8s-ressurs-fila

### Outputs
* `cluster-name` Navn på cluster utledet fra input

Eksempel på bruk:
```yaml
jobs:
  set-cluster:
    name: 'Sett cluster basert på config-file'
    runs-on: ubuntu-latest
    outputs:
      cluster-name: ${{ steps.set-cluster-name.outputs.cluster-name }}
    steps:
      - name: Sett cluster basert på config-file
        id: set-cluster-name
        uses: navikt/sosialhjelp-ci/actions/set-cluster-name@v2
        with:
          resource-name: ${{ inputs.config-file-name }}
```