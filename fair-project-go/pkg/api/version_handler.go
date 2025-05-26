package api

import (
	"net/http"

	"github.com/sadeq/fair-project-go/pkg/errors"
	"github.com/sadeq/fair-project-go/pkg/version"
)

// VersionHandler handles requests to the /version endpoint.
// It returns the current version of the application.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, errors.New(errors.ErrMethodNotAllowed, "Only GET method is allowed for version endpoint"))
		return
	}

	RespondJSON(w, http.StatusOK, version.Map())
}
