package db

import (
	"bytes"

	log "github.com/sirupsen/logrus"
)

type CurrencyExchangeRepository struct {
	db *DB
}

func NewCurrencyExchangeRepository(d *DB) *CurrencyExchangeRepository {
	return &CurrencyExchangeRepository{db: d}
}

type CurrencyExchange struct {
	Change24Hour    float64 `json:"CHANGE24HOUR"`
	ChangePtc24Hour float64 `json:"CHANGEPCT24HOUR"`
	Open24Hour      float64 `json:"OPEN24HOUR"`
	Volume24Hour    float64 `json:"VOLUME24HOUR"`
	Volume24HourTo  float64 `json:"VOLUME24HOURTO"`
	Low24Hour       float64 `json:"LOW24HOUR"`
	High24Hour      float64 `json:"HIGH24HOUR"`
	Price           float64 `json:"PRICE"`
	Supply          float64 `json:"SUPPLY"`
	Mktcap          float64 `json:"MKTCAP"`
}

type CurrencyExchangeDisplay struct {
	Change24Hour    string `json:"CHANGE24HOUR"`
	ChangePtc24Hour string `json:"CHANGEPCT24HOUR"`
	Open24Hour      string `json:"OPEN24HOUR"`
	Volume24Hour    string `json:"VOLUME24HOUR"`
	Volume24HourTo  string `json:"VOLUME24HOURTO"`
	Low24Hour       string `json:"LOW24HOUR"`
	High24Hour      string `json:"HIGH24HOUR"`
	Price           string `json:"PRICE"`
	Supply          string `json:"SUPPLY"`
	Mktcap          string `json:"MKTCAP"`
}

type Exchange struct {
	FromCurrency string
	ToCurrency   string
	Ce           CurrencyExchange
	CeDisplay    CurrencyExchangeDisplay
}

func (cr *CurrencyExchangeRepository) Store(e []Exchange) bool {

	if len(e) > 0 {
		stmt := `insert into currency_exchanges (
		from_currency,
		to_currency,
		change24hour,
		change24hour_display,
		change_ptc24hour,
		change_ptc24hour_display,
		open24hour,
		open24hour_display,
		volume24hour,
		volume24hour_display,
		volume24hour_to,
		volume24hour_to_display,
		low24hour,
		low24hour_display,
		high24hour,
		high24hour_display,
		price,
		price_display,
		supply,
		supply_display,
		mktcap,
		mktcap_display
		) values `
		b := bytes.Buffer{}
		b.WriteString(stmt)
		for i := 0; i < len(e); i++ {
			b.WriteString(`(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
			if i != len(e)-1 {
				b.WriteString(",")
			}
		}
		stmt = b.String()
		args := make([]interface{}, 0)
		for _, v := range e {
			args = append(
				args,
				v.FromCurrency,
				v.ToCurrency,
				v.Ce.Change24Hour,
				v.CeDisplay.Change24Hour,
				v.Ce.ChangePtc24Hour,
				v.CeDisplay.ChangePtc24Hour,
				v.Ce.Open24Hour,
				v.CeDisplay.Open24Hour,
				v.Ce.Volume24Hour,
				v.CeDisplay.Volume24Hour,
				v.Ce.Volume24HourTo,
				v.CeDisplay.Volume24HourTo,
				v.Ce.Low24Hour,
				v.CeDisplay.Low24Hour,
				v.Ce.High24Hour,
				v.CeDisplay.High24Hour,
				v.Ce.Price,
				v.CeDisplay.Price,
				v.Ce.Supply,
				v.CeDisplay.Supply,
				v.Ce.Mktcap,
				v.CeDisplay.Mktcap,
			)
		}

		t, err := cr.db.Conn.Begin()

		if err != nil {
			log.WithField("err", err).Warn("Error opening db transaction")
			return false
		}

		_, err = t.Exec(`delete from currency_exchanges`)

		if err != nil {
			log.WithField("err", err).Warn("Error deleting previous exchanges")
			return false
		}

		_, err = t.Exec(stmt, args...)

		if err != nil {
			log.WithField("err", err).Warn("Error inserting new exchanges")
			return false
		}

		err = t.Commit()

		if err != nil {
			log.WithField("err", err).Warn("Error commiting db transaction")
			return false
		}
	}
	return true
}

func (cr *CurrencyExchangeRepository) GetByCurrenciesList(from []string, to []string) ([]Exchange, error) {
	q := `select
	from_currency,
	to_currency,
	change24hour,
	change24hour_display,
	change_ptc24hour,
	change_ptc24hour_display,
	open24hour,
	open24hour_display,
	volume24hour,
	volume24hour_display,
	volume24hour_to,
	volume24hour_to_display,
	low24hour,
	low24hour_display,
	high24hour,
	high24hour_display,
	price,
	price_display,
	supply,
	supply_display,
	mktcap,
	mktcap_display from currency_exchanges where from_currency in (`

	var buff bytes.Buffer

	buff.WriteString(q)

	for i := 0; i < len(from); i++ {
		buff.WriteString("?")
		if i != len(from)-1 {
			buff.WriteString(",")
		}
	}

	buff.WriteString(")")
	buff.WriteString(" and to_currency in (")

	for i := 0; i < len(to); i++ {
		buff.WriteString("?")
		if i != len(to)-1 {
			buff.WriteString(",")
		}
	}
	buff.WriteString(")")

	q = buff.String()

	var args = make([]interface{}, 0)

	for _, v := range from {
		args = append(args, v)
	}

	for _, v := range to {
		args = append(args, v)
	}

	rows, err := cr.db.Conn.Query(q, args...)

	if err != nil {
		log.WithField("err", err).Warn("Error selecting rows from exchanges")
		return nil, err
	}

	defer rows.Close()

	exchanges := make([]Exchange, 0)

	for rows.Next() {
		e := Exchange{}

		err = rows.Scan(
			&e.FromCurrency,
			&e.ToCurrency,
			&e.Ce.Change24Hour,
			&e.CeDisplay.Change24Hour,
			&e.Ce.ChangePtc24Hour,
			&e.CeDisplay.ChangePtc24Hour,
			&e.Ce.Open24Hour,
			&e.CeDisplay.Open24Hour,
			&e.Ce.Volume24Hour,
			&e.CeDisplay.Volume24Hour,
			&e.Ce.Volume24HourTo,
			&e.CeDisplay.Volume24HourTo,
			&e.Ce.Low24Hour,
			&e.CeDisplay.Low24Hour,
			&e.Ce.High24Hour,
			&e.CeDisplay.High24Hour,
			&e.Ce.Price,
			&e.CeDisplay.Price,
			&e.Ce.Supply,
			&e.CeDisplay.Supply,
			&e.Ce.Mktcap,
			&e.CeDisplay.Mktcap,
		)
		if err != nil {
			log.WithField("err", err).Warn("Error scanning rows from exchanges")
			return nil, err
		}

		exchanges = append(exchanges, e)
	}

	return exchanges, nil
}
