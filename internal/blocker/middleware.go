package blocker

import (
	"net/http"
	"net/http/httputil"

	"github.com/Nexadis/TCPTools/internal/logger"
)

func WithLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dumped, err := httputil.DumpRequest(r, true)
		if err == nil {
			logger.Log.Infoln(string(dumped))
		}
		h.ServeHTTP(w, r)
	}
}
