package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

var ConnexionString string
var icons []Icon

func oneTimeCreateLanguages() {
	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()
	var cat Language

	cat.Name = "Français"
	cat.Short = "fr"
	db.Create(&cat)
	cat.Name = "English"
	cat.Short = "en"
	db.Create(&cat)

}
func oneTimeCreateAccounts(adminMDP string) {
	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()
	var admin Account

	admin.Name = "Oro"
	admin.Login = "Orosius"
	admin.IsAdmin = true
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminMDP), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	admin.Password = string(hashedPassword)
	admin.ProfilePic = "adminPic"
	db.Create(&admin)
}
func oneTimeCreateTheme() {
	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()
	var theme Theme
	theme.BackColor = "#2b2a2a"
	theme.FontColor = "#ffffff"
	theme.UnderColor = "#f37272"
	theme.SecondColor = "#584e4a"
	theme.Name = "Default Dark"
	theme.Active = true
	db.Create(&theme)
}

func oneTimeCreateConfig(url string) {
	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()
	var bc BlogConfig
	bc.Comments = true
	bc.Title = "Oro Blog"
	bc.Meta = "Blog orolol golang vuejs"
	bc.Grid = false
	bc.Url = url
	bc.Description = "Blog description !"
	db.Create(&bc)
}

func initIcons() []Icon {

	root := "./icons"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		name := info.Name()

		fmt.Println(name, name[len(name)-3:])
		if name[len(name)-3:] == "png" {
			icons = append(icons, Icon{
				Name: name[4 : len(name)-4],
				Url:  path,
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(icons)
	// for _, file := range files {
	// 	// file.
	// 	// fmt.Println(file)
	// }
	return icons
}

func main() {

	var configuration Configuration
	var filename = "config.json"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println(err)
	}

	ConnexionString = configuration.Connection_String

	db, err := gorm.Open("mysql", ConnexionString)
	db.Set("gorm:table_options", "charset=utf8mb4")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()
	initIcons()
	db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&Post{})
	db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&Page{})
	db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&ContentPost{})
	db.AutoMigrate(&Note{})
	db.AutoMigrate(&BlogConfig{})
	db.AutoMigrate(&Link{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Theme{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&CategoryTrad{})
	db.AutoMigrate(&Language{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Picture{})
	db.Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(&History{})
	db.AutoMigrate(&Comment{})

	db.Model(&Post{}).Related(&ContentPost{})
	db.Model(&Page{}).Related(&ContentPost{})
	db.Model(&Category{}).Related(&CategoryTrad{})
	db.Model(&ContentPost{}).Related(&Note{})
	db.Model(&Post{}).Related(&Tag{})
	db.Model(&Post{}).Related(&Comment{})
	var users []Account
	db.Find(&users)
	var themes []Theme
	db.Find(&themes)
	var bc []BlogConfig
	db.Find(&bc)

	if len(users) == 0 {
		oneTimeCreateAccounts(configuration.Admin_MDP)
	}
	if len(themes) == 0 {
		oneTimeCreateTheme()
	}
	if len(bc) == 0 {
		oneTimeCreateConfig(configuration.Url)
	}
	var lg []Language
	db.Find(&lg)

	if len(lg) == 0 {
		oneTimeCreateLanguages()
	}

	initRoutes()

}
