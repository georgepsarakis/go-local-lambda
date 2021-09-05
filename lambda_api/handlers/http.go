package handlers

import (
	"github.com/djhworld/go-lambda-invoke/golambdainvoke"
	"github.com/georgepsarakis/go-local-lambda/lambda_api/configuration"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)


var logger = log.New(os.Stdout, "go-local-lambda-server", log.Lmsgprefix)


var invokeURLPattern = regexp.MustCompile("(?i)^/2015-03-31/functions/(?P<functionName>[a-z0-9_-]+)/invocations$")

func functionNameFromURL(requestURL string) string {
	matches := invokeURLPattern.FindStringSubmatch(requestURL)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}


type InvokeHandler struct {
	Configuration *configuration.LocalLambdaConfiguration
}

// POST /2015-03-31/functions/:functionName/invocations
func (h *InvokeHandler) Run(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_, _ = w.Write([]byte("Only POST method is allowed"))
		return
	}

	functionName := functionNameFromURL(r.URL.RequestURI())
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		logError(err)
		return
	}
	input := golambdainvoke.Input{
		Payload: payload,
		Port: int(h.Configuration.FindPort(functionName)),
	}
	response, err := golambdainvoke.Run(input)
	if response != nil {
		_, err := w.Write(response)
		logError(err)
	}
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		logError(err)
	}
}

func logError(e error) {
	if e != nil {
		logger.Println(e)
	}
}
