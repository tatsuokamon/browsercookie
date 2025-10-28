package browsercookie

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"gopkg.in/ini.v1"
)

func newFireFox() *firefox {
	return &firefox{browserCookieLoader{}}
}

type firefox struct{ browserCookieLoader }

func (f *firefox) String() string {
	return "firefox"
}

func (f *firefox) parseProfile(profile string) (string, error) {
	profileDir := path.Dir(profile)
	cfg, err := ini.Load(profile)

	filePath := ""

	if err != nil {
		return "", err
	}
	for _, section := range cfg.Sections() {
		if strings.HasPrefix(section.Name(), "Install") {
			if section.HasKey("Default") {
				key, _ := section.GetKey("Default")
				defaultPath := key.String()
				filePath = filepath.Join(profileDir, defaultPath)
				break
			}
		}
	}

	return filePath, nil
}

func (f *firefox) findDefaultProfile() (string, error) {
	if runtime.GOOS == "darwin" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s/%s", home, "Library/Application Support/Firefox/profiles.ini"), nil
	}

	return "", ErrNotImplemented
}

func (f *firefox) findCookieFilesIter() (func(func(string) bool), error) {
	profile, err := f.findDefaultProfile()
	if err != nil {
		return func(func(string) bool) {}, err
	}
	_path, err := f.parseProfile(profile)
	fmt.Println(_path)
	if err != nil {
		return func(func(string) bool) {}, err
	}
	cookieFile := path.Join(_path, "cookies.sqlite")
	cookieDir := path.Dir(cookieFile)

	return func(yield func(string) bool) {
		for _, file := range []string{
			path.Join(cookieDir, "sessionstore-backups", "recovery.js"),
			path.Join(cookieDir, "sessionstore-backups", "recovery.json"),
			path.Join(cookieDir, "sessionstore-backups", "recovery.jsonlz4"),
			path.Join(cookieDir, "sessionstore.js"),
			path.Join(cookieDir, "sessionstore.json"),
			path.Join(cookieDir, "sessionstore.jsonlz4"),
			cookieFile,
		} {
			_, err := os.Stat(file)
			if err == nil || !os.IsNotExist(err) {
				if !yield(file) {
					return
				}
			}
		}
	}, nil
}

func (f *firefox) getCookiesIter() (func(func(*http.Cookie) bool), error) {
	// hasSessionFiles := false

	cIter, err := f.findCookieFilesIter()
	if err != nil {
		return func(func(*http.Cookie) bool) {}, err
	}

	return func(yield func(*http.Cookie) bool) {
		for cookieFile := range cIter {
			file, err := createLocalCopy(cookieFile, path.Ext(cookieFile))
			if err != nil {
				continue
			}
			defer os.Remove(file)

			if strings.HasSuffix(file, ".sqlite") {
				db, err := sql.Open("sqlite3", file)
				if err != nil {
					fmt.Println(err)
					continue
				}
				rows, err := db.Query("SELECT host, path, isSecure, expiry, name, value FROM moz_cookies")
				if err != nil {
					fmt.Println(err)
					continue
				}

				for rows.Next() {
					var cookie http.Cookie
					var i int64
					err := rows.Scan(&cookie.Domain, &cookie.Path, &cookie.Secure, &i, &cookie.Name, &cookie.Value)
					if err != nil {
						continue
					}
					cookie.Expires = time.Unix(0, i*int64(time.Millisecond))
					if !yield(&cookie) {
						return
					}
				}
			} else {
				_, err := os.ReadFile(file)
				if err != nil {
					continue
				}
			}
		}
	}, nil
}
