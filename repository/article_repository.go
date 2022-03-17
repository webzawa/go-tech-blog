package repository

import (
	"math"
	"os"
	"time"

	"go-tech-blog/model"
)

// ArticleList ...
func ArticleListByCursor(cursor int) ([]*model.Article, error) {
	// 引数で渡されたカーソルの値が0以下の場合は代わりにint型の最大値で置き換える
	if cursor <= 0 {
		cursor = math.MaxInt32
	}

	// クエリ結果格納スライス作成、、10件取得と決め、サイズとキャパシティを指定
	articles := make([]*model.Article, 0, 10)

	//クエリ結果格納変数、クエリ文字列、パラメータを指定してクエリを実行する
	err := db.Where("id < ?", cursor).Order("id desc").Find(&articles).Limit(10).Error
	if err != nil {
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

func ArticleGetByID(id int) (*model.Article, error) {
	var article model.Article

	//クエリ結果格納変数、クエリ文字列、パラメータを指定してクエリを実行する
	//複数件取得はdb.Selectだが一件取得の場合はdb.Getを使用する
	if err := db.Find(&article, id).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

func ArticleCreate(article *model.Article, dst *os.File) (*model.Article, error) {
	now := time.Now()
	article.Created = now
	article.Updated = now
	article.ImagePath = dst.Name()
	article.ThumbPath = ""

	// トランザクションを開始
	tx := db.Begin()

	// クエリ文字列と構造体を引数に渡してSQL実行、クエリ文字列の「:title」などは構造体の値で置換される
	// 構造体タグで指定してあるフィールドが対象となる。`db: "title"`など
	err := tx.Create(&article).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return article, nil
}

func ArticleUpdate(article_org *model.Article) (*model.Article, error) {
	var article model.Article
	db.Find(&article, article_org.ID)

	now := time.Now()
	article_org.Updated = now

	// トランザクションを開始
	tx := db.Begin()

	err := tx.Model(article).Updates(model.Article{Title: article_org.Title, Body: article_org.Body, Created: article_org.Created, Updated: now}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &article, nil
}

func ArticleDelete(id int) error {

	var article model.Article

	// トランザクションを開始
	tx := db.Begin()
	tx.Find(&article, id)
	err := tx.Delete(&article).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
