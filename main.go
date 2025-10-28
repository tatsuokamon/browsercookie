package browsercookie

import "net/http"

func FireFox() ([]*http.Cookie, error) {
	f := newFireFox()
	return f.load()
}
