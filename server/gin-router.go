package main

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	UserName  string
	FirstName string
	LastName  string
}

var identityKey = "id"

func initRoutes() {
	// Disable Console Color
	// gin.DisableConsoleColor()

	// Creates a gin r with default middleware:
	// logger and recovery (crash-free) middleware

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Use(LiberalCORS)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Static("/pics", "./pics")
	r.Static("/icons", "./icons")
	r.Static("/js", "../dist/js")
	r.Static("/css", "../dist/css")
	// r.StaticFS("/pics", http.Dir("./pics"))
	r.MaxMultipartMemory = 8 << 20

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("super secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(AccountApi); ok {
				return jwt.MapClaims{
					identityKey: v.Login,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			var a *AccountApi

			if _, ok := claims[identityKey].(string); ok {

				return claims[identityKey]
			}
			return a
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			var acc Account
			var accApi AccountApi
			db, _ := gorm.Open("mysql", ConnexionString)
			db.First(&acc, "Login = ?", userID)
			errPass := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password))

			if errPass != nil {
				return nil, jwt.ErrFailedAuthentication
			} else if acc.ID == 0 {
				return nil, jwt.ErrFailedAuthentication
			}

			accApi.ID = acc.ID
			accApi.Login = acc.Login
			return accApi, nil
		},
		Authorizator: func(user interface{}, c *gin.Context) bool {
			var acc Account
			db, _ := gorm.Open("mysql", ConnexionString)
			db.First(&acc, "is_admin = true")

			if v, ok := user.(string); ok && v == acc.Login {
				return true
			} else {
				return false
			}

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	auth := r.Group("/auth")
	admin := r.Group("/api/admin")

	r.NoRoute(GlobalHandler)
	r.GET("/", GlobalHandler)
	r.GET("/admin2", GlobalHandler)
	r.GET("/post/:title", PostHandler)
	r.GET("/page/:title", PostHandler)
	r.POST("/api/Login", authMiddleware.LoginHandler)
	r.POST("/api/SignUp", SignUp)
	r.GET("/api/post/:title", GetPost)
	r.GET("/api/page/:title", GetPage)
	r.GET("/api/pages", GetPages)
	r.POST("/api/comment", PostComment)
	r.GET("/api/posts/*cat", GetPostByCats)
	r.GET("/api/cats/*parent", GetCategories)
	r.GET("/api/languages", GetLanguages)
	r.GET("/api/config", GetConfig)
	r.GET("/api/theme", GetActiveTheme)
	r.GET("/api/themes", GetAllThemes)
	r.GET("/api/links", GetLinks)

	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/api/RefreshToken", authMiddleware.RefreshHandler)
	}
	admin.Use(authMiddleware.MiddlewareFunc())
	{
		admin.POST("/post", NewPost)
		admin.POST("/page", NewPage)
		admin.POST("/acomment", CommentStatus)
		admin.POST("/tag", NewTag)
		admin.POST("/cat", NewCat)
		admin.POST("/dcat", DeleteCat)
		admin.POST("/dpost", DeletePost)
		admin.POST("/aposts", GetPostsAdmin)
		admin.POST("/apic", AddPictures)
		admin.GET("/icons", GetAllIcons)
		admin.POST("/atheme", NewTheme)
		admin.POST("/etheme", EditTheme)
		admin.POST("/dtheme", DeleteTheme)
		admin.POST("/alink", NewLink)
		admin.POST("/dlink", DeleteLink)
		admin.POST("/activet", SetActiveTheme)
		admin.POST("/sconfig", EditConfig)
		admin.GET("/pics", GetAllPictures)
		admin.GET("/preview/:title", GetPostPreview)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run(":3010")
	// r.Run(":3000") for a hard coded port

}

func LiberalCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	if c.Request.Method == "OPTIONS" {
		if len(c.Request.Header["Access-Control-Request-Headers"]) > 0 {
			c.Header("Access-Control-Allow-Headers", c.Request.Header["Access-Control-Request-Headers"][0])
		}
		c.AbortWithStatus(http.StatusOK)
	}
}
