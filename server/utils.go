package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/jinzhu/gorm"
)

//Account Account Model
type Account struct {
	gorm.Model
	Name       string `gorm:"not null;unique"`
	Login      string `gorm:"not null;unique"`
	Password   string
	IsAdmin    bool
	ProfilePic string `gorm:"default:'pp1'"`
}

type BlogConfig struct {
	gorm.Model
	Title       string
	Meta        string
	Grid        bool
	Comments    bool
	BlogPic     string
	Description string
	Url         string
}

type Link struct {
	gorm.Model
	Side  string
	Url   string
	Icon  string
	Label string
}

//Account Account Model
type AccountApi struct {
	ID         uint
	Login      string
	Name       string
	ProfilePic string
}

//Model for languges
type Language struct {
	Name   string
	Short  string
	PicUrl string
}

//Category cat
type Category struct {
	gorm.Model
	Name   string
	Parent *Category
	Trads  []CategoryTrad
}

type CategoryTrad struct {
	gorm.Model
	CategoryID int
	LanguageID string
	Name       string
}

type Theme struct {
	gorm.Model
	Name        string
	BackColor   string
	FontColor   string
	UnderColor  string
	SecondColor string
	Active      bool
}

type Post struct {
	gorm.Model

	Public   bool
	Category string
	Tags     []Tag `gorm:"many2many:post_tags;"`

	Comments    []Comment
	Date        time.Time
	Views       int
	UniqueViews int
	Contents    []ContentPost
	MainPicture string
}
type Page struct {
	gorm.Model

	Public bool
	// Tags     []Tag `gorm:"many2many:post_tags;"`
	// Comments    []Comment
	Date        time.Time
	Views       int
	UniqueViews int
	Contents    []ContentPost
}

type ContentPost struct {
	gorm.Model
	PostID     uint
	PageID     uint
	URL        string `gorm:"not null;unique"`
	Abstract   string `gorm:"size:2048"`
	Content    string `sql:"type:longText"`
	Title      string `gorm:"not null;unique"`
	Notes      []Note
	LanguageID string
}

type History struct {
	Identity string `sql:"type:text"`
	HitDate  time.Time
	PostID   uint
}

type Tag struct {
	gorm.Model
	PostID  uint
	TagName string
}

type Comment struct {
	gorm.Model
	PostID   uint
	Content  string `gorm:"size:2048"`
	Status   string
	UserName string
}

type Note struct {
	gorm.Model
	ContentPostID uint
	Number        int
	Content       string `gorm:"size:2048"`
}

type Configuration struct {
	Connection_String string
	Admin_MDP         string
	Url               string
}

type Picture struct {
	Name     string
	Url      string
	UrlThumb string
}

type Icon struct {
	Name string
	Url  string
}

func hashIdString(text string) string {
	converted := []byte(text)
	hasher := sha256.New()
	hasher.Write(converted)
	return hex.EncodeToString(hasher.Sum(nil))
}
