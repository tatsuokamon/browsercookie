package browsercookie

import "net/http"

func newBrowserCookieLoader() {
}

type browserCookieLoader struct {
	CookieFiles []string
}

func (b *browserCookieLoader) findCookieFiles() ([]string, error) {
	return []string{}, ErrNotImplemented
}

func (b *browserCookieLoader) getCookies() ([]*http.Cookie, error) {
	return []*http.Cookie{nil}, ErrNotImplemented
}

func (b *browserCookieLoader) load() ([]*http.Cookie, error) {
	return []*http.Cookie{nil}, ErrNotImplemented
}
