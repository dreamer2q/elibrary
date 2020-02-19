package douban

type SuggestResult []struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Pic        string `json:"pic"`
	AuthorName string `json:"author_name"`
	Year       string `json:"year"`
	Type       string `json:"type"`
	ID         string `json:"id"`
}

type SearchResult struct {
	Count int    `json:"count"`
	Start int    `json:"start"`
	Total int    `json:"total"`
	Books []Book `json:"books"`
}

type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      []string  `json:"author"`
	AuthorIntro string    `json:"author_intro"`
	Pubdate     string    `json:"pubdate"`
	Tags        []BookTag `json:"tags"`
	Rating      BookRate  `json:"rating"`
	Image       string    `json:"image"`
	Catalog     string    `json:"catalog"`
	Pages       string    `json:"pages"`
	Publisher   string    `json:"publisher"`
	Isbn13      string    `json:"isbn13"`
	Summary     string    `json:"summary"`
	Price       string    `json:"price"`
}

type BookTag struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

type BookRate struct {
	Max       int    `json:"max"`
	NumRaters int    `json:"numRaters"`
	Average   string `json:"average"`
	Min       int    `json:"min"`
}
type BookImage struct {
	Small  string `json:"small"`
	Large  string `json:"large"`
	Medium string `json:"medium"`
}

type ErrorMsg struct {
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
	Request string `json:"request"`
}
