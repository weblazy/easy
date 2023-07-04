package eprometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunPrometheus(cfg *Config) {
	http.Handle(cfg.Path, promhttp.Handler())
	http.ListenAndServe(cfg.Port, nil)
}
