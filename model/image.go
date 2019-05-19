package model

type ImageModel struct {
	Id        int64  `db:"id"`
	FullUrl   string `db:"full_url"`
	FullSize  int64  `db:"full_size"`
	FullFile  string `db:"full_file"`
	ThumbUrl  string `db:"thumb_url"`
	ThumbSize int64  `db:"thumb_size"`
	ThumbFile string `db:"thumb_file"`
}
