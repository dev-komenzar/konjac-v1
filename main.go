package main

import (
	"github.com/gin-gonic/gin"

	"github.com/tuckKome/konjac-v1/db"
	"github.com/tuckKome/konjac-v1/handler"
)

func main() {
	db.Init()
	router := gin.Default()

	// router.Use(cors.New(cors.Config{
	// 	// 許可したいアクセス元の一覧
	// 	AllowOrigins: []string{
	// 		"*",
	// 	},
	// }))

	// router.Static("/templates", "./templates")
	// router.Static("/bootstrap-4", "./bootstrap-4")
	// router.Static("/jsmind", "./jsmind")

	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/apple-touch-icon.png", "./resources/apple-touch-icon.png")
	router.StaticFile("/favicon-32x32.png", "./resources/favicon-32x32.png")
	router.StaticFile("/favicon-16x16.png", "./resources/favicon-16x16.png")
	router.StaticFile("/manifest.json", "./resources/manifest.json")
	router.StaticFile("/safari-pinned-tab.svg", "./resources/safari-pinned-tab.svg")
	router.StaticFile("/android-chrome-192x192.png", "./resources/android-chrome-192x192.png")
	router.StaticFile("/android-chrome-256x256.png", "./resources/android-chrome-256x256.png")

	router.LoadHTMLGlob("./react/*.html")
	router.Static("/static", "./react/static")

	// router.HTMLRender = render.LoadTemplates("./templates") // 事前にテンプレートをロード multitemplateで
	// store := cookie.NewStore([]byte("secret"))
	// router.Use(sessions.Sessions("mysession", store))

	router.GET("/", handler.Index)
	// router.Use(static.Serve("/", static.LocalFile("./react", true)))
	router.GET("/translate/:text", handler.Index)

	router.GET("/api/:text", handler.GetTranslation)

	router.GET("/error", handler.GetError)

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(308, "/")
	})

	router.Run()
}
