package douban

import (
	"fmt"
	cookiejar2 "github.com/juju/persistent-cookiejar"
	"net/http"
	"net/url"
	"strings"
)

const (
	reqPhoneCode = "https://accounts.douban.com/j/mobile/login/request_phone_code" //post
	//ck=&area_code=%2B86&number=12312332112
	verPhoneCode = "https://accounts.douban.com/j/mobile/login/verify_phone_code" //post
	//ck=&area_code=%2B86&number=&code=5753&remember=false&ticket=
)

type VerifyCode func(code string) error

func LoginByPhone(num string) (VerifyCode, error) {
	val := url.Values{
		"ck":        {},
		"area_code": {"+86"},
		"number":    {num},
	}
	//resp, err := http.PostForm(reqPhoneCode, val)
	tmpJar, _ := cookiejar2.New(nil)
	mclient := &http.Client{
		Jar: tmpJar,
	}
	req, _ := http.NewRequest("POST", reqPhoneCode, strings.NewReader(val.Encode()))
	//resp, err := mclient.PostForm(reqPhoneCode, val)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	resp, err := mclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}
	return func(code string) error {
		val.Add("code", code)
		val.Add("remember", "false")
		val.Add("ticket", "")
		req, _ = http.NewRequest("POST", verPhoneCode, strings.NewReader(val.Encode()))
		//resp, err := mclient.PostForm(verPhoneCode, val)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", userAgent)
		resp, err := mclient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		//io.Copy(os.Stdout, resp.Body)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("code might error")
		}
		client = mclient
		//parse, err := url.Parse("https://douban.com")
		//fmt.Println(mclient.Jar.Cookies(parse))
		return nil
	}, nil
}
