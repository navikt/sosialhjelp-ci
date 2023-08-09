# Build and push docker image

### Beskrivelse
Samler alle steg som brukes for å bygge og pushe et docker image til ghcr.io. 
Dette inkluderer også login til ghcr.io, samt tagge imaget.

### Inputs
* `image-name` Navn på Docker Image
* `artifact-version` Tag. F.eks. fra 'create-artifact-version'-action.
* `github-token` Trenger github-token for innlogging mot ghcr.io

### Eksempel på bruk
```yaml
steps:
  - name: Create and release tag
    uses: navikt/sosialhjelp-ci/actions/build-and-push-docker-image@master
    with:
      image-name: ghcr.io/${{ github.repository }}/${{ github.event.repository.name }}
      github-token: ${{ secrets.GITHUB_TOKEN }}
      artifact-version: ${{ steps.artifact-version.outputs.artifact-version }}
```