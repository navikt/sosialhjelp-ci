# Deploy Application
Felles-logikk for å deploye en app med NAIS Deploy.

Jobb 1 forbereder deploy. Denne jobben henter
ut siste tag (versjon fra commit) og sjekker om image finnes
ifra før. Hvis ikke bygges det.

Jobb 2 kjører NAIS Deploy med inputs og tag fra jobb 1.

### Inputs:
* `cluster-name` Hvilket nais-cluster skal appen deployes til.
* `resource-name` Hva heter k8s resource fila som skal brukes.
* `resource-folder` Hvilken folder: (dev/prod). Default: `dev`
* `build-always` Hvis image alltid skal bygges. Default: `false`

Eksempel på bruk (deploy til produksjon):
```yaml
jobs:
  deploy-to-prod:
    name: 'Deploy to prod-fss'
    permissions:
      packages: write
      contents: write
    uses: navikt/sosialhjelp-ci/.github/workflows/deploy_app.yml@v2
    with:
      cluster-name: 'prod-fss'
      resource-folder: 'prod'
      resource-name: 'prod-fss'
      build-always: true
    secrets: inherit
```
Eksempel (deploy til dev):
```yaml
  deploy-to-dev:
    name: 'Deploy to dev'
    permissions:
      packages: write
      contents: write
    uses: navikt/sosialhjelp-ci/.github/workflows/deploy_app.yml@v2
    with:
      cluster-name: <cluster name>
      resource-name: <resource name>
    secrets: inherit
```