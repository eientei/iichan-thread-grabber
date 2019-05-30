package model

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"
)

type ThreadModel struct {
	Url     string
	Created time.Time
	Images  []*ImageModel
}

type ImageModel struct {
	Md5       string
	Filepath  string
	Thumbpath string
	Postid    int
	Tags      string
	Rating    string
	ParentMd5 string
	Group     int
	Order     int
	Index     int
}

func (img *ImageModel) Update(thread string) error {
	_, err := psql.Update("grabber_images").SetMap(map[string]interface{}{
		"parent_md5": img.ParentMd5,
		"rating":     img.Rating,
		"tags":       img.Tags,
	}).Where(squirrel.Eq{"md5": img.Md5}).Exec()
	fmt.Println(psql.Update("grabber_images").SetMap(map[string]interface{}{
		"parent_md5": img.ParentMd5,
		"rating":     img.Rating,
		"tags":       img.Tags,
	}).Where(squirrel.Eq{"md5": img.Md5}).ToSql())
	if err != nil {
		return err
	}
	_, err = psql.Update("grabber_threads_images").SetMap(map[string]interface{}{
		"pgroup": img.Group,
		"porder": img.Order,
	}).Where(squirrel.Eq{"image_md5": img.Md5, "thread_url": thread}).Exec()
	return err
}

func ListThreads() ([]string, error) {
	rows, err := psql.Select("url").From("grabber_threads").Query()
	if err != nil {
		return nil, err
	}
	var res []string
	for rows.Next() {
		var r string
		err = rows.Scan(&r)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func GetThread(url string) (*ThreadModel, error) {
	row := psql.Select("url", "created").From("grabber_threads").Where(squirrel.Eq{"url": url}).QueryRow()
	model := &ThreadModel{}
	err := row.Scan(&model.Url, &model.Created)
	if err != nil {
		return nil, err
	}
	rows, err := psql.Select(
		"grabber_images.md5",
		"grabber_images.filepath",
		"grabber_images.thumbpath",
		"grabber_images.postid",
		"grabber_images.tags",
		"grabber_images.rating",
		"grabber_images.parent_md5",
		"grabber_threads_images.pgroup",
		"grabber_threads_images.porder",
		"grabber_threads_images.norder",
	).
		From("grabber_threads_images left join grabber_images on grabber_threads_images.image_md5 = grabber_images.md5").
		Where(squirrel.Eq{"grabber_threads_images.thread_url": url}).
		OrderBy("grabber_threads_images.pgroup asc", "grabber_threads_images.porder asc").Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		img := &ImageModel{}
		err = rows.Scan(&img.Md5, &img.Filepath, &img.Thumbpath, &img.Postid, &img.Tags, &img.Rating, &img.ParentMd5, &img.Group, &img.Order, &img.Index)
		if err != nil {
			return nil, err
		}
		model.Images = append(model.Images, img)
	}

	return model, nil
}

func CreateThread(url string) error {
	_, err := psql.Insert("grabber_threads").
		Columns("url", "created").Values(url, time.Now()).Suffix(`on conflict(url) do update set created = now()`).Exec()
	return err
}

func CreateImage(url string, md5, filepath, thumbpath string) error {
	postid := 0
	rating := "q"
	parentmd5 := ""
	tags := ""
	rows, err := DB.Query(`
select 
  foo.id,  
  foo.rating, 
  posts.md5 as parent_md5, 
  foo.tags
from (
  select 
    md5,
    id,
    parent_id,
    rating,array_to_string(tags_array, ' ') as tags 
  from posts
) as foo left join posts on 
  posts.id = foo.parent_id 
where foo.md5 = $1
  `, md5)
	if rows != nil && rows.Next() {
		_ = rows.Scan(&postid, &rating, &parentmd5, &tags)
		_ = rows.Close()
	}

	_, err = psql.Insert("grabber_images").
		Columns("md5", "filepath", "thumbpath", "postid", "rating", "parent_md5", "tags").
		Values(md5, filepath, thumbpath, postid, rating, parentmd5, tags).Suffix(`on conflict(md5) do update set filepath = $2, thumbpath = $3`).Exec()
	if err != nil {
		return err
	}
	_, err = psql.Insert("grabber_threads_images").
		Columns("thread_url", "image_md5").
		Values(url, md5).Suffix(`on conflict(thread_url, image_md5) do nothing`).Exec()
	return err
}

func UpdateImage(md5 string, rating, parentmd5, tags *string) error {
	updmap := make(map[string]interface{})
	if rating != nil {
		updmap["rating"] = *rating
	}
	if parentmd5 != nil {
		updmap["parent_md5"] = *parentmd5
	}
	if tags != nil {
		updmap["tags"] = *tags
	}
	_, err := psql.Update("grabber_images").SetMap(updmap).Where(squirrel.Eq{"md5": md5}).Exec()
	return err
}
