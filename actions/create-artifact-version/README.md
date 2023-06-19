# Create artifact version

### Beskrivelse
Action som genererer en streng satt sammen av `1.1_${GIT_COMMIT_DATE}_${GIT_COMMIT_HASH}`.
`GIT_COMMIT_DATE` er på formatet `yyyymmdd_hhmm`.

Eksempel: `1.1_20230616.1102_9c8933b`

### Outputs
* `version` Den genererte strengen

### Eksempel på bruk
```yaml
steps:
  - name: Create artifact version
    id: artifact-version
    uses: navikt/sosialhjelp-ci/actions/create-artifact-version@master
```
#### Bruke output (avhengig av `steps.id`):
`${{ steps.artifact-version.outputs.version }}`