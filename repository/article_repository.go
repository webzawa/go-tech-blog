package repository

import (
	"database/sql"
	"time"

	"go-tech-blog/model"
)

// ArticleList ...
func ArticleList() ([]*model.Article, error) {
	query := `SELECT * FROM articles;`

	var articles []*model.Article
	if err := db.Select(&articles, query); err != nil {
		return nil, err
	}

	return articles, nil
}

func ArticleCreate(article *model.Article) (sql.Result, error) {
	now := time.Now()
	article.Created = now
	article.Updated = now

	query := `INSERT INTO articles (title, body, created, updated)
	VALUES (:title, :body, :created, :updated);`

	// トランザクションを開始
	tx := db.MustBegin()

	// クエリ文字列と構造体を引数に渡してSQL実行、クエリ文字列の「:title」などは構造体の値で置換される
	// 構造体タグで指定してあるフィールドが対象となる。`db: "title"`など
	res, err := tx.NamedExec(query, article)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return res, nil
}
