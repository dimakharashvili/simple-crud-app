package handler

import (
	"bytes"
	"context"
	"dmmak/simple-rest-crud/internal/entity"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

type (
	RedditPostsRepo interface {
		Save(ctx context.Context, post *entity.RedditPost) (err error)
		Get(ctx context.Context, postID string) (post *entity.RedditPost, found bool, err error)
		Delete(ctx context.Context, postID string) (found bool, err error)
	}

	HttpHandler struct {
		postsRepo RedditPostsRepo
	}
)

func New(postsRepo RedditPostsRepo) *HttpHandler {
	return &HttpHandler{postsRepo: postsRepo}
}

func (h *HttpHandler) SavePost(w http.ResponseWriter, r *http.Request) {
	redditPost := &entity.RedditPost{}
	if err := json.NewDecoder(r.Body).Decode(redditPost); err != nil {
		log.Printf("can't decode save post request body: %s\n", err)
		http.Error(w, "Something gone wrong", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.postsRepo.Save(ctx, redditPost); err != nil {
		log.Printf("error while saving post: %s\n", err)
		http.Error(w, "Something gone wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HttpHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	redditPost := &entity.RedditPost{}

	re := regexp.MustCompile(`[a-z0-9]{8}`)
	postID := re.FindString(r.URL.Path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	redditPost, found, err := h.postsRepo.Get(ctx, postID)
	if err != nil {
		log.Printf("error while getting post: %s\n", err)
		http.Error(w, "Something gone wrong", http.StatusInternalServerError)
		return
	}
	if !found {
		errMsg := fmt.Sprintf("can't find post with uuid=%v\n", postID)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusNotFound)
	}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(redditPost)
	if err != nil {
		log.Printf("can't encode get post response body: %s\n", err)
		http.Error(w, "Something gone wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b.Bytes())
	if err != nil {
		log.Printf("can't write get post response body: %s\n", err)
		http.Error(w, "Something gone wrong", http.StatusInternalServerError)
	}
}

func (h *HttpHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`[a-z0-9]{8}`)
	postID := re.FindString(r.URL.Path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	found, err := h.postsRepo.Delete(ctx, postID)
	if err != nil {
		log.Printf("error while deleting post: %s\n", err)
		http.Error(w, "Something gone wrong", http.StatusInternalServerError)
		return
	}
	if !found {
		errMsg := fmt.Sprintf("can't find post with uuid=%v\n", postID)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HttpHandler) NotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Something gone wrong", http.StatusMethodNotAllowed)
}

func (h *HttpHandler) Route(allowedMethods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var actualMethod string
		for _, method := range allowedMethods {
			if r.Method == method {
				actualMethod = r.Method
			}
		}
		switch actualMethod {
		case http.MethodPost:
			h.SavePost(w, r)
		case http.MethodGet:
			h.GetPost(w, r)
		case http.MethodDelete:
			h.DeletePost(w, r)
		default:
			h.NotAllowed(w, r)
		}
	}
}
