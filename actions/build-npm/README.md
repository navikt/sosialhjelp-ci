# Build npm

### Beskrivelse
Samler alle steg for å bygge et genereisk npm-prosjekt, samt kjøre tester.

### Inputs
* `reader-token` Trenger secrets.READER_TOKEN for å hente dependencies fra Github Package Registry

### Eksempel på bruk
```yaml
steps:
  - uses: navikt/sosialhjelp-ci/actions/build-npm@master
    with:
      reader-token: ${{ secrets.READER_TOKEN }}
```