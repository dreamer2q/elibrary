package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func (s *service) index(c *gin.Context) {
	type IndexTmpl struct {
		Books    []Book
		PrePage  string
		NextPage string
		Pages    []int16
	}
	pg := c.DefaultQuery("page", "0")
	page, err := strconv.ParseUint(pg, 10, 64)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		s.tmpl.ExecuteTemplate(c.Writer, "error", err)
		return
	}

	var counter int64
	var lastRecord = &Book{}
	s.db.Model(&Book{}).Count(&counter)
	s.db.Last(lastRecord)
	pagelen := counter / 20
	if counter%20 != 0 {
		pagelen++
	}
	if page > uint64(pagelen) {
		c.Writer.WriteHeader(http.StatusBadRequest)
		s.tmpl.ExecuteTemplate(c.Writer, "error", "bad payload")
		return
	}
	page = uint64(pagelen) - page - 1
	books := make([]Book, 0, 25)
	s.db.Model(&Book{}).Where("id < ? AND id >= ?", page*20+20, page*20).Find(&books)
	if len(books) == 0 {
		c.Writer.WriteHeader(http.StatusNotFound)
		s.tmpl.ExecuteTemplate(c.Writer, "404", nil)
		return
	}
	var prePage, nextPage string
	if page != 0 {
		prePage = fmt.Sprintf("/index?page=%d", page-1)
	}
	if page < uint64(pagelen) {
		nextPage = fmt.Sprintf("/index?page=%d", page+1)
	}
	var pages = make([]int16, pagelen, pagelen)
	for i := range pages {
		pages[i] = int16(i)
	}
	data := &IndexTmpl{
		Books:    books,
		PrePage:  prePage,
		NextPage: nextPage,
		Pages:    pages,
	}
	//c.Writer.WriteHeader(http.StatusOK)
	//s.tmpl.ExecuteTemplate(c.Writer, "index", &data)
	c.HTML(http.StatusOK, "index", &data)
}

func (s *service) bookById(c *gin.Context) {
	type BookItem struct {
		Book
		Prev, Next         string
		PrevName, NextName string
		Views              int
		Douban             string
	}
	id := c.Param("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "error", "missing id param")
		return
	}
	ID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error", err)
		return
	}
	var book Book
	s.db.Where(&Book{Model: gorm.Model{ID: uint(ID)}}).Limit(1).Find(&book)
	if book.ID == 0 {
		c.Writer.WriteHeader(http.StatusNotFound)
		s.tmpl.ExecuteTemplate(c.Writer, "404", nil)
		return
	}
	bookItem := &BookItem{
		Book:   book,
		Douban: fmt.Sprintf("https://book.douban.com/subject/%s/", book.BookId),
	}
	var preBook, nextBook Book
	s.db.Find(&preBook, book.ID-1)
	s.db.Find(&nextBook, book.ID+1)
	if preBook.ID != 0 {
		bookItem.Prev = strconv.Itoa(int(preBook.ID))
		bookItem.PrevName = preBook.Title
	}
	if nextBook.ID != 0 {
		bookItem.Next = strconv.Itoa(int(nextBook.ID))
		bookItem.NextName = nextBook.Title
	}
	c.Writer.WriteHeader(http.StatusOK)
	s.tmpl.ExecuteTemplate(c.Writer, "book", &bookItem)
}

func (s *service) pageNotFound(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotFound)
	s.tmpl.ExecuteTemplate(c.Writer, "404", "")
}
