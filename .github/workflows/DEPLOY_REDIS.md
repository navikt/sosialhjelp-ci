# Deployer REDIS til aktuelle system

Deployer REDIS til system angitt i input. 
K8s-ressursene er hardkodet, og krever like filnavn i repo som 
bruker workflowen.

### Inputs (alle er default `false`)
* `to-prod-fss` Til prod-fss
* `to-dev-fss` Til dev-fss
* `to-dev-gcp` Til dev-gcp 

Eksempel på bruk:
```yaml
name: Auto-deploy redis dev
on:
  push: # Deploys automatically to dev-fss and dev-gcp if there is a change in any redis files.
    paths:
      - "nais/redis-config.yml"
      - "nais/redisexporter.yml"
      - ".github/workflows/autodeploy_dev_redis.yml"
    branches-ignore:
      - 'master'
    tags-ignore:
      - "**" # Don't build any tags

jobs:
  deploy-redis:
    name: 'Deploy Redis to Dev'
    uses: navikt/sosialhjelp-ci/.github/workflows/deploy_redis.yml@v2
    secrets: inherit
    with:
      to-dev-fss: true
      to-dev-gcp: true
```