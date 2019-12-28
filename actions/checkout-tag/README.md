# checkout-tag

Check out a specific git tag, and optionally assert that the tag is on origin/master

## Usage

### nais/deploy action

Deploy to prod-sbs:

```yaml
  deploy_prod:
    runs-on: ubuntu-latest
    if: github.event.action == 'deploy_prod_tag'
    steps:
    - uses: navikt/sosialhjelp-ci/actions/checkout-tag@master
      with:
        tag: ${{ github.event.client_payload.TAG }}
        vars: "nais/prod/default.json"
        assert_master: "true"
    - uses: nais/deploy/actions/deploy@v1
      env:
        APIKEY: ${{ secrets.NAIS_DEPLOY_APIKEY }}
        CLUSTER: prod-sbs
        RESOURCE: "nais.yaml"
        VARS: "nais/prod/default.json"
```

Deploy to dev-sbs:

```yaml
  deploy_miljo:
    runs-on: ubuntu-latest
    if: github.event.action == 'deploy_miljo_tag'
    steps:
    - uses: navikt/sosialhjelp-ci/actions/checkout-tag@master
      with:
        tag: ${{ github.event.client_payload.TAG }}
        vars: "nais/dev/${{ github.event.client_payload.MILJO }}.json"
    - uses: nais/deploy/actions/deploy@v1
      env:
        APIKEY: ${{ secrets.NAIS_DEPLOY_APIKEY }}
        CLUSTER: prod-sbs
        RESOURCE: "nais.yaml"
        VARS: "nais/dev/${{ github.event.client_payload.MILJO }}.json"
```

### deployment-cli Docker container

Deploy to prod-sbs:

```yaml
  deploy_prod:
    runs-on: ubuntu-latest
    if: github.event.action == 'deploy_prod_tag'
    container: navikt/deployment-cli:0.4.0
    steps:
    - uses: navikt/sosialhjelp-ci/actions/checkout-tag@master
      with:
        tag: ${{ github.event.client_payload.TAG }}
        assert_master: "true"
    - run: |
        deployment-cli deploy create \
          --cluster=prod-sbs --repository=${GITHUB_REPOSITORY} --team=digisos \
          -r=nais.yaml --var version=${VERSION} --vars=nais/prod/default.json
      env:
        VERSION: ${{ github.event.client_payload.TAG }}
        DEPLOYMENT_USERNAME: ${{ secrets.DEPLOYMENT_USERNAME }}
        DEPLOYMENT_PASSWORD: ${{ secrets.DEPLOYMENT_PASSWORD }}
```

Deploy to dev-sbs:

```yaml
  deploy_miljo:
    runs-on: ubuntu-latest
    if: github.event.action == 'deploy_miljo_tag'
    container: navikt/deployment-cli:0.4.0
    steps:
    - uses: navikt/sosialhjelp-ci/actions/checkout-tag@master
      with:
        tag: ${{ github.event.client_payload.TAG }}
    - run: |
        deployment-cli deploy create \
          --cluster=dev-sbs --repository=${GITHUB_REPOSITORY} --team=digisos \
          -r=nais.yaml --var version=${VERSION} --vars=nais/dev/${MILJO}.json
      env:
        VERSION: ${{ github.event.client_payload.TAG }}
        MILJO: ${{ github.event.client_payload.MILJO }}
        DEPLOYMENT_USERNAME: ${{ secrets.DEPLOYMENT_USERNAME }}
        DEPLOYMENT_PASSWORD: ${{ secrets.DEPLOYMENT_PASSWORD }}
```
