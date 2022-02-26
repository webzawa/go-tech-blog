package repository

import (
	"database/sql"
	"math"
	"time"

	"go-tech-blog/model"
)

// ArticleList ...
func ArticleListByCursor(cursor int) ([]*model.Article, error) {
	// 引数で渡されたカーソルの値が0以下の場合は代わりにint型の最大値で置き換える
	if cursor <= 0 {
		cursor = math.MaxInt32
	}

	// IDの降順に記事データを10件取得
	query := `SELECT *
	FROM articles
	Where ID < ?
	ORDER BY id desc
	LIMIT 10;`

	// クエリ結果格納スライス作成、、10件取得と決め、サイズとキャパシティを指定
	articles := make([]*model.Article, 0, 10)

	//クエリ結果格納変数、クエリ文字列、パラメータを指定してクエリを実行する
	if err := db.Select(&articles, query, cursor); err != nil {
		return nil, err
	}

	return articles, nil
}

// // ArticleList ...
// func ArticleList() ([]*model.Article, error) {
// 	query := `SELECT * FROM articles;`

// 	var articles []*model.Article
// 	if err := db.Select(&articles, query); err != nil {
// 		return nil, err
// 	}

// 	return articles, nil
// }

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
