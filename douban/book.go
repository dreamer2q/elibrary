package douban

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	bookInfoById = "/v2/book/"                                 //:id
	bookSearch   = "/v2/book/search"                           //q=""
	bookSuggest  = "https://book.douban.com/j/subject_suggest" //q=""
)

var (
	baseUrl = [...]string{
		"https://douban-api.uieee.com",
		"https://douban.uieee.com",
		"https://douban-api.now.sh",
		"https://douban-api.zce.now.sh",
		"https://douban-api-git-master.zce.now.sh",
	}
	index int
)

func getBaseUrl() string {
	index = (index + 1) % len(baseUrl)
	return baseUrl[index]
}

func BookSuggestion(q string) (SuggestResult, error) {
	_, body, err := get(bookSuggest, url.Values{"q": []string{q}})
	if err != nil {
		return nil, err
	}
	if bytes.Contains(body, []byte("登录")) {
		return nil, fmt.Errorf("%s", body)
	}
	retResult := make(SuggestResult, 0, 5)
	_ = json.Unmarshal(body, &retResult)
	if len(retResult) == 0 {
		return nil, fmt.Errorf("emtpy result")
	}
	return retResult, nil
}

func BookSearch(q string) (SearchResult, error) {
	searchResult := new(SearchResult)
	err := do(bookSearch, url.Values{"q": []string{q}}, searchResult)
	if err == nil && searchResult.Count == 0 {
		return *searchResult, fmt.Errorf("empty result")
	}
	return *searchResult, err
}

func BookInfo(id string) (Book, error) {
	book := new(Book)
	err := do(bookInfoById+id, nil, book)
	return *book, err
}

func do(relatedUrl string, query url.Values, ret interface{}) error {
	resp, body, err := get(getBaseUrl()+relatedUrl, query)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		errMsg := new(ErrorMsg)
		_ = json.Unmarshal(body, errMsg)
		return fmt.Errorf("do %s: %s", errMsg.Request, errMsg.Msg)
	}
	return json.Unmarshal(body, ret)
}
