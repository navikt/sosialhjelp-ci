# Create artifact version

### Beskrivelse
Action som genererer en streng satt sammen av `yyyymmdd.hhmm`, 
pluss første segment av en uuid.

### Outputs
* `version` Strengen som er generert.

### Eksempel på bruk
```yaml
steps:
  - name: Create artifact version
    id: artifact-version
    uses: navikt/sosialhjelp-ci/actions/create-artifact-version@master
```
Bruke output (avhengig av `steps.id`):
`${{ steps.artifact-version.outputs.version }}`