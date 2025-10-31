package browsercookie

import (
	"testing"
)

func TestSafari(t *testing.T) {
	s := newSafari()
	fIter, err := s.findCookieFilesIter()
	if err != nil {
		t.Log(err)
	}

	for file := range fIter {
		t.Log(file)
	}

	cIter, err := s.getCookiesIter()
	if err != nil {
		t.Log(err)
	}

	for _ = range cIter {}
}
