# Reusable Workflows

## - CodeQL
### Beskrivelse
Kan kjøres som en uavhengig jobb. Sjekker ut kode, setter opp og bygger java, 
og gjør selve kodeskanningen.

### Permissions
Trenger følgende permissions fra caller-workflow:
* `actions: read`
* `security-events: write` For å generere security-alerts
* `contents: read`

### Inputs
* `githubUser` For å sjekke ut repo. Default: `x-access-token`

### Eksempel på bruk
```yaml
jobs:
  analyze_code:
    name: Analyze Code
    permissions:
      actions: read
      security-events: write
      contents: read
    uses: navikt/sosialhjelp-ci/.github/workflows/codeql_java.yml@master
    with:
      githubUser: x-access-token
```

## - Dependency Submission Gradle
### Beskrivelse
Hvis avhengigheten ikke er søkbar under `<repo> -> Insight -> Dependency Graph`.
Kan kjøre som en uavhengig jobb for å sende inn avhengigheter ved bruk av 
Dependency Submission API. For gradle-prosjekter.

### Permissions
Trenger følgende permissions fra caller-workflow:
* `contents: write`

### Inputs
* `githubUser` For å sjekke ut repo. Default: `x-access-token`

### Eksempel på bruk
```yaml
jobs:
  dependency_submission:
    name: Dependency Submission Gradle
    permissions:
      contents: write
    uses: navikt/sosialhjelp-ci/.github/workflows/dependency_submission_gradle.yml@master
    with:
      githubUser: x-access-token
```