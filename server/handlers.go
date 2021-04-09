package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func GlobalHandler(c *gin.Context) {
	fmt.Println("Global", c.Request.UserAgent(), c.Request.URL)
	var isBot = false
	for _, b := range Bots {
		if strings.Contains(c.Request.UserAgent(), b) {
			isBot = true
			break
		}
	}

	if isBot {
		var post Post
		var cpost ContentPost
		var config BlogConfig
		db, err := gorm.Open("mysql", ConnexionString)
		if err != nil {
			panic("failed to connect database")
		}
		defer db.Close()

		title := c.Param("title")
		title, _ = url.PathUnescape(title)
		db.Where("Url = ?", title).First(&cpost)
		db.Where("ID = ?", cpost.PostID).First(&post)
		db.First(&config)

		// c.Request.RemoteAddr
		// c.Request.URL.Path
		fmt.Println(c.Request.URL.ResolveReference(c.Request.URL))
		c.HTML(http.StatusOK, "robots.tmpl", BotInfos{
			Title:       config.Title,
			Description: config.Description,
			Pic:         config.Url + "/" + config.BlogPic,
			Url:         config.Url,
			BlogName:    config.Title})
	} else {
		c.File("../dist/index.html")
	}
}

type BotInfos struct {
	Title       string
	Description string
	Author      string
	Pic         string
	Url         string
	BlogName    string
}

var Bots []string = []string{"facebookexternalhit", "twitter", "LinkedInBot", "skype"}

func PostHandler(c *gin.Context) {
	fmt.Println("Post", c.Request.UserAgent(), c.Request.URL)
	var isBot = false
	for _, b := range Bots {
		if strings.Contains(c.Request.UserAgent(), b) {
			isBot = true
			break
		}
	}

	if isBot {

		var post Post
		var cpost ContentPost
		var config BlogConfig
		db, err := gorm.Open("mysql", ConnexionString)
		if err != nil {
			panic("failed to connect database")
		}
		defer db.Close()

		title := c.Param("title")
		title, _ = url.PathUnescape(title)
		db.Where("Url = ?", title).First(&cpost)
		db.Where("ID = ?", cpost.PostID).First(&post)
		db.First(&config)

		c.HTML(http.StatusOK, "robots.tmpl", BotInfos{
			Title:       cpost.Title,
			Description: cpost.Abstract,
			Pic:         config.Url + "/" + post.MainPicture,
			Url:         config.Url + c.Request.URL.String(),
			BlogName:    config.Title})

	} else {
		c.File("../dist/index.html")
	}

}

func NewPost(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var post Post
	c.ShouldBind(&post)
	for i, _ := range post.Contents {
		if post.Contents[i].Title == "" {
			post.Contents = append(post.Contents[:i], post.Contents[i+1:]...)
		} else {
			post.Contents[i].URL = strings.Replace(post.Contents[i].Title, " ", "-", -1)
			post.Contents[i].URL = strings.Replace(post.Contents[i].URL, "?", "-", -1)
			post.Contents[i].URL = strings.Replace(post.Contents[i].URL, "&", "-", -1)
			post.Contents[i].URL = strings.Replace(post.Contents[i].URL, "=", "-", -1)
			post.Contents[i].URL = strings.Replace(post.Contents[i].URL, "+", "-", -1)
		}
	}

	if post.ID != 0 {
		db.Save(&post)
		c.Status(http.StatusCreated)
	} else {
		post.Date = time.Now()
		if err := db.Create(&post).Error; err != nil {
			c.String(http.StatusInternalServerError, "Error during post creation")
			return
		} else {
			c.Status(http.StatusCreated)
		}
	}

}

func NewPage(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var post Page
	c.ShouldBind(&post)
	for i, _ := range post.Contents {
		if post.Contents[i].Title == "" {
			post.Contents = append(post.Contents[:i], post.Contents[i+1:]...)
		} else {
			post.Contents[i].URL = strings.Replace(post.Contents[i].Title, " ", "-", -1)
		}
	}

	if post.ID != 0 {
		db.Save(&post)
		c.Status(http.StatusCreated)
	} else {
		post.Date = time.Now()
		if err := db.Create(&post).Error; err != nil {
			c.String(http.StatusInternalServerError, "Error during post creation")
			return
		} else {
			c.Status(http.StatusCreated)
		}
	}

}
func NewTag(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var tag Tag
	c.ShouldBind(&tag)
	db.Save(&tag)
	c.Status(http.StatusCreated)

}
func NewCat(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var tag Category
	c.ShouldBind(&tag)
	fmt.Println("Cat", tag)
	fmt.Println("Trads", tag.Trads)
	if tag.ID != 0 {
		db.Debug().Save(&tag)
	} else {
		db.Debug().Create(&tag)
	}

	c.Status(http.StatusCreated)

}
func PostComment(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var comment Comment
	c.ShouldBind(&comment)
	comment.Status = "pending"
	db.Save(&comment)
	c.Status(http.StatusCreated)

}

func GetCategories(c *gin.Context) {
	var cats []Category

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Find(&cats)
	for i, _ := range cats {
		db.Model(cats[i]).Related(&cats[i].Trads)
	}

	c.JSON(http.StatusOK, cats)
}
func GetLanguages(c *gin.Context) {
	var lg []Language

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Find(&lg)

	c.JSON(http.StatusOK, lg)
}

func GetPost(c *gin.Context) {

	var cpost ContentPost
	var post Post

	var hit History
	var otherHit []History

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	hash := hashIdString(c.ClientIP() + c.GetHeader("User-Agent"))
	hit.Identity = hash
	hit.HitDate = time.Now()

	title := c.Param("title")
	title, _ = url.PathUnescape(title)
	// title = strings.Replace(title, "-", " ", -1)
	db.Where("URL = ?", title).First(&cpost)
	db.Model(cpost).Related(&cpost.Notes)
	db.Where("ID = ? AND Public = true", cpost.PostID).First(&post)

	db.Model(post).Where("Status = 'Approved'").Related(&post.Comments)

	post.Contents = append(post.Contents, cpost)

	db.Where("Post_id = ? AND Identity = ? AND DATE(hit_date) = DATE(NOW())", post.ID, hit.Identity).Find(&otherHit)

	if len(otherHit) == 0 {
		post.UniqueViews++
	}

	post.Views++
	db.Save(&post)
	hit.PostID = post.ID
	db.Save(&hit)

	c.JSON(http.StatusOK, post)
}
func GetPostPreview(c *gin.Context) {

	var cpost ContentPost
	var post Post

	var hit History

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	hash := hashIdString(c.ClientIP() + c.GetHeader("User-Agent"))
	hit.Identity = hash
	hit.HitDate = time.Now()

	title := c.Param("title")
	title, _ = url.PathUnescape(title)
	// title = strings.Replace(title, "-", " ", -1)
	db.Where("URL = ?", title).First(&cpost)
	db.Model(cpost).Related(&cpost.Notes)
	db.Where("ID = ?", cpost.PostID).First(&post)
	db.Model(post).Where("Status = 'Approved'").Related(&post.Comments)
	post.Contents = append(post.Contents, cpost)

	c.JSON(http.StatusOK, post)
}

