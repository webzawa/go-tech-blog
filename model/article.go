package model

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
)

// Article ...
type Article struct {
	ID        int       `db:"id" form:"id" json:"id"`
	Title     string    `db:"title" form:"title" validate:"required,max=50" json:"title"`
	Body      string    `db:"body" form:"body" validate:"required" json:"body"`
	ImagePath string    `db:"image_path" form:"image_path" json:"image_path"`
	ThumbPath string    `db:"thumb_path" form:"thumb_path" json:"thumb_path"`
	Created   time.Time `db:"created" json:"created"`
	Updated   time.Time `db:"updated" json:"updated"`
}

func (a *Article) ValidationErrors(err error) []string {
	// メッセージ格納用スライス
	var errMessages []string
	//複数エラーが発生する場合があるのでループ処理
	for _, err := range err.(validator.ValidationErrors) {
		//メッセージ格納変数を用意
		var message string
		//エラーフィールドを特定
		switch err.Field() {
		case "Title":
			//エラーになったバリデーションルールの特定
			switch err.Tag() {
			case "required":
				message = "タイトルは必須です。"
			case "max":
				message = "タイトルは50文字までです。"
			}
		case "Body":
			message = "本文は必須です。"
		}
		// messageをスライスに追加
		if message != "" {
			errMessages = append(errMessages, message)
		}
	}
	return errMessages
}
