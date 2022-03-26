package model

type User struct {
	ID   int    `db:"id" form:"id" json:"id"`
	Uid  string `db:"uid" form:"uid" json:"uid"`
	Name string `db:"name" form:"name" json:"name"`
}
