package browsercookie

import (
	"io"
	"os"
)

func createLocalCopy(cookiePath string) (string, error) {
	src, err := os.Open(cookiePath)
	if err != nil {
		return "", err
	}

	tmpCookieFile, err := os.CreateTemp(".", "")
	if err != nil {
		return "", err
	}
	defer func() {
		tmpCookieFile.Close()
	}()

	tmpCookieName := tmpCookieFile.Name()
	_, err = io.Copy(tmpCookieFile, src)

	if err != nil {
		return tmpCookieName, err
	}

	return tmpCookieName, nil
}
