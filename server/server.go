package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stlasos/cryptocompare-collector/api"
	"github.com/stlasos/cryptocompare-collector/config"
	"github.com/stlasos/cryptocompare-collector/db"
)

type Server struct {
	api        *api.Api
	repo       *db.CurrencyExchangeRepository
	httpServer *http.Server
}

func NewServer(a *api.Api, r *db.CurrencyExchangeRepository) *Server {
	return &Server{
		api:  a,
		repo: r,
	}
}

func (s *Server) Init(port int) {
	m := http.NewServeMux()

	m.HandleFunc("/service/price", s.HandleRequest)

	s.httpServer = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: m}

	log.Info("Server started")

	err := s.httpServer.ListenAndServe()

	if err != nil {
		log.Panic("Error starting server")
	}
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	fsyms, fsymsExists := r.URL.Query()["fsyms"]
	tsyms, tsymsExists := r.URL.Query()["tsyms"]
	if fsymsExists && tsymsExists {
		resp, err := s.api.GetByCurrenciesRaw(fsyms[0], tsyms[0])
		if err == nil {
			w.WriteHeader(200)
			w.Write(resp)
			return
		} else {
			s.generateResponse(w, fsyms[0], tsyms[0])
			return
		}
	}
	w.WriteHeader(400)
}

func (s *Server) generateResponse(w http.ResponseWriter, fsyms string, tsyms string) {
	res, err := s.repo.GetByCurrenciesList(strings.Split(fsyms, ","), strings.Split(tsyms, ","))
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var response struct {
		RAW     map[string]map[string]db.CurrencyExchange        `json:"RAW"`
		DISPLAY map[string]map[string]db.CurrencyExchangeDisplay `json:"DISPLAY"`
	}

	response.RAW = make(map[string]map[string]db.CurrencyExchange)
	response.DISPLAY = make(map[string]map[string]db.CurrencyExchangeDisplay)

	for _, e := range res {
		if _, ok := response.RAW[e.FromCurrency]; !ok {
			response.RAW[e.FromCurrency] = make(map[string]db.CurrencyExchange)
		}
		response.RAW[e.FromCurrency][e.ToCurrency] = e.Ce

		if _, ok := response.DISPLAY[e.FromCurrency]; !ok {
			response.DISPLAY[e.FromCurrency] = make(map[string]db.CurrencyExchangeDisplay)
		}
		response.DISPLAY[e.FromCurrency][e.ToCurrency] = e.CeDisplay
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		log.WithField("err", err).Warn("Error encoding server response")
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
	w.Write(jsonResponse)
}

func (s *Server) InitBgCollector(conf config.Collector) {
	go func() {
		ticker := time.NewTicker(conf.Scheduler)
		log.Info("Background collector started")
		for {
			<-ticker.C
			res, err := s.api.GetByCurrencies(conf.Currencies.From, conf.Currencies.To)
			if err != nil {
				continue
			}
			s.repo.Store(res.Exchanges)
		}
	}()
}
