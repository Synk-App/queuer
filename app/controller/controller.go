package controller

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"synk/gateway/app/util"
	"time"

	"github.com/getsentry/sentry-go"
)

type Response struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	Info  any    `json:"info"`
	List  []any  `json:"list"`
}

type ResponseHeader struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

const AUTH_TIMEOUT = time.Second * 5
const SENTRY_LOG_TIMEOUT = time.Second * 5

func WriteErrorResponse(w http.ResponseWriter, response any, route string, message string, status int) {
	defer sentry.Flush(SENTRY_LOG_TIMEOUT)

	util.LogRoute(route, message)

	sentry.CaptureMessage("error(@queuer" + route + "): " + message)

	jsonResp, _ := json.Marshal(response)

	w.WriteHeader(status)
	w.Write(jsonResp)
}

func WriteSuccessResponse(w http.ResponseWriter, response any) {
	defer sentry.Flush(SENTRY_LOG_TIMEOUT)

	jsonResp, _ := json.Marshal(response)

	sentry.CaptureMessage("success(@queuer): " + string(jsonResp))

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func WriteStatusResponse(w http.ResponseWriter, response any) {
	defer sentry.Flush(SENTRY_LOG_TIMEOUT)

	jsonResp, _ := json.Marshal(response)

	sentry.CaptureMessage("status(@queuer): " + string(jsonResp))

	w.WriteHeader(http.StatusMultiStatus)
	w.Write(jsonResp)
}

func SetJsonContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func NewPublisherServiceClient() *http.Client {
	env := os.Getenv("ENV")

	if env == "production" {
		return &http.Client{
			Timeout: AUTH_TIMEOUT,
		}
	}

	caCert, caErr := os.ReadFile("/cert/rootCA.pem")
	if caErr != nil {
		util.LogRoute("/", "error reading root CA file: "+caErr.Error())
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   AUTH_TIMEOUT,
	}

	return client
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allowedOriginsMap = map[string]struct{}{
			strings.TrimSuffix(os.Getenv("QUEUER_ENDPOINT"), "/"): {},
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		origin := r.Header.Get("Origin")
		if _, ok := allowedOriginsMap[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
