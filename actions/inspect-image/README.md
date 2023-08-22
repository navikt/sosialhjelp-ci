# Inspect Image 

Liten greie som brukes til å sjekke om image med aktuell tag
finnes ifra før.

### Inputs
* `docker-tag` Aktuell tag på image

### Outputs
* `image-manifest` Hvis _**manifest unknown**_, finnes ikke imaget

Eksempel på bruk:
```yaml
- name: 'Check if Image Exists'
  id: inspect-image
  if: inputs.build-always == false
  uses: navikt/sosialhjelp-ci/actions/inspect-image@v2
  with:
    docker-tag: ${{ steps.create-docker-tag.outputs.docker-tag }}

- name: 'If Needed - Release tag and Build Image'
  if: steps.inspect-image.outputs.image-manifest == 'manifest unknown'
  uses: navikt/sosialhjelp-ci/actions/build-for-deploy-kotlin@v2
  with:
    artifact-version: ${{ steps.artifact-version.outputs.version }}
```