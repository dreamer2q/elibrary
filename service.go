package main

import (
	"ebook/douban"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type service struct {
	db   *gorm.DB
	gin  *gin.Engine
	tmpl *template.Template
}

func (s *service) init() {
	s.initDB()
	s.initTmpl()
	s.initServer()
}

func (s *service) initTmpl() {
	s.tmpl = template.New("myTpl")
	s.tmpl = template.Must(s.tmpl.ParseGlob("tpl/*.html"))
}

func (s *service) initServer() {
	g := gin.Default()
	g.GET("/", s.index)
	g.GET("/index", s.index)
	g.GET("/book/:id", s.bookById)
	g.GET("/404", s.pageNotFound)
	g.Static("/static", "html/")
	g.HTMLRender = &render.HTMLDebug{
		Glob: "tpl/*.html",
	}
	s.gin = g
}

func (s *service) run() {
	log.Panicln(s.gin.Run())
}

func (s *service) initDB() {
	db, err := gorm.Open("mysql", "root:admin@/library?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&Book{})
	db.AutoMigrate(&BookTag{})
	db.AutoMigrate(&BookPath{})
	//db.Model(&Book{}).Related(&BookTag{}).Related(&BookPath{})
	s.db = db
	//s.scanBook(`F:\学习\电子书\ebook`)
}

var filter = map[string]bool{
	".pdf":  true,
	".epub": true,
	".azw":  true,
	".azw3": true,
	".mobi": true,
}

func (s *service) scanBook(dir string) {
	var pathCh = make(chan string, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for d := range pathCh {
				s.addBook(d)
			}
		}()
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filter[filepath.Ext(info.Name())] {
			fmt.Println(path)
			pathCh <- path
		}
		return nil
	})
	close(pathCh)
}

func (s *service) addBook(path string) error {
	ext := filepath.Ext(path)[1:]
	name := filepath.Base(path)
	name = name[:strings.IndexRune(name, '.')]

	//bookInfo, err := douban.BookSearch(name)
	bookInfo, err := douban.BookSuggestion(name)
	if err != nil {
		//log.Panic(err)
		log.Println(err, "sleeping for 0 min")
		//time.Sleep(10 * time.Minute)
		return err
	}
	//tmp := bookInfo.Books[0]
	tmp := bookInfo[0]
	books := make([]Book, 0)
	fmt.Println(tmp.Title, tmp.ID, tmp.AuthorName)
	s.db.Where(&Book{BookId: tmp.ID}).Find(&books)
	if len(books) != 0 {
		log.Println("already exist")
		return fmt.Errorf("book %s (%s) already exist", books[0].Title, books[0].BookId)
	}

	book, err := douban.BookInfo(tmp.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	b := book
	bookItem := &Book{
		BookId:      b.ID,
		Title:       b.Title,
		Author:      strings.Join(b.Author, ","),
		AuthorIntro: b.AuthorIntro,
		PublicDate:  parseTime(b.Pubdate),
		ImageUrl:    b.Image,
		Pages:       intWrapper(strconv.Atoi(b.Pages)),
		Publisher:   b.Publisher,
		Isbn:        b.Isbn13,
		Summary:     b.Summary,
		Catalog:     b.Catalog,
		Tags:        wrapBookTags(b.Tags),
		BookPaths: []BookPath{{
			Ext:      ext,
			Filepath: path,
		}},
	}

	s.db.Create(&bookItem)
	return nil
}

func (s *service) clearDB() {
	s.db.Unscoped().Delete(&Book{})
	s.db.Unscoped().Delete(&BookTag{})
	s.db.Unscoped().Delete(&BookPath{})
}
