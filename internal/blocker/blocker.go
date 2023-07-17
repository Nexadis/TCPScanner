package blocker

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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
	return http.ListenAndServe(bl.Addr, http.HandlerFunc(bl.Block))
}

func (bl *Blocker) Block(w http.ResponseWriter, r *http.Request) {
	for _, blocked := range bl.blocklist {
		if r.Host == blocked {
			http.Error(w, fmt.Sprintf("It's blocked %s=%s", r.Host, blocked), http.StatusLocked)
			return
		}
	}
	w.Write([]byte(`request allowed`))
}

func ReadBlocklist(fileBlocklist string) (blocklist, error) {
	blocklist, err := os.ReadFile(fileBlocklist)
	if err != nil {
		return nil, fmt.Errorf("blocklist %w", err)
	}
	return strings.Split(strings.Trim(string(blocklist), "\n"), "\n"), nil
}
