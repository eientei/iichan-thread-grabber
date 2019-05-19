package model

type ImageRating int32

const (
	ImageRatingQuestionable ImageRating = iota
	ImageRatingSafe
	ImageRatingExplicit
)

type PostModel struct {
	Id          int64       `db:"id"`
	ImageId     int64       `db:"image_id"`
	ImageTags   string      `db:"image_tags"`
	ImageRating ImageRating `db:"image_rating"`
}

type PostgroupModel struct {
	Id int64 `db:"id"`
}

type PostgroupImageModel struct {
	Id      int64  `db:"id"`
	Order   int64  `db:"order"`
	GroupId int64  `db:"group_id"`
	ImageId string `db:"image_id"`
}
