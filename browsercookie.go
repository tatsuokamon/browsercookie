package browsercookie

import "net/http"

func Firefox() ([]*http.Cookie, error) {
	f := newFireFox()
	return f.load()
}

func Safari() ([]*http.Cookie, error) {
	s := newSafari()
	return s.load()
}
