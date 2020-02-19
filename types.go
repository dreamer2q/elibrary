package main

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Book struct {
	gorm.Model
	BookId      string
	Title       string
	Author      string
	AuthorIntro string `gorm:"type:text"`
	PublicDate  time.Time
	ImageUrl    string
	Pages       int
	Publisher   string
	Isbn        string
	Summary     string     `gorm:"type:text"`
	Catalog     string     `gorm:"type:text"`
	Tags        []BookTag  `gorm:"ForeignKey:BookId;association_foreignKey:BookId"`
	BookPaths   []BookPath `gorm:"ForeignKey:BookId;association_foreignKey:BookId"`
}

type BookTag struct {
	gorm.Model
	BookId string
	Title  string
	Name   string
	Count  int
}

type BookPath struct {
	gorm.Model
	BookId   string
	Ext      string
	Filepath string
}
