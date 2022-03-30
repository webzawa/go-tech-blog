package main

import (
	"log"

	"go-tech-blog/handler"
	"go-tech-blog/repository"

	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var db *gorm.DB
var e = createMux()

func main() {
	db = connectDB()
	repository.SetDB(db)

	//TOPページに記事一覧を表示
	e.GET("/", handler.ArticleIndex)

	//記事関連ページは"/articles"で開始することとする
	//記事一覧画面には"/"と"articles"の療法でアクセスできるようにする
	//パスパラメータの":id"も":articleID"と明確にしている
	e.GET("/articles", handler.ArticleIndex)                // 一覧画面
	e.GET("/articles/new", handler.ArticleNew)              // 新規作成画面
	e.GET("/articles/:articleID", handler.ArticleShow)      // 詳細画面
	e.GET("/articles/:articleID/edit", handler.ArticleEdit) // 編集画面

	//JSON返却処理は"/api"で開始する、記事関連なので"/articles"を続ける
	e.GET("/api/articles", handler.ArticleList)                 // 一覧
	e.POST("/api/articles", handler.ArticleCreate)              // 作成
	e.DELETE("/api/articles/:articleID", handler.ArticleDelete) // 削除
	e.PATCH("/api/articles/:articleID", handler.ArticleUpdate)  // 更新

	e.POST("/api/users", handler.UserCreate) // 作成

	// e.Validator = &CustomValidator{validator: validator.New()}

	ctx := context.Background()
	opt := option.WithCredentialsFile("./seisou-75737-firebase-adminsdk-ptcqo-12d01b7869.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		// return nil, fmt.Errorf("error initializing app: %v", err)
		log.Printf("not Successfully fetched user data")
	}
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	uid := "jpAemF4fXWSNqD4HsPtUvkJP8o92"
	u, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Fatalf("error getting user %s: %v\n", uid, err)
	}
	fmt.Println("============================")
	log.Printf("user data: %#v\n", u)
	log.Printf("Successfully fetched user data: %#v\n", u.UserInfo)
	fmt.Println("============================")

	idToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjQ2NDExN2FjMzk2YmM3MWM4YzU5ZmI1MTlmMDEzZTJiNWJiNmM2ZTEiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoic2VyaXphd2EiLCJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vc2Vpc291LTc1NzM3IiwiYXVkIjoic2Vpc291LTc1NzM3IiwiYXV0aF90aW1lIjoxNjQ4NjMxODA2LCJ1c2VyX2lkIjoianBBZW1GNGZYV1NOcUQ0SHNQdFV2a0pQOG85MiIsInN1YiI6ImpwQWVtRjRmWFdTTnFENEhzUHRVdmtKUDhvOTIiLCJpYXQiOjE2NDg2MzI1MjQsImV4cCI6MTY0ODYzNjEyNCwiZW1haWwiOiJtLnNlcml6YXdhMjA2NEBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZW1haWwiOlsibS5zZXJpemF3YTIwNjRAZ21haWwuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoicGFzc3dvcmQifX0.Zl3vtbQbefHcdJB-uzAbnMzKkUGrrtQy_WCObiAmhJWFSBMQbncBQUiod6AVEKP45W-EhjBJbHuogRx-VXjF3cKSg25ziqyBHuawEc4r5F0F8NOsXuF7aULUFqNr__9PwyjjjVUngXs2Z9rz2XlfqKl0AEIMEAVHMiBSnaZGkBeJII3oaLQnCc7jkdfEl6okfG7N11gaTNh9xhiV1p13lFaxgo6x61tPalp4DCVg6QBLpv5ElfPGGz6cDvjBN1VYPeMVkKsWNfhN78TnE8Jf7hKYOfDm1405HHowOdL3XaO9CSNPKW9NS7TpdO8WSY7roOa6-O__Ej6O7f8KWlKC-A"
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)

	e.Logger.Fatal(e.Start(":8080"))
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	// e.Use(middleware.CSRF())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"authorization", "Content-Type"},
		AllowCredentials: true,
		AllowMethods:     []string{echo.OPTIONS, echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Static("/css", "src/css")
	e.Static("/js", "src/js")

	return e
}

func connectDB() *gorm.DB {
	// dsn := os.Getenv("DSN")
	dsn := "workuser:Passw0rd!@tcp(127.0.0.1:3306)/techblog?parseTime=true&autocommit=0&sql_mode=%27TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY%27"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal(err)
	}
	// if err := db.Ping(); err != nil {
	// 	e.Logger.Fatal(err)
	// }
	log.Println("db connection succeeded")
	return db
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
