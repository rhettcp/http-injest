package main

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
}

func NewHandler(ctx context.Context) (*Handler, error) {

	return &Handler{}, nil
}

func (h *Handler) Netog(w http.ResponseWriter, r *http.Request) {
	validateRequests.Inc()

	// Get Request from http
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !h.authValid(r.Header.Get("Authorization")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Info("Request Received for Netog!")

	for k, vals := range r.Header {
		log.Info("Header: ", k, ": ", strings.Join(vals, ", "))
	}

	log.Info("Body: ", string(requestBody))

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) authValid(authString string) bool {
	if !strings.HasPrefix(authString, "Basic ") {
		return false
	}

	authString = strings.Replace(authString, "Basic ", "", -1)
	baseDecode, err := base64.StdEncoding.DecodeString(authString)
	if err != nil {
		log.Error("authValid, decode: ", err)
		return false
	}

	baseString := string(baseDecode)
	if !strings.Contains(baseString, ":") {
		return false
	}

	creds := strings.Split(baseString, ":")

	if len(creds) != 2 {
		return false
	}

	if creds[0] != "netog" || creds[1] != "netog11" {
		log.Warn("Invalid auth: ", creds)
		return false
	}
	return true
}
