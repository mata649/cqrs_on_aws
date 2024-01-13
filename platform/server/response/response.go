package response

import (
	"encoding/json"
	"net/http"
	"strings"
)

func WriteResponse(status int, content interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	switch content.(type) {
	case string:
		w.Write([]byte(`{"message":"` + content.(string) + `"`))
	default:
		resp, err := json.Marshal(content)
		if err != nil {
			resp = []byte(`{"message":"null"}`)
		}
		w.Write(resp)
	}
}

func WriteBindingErrorResponse(error string, w http.ResponseWriter) {

	errorsParsed := map[string]string{}
	errorSplited := strings.Split(error, ":")
	key := strings.Trim(errorSplited[0], " ")

	key = strings.ToLower(key[0:1]) + key[1:]
	errorsParsed[key] = errorSplited[1]

	resp, err := json.Marshal(errorsParsed)
	if err != nil {
		resp = []byte(`{"message":"request error"}`)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}
