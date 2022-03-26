package handler

import (
	"net/http"

	"go-tech-blog/model"
	"go-tech-blog/repository"

	"github.com/labstack/echo/v4"
)

type UserCreateOutput struct {
	User             *model.User
	Message          string
	ValidationErrors []string
}

func UserCreate(c echo.Context) error {

	// 送信されてくるフォームの内容を格納する構造体を宣言
	var user model.User

	// レスポンスとして返却する構造体を宣言
	var out UserCreateOutput

	// フォームの内容を構造体に埋め込む
	if err := c.Bind(&user); err != nil {
		// エラーの内容をログ出力
		c.Logger().Error(err.Error())
		// リクエストの解釈に失敗した場合は400エラーを返却
		return c.JSON(http.StatusBadRequest, out)
	}

	// repositoryを呼び出して保存処理を実行する
	res, err := repository.UserCreate(&user)
	if err != nil {
		// エラーの内容をログ出力
		c.Logger().Error(err.Error())
		// 500エラーを返却
		return c.JSON(http.StatusInternalServerError, out)
	}

	// SQL実行結果から作成されたレコードのIDを取得
	id := res.ID

	//構造体にIDをセット
	user.ID = int(id)

	// レスポンスの構造体に保存した記事のデータを格納
	out.User = &user

	// 処理成功時にステータスコード200でレスポンスを返却する
	return c.JSON(http.StatusOK, out)
}
