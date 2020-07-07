package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/tuckKome/konjac/db"
	"github.com/tuckKome/konjac/handler"
	"github.com/tuckKome/konjac/render"
)

func main() {
	db.Init()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		// 許可したいアクセス元の一覧
		AllowOrigins: []string{
			"localhost:8080",
			"localhost:3000",
		},
	}))

	router.Static("/templates", "./templates")
	router.Static("/bootstrap-4", "./bootstrap-4")
	router.Static("/jsmind", "./jsmind")
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	router.StaticFile("/apple-touch-icon.png", "./resources/apple-touch-icon.png")
	router.StaticFile("/favicon-32x32.png", "./resources/favicon-32x32.png")
	router.StaticFile("/favicon-16x16.png", "./resources/favicon-16x16.png")
	router.StaticFile("/manifest.json", "./resources/manifest.json")
	router.StaticFile("/safari-pinned-tab.svg", "./resources/safari-pinned-tab.svg")
	router.StaticFile("/android-chrome-192x192.png", "./resources/android-chrome-192x192.png")
	router.StaticFile("/android-chrome-256x256.png", "./resources/android-chrome-256x256.png")

	router.HTMLRender = render.LoadTemplates("./templates") // 事前にテンプレートをロード multitemplateで
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/", handler.Index)

	router.GET("/translate/:text", handler.GetTranslation)

	router.GET("/error", handler.GetError)

	router.Run()
}
