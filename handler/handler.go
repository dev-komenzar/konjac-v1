package handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	translate "cloud.google.com/go/translate/apiv3"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"

	"github.com/tuckKome/konjac/db"
)

type sessionInfo struct {
	UserID         int
	Name           string //ログインしているユーザの表示名
	IsSessionAlive bool   //セッションが生きているかどうか
}

type pair struct {
	Language string
	Code     string
	Parent   string
	Response string
}

type tr struct {
	Translated string `json:"translated"`
	Language   string `json:"language"`
	Lcode      string `json:"lcode"`
	Random     int    `json:"random"`
}

// Login writes user info to session
func Login(c *gin.Context, u db.User) {
	session := sessions.Default(c)
	session.Set("Name", u.Name)
	session.Set("Email", u.Email)
	session.Save()
}

//getSessionInfo get current user info
func getSessionInfo(c *gin.Context) sessionInfo {
	var info sessionInfo
	session := sessions.Default(c)
	userID := session.Get("UserID")
	userName := session.Get("Name")
	if userID == nil {
		info = sessionInfo{
			UserID: -1, Name: "", IsSessionAlive: false,
		}
	} else {
		info = sessionInfo{
			UserID:         userID.(int),
			Name:           userName.(string),
			IsSessionAlive: true,
		}
		log.Println(info)
	}
	return info
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//Index show index page
func Index(c *gin.Context) {
	info := getSessionInfo(c)
	c.HTML(200, "index.html", gin.H{
		"sessionInfo": info,
	})
}

//GetTranslation translate word using google cloud translate.
//If you update languages, update template file also.
func GetTranslation(c *gin.Context) {
	pairs := []pair{
		{"日本語", "ja", "root", ""},
		{"英語", "en", "germanic", ""},
		{"ドイツ語", "de", "germanic", ""},
		{"アイルランド語", "ga", "celtic", ""},
		{"アラビア語", "ar", "afro-asiatic", ""},
		{"ギリシャ語", "el", "hellenic", ""},
		{"エスペラント", "eo", "root", ""},
		{"スペイン語", "es", "la", ""},
		{"フランス語", "fr", "la", ""},
		{"イタリア語", "it", "la", ""},
		{"オランダ語", "nl", "germanic", ""},
		{"アフリカーンス語", "af", "germanic", ""},
		{"フィンランド語", "fi", "uralic", ""},
		{"スウェーデン語", "sv", "germanic", ""},
		{"ノルウェー語", "no", "germanic", ""},
		{"アイスランド語", "is", "germanic", ""},
		{"ロシア語", "ru", "slavic", ""},
		{"ポーランド語", "pl", "slavic", ""},
		{"ブルガリア語", "bg", "slavic", ""},
		{"ウクライナ語", "uk", "slavic", ""},
		{"キルギス語", "ky", "turkic", ""},
		{"トルコ語", "tr", "turkic", ""},
		{"スワヒリ語", "sw", "root", ""},
		{"中国語", "zh", "sino-tibetan", ""},
		{"ミャンマー語", "my", "sino-tibetan", ""},
		{"ベトナム語", "vi", "austroasiatic", ""},
		{"タイ語", "th", "tai-kadai", ""},
		{"ラーオ語", "lo", "tai-kadai", ""},
		{"ラテン語", "la", "italic", ""},
		{"リトアニア語", "lt", "baltic", ""},
		{"ラトビア語", "lv", "baltic", ""},
		{"ヒンディー語", "hi", "indo", ""},
		{"チェコ語", "cs", "slavic", ""},
		{"モンゴル語", "mn", "root", ""},
		{"ゲール語", "gd", "celtic", ""},
		{"ネパール語", "ne", "indo", ""},
		{"ウルドゥー語", "ur", "indo", ""},
		{"クメール語", "km", "austroasiatic", ""},
		{"ヘブライ語", "he", "afro-asiatic", ""},
		{"ペルシア語", "fa", "iranian", ""},
		{"クルド語", "ku", "iranian", ""},
		{"エストニア語", "et", "uralic", ""},
		{"ハンガリー語", "hu", "uralic", ""},
		{"韓国語", "ko", "root", ""},
		{"パンジャーブ語", "pa", "indo", ""},
		{"ベンガル語", "bn", "indo", ""},
		{"タジク語", "tg", "iranian", ""},
		{"カザフ語", "kk", "turkic", ""},
		{"ジョージア語", "ka", "root", ""},
		{"インドネシア語", "id", "austronesian", ""},
		{"ハワイ語", "haw", "austronesian", ""},
		{"サモア語", "sm", "austronesian", ""},
		{"アゼルバイジャン語", "az", "turkic", ""},
	}

	text := c.Param("text")
	texts := []string{text}
	fmt.Println(pairs)
	// Clientをつくる
	client, err := translate.NewTranslationClient(c)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	// ジョブ数をあらかじめ登録
	wg.Add(len(pairs))

	t := func(i int) {
		request := &translatepb.TranslateTextRequest{
			Contents:           texts,
			TargetLanguageCode: pairs[i].Code,
			Parent:             "projects/honyac-konjac",
		}

		translation, err := client.TranslateText(c, request)
		if err != nil {
			log.Fatalf("Failed to translate text: %v", err)
		}
		fmt.Println(translation)
		t := translation.Translations[0].GetTranslatedText()
		pairs[i].Response = strings.ReplaceAll(t, "&#39;", "'") //フランス語でL'EgypteがL&#39;Egupteとなる問題
		fmt.Println(pairs[i])                                   //デバッグ用
		wg.Done()                                               // Doneで完了を通知
	}

	for i := range pairs {
		go t(i)
	}
	wg.Wait()

	//example struct whitch is returned
	// [{response:"日本語\nこんにちは"}]

	var tt []tr = make([]tr, len(pairs), len(pairs))
	for i := range tt {
		tt[i].Translated = pairs[i].Response
		tt[i].Language = pairs[i].Language
		tt[i].Lcode = pairs[i].Code

		rand.Seed(time.Now().UnixNano())
		tt[i].Random = rand.Intn(3)
	}

	c.JSONP(http.StatusOK, tt)
}

func GetError(c *gin.Context) {
	info := getSessionInfo(c)

	c.HTML(200, "error.html", gin.H{
		"SessionInfo": info,
		"Error":       c.Errors,
	})
}
