# Kodeskanning med CodeQL
Gjør alle nødvendige steg for å kjøre en kodeskanning for 
Java/Kotlin.

Eksempel på bruk:
```yaml
jobs:
  analyze_code:
    name: Analyze Code
    permissions:
      actions: read
      security-events: write
      contents: read
    uses: navikt/sosialhjelp-ci/.github/workflows/codeql_java.yml@v2
```