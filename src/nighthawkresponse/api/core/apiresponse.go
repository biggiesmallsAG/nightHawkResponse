package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nighthawkresponse/log"
	"strings"
)

// ApiResponse is comment Nighthawk Response API message
// sent by server
type ApiResponse struct {
	Response string      `json:"response"`
	Reason   string      `json:"reason"`
	Data     interface{} `json:"data"`
}

func ApiResponseMessage(response, reason string, data interface{}) string {
	apiresponse := ApiResponse{
		Response: response,
		Reason:   reason,
		Data:     data,
	}

	jsonResponse, _ := json.MarshalIndent(&apiresponse, " ", "  ")
	return string(jsonResponse)
}

func HttpResponseReturn(w http.ResponseWriter, r *http.Request, response string, reason string, data interface{}) {
	LogDebug(DEBUG, fmt.Sprintf("[+] %s %s, %s - %s ", r.Method, r.RequestURI, response, reason))

	// Log failed messages to nighthawk/log
	if strings.ToLower(response) == "failed" {
		log.LogMessage(r.RequestURI, "NOTICE", reason)
	}

	fmt.Fprintln(w, ApiResponseMessage(response, reason, data))
}
