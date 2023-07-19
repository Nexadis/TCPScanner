package blocker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/Nexadis/TCPTools/internal/logger"
)

type blocklist []string

type Blocker struct {
	blocklist blocklist
	Addr      string
}

func New(fileBlocklist string, addr string) (*Blocker, error) {
	blocklist, err := ReadBlocklist(fileBlocklist)
	if err != nil {
		return nil, err
	}
	return &Blocker{
		blocklist: blocklist,
		Addr:      addr,
	}, nil
}

func (bl *Blocker) Run() error {
	return http.ListenAndServe(bl.Addr, http.HandlerFunc(
		WithLog(bl.Block),
	))
}

func (bl *Blocker) Block(w http.ResponseWriter, r *http.Request) {
	if bl.IsBlocked(r) {
		http.Error(w, fmt.Sprintf("Blocked host: %s", r.Host), http.StatusLocked)
		return
	}
	req, err := copyRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf(`error while redirect request: %v`, err), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`error while redirect request: %v`, err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func ReadBlocklist(fileBlocklist string) (blocklist, error) {
	blocklist, err := os.ReadFile(fileBlocklist)
	if err != nil {
		return nil, fmt.Errorf("blocklist %w", err)
	}
	return strings.Split(strings.Trim(string(blocklist), "\n"), "\n"), nil
}

func (bl *Blocker) IsBlocked(r *http.Request) bool {
	for _, blocked := range bl.blocklist {
		if r.Host == blocked {
			logger.Log.Infof("Blocked %v", blocked)
			return true
		}
	}
	return false
}

func copyRequest(r *http.Request) (*http.Request, error) {
	req, err := http.NewRequestWithContext(r.Context(), r.Method, r.RequestURI, r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	copyHeaders(req, r)
	dumped, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	logger.Log.Infof(
		"Send request:\n%v", string(dumped),
	)
	return req, nil
}

func copyHeaders(to, from *http.Request) {
	for k, v := range from.Header {
		if strings.Contains(k, "Proxy") {
			continue
		}
		to.Header.Add(k, v[0])
	}
}
