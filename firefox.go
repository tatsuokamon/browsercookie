package browsercookie

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/ini.v1"
)

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

func (f *firefox) findCookieFiles(yield func(string)bool) {
	profile, _ := f.findDefaultProfile()
	_path, _ := f.parseProfile(profile)
	cookieFile := fmt.Sprintf("%s/%s", _path, "cookies.sqlite")
	cookieDir := path.Dir(cookieFile)

	for _, file := range []string{
		path.Join(cookieDir, "sessionstrore-backups", "recovery.js"),
		path.Join(cookieDir, "sessionstrore-backups", "recovery.json"),
		path.Join(cookieDir, "sessionstrore-backups", "recovery.jsonlz4"),
		path.Join(cookieDir, "sessionstrore.js"),
		path.Join(cookieDir, "sessionstrore.json"),
		path.Join(cookieDir, "sessionstrore.jsonlz4"),
	} {
		_, err := os.Open(file)
		if err == nil || !os.IsNotExist(err) {
			if !yield(file) {
				return
			}
		}
	}
}

func (f *firefox)getCookies() {
	hasSessionFiles := false
	for cookieFile := range f.findCookieFiles {
		file, err := createLocalCopy(cookieFile)
	}
}
