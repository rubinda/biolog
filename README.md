## BioLog

Spletna aplikacija, ki nudi podporo pri popisu vrst po Sloveniji.

## Namestitev

Za delovanje aplikacije potrebujete [golang](https://golang.org/dl/), [dep](https://github.com/golang/dep) in [PostgreSQL](https://www.postgresql.org/download/). Podatke za povezljivost na Postgres podatkovno bazo je potrebno dodati v `config\config.yaml`. V mapi `certs\` se morata nahajati tudi SSL certifikat in ključ. Za dostop do [Google APIs](https://developers.google.com/identity/protocols/OAuth2) storitev potrebujete tudi client ID in client secret.

Za vzpostavitev podatkovne baze uporabite `pg_restore` in datoteko `biolog.dump`:
```sh
$ pg_restore -U postgres --schema-only -d ime_baze biolog.dump
```

Za namestitev odvisnih paketov uporabite `dep`:
```sh
$ dep ensure
```
Aplikacijo lahko prevedete in poženete preko ukazne vrstice:
```sh
$ go run cmd/biolog/main.go
```

#### Opomba

Delovanje aplikacije je trenutno preverjeno le na operacijskem sistemu macOS.