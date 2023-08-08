package repo_it_test

import (
	"context"
	"database/sql"
	"dmmak/simple-rest-crud/internal/entity"
	"dmmak/simple-rest-crud/internal/handler"
	"dmmak/simple-rest-crud/internal/repo"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

var pg handler.RedditPostsRepo

func TestMain(m *testing.M) {
	pgUrl := os.Getenv("PG_URL")
	db, err := sql.Open("pgx", pgUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//TODO make migrate conf
	_, err = db.Exec("drop schema if exists test cascade")
	if err != nil {
		log.Fatalf("can't drop test schema: %s", err)
	}
	_, err = db.Exec("create schema test")
	if err != nil {
		log.Fatalf(("can't create test schema"))
	}
	_, err = db.Exec("create table test.posts (uuid varchar(8) primary key, title varchar(256), likes smallint)")
	if err != nil {
		log.Fatalf(("can't create table posts"))
	}
	_, err = db.Exec(`
		create table test.comments (
			uuid varchar(8) primary key, 
			post_uuid varchar(8) CONSTRAINT post_fk REFERENCES test.posts(uuid) ON DELETE CASCADE, 
			body varchar(4000), 
			likes smallint);
	`)
	if err != nil {
		log.Fatalf(("can't create table comments"))
	}
	defer db.Exec("drop schema test cascade")

	pg = repo.NewPGRepo(db)

	code := m.Run()
	os.Exit(code)
}

func TestSavePostWithCommentsThenGet(t *testing.T) {
	expectedPost := &entity.RedditPost{
		UUID:     "p1",
		Title:    "Test reddit post",
		Likes:    5,
		Comments: []*entity.Comment{},
	}
	comments := []*entity.Comment{
		{
			UUID:  "c1",
			Body:  "Comment body 1",
			Likes: 1,
		},
		{
			UUID:  "c2",
			Body:  "Comment body 2",
			Likes: 2,
		},
	}
	expectedPost.Comments = comments
	err := pg.Save(context.Background(), expectedPost)
	if err != nil {
		t.Fatalf("error while saving post: %s", err)
	}

	actualPost, found, err := pg.Get(context.Background(), expectedPost.UUID)
	if err != nil {
		t.Fatalf("error while getting post: %s", err)
	}
	assert.True(t, found)
	assert.Equal(t, expectedPost, actualPost)
}

func TestSavePostWithCommentsThenDelete(t *testing.T) {
	post := &entity.RedditPost{
		UUID:     "pDelete",
		Title:    "Test reddit post for deletion",
		Likes:    1,
		Comments: []*entity.Comment{},
	}
	comments := []*entity.Comment{
		{
			UUID:  "cDelete",
			Body:  "Comment body for deletion",
			Likes: 1,
		},
	}
	post.Comments = comments

	err := pg.Save(context.Background(), post)
	if err != nil {
		t.Fatalf("error while saving post: %s", err)
	}
	found, err := pg.Delete(context.Background(), "pDelete")
	if err != nil {
		t.Fatalf("error while deleting post, uuid=%v: %s", post.UUID, err)
	}
	assert.True(t, found)

	_, found, err = pg.Get(context.Background(), "pDelete")
	if err != nil {
		t.Fatalf("error while getting post, uuid=%v: %s", post.UUID, err)
	}
	assert.False(t, found)
}
