package errorcodeutil

import (
	"log"
	"net/http"

	"github.com/savannahghi/serverutils"
)

// ReportErr writes the indicated error to supplied response writer and also logs it
func ReportErr(w http.ResponseWriter, err error, status int) {
	if serverutils.IsDebug() {
		log.Printf("%s", err)
	}
	serverutils.WriteJSONResponse(w, ErrorMap(err), status)
}

// ErrorMap turns the supplied error into a map with "error" as the key
func ErrorMap(err error) map[string]string {
	errMap := make(map[string]string)
	errMap["error"] = err.Error()
	return errMap
}

// RespondWithError writes an error response
func RespondWithError(w http.ResponseWriter, code int, err error) {
	errMap := ErrorMap(err)
	errBytes, err := json.Marshal(errMap)
	if err != nil {
		errBytes = []byte(fmt.Sprintf("error: %s", err))
	}
	RespondWithJSON(w, code, errBytes)
}

// RespondWithJSON writes a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(payload)
	if err != nil {
		log.Printf(
			"unable to write payload `%s` to the http.ResponseWriter: %s",
			string(payload),
			err,
		)
	}
}