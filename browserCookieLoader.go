package browsercookie

import "net/http"

func newBrowserCookieLoader() {
}

type browserCookieLoader struct {
	CookieFiles []string
}

func (b *browserCookieLoader) findCookieFilesIter() (func(func(string) bool), error) {
	return func(func(string) bool) {}, ErrNotImplemented
}

func (b *browserCookieLoader) getCookiesIter() (func(func(*http.Cookie) bool), error) {
	return func(func(*http.Cookie) bool) {}, ErrNotImplemented
}

func (b *browserCookieLoader) load() ([]*http.Cookie, error) {
	return []*http.Cookie{nil}, ErrNotImplemented
}
