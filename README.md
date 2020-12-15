# sosialhjelp-ci-scripts
Hjelpe-script og verktøy som brukes til å bygge og kontrollere byggejobber for Digisos.

## Nedlasting
[Windows](https://github.com/navikt/sosialhjelp-ci/raw/master/cistatus/release/cistatus.exe)

[Linux](https://github.com/navikt/sosialhjelp-ci/raw/master/cistatus/release/linux.tar.gz)

[macOS](https://github.com/navikt/sosialhjelp-ci/raw/master/cistatus/release/osx.tar.gz)

### CLI

[Windows](https://github.com/navikt/sosialhjelp-ci/raw/master/deploy/deploy.exe)

[macOS](https://github.com/navikt/sosialhjelp-ci/raw/master/deploy/deploy)

[Linux](https://github.com/navikt/sosialhjelp-ci/raw/master/deploy/deploy-linux)

Husk å legge til Circle ci token i `~/.cistatus.json`

## Komme i gang

Ved første oppstart av applikasjonen vil man bli bedt om å legge inn et token i `.cistatus.json`-fila
i hjemmekatalogen sin. Dette tokenet opprettes på [CircleCI | Personal API Tokens](https://circleci.com/account/api),
og vil kun være synlig når du oppretter det, så det er lurt å lagre en kopi på et sikkert sted.

Applikasjonen brukes til å deploye prosjekter som har blitt migrert til CircleCI. I skrivende stund
gjelder dette applikasjonene:
- [sosialhjelp-soknad](https://circleci.com/gh/navikt/sosialhjelp-soknad)
- [sosialhjelp-soknad-api](https://circleci.com/gh/navikt/sosialhjelp-soknad-api)
- [sosialhjelp-innsyn](https://circleci.com/gh/navikt/sosialhjelp-innsyn)
- [sosialhjelp-innsyn-api](https://circleci.com/gh/navikt/sosialhjelp-innsyn-api)

For at de skal være tilgjengelige i applikasjonen, må man gå inn på prosjektene i CircleCI (lenkene over)
og trykke på den blå `Follow Prosjekt`-knappen oppe til høyre.

## Bruk: deploy-app

For å deploye den sist bygde versjonen av en applikasjon i et miljø, trykker man på knappen for det
miljøet man ønsker å deploye til på linja med applikasjonen man ønsker å deploye. Hvis du ikke er sikker
på om det er din versjon av applikasjonen som er bygd sist, kan du gå inn på prosjektsiden på CircleCI,
og sjekke. For å bygge din versjon på nytt, kan du trykke deg inn på den siste vellykkede jobben på din branch,
og trykke på `Rerun workflow`-knappen oppe til høyre. Når denne jobben er ferdig, er prosjektet klart
til deploy i ønsket miljø.

## Bruk: cli
La currect working directory være det repositoriet du ønsker å deploye. Pass på at HEAD er den commiten du 
ønsker å deploye. Skriv:

`deploy <miljø> <configfilnavn>`

der <miljø> er dev-sbs, dev-fss, dev-gcp, labs-gcp eller prod.

Configfilnavn er optional og trenger bare å brukes hvis filnavnet på configfila som skal brukes heter noe annet enn `miljø`. Eksempel på dette er `dev-sbs-intern` som brukes for intern-appene våre.

Pass på at image-taggen du ønsker å deploye er knyttet til commiten du står i. Det kan være du må kjøre en `git fetch`.

## Hvordan lage deploy-cli tool til ulike OS fra windows: 
```
set GOOS=windows
go build
```
```
set GOOS=darwin
go build 
```
```
set GOOS=linux
go build -o deploy-linux
```

## Henvendelser
Henvendelser kan sendes via Slack i kanalen #digisos.
