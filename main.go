package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shadowtactical/shadow-common/env"
	"github.com/shadowtactical/shadow-common/logging"
	"github.com/shadowtactical/shadow-common/signals"
	"github.com/shadowtactical/shadow-common/status"
)

var (
	log = logging.Log

	validateRequests prometheus.Counter
	validFfl         prometheus.Counter
	invalidFfl       prometheus.Counter
)

func init() {
	envName := env.GetEnvironmentWithDefault("APP_ENV", "dev")
	namespace := fmt.Sprintf("sct_%s", envName)

	validateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "sct_gb_injest",
		Name:      "validate_requests",
		Help:      "counter for validate requests received",
	})
	validFfl = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "sct_gb_injest",
		Name:      "valid_ffl",
		Help:      "counter for valid ffl validations",
	})
	invalidFfl = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "sct_gb_injest",
		Name:      "invalid_ffl",
		Help:      "counter for invalid ffl validations",
	})
}

func main() {
	logging.SetLogLevel(env.GetLogLevel())
	log.Info("http-injest Starting...")
	prometheus.MustRegister(validateRequests, validFfl, invalidFfl)

	can := status.ServeStatusEndpoint(context.Background())

	h, err := NewHandler(context.Background())
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/v0/http-injest/netog", h.Netog)
	http.Handle("/", r)

	log.Info("http-injest Running...")
	http.ListenAndServe(":"+env.GetEnvironmentWithDefault("SERVICE_PORT", "8080"), r)
	log.Error(signals.SignalHandlerLoop())
	can()
}
