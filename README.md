# sosialhjelp-ci

Repository for felles logikk knyttet til Workflows for Team Digisos.
Dette repoet inneholder Reusable Workflows (RW) og Custom Actions (CA).

RW brukes typisk i andre repos for å slippe duplisere for mye kode pr. 
repo, og for at Workflows skal være like på tvers av repos for teamet.
RW må p.t. ligge under `.github/workflows`-mappen. Alle RW skal ha en 
README.md med navnet i store bokstaver. 

CA brukes for å samle de logiske stegene ved en operasjon i en action, 
slik at RW er delt i logiske (del-)steg slik at Workflow skal være 
enklere å forstå. F.eks. nødvendige steg for å bygge kode, eller bygge 
og pushe docker image kan være samlet i en action.
CA ligger under `actions`-mappen. Hver action må ha sin egen mappe, som
angir navnet på actionet. Logikken må ligge i en `action.yml` fil under 
denne mappen. For dokumentasjon ligger det en README.md i tilknytning
til hver action.

Les om Reusble Workflows og Custom Actions i GitHub-dokumentasjonen for
mer detaljert informasjon og eksempler.

## Henvendelser
Henvendelser kan sendes via Slack i kanalen #team_digisos.
