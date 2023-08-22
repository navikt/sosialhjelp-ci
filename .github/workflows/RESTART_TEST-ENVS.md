# Restarter aktuelle test-system

Reusable workflow som deployer på nytt til alle testsystemer.
Brukes for det meste med `On: Schedule`. Det vil sørge for at alle 
testsystemer kjører med master-branch igjen.

Workflow er delt opp i 4 jobber, hvor jobb nummer 1 henter ut riktig
docker tag basert på docker tag fra siste commit på gjeldende branch.

De 3 resterende er deploy for hvert testsystem, og vil kjøre parallelt.
(Hvis de er satt til `true` i inputs).

### Inputs (alle har default false)
* `to-mock` Skal det deployes til mock (dev-gcp)
* `to-dev-gcp` Testsystem dev-gcp
* `to-dev-fss` Testsystem dev-fss

Eksempel på bruk:
```yaml
jobs:
  restart-testenvs:
    name: 'Restart test-envs'
    uses: navikt/sosialhjelp-ci/.github/workflows/restart_test-envs.yml@v2
    secrets: inherit
    with:
      to-mock: true
      to-dev-fss: true
```