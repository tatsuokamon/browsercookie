package browsercookie

import (
	"net/http"
	"os"
	"path"
	"runtime"
)

func newSafari() *safari {
	return &safari{browserCookieLoader{}}
}

type safari struct{ browserCookieLoader }

func (s *safari) String() string {
	return "safari"
}

func (s *safari) load() ([]*http.Cookie, error) {
	result := []*http.Cookie{}
	cIter, err := s.getCookiesIter()
	if err != nil {
		return result, err
	}

	for c := range cIter {
		result = append(result, c)
	}

	return result, nil
}

func (s *safari) findCookieFilesIter() (func(func(string) bool), error) {
	if runtime.GOOS != "darwin" {
		return func(func(string) bool) {}, ErrSafariOnlyOnOSX
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return func(func(string) bool) {}, err
	}
	cookiePath := path.Join(home, "Library/Cookies")

	return func(yield func(string) bool) {
		for _, file := range []string{
			cookiePath,
		} {
			_, err := os.Stat(file)
			if err == nil || !os.IsNotExist(err) {
				if !(yield(file)) {
					return
				}
			}
		}
	}, nil
}

func (s *safari) getCookiesIter() (func(func(*http.Cookie) bool), error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return func(func(*http.Cookie) bool) {}, err
	}
	filePath := path.Join(home, "Library/Cookies/Cookies.binarycookies")
	_, err = os.Stat(filePath)
	if err != nil {
		return func(func(*http.Cookie) bool) {}, err
	}

	return func(yield func(*http.Cookie) bool) {
		src, err := os.ReadFile(filePath)
		if err != nil {
			return 
		}

		var numPage int8
	}, nil
}
