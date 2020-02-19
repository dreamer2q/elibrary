package douban

import (
	"fmt"
	cookiejar2 "github.com/juju/persistent-cookiejar"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"
)

var (
	client *http.Client
)

func init() {
	jar, err := cookiejar2.New(nil)
	if err != nil {
		log.Panic(err)
	}
	client = &http.Client{
		Jar: jar,
	}
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Kill, os.Interrupt)
	go func() {
		k := <-kill
		_ = jar.Save()
		log.Println("signal", k)
		os.Exit(1)
	}()
}

func get(url string, query url.Values) (*http.Response, []byte, error) {
	urlQuery := fmt.Sprintf("%s?%s", url, query.Encode())
	req, _ := http.NewRequest("GET", urlQuery, nil)
	req.Header.Add("User-Agent", userAgent)
	//resp, err := (&http.Client{}).Do(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return resp, body, err
}
