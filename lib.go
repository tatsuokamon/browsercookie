package browsercookie

import (
	"fmt"
	"io"
	"os"
)

func createLocalCopy(cookiePath, suffix string) (string, error) {
	src, err := os.Open(cookiePath)
	if err != nil {
		return "", err
	}

	tmpCookieFile, err := os.CreateTemp(".", fmt.Sprintf("*%s", suffix))
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
