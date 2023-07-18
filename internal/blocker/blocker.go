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
	for _, blocked := range bl.blocklist {
		if r.Host == blocked {
			logger.Log.Infof("Blocked %v", blocked)
			http.Error(w, fmt.Sprintf("It's blocked %s=%s", r.Host, blocked), http.StatusLocked)
			return
		}
	}
	req, err := http.NewRequest(r.Method, r.RequestURI, r.Body)
	for k, v := range r.Header {
		if strings.Contains(k, "Proxy") {
			continue
		}
		req.Header.Add(k, v[0])
	}
	defer r.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf(`error while redirect request: %v`, err), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	dumped, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return
	}
	logger.Log.Infof(
		"Send request:\n%v", string(dumped),
	)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`error while redirect request: %v`, err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	w.Write(body)
}

func ReadBlocklist(fileBlocklist string) (blocklist, error) {
	blocklist, err := os.ReadFile(fileBlocklist)
	if err != nil {
		return nil, fmt.Errorf("blocklist %w", err)
	}
	return strings.Split(strings.Trim(string(blocklist), "\n"), "\n"), nil
}
