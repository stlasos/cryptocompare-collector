package main

import (
	"flag"

	"github.com/stlasos/cryptocompare-collector/api"
	"github.com/stlasos/cryptocompare-collector/config"
	"github.com/stlasos/cryptocompare-collector/db"
	"github.com/stlasos/cryptocompare-collector/server"
)

func main() {
	mg := flag.String("migrate", "data/migration.sql", "run migrations instead of launching app to setup database")
	configLocation := "./cryptocompare-collector.yml"
	if len(flag.Args()) > 0 {
		configLocation = flag.Arg(0)
	}
	flag.Parse()
	conf := config.NewConfig(configLocation)
	d := db.NewDB(conf.DB)
	mgSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "migrate" {
			mgSet = true
		}
	})
	if mgSet {
		if *mg == "" {
			*mg = "data/migration.sql"
		}
		d.Migrate(*mg)
		return
	}

	a := api.NewApi(conf.Collector.ServiceUrl)
	repo := db.NewCurrencyExchangeRepository(d)

	s := server.NewServer(a, repo)

	s.InitBgCollector(conf.Collector)
	s.Init(conf.Port)
}
