# Cryptocompare collector

Service that collects data from cryptocompare API.

## Installation

Before running service you should run database migrations (default file location is [`data/migration.sql`](data/migration.sql))

``` shell
go build
# if current dir is project root, otherwise you should specify the location of migration file in the -migrate option
./cryptocompare-collector -migrate=
```

## Configuration

``` yaml
# application port (default 9821)
port: 9821
db:
  # db params
  host: 127.0.0.1
  port: 3306
  user: root
  password:
  db_name: cryptocompare
collector:
  # currencies list to collect
  from: 
    - BTC
    - XRP
  to:
    - USD
    - EUR
  # collector time interval (default 1m)
  schedule: 5s
  # service api url (default https://min-api.cryptocompare.com)
  service_url: "https://min-api.cryptocompare.com"
```

## Run

Default config path is ./cryptocompare-collector.yml (will be used if not specified directly)

``` shell
./cryptocompare-collector <path/to/config>
```