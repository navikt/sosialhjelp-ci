# Build kotlin

### Beskrivelse
Samler stegene som gjøres for å bygge et generisk kotlin(eller java)-prosjekt.
Denne inkluderer også `ktlintCheck`

### Inputs
* `task-type` (Ikke required) Kan brukes hvis man ikke ønsker å kjøre full build (ikke tester).
Se under.

### Eksempel på bruk
```yaml
steps:
  - name: Build kotlin, and run tests/lint
    uses: navikt/sosialhjelp-ci/actions/build-kotlin@master
    with:
      task-type: 'assemble'
```