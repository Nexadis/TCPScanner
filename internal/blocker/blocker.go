package blocker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/Nexadis/TCPTools/internal/blocker/config"
	"github.com/Nexadis/TCPTools/internal/logger"
)

type (
	addrlist []string
)

type Blocker struct {
	blacklist addrlist
	whitelist addrlist
	c         *config.Config
}

func New(c *config.Config) (*Blocker, error) {
	blacklist, err := ReadList(c.BlackList)
	if err != nil {
		return nil, err
	}
	whitelist, err := ReadList(c.WhiteList)
	if err != nil {
		return nil, err
	}
	return &Blocker{
		blacklist: blacklist,
		whitelist: whitelist,
		c:         c,
	}, nil
}

func (bl *Blocker) Run() error {
	return http.ListenAndServe(bl.c.Address, http.HandlerFunc(
		WithLog(bl.Block),
	))
}

func (bl *Blocker) Block(w http.ResponseWriter, r *http.Request) {
	if bl.IsBlocked(r) || !bl.IsAllowed(r) {
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

func ReadList(fileWithList string) (addrlist, error) {
	if fileWithList == "" {
		return addrlist{}, nil
	}
	blocklist, err := os.ReadFile(fileWithList)
	if err != nil {
		return nil, fmt.Errorf("addrlist %w", err)
	}
	return strings.Split(strings.Trim(string(blocklist), "\n"), "\n"), nil
}

func (bl *Blocker) IsBlocked(r *http.Request) bool {
	if len(bl.blacklist) == 0 {
		return false
	}
	hostname := r.URL.Hostname()
	for _, blocked := range bl.blacklist {
		if strings.Contains(hostname, blocked) {
			logger.Log.Infof("Blocked %v", blocked)
			return true
		}
	}
	return false
}

func (bl *Blocker) IsAllowed(r *http.Request) bool {
	if len(bl.whitelist) == 0 {
		return true
	}
	hostname := r.URL.Hostname()
	for _, allowed := range bl.whitelist {
		allowedURL, _ := url.Parse(allowed)
		if hostname == allowedURL.Hostname() {
			return true
		}
	}

	logger.Log.Infof("Is not allowed %v", hostname)
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
