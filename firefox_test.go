package browsercookie

import (
	"strings"
	"testing"
)

func TestFireFox(t *testing.T) {
	f := newFireFox()
	fIter, err := f.findCookieFilesIter()
	if err != nil {
		t.Log(err)
	}

	for file := range fIter {
		t.Log(file)
	}

	cIter, err := f.getCookiesIter()
	if err != nil {
		t.Log(err)
	}

	for c := range cIter {
		if strings.Contains(c.Domain, "utol") {
			t.Log(c)
		}
	}
}
