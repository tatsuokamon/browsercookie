package browsercookie

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
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
	cookiePath := path.Join(home, "Library/Containers/com.apple.Safari/Data/Library/Cookies")

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
	filePath := path.Join(home, "Library/Containers/com.apple.Safari/Data/Library/Cookies/Cookies.binarycookies")
	_, err = os.Stat(filePath)
	if err != nil {
		return func(func(*http.Cookie) bool) {}, err
	}

	return parseSafariBinaryCookiesFileToCooiesIter(filePath)
}

func parseSafariBinaryCookiesFileToCooiesIter(file string) (func(func(*http.Cookie) bool), error) {
	return func(yield func(*http.Cookie) bool) {
		src, err := os.Open(file)
		if err != nil {
			fmt.Println(err, 0)
		}
		defer src.Close()

		src.Seek(4, io.SeekStart)
		numPages, err := unpack[int32](src, binary.BigEndian, 4)
		if err != nil {
			fmt.Println(err, 1)
		}

		pageSizes := []int32{}
		for range numPages {
			ps, err := unpack[int32](src, binary.BigEndian, 4)
			if err != nil {
				fmt.Println(err, 2)
			}
			pageSizes = append(pageSizes, ps)
		}

		pages := [][]byte{}
		for _, ps := range pageSizes {
			b := make([]byte, ps)
			_, err := src.Read(b)
			if err != nil {
				fmt.Println(err, 3)
			}
			pages = append(pages, b)
		}

		for _, page := range pages {
			pageReader := bytes.NewReader(page)
			pageReader.Seek(4, io.SeekStart)
			numCookies, err := unpack[int32](pageReader, binary.LittleEndian, 4)
			if err != nil {
				fmt.Println(err, 4)
			}

			cookieOffsets := []int32{}
			for range numCookies {
				cookieOffset, err := unpack[int32](pageReader, binary.LittleEndian, 4)
				if err != nil {
					fmt.Println(err, 5)
				}
				cookieOffsets = append(cookieOffsets, cookieOffset)
			}

			pageReader.Seek(4, io.SeekStart)

			for _, offset := range cookieOffsets {
				pageReader.Seek(int64(offset), io.SeekStart)
				cookieSize, err := unpack[int32](pageReader, binary.LittleEndian, 4)
				if err != nil {
					fmt.Println(err, 6)
				}

				cookieBytes := make([]byte, cookieSize)
				pageReader.Read(cookieBytes)
				cookieReader := bytes.NewReader(cookieBytes)

				cookieReader.Seek(4, io.SeekStart)

				flags, err := unpack[int32](cookieReader, binary.LittleEndian, 4)
				// ^^^ok^^^
				if err != nil {
					fmt.Println(err, 7)
				}
				cookieFlags := fromFlagsToCookieFlag(flags)
				cookieReader.Seek(4, io.SeekCurrent)

				var (
					urlOffset   int32
					nameOffset  int32
					pathOffset  int32
					valueOffset int32
				)

				for _, intPtr := range []*int32{
					&urlOffset,
					&nameOffset,
					&pathOffset,
					&valueOffset,
				} {
					value, err := unpack[int32](cookieReader, binary.LittleEndian, 4)
					if err != nil {
						fmt.Println(err, 8)
					}
					*intPtr = value
				}

				expiryFloat, err := unpack[float64](cookieReader, binary.LittleEndian, 8)
				if err != nil {
					fmt.Println(err, 9)
				}

				host := ""
				name := ""
				path := ""
				value := ""

				for strPtr, eachOffset := range map[*string]int32{
					&host:  urlOffset,
					&name:  nameOffset,
					&path:  pathOffset,
					&value: valueOffset,
				} {
					cookieReader.Seek(int64(eachOffset-4), io.SeekStart)
					result := ""

					for {
						b := make([]byte, 1)
						_, err := cookieReader.Read(b)
						if err != nil {
							fmt.Println(err, 10)
							break
						}

						var i int8
						err = binary.Read(bytes.NewReader(b), binary.LittleEndian, &i)
						if err != nil {
							fmt.Println(err, 11)
						}

						if i == 0 {
							break
						}

						result += string(b)
					}

					*strPtr = result
					// fmt.Println(result)
				}

				if !yield(&http.Cookie{
					Domain:  host,
					Name:    name,
					Path:    path,
					Value:   value,
					Secure:  cookieFlags,
					Expires: fromFloatToTime(expiryFloat),
				}) {
					fmt.Println("hello")
					return
				}
			}
		}
	}, nil
}

func fromFloatToTime(float float64) time.Time {
	float += 978307200
	return time.Unix(int64(float), 0)
}
func fromFlagsToCookieFlag(flags int32) bool {
	cookieFlags := false
	switch flags {
	case 0:
		{
			cookieFlags = false
		}
	case 1:
		{
			cookieFlags = true
		}
	case 4:
		{
			cookieFlags = false
		}
	case 5:
		{
			cookieFlags = true
		}
	default:
		{
			cookieFlags = false
		}
	}

	return cookieFlags
}

func unpack[T any](src io.Reader, endian binary.ByteOrder, byteLength int) (T, error) {
	b := make([]byte, byteLength)
	dummy := new(T)

	_, err := src.Read(b)
	if err != nil {
		return *dummy, err
	}

	err = binary.Read(bytes.NewReader(b), endian, dummy)
	if err != nil {
		return *dummy, err
	}

	return *dummy, err
}
