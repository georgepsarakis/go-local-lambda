package handlers

import (
	"encoding/json"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/awslambda"
	"github.com/georgepsarakis/go-local-lambda/local_lambda/configuration"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"regexp"
)

var invokeURLPattern = regexp.MustCompile("(?i)^/2015-03-31/functions/(?P<functionName>[a-z0-9_-]+)/invocations$")

func functionNameFromURL(requestURL string) string {
	matches := invokeURLPattern.FindStringSubmatch(requestURL)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

type InvokeHandler struct {
	Configuration configuration.Configuration
	Logger *zap.Logger
}

// POST /2015-03-31/functions/:functionName/invocations
func (h *InvokeHandler) Run(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Only POST method is allowed"))
		return
	}

	functionName := functionNameFromURL(r.URL.RequestURI())
	payload, err := ioutil.ReadAll(r.Body)
	h.Logger.Info("received request", zap.String("payload", string(payload)))
	if logError(h.Logger, w, err) {
		return
	}
	req := awslambda.Request{
		Payload: payload,
	}
	client, err := awslambda.NewClient(functionName, h.Configuration.FindPort(functionName))
	if logError(h.Logger, w, err) {
		return
	}
	response, err := client.Invoke(req)
	if logError(h.Logger, w, err) {
		return
	}

	if len(response.Payload) > 0 {
		_, err := w.Write(response.Payload)
		if logError(h.Logger, w, err) {
			return
		}
	}
	if response.Error != nil {
		invocationError, err := json.Marshal(response.Error)
		if logError(h.Logger, w, err) {
			return
		}
		if _, err := w.Write(invocationError); err != nil {
			logError(h.Logger, w, err)
		}
	}
}

func logError(logger *zap.Logger, w http.ResponseWriter, e error) bool {
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error: "+e.Error())); err != nil {
			logger.Error("response write error", zap.Error(err))
		}
		logger.Error("server error", zap.Error(e))
		return true
	}
	return false
}
