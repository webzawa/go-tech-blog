package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"go-tech-blog/model"
	"go-tech-blog/repository"

	"github.com/labstack/echo/v4"
)

type ArticleCreateOutput struct {
	Article          *model.Article
	Message          string
	ValidationErrors []string
}

func ArticleIndex(c echo.Context) error {
	articles, err := repository.ArticleList()

	if err != nil {
		log.Println(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Message":  "Article Index",
		"Now":      time.Now(),
		"Articles": articles,
	}
	return render(c, "article/index.html", data)
}

func ArticleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}
	return render(c, "article/new.html", data)
}

func ArticleShow(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	data := map[string]interface{}{
		"Message": "Article Show",
		"Now":     time.Now(),
		"ID":      id,
	}
	return render(c, "article/show.html", data)
}

func ArticleEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	data := map[string]interface{}{
		"Message": "Article Edit",
		"Now":     time.Now(),
		"ID":      id,
	}
	return render(c, "article/edit.html", data)
}

func ArticleCreate(c echo.Context) error {
	// 送信されてくるフォームの内容を格納する構造体を宣言
	var article model.Article

	//　レスポンスとして返却する構造体を宣言
	var out ArticleCreateOutput

	// フォームの内容を構造体に埋め込む
	if err := c.Bind(&article); err != nil {
		// エラーの内容をログ出力
		c.Logger().Error(err.Error())
		// リクエストの解釈に失敗した場合は400エラーを返却
		return c.JSON(http.StatusBadRequest, out)
	}

	// repositoryを呼び出して保存処理を実行する
	res, err := repository.ArticleCreate(&article)
	if err != nil {
		// エラーの内容をログ出力
		c.Logger().Error(err.Error())
		// 500エラーを返却
		return c.JSON(http.StatusInternalServerError, out)
	}

	// SQL実行結果から作成されたレコードのIDを取得
	id, _ := res.LastInsertId()

	//構造体にIDをセット
	article.ID = int(id)

	// レスポンスの構造体に保存した記事のデータを格納
	out.Article = &article

	// 処理成功時にステータスコード200でレスポンスを返却する
	return c.JSON(http.StatusOK, out)
}
