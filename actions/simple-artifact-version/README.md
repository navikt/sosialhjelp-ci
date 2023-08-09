# simple-artifact-version

### Beskrivelse
Generisk action for lage en artifact version basert på timestamp og hash fra commit.

### Outputs
* `artifact-version` Sammensatt versjon fra commit

### Eksempel på bruk
```yaml
steps:
  - uses: navikt/sosialhjelp-ci/actions/simple-artifact-version@master
    id: artifact-version
```
Output kan hentes fra variabel (avhengig av `steps.id`):
```yaml
run: echo ${{ steps.{id}.outputs.artifact-version }} >> $GITHUB_ENV
```