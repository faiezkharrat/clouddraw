package uri

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"makeit.imfr.cgi.com/makeit2/scm/lino/pimo/pkg/model"
)

func Read(uri string) ([]model.Entry, error) {
	var result []model.Entry
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "file" {
		file, err := os.Open(u.Host + u.Path)
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}
		return result, nil
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		/* #nosec */
		rep, err := http.Get(uri)
		if err != nil {
			return nil, err
		}
		defer rep.Body.Close()
		scanner := bufio.NewScanner(rep.Body)
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}
		return result, nil
	}

	return nil, fmt.Errorf("Not a valid scheme")
}
