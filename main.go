package browsercookie

import "net/http"

func Firefox() ([]*http.Cookie, error) {
	f := newFireFox()
	return f.load()
}
