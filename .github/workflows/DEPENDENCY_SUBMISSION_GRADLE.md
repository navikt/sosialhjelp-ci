# Dependency Submission (Gradle)

Det er ikke alltid Dependabot klarer å detecte alle dependencies
automatisk. Denne workflowen henter ut dependencies fra build.gradle og 
sender det inn via Dependency Submission API.

Eksempel på bruk:
```yaml
jobs:
  dependency_submission:
    name: Dependency Submission Gradle
    permissions:
      contents: write
    uses: navikt/sosialhjelp-ci/.github/workflows/dependency_submission_gradle.yml@v2
```