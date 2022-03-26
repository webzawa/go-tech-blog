package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

type ArticleUpdateOutput struct {
	Article          *model.Article
	Message          string
	ValidationErrors []string
}

func ArticleIndex(c echo.Context) error {

	// "/articles"のパスでリクエスト時に"/"でリダイレクトする
	// GoogleAnalyticsなどでのアクセス解析時にパスが統一され分析しやすくなる
	if c.Request().URL.Path == "/articles" {
		c.Redirect(http.StatusPermanentRedirect, "/")
	}

	articles, err := repository.ArticleListByCursor(0)

	if err != nil {
		log.Println(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	//取得できた最後のArticleIDをカーソルとして設定
	var cursor int
	if len(articles) != 0 {
		cursor = articles[len(articles)-1].ID
	}

	data := map[string]interface{}{
		"Articles": articles,
		"Cursor":   cursor,
	}

	return render(c, "article/index.html", data)
}

func ArticleList(c echo.Context) error {
	//クエリパラメータからカーソル値を取得、数値型にキャスト
	cursor, _ := strconv.Atoi(c.QueryParam("cursor"))
	//リポジトリ処理を呼び出して記事一覧データを取得、引数にカーソル値を渡して、
	//IDのどの位置から10件取得するか指定する
	articles, err := repository.ArticleListByCursor(cursor)
	//エラー処理
	if err != nil {
		c.Logger().Error(err.Error())
		//クライアントに500エラーでレスポンスを返却
		//JSON形式でデータのみ返却するので、c.HTMLBlob()でなくc.JSON()を呼び出す
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, articles)
}

func ArticleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}
	return render(c, "article/new.html", data)
}

func ArticleShow(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleID"))

	article, err := repository.ArticleGetByID(id)

	if err != nil {
		c.Logger().Error(err.Error())
		//クライアントに500エラーでレスポンスを返却
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Article": article,
	}

	return render(c, "article/show.html", data)
}

func ArticleEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleID"))

	article, err := repository.ArticleGetByID(id)

	if err != nil {
		c.Logger().Error(err.Error())
		//クライアントに500エラーでレスポンスを返却
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Article": article,
	}

	return render(c, "article/edit.html", data)
}

func ArticleCreate(c echo.Context) error {
	// file, err := c.FormFile("image")
	// if err != nil {
	// 	return err
	// }
	// src, err := file.Open()
	// if err != nil {
	// 	return err
	// }
	// defer src.Close()

	// // Destination
	// dst, err := os.Create("uploads" + "/" + file.Filename)
	// if err != nil {
	// 	return err
	// }
	// defer dst.Close()

	// // Copy
	// if _, err = io.Copy(dst, src); err != nil {
	// 	return err
	// }

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

	// if err := c.Validate(&article); err != nil {
	// 	// エラーの内容をログ出力
	// 	c.Logger().Error(err.Error())
	// 	// // エラーを内容をレスポンスの構造体に格納する
	// 	// out.Message = err.Error()
	// 	// エラー内容を検査してカスタムエラーメッセージを取得します
	// 	out.ValidationErrors = article.ValidationErrors(err)
	// 	// 解釈したParamが許可していない値の場合は422エラーを返却
	// 	return c.JSON(http.StatusUnprocessableEntity, out)
	// }

	// repositoryを呼び出して保存処理を実行する
	// res, err := repository.ArticleCreate(&article, dst)
	res, err := repository.ArticleCreate(&article)
	if err != nil {
		// エラーの内容をログ出力
		c.Logger().Error(err.Error())
		// 500エラーを返却
		return c.JSON(http.StatusInternalServerError, out)
	}

	// SQL実行結果から作成されたレコードのIDを取得
	id := res.ID

	//構造体にIDをセット
	article.ID = int(id)

	// レスポンスの構造体に保存した記事のデータを格納
	out.Article = &article

	// 処理成功時にステータスコード200でレスポンスを返却する
	return c.JSON(http.StatusOK, out)
}

func ArticleUpdate(c echo.Context) error {
	//リクエスト送信元のパスを取得
	ref := c.Request().Referer()
	//リクエスト送信元のパスから記事IDを抽出
	refID := strings.Split(ref, "/")[4]
	//リクエストURLのパスパラメータのから記事IDを抽出する
	reqID := c.Param("articleID")

	if reqID != refID {
		return c.JSON(http.StatusBadRequest, "")
	}

	// 送信されてくるフォームの内容を格納する構造体を宣言
	var article model.Article

	//　レスポンスとして返却する構造体を宣言
	var out ArticleCreateOutput

	// フォームで送信されたデータを変数に格納する
	if err := c.Bind(&article); err != nil {
		// エラーの内容をログ出力
		c.Logger().Error(err.Error())
		// リクエストの解釈に失敗した場合は400エラーを返却
		return c.JSON(http.StatusBadRequest, out)
	}

	if err := c.Validate(&article); err != nil {
		// エラーの内容をログ出力
		c.Logger().Error(err.Error())
		// // エラーを内容をレスポンスの構造体に格納する
		// out.Message = err.Error()
		// エラー内容を検査してカスタムエラーメッセージを取得します
		out.ValidationErrors = article.ValidationErrors(err)
		// 解釈したParamが許可していない値の場合は422エラーを返却
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	//文字列型のIDを数値型にキャスト
	articleID, _ := strconv.Atoi(reqID)
	//フォームデータを格納した構造体にIDをセット
	article.ID = articleID
	//記事を更新する処理を呼び出し
	_, err := repository.ArticleUpdate(&article)

	if err != nil {
		//レスポンスの構造体にエラー内容をセット
		out.Message = err.Error()
		//リクエストが正しいがサーバ側でエラー発生時には500エラー返却
		return c.JSON(http.StatusInternalServerError, out)
	}

	//レスポンスの構造体に記事データをセット
	out.Article = &article

	//処理成功時はステータスコード200でレスポンスを返却
	return c.JSON(http.StatusOK, out)
}

func ArticleDelete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleID"))

	if err := repository.ArticleDelete(id); err != nil {
		// エラーの内容をログ出力
		c.Logger().Error(err.Error())
		// 500エラーを返却
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Article %d id deleted", id))
}
