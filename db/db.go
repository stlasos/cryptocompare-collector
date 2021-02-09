package db

import (
	"database/sql"
	"io/ioutil"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/stlasos/cryptocompare-collector/config"
)

type DB struct {
	Conn *sql.DB
}

func NewDB(dbConf config.DB) *DB {
	d := DB{}
	d.Connect(dbConf)
	return &d
}

func (d *DB) Connect(dbConf config.DB) {
	connStr := dbConf.User
	if dbConf.Password != "" {
		connStr += ":" + dbConf.Password
	}
	connStr += "@(" + dbConf.Host + ":" + strconv.Itoa(dbConf.Port) + ")/" + dbConf.Database
	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		log.WithField("err", err).Panic("Can't open db connection")
	}
	d.Conn = conn
}

func (d *DB) Migrate(filename string) {
	migration, err := ioutil.ReadFile(filename)

	if err != nil {
		log.WithFields(map[string]interface{}{"err": err, "location": filename}).Panic("Can't open migration file")
	}

	_, err = d.Conn.Exec(string(migration))

	if err != nil {
		log.WithField("err", err).Panic("Error during database migration")
	}

	log.WithField("location", filename).Info("Migration completed sucessfully")
}