func GetPage(c *gin.Context) {

	var cpost ContentPost
	var post Page

	var hit History
	var otherHit []History

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	hash := hashIdString(c.ClientIP() + c.GetHeader("User-Agent"))
	hit.Identity = hash
	hit.HitDate = time.Now()

	title := c.Param("title")
	title, _ = url.PathUnescape(title)
	title = strings.Replace(title, "-", " ", -1)
	fmt.Println("TITLE", title)
	db.Where("Title = ?", title).First(&cpost)
	db.Model(cpost).Related(&cpost.Notes)
	db.Where("ID = ?", cpost.PostID).First(&post)

	post.Contents = append(post.Contents, cpost)

	db.Where("Post_id = ? AND Identity = ? AND DATE(hit_date) = DATE(NOW())", post.ID, hit.Identity).Find(&otherHit)

	if len(otherHit) == 0 {
		post.UniqueViews++
	}

	post.Views++
	db.Save(&post)
	hit.PostID = post.ID
	db.Save(&hit)

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {

	var post Post

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	c.ShouldBind(&post)
	db.Delete(&post)

	c.Status(http.StatusOK)
}
func DeleteCat(c *gin.Context) {

	var post Category

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	c.ShouldBind(&post)
	db.Where("ID = ?", post.ID).First(&post)
	fmt.Println(post)
	db.Debug().Delete(&post)

	c.Status(http.StatusOK)
}

func CommentStatus(c *gin.Context) {

	var com Comment

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	c.ShouldBind(&com)
	db.Save(&com)

	c.Status(http.StatusOK)
}

func GetPostByCats(c *gin.Context) {

	var posts []Post

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	claims := jwt.ExtractClaims(c)
	if claims["cats"] != nil {
		db.Where("Category = ? AND Public = true", claims["cats"]).Order("date desc").Find(&posts)
	} else {
		db.Where("Public = true").Order("date desc").Find(&posts)
	}
	for i, _ := range posts {
		db.Model(posts[i]).Related(&posts[i].Contents)
	}
	c.JSON(http.StatusOK, posts)
}

func GetPages(c *gin.Context) {

	var posts []Page

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Order("date desc").Find(&posts)

	for i, _ := range posts {
		db.Model(posts[i]).Related(&posts[i].Contents)
	}
	c.JSON(http.StatusOK, posts)
}

func GetPostByTags(c *gin.Context) {

	var posts []Post

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	claims := jwt.ExtractClaims(c)
	if claims["cats"] != nil {
		db.Where("Category = ?  AND Public = true", claims["cats"]).Order("date desc").Find(&posts)
	} else {
		db.Where("Public = true").Order("date desc").Find(&posts)
	}
	for i, _ := range posts {
		db.Model(posts[i]).Related(&posts[i].Contents)
	}
	c.JSON(http.StatusOK, posts)
}
func GetPostsAdmin(c *gin.Context) {
	var posts []Post

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Order("date desc").Find(&posts)
	for i, _ := range posts {
		db.Model(posts[i]).Related(&posts[i].Contents).Related(&posts[i].Comments)
		for j, _ := range posts[i].Contents {
			db.Model(posts[i].Contents[j]).Related(&posts[i].Contents[j].Notes)
		}
		// db.Model(posts[i])
	}
	c.JSON(http.StatusOK, posts)
}

func AddPictures(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var pic Picture

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filename := filepath.Join("pics", filepath.Base(file.Filename))

	// filepath.Join(dir, filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	pic.Name = file.Filename
	pic.Url = filename
	db.Save(&pic)

	c.Status(http.StatusOK)
}
func GetAllPictures(c *gin.Context) {

	var lg []Picture

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Find(&lg)

	c.JSON(http.StatusOK, lg)
}
func GetAllIcons(c *gin.Context) {

	c.JSON(http.StatusOK, icons)
}

func GetActiveTheme(c *gin.Context) {

	var lg Theme

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Where("active is true").First(&lg)
	fmt.Println(lg)

	c.JSON(http.StatusOK, lg)
}
func GetAllThemes(c *gin.Context) {

	var lg []Theme

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Find(&lg)

	c.JSON(http.StatusOK, lg)
}

func NewTheme(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var theme Theme
	c.ShouldBind(&theme)
	db.Save(&theme)
	c.Status(http.StatusCreated)

}

func DeleteTheme(c *gin.Context) {

	var theme Theme

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	c.ShouldBind(&theme)
	db.Where("Name = ?", theme.Name).First(&theme)
	db.Delete(&theme)

	c.Status(http.StatusOK)
}
func EditTheme(c *gin.Context) {

	var theme Theme

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	c.ShouldBind(&theme)
	db.Save(&theme)

	c.Status(http.StatusOK)
}
func SetActiveTheme(c *gin.Context) {

	var theme Theme

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	c.ShouldBind(&theme)
	db.Where("ID = ?", theme.ID).First(&theme)
	db.Model(Theme{}).Update("Active", false)
	theme.Active = true
	db.Save(&theme)

	c.Status(http.StatusOK)
}

func GetConfig(c *gin.Context) {

	var lg BlogConfig

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.First(&lg)

	c.JSON(http.StatusOK, lg)
}

func EditConfig(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var theme BlogConfig
	c.ShouldBind(&theme)
	db.Save(&theme)
	c.Status(http.StatusCreated)

}

func GetLinks(c *gin.Context) {

	var lg []Link

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Find(&lg)

	c.JSON(http.StatusOK, lg)
}

func NewLink(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var theme Link
	c.ShouldBind(&theme)
	db.Save(&theme)
	c.Status(http.StatusCreated)

}

func DeleteLink(c *gin.Context) {

	var theme Link

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	c.ShouldBind(&theme)
	db.Where("ID = ?", theme.ID).First(&theme)
	db.Delete(&theme)

	c.Status(http.StatusOK)
}

func SignUp(c *gin.Context) {

	db, err := gorm.Open("mysql", ConnexionString)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var acc Account
	c.ShouldBind(&acc)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	acc.Password = string(hashedPassword)

	if err := db.Create(&acc).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error during account creation")
		return
	} else {
		c.Status(http.StatusCreated)
	}

}
