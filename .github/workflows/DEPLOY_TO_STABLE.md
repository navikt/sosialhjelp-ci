# Deploy Application
Felles-logikk for å deploye prod og prod-test.

Jobb 1 forbereder deploy. Denne jobben henter
ut siste tag (versjon fra commit) og bygger image.

Jobb 2 og 3 deployer imaget til prod og prod-test.

Eksempel på bruk (deploy til produksjon):
```yaml
jobs:
  deploy-to-stable:
    name: 'Deploy to prod og prod-test'
    permissions:
      packages: write
      contents: write
    uses: navikt/sosialhjelp-ci/.github/workflows/deploy_app.yml@v2
    secrets: inherit
```