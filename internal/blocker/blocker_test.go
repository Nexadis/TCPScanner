package blocker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testbl = addrlist{
	"google.com",
	"ya.ru",
	"12.163.34.11",
}

type proxyReq struct {
	name      string
	blocklist addrlist
	hostname  string
	isBlocked bool
}

var blackSuites = []proxyReq{
	{
		name:      "Valid url",
		blocklist: testbl,
		hostname:  "ya.ru",
		isBlocked: true,
	},
	{
		name:      "Invalid url",
		blocklist: testbl,
		hostname:  "yaru",
		isBlocked: false,
	},
	{
		name:      "Valid subdomain",
		blocklist: testbl,
		hostname:  "dsen.ya.ru.com",
		isBlocked: false,
	},
	{
		name:      "Valid subdomain",
		blocklist: testbl,
		hostname:  "dsen.ya.ru",
		isBlocked: true,
	},
	{
		name:      "Valid subdomain with port",
		blocklist: testbl,
		hostname:  "dsen.ya.ru:8080",
		isBlocked: true,
	},
	{
		name:      "Valid url",
		blocklist: testbl,
		hostname:  "ru",
		isBlocked: false,
	},
	{
		name:      "Long random invalid hostname",
		blocklist: testbl,
		hostname:  "aaksjdflasjdf;lkjsaldfkjalsk;jdfjhalkjher-[sfafj234ja;sdfkj]",
		isBlocked: false,
	},
	{
		name:      "Valid IP with port",
		blocklist: testbl,
		hostname:  "12.163.34.11:23432",
		isBlocked: true,
	},
}

func TestBlock(t *testing.T) {
	b := &Blocker{
		blacklist: testbl,
	}
	for _, test := range blackSuites {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://"+test.hostname, nil)
			t.Log(test.hostname)
			assert.Equal(t, test.isBlocked, b.IsBlocked(r))
		})
	}
}

var whiteSuites = []proxyReq{
	{
		name:      "Valid url",
		blocklist: testbl,
		hostname:  "ya.ru",
		isBlocked: false,
	},
	{
		name:      "Invalid url",
		blocklist: testbl,
		hostname:  "yaru",
		isBlocked: true,
	},
	{
		name:      "Valid subdomain",
		blocklist: testbl,
		hostname:  "dsen.ya.ru.com",
		isBlocked: true,
	},
	{
		name:      "Valid subdomain",
		blocklist: testbl,
		hostname:  "dsen.ya.ru",
		isBlocked: false,
	},
	{
		name:      "Valid subdomain with port",
		blocklist: testbl,
		hostname:  "dsen.ya.ru:8080",
		isBlocked: false,
	},
	{
		name:      "Valid url",
		blocklist: testbl,
		hostname:  "ru",
		isBlocked: true,
	},
	{
		name:      "Long random invalid hostname",
		blocklist: testbl,
		hostname:  "aaksjdflasjdf;lkjsaldfkjalsk;jdfjhalkjher-[sfafj234ja;sdfkj]",
		isBlocked: true,
	},
	{
		name:      "Valid IP with port",
		blocklist: testbl,
		hostname:  "12.163.34.11:23432",
		isBlocked: false,
	},
}

func TestAllower(t *testing.T) {
	b := &Blocker{
		whitelist: testbl,
	}
	for _, test := range whiteSuites {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://"+test.hostname, nil)
			assert.Equal(t, test.isBlocked, !b.IsAllowed(r))
		})
	}
}
