# Build for Deploy
Et typisk bygg for deploy. Inkluderer:
* Bygger java/kotlin uten tester (=assemble)
* Lager tag fra siste commit - og releaser
* Bygger og pusher Docker Image 
* Returnerer tag som ble laget

### Inputs
* `github-token`
* `build-image` Skal image bygges og pushes, eller finnes det ifra før?

### Outputs
* `image-tag` Endelig tag for docker image.

### Eksempel på bruk
```yaml
steps:
  - name: Build and Tag for Deploy
    id: build-and-release
    uses: navikt/sosialhjelp-ci/actions/build-for-deploy@master
    with:
      github-token: ${{ secrets.GITHUB_TOKEN }}
      build-image: ${{ inputs.build-image }}
```

Output kan brukes slik (avhengig av `steps.id`):
#### `${{ steps.{id}.outputs.image-tag }}`
