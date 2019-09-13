# sosialhjelp-ci-scripts
Hjelpe-script og verktøy som brukes til å bygge og kontrollere byggejobber for Digisos.

## Nedlasting
[Windows](https://github.com/navikt/sosialhjelp-ci/raw/master/cistatus/release/cistatus.exe)

[Linux](https://github.com/navikt/sosialhjelp-ci/raw/master/cistatus/release/linux.tar.gz)

[macOS](https://github.com/navikt/sosialhjelp-ci/raw/master/cistatus/release/osx.tar.gz)

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

## Bruk

For å deploye den sist bygde versjonen av en applikasjon i et miljø, trykker man på knappen for det
miljøet man ønsker å deploye til på linja med applikasjonen man ønsker å deploye. Hvis du ikke er sikker
på om det er din versjon av applikasjonen som er bygd sist, kan du gå inn på prosjektsiden på CircleCI,
og sjekke. For å bygge din versjon på nytt, kan du trykke deg inn på den siste vellykkede jobben på din branch,
og trykke på `Rerun workflow`-knappen oppe til høyre. Når denne jobben er ferdig, er prosjektet klart
til deploy i ønsket miljø.

## Henvendelser
Henvendelser kan sendes via Slack i kanalen #digisos.
