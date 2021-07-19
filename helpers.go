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
