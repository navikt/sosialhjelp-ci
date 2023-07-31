# Create and release tag

### Beskrivelse
Generisk action for å generere og release tag i kjørende repo.
Bruker `create-artifact-version`-action for å lage en tag basert på dato og klokkeslett.

### Inputs
* `prefix` F.eks. 'sosialhjelp-soknad' bruker dev/prod-prefix før generert tag

### Outputs
* `artifact-version` Generert tag/version

### Eksempel på bruk
```yaml
steps:
  - uses: navikt/sosialhjelp-ci/actions/create-and-release-tag@test_prefix
    id: artifact-version
    with:
      prefix: dev-sbs-
```
Output kan brukes slik (avhengig av `steps.id`):
#### `${{ steps.{id}.outputs.artifact-version }}`