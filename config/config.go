package config

import (
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port      int       `yaml:"port"`
	DB        DB        `yaml:"db"`
	Collector Collector `yaml:"collector"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Database string `yaml:"db_name"`
	Password string `yaml:"password"`
}

type Collector struct {
	Currencies struct {
		From []string `yaml:"from"`
		To   []string `yaml:"to"`
	} `yaml:"currencies"`
	Scheduler  time.Duration `yaml:"schedule"`
	ServiceUrl string        `yaml:"service_url"`
}

func NewConfig(filename string) *Config {
	c := Config{}
	c.Load(filename)
	return &c
}

func (c *Config) Load(filename string) {
	confData, err := ioutil.ReadFile(filename)

	if err != nil {
		log.WithField("err", err).Panic("Error loading config file")
	}

	err = yaml.Unmarshal(confData, &c)

	if err != nil {
		log.WithField("err", err).Panic("Error converting config file to yaml")
	}

	if c.Port == 0 {
		c.Port = 9821
	}

	if c.DB.Host == "" {
		c.DB.Host = "localhost"
	}

	if c.DB.Port == 0 {
		c.DB.Port = 3306
	}

	if c.DB.User == "" {
		c.DB.User = "root"
	}

	if len(c.Collector.Currencies.From) == 0 {
		c.Collector.Currencies.From = []string{"BTC", "XRP", "ETH", "BCH", "EOS", "LTC", "XMR", "DASH"}
	}

	if len(c.Collector.Currencies.To) == 0 {
		c.Collector.Currencies.To = []string{"USD", "EUR", "GBP", "JPY", "RUR"}
	}

	if c.Collector.Scheduler == 0 {
		c.Collector.Scheduler = time.Minute
	}

	if c.Collector.ServiceUrl == "" {
		c.Collector.ServiceUrl = "https://min-api.cryptocompare.com"
	}

	log.WithField("location", filename).Info("Config loaded")
}
