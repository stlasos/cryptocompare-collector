package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/stlasos/cryptocompare-collector/db"
)

type Api struct {
	url string
}

func NewApi(url string) *Api {
	return &Api{url: url}
}

type ApiResponse struct {
	Exchanges []db.Exchange
}

func (a *Api) GetByCurrenciesRaw(fsyms string, tsyms string) ([]byte, error) {
	resp, err := http.Get(a.url + "/data/pricemultifull?fsyms=" + fsyms + "&tsyms=" + tsyms)

	if err != nil {
		log.WithField("err", err).Warn("Error during get currencies request")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.WithField("err", err).Warn("Error reading get currencies api response")
		return nil, err
	}

	return body, nil
}

func (a *Api) GetByCurrencies(from []string, to []string) (*ApiResponse, error) {
	fsyms := a.convertCurrenciesToGetParams(from)
	tsyms := a.convertCurrenciesToGetParams(to)
	body, err := a.GetByCurrenciesRaw(fsyms, tsyms)

	if err != nil {
		return nil, err
	}

	var decodedResp struct {
		RAW     map[string]map[string]db.CurrencyExchange        `json:"RAW"`
		DISPLAY map[string]map[string]db.CurrencyExchangeDisplay `json:"DISPLAY"`
	}

	err = json.Unmarshal(body, &decodedResp)

	if err != nil {
		log.WithField("err", err).Warn("Error while json decoding api response")
		return nil, err
	}

	formattedRes := ApiResponse{}

	for from, toMany := range decodedResp.RAW {
		for to, exchange := range toMany {
			d := db.CurrencyExchangeDisplay{}
			if eD, exists := decodedResp.DISPLAY[from][to]; exists {
				d = eD
			}
			e := db.Exchange{
				FromCurrency: from,
				ToCurrency:   to,
				Ce:           exchange,
				CeDisplay:    d,
			}
			formattedRes.Exchanges = append(formattedRes.Exchanges, e)
		}
	}

	return &formattedRes, nil
}

func (a *Api) convertCurrenciesToGetParams(currencies []string) string {
	res := ""
	for i, v := range currencies {
		res += v
		if i != len(currencies)-1 {
			res += ","
		}
	}
	return res
}
