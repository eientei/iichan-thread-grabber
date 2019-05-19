package model

import "time"

type ThreadModel struct {
	Id      int64     `db:"id"`
	Url     string    `db:"url"`
	Created time.Time `db:"created"`
}

type ThreadImageModel struct {
	Id       int64 `db:"id"`
	ThreadId int64 `db:"thread_id"`
	ImageId  int64 `db:"image_id"`
}
