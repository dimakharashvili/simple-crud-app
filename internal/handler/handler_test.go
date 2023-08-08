package handler_test

import (
	"bytes"
	"context"
	"dmmak/simple-rest-crud/internal/entity"
	"dmmak/simple-rest-crud/internal/handler"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSavePost(t *testing.T) {
	tests := []struct {
		name          string
		expStatusCode int
	}{
		{
			name:          "success",
			expStatusCode: http.StatusOK,
		},
		{
			name:          "error",
			expStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestEntity := &entity.RedditPost{
				UUID:  "p1000000",
				Title: "Some title",
				Likes: 5,
				Comments: []*entity.Comment{
					{
						UUID:  "c0000001",
						Body:  "Some comment text 1",
						Likes: 1,
					},
					{
						UUID:  "c0000002",
						Body:  "Some comment text 2",
						Likes: 2,
					},
				},
			}
			mockCtrl := gomock.NewController(t)
			mockRepo := NewMockRedditPostsRepo(mockCtrl)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			mockRecorder := mockRepo.EXPECT().Save(gomock.AssignableToTypeOf(ctx), gomock.Eq(requestEntity))

			switch test.name {
			case "success":
				mockRecorder.DoAndReturn(
					func(context.Context, *entity.RedditPost) error {
						return nil
					},
				)
			case "error":
				mockRecorder.DoAndReturn(
					func(context.Context, *entity.RedditPost) error {
						return fmt.Errorf("some repo internal error")
					},
				)
			}

			h := handler.New(mockRepo)
			rr := httptest.NewRecorder()
			reqBody := new(bytes.Buffer)
			err := json.NewEncoder(reqBody).Encode(requestEntity)
			if err != nil {
				log.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodGet, "/posts/", reqBody)
			if err != nil {
				log.Fatal(err)
			}
			h.SavePost(rr, req)

			assert.Equal(t, test.expStatusCode, rr.Result().StatusCode)
		})
	}
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		expStatusCode  int
		responseEntity *entity.RedditPost
	}{
		{
			name:          "success",
			postID:        "p1000000",
			expStatusCode: http.StatusOK,
			responseEntity: &entity.RedditPost{
				UUID:  "p1000000",
				Title: "Some title",
				Likes: 5,
				Comments: []*entity.Comment{
					{
						UUID:  "c0000001",
						Body:  "Some comment text 1",
						Likes: 1,
					},
					{
						UUID:  "c0000002",
						Body:  "Some comment text 2",
						Likes: 2,
					},
				},
			},
		},
		{
			name:          "not found",
			postID:        "p1000001",
			expStatusCode: http.StatusNotFound,
		},
		{
			name:          "error",
			postID:        "p1000002",
			expStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockRepo := NewMockRedditPostsRepo(mockCtrl)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			mockRecorder := mockRepo.EXPECT().Get(gomock.AssignableToTypeOf(ctx), gomock.Eq(test.postID))

			switch test.name {
			case "success":
				mockRecorder.DoAndReturn(
					func(context.Context, string) (*entity.RedditPost, bool, error) {
						return test.responseEntity, true, nil
					},
				)
			case "not found":
				mockRecorder.DoAndReturn(
					func(context.Context, string) (*entity.RedditPost, bool, error) {
						return nil, false, nil
					},
				)
			case "error":
				mockRecorder.DoAndReturn(
					func(context.Context, string) (*entity.RedditPost, bool, error) {
						return nil, false, fmt.Errorf("some repo internal error")
					},
				)
			}

			h := handler.New(mockRepo)
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/posts/"+test.postID, nil)
			if err != nil {
				log.Fatal(err)
			}
			h.GetPost(rr, req)

			assert.Equal(t, test.expStatusCode, rr.Result().StatusCode)
		})
	}
}

func TestDeletePost(t *testing.T) {
	tests := []struct {
		name          string
		postID        string
		expStatusCode int
	}{
		{
			name:          "success",
			postID:        "p1000000",
			expStatusCode: http.StatusOK,
		},
		{
			name:          "not found",
			postID:        "p1000001",
			expStatusCode: http.StatusNotFound,
		},
		{
			name:          "error",
			postID:        "p1000002",
			expStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockRepo := NewMockRedditPostsRepo(mockCtrl)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			mockRecorder := mockRepo.EXPECT().Delete(gomock.AssignableToTypeOf(ctx), gomock.Eq(test.postID))

			switch test.name {
			case "success":
				mockRecorder.DoAndReturn(
					func(context.Context, string) (bool, error) {
						return true, nil
					},
				)
			case "not found":
				mockRecorder.DoAndReturn(
					func(context.Context, string) (bool, error) {
						return false, nil
					},
				)
			case "error":
				mockRecorder.DoAndReturn(
					func(context.Context, string) (bool, error) {
						return false, fmt.Errorf("some repo internal error")
					},
				)
			}

			h := handler.New(mockRepo)
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/posts/"+test.postID, nil)
			if err != nil {
				log.Fatal(err)
			}
			h.DeletePost(rr, req)

			assert.Equal(t, test.expStatusCode, rr.Result().StatusCode)
		})
	}

}

func TestHandlerRoute(t *testing.T) {
	tests := []struct {
		name           string
		allowedMethods []string
		request        *http.Request
		expStatusCode  int
	}{
		{
			name:           "success GET",
			allowedMethods: []string{http.MethodGet, http.MethodDelete},
			request:        httptest.NewRequest(http.MethodGet, "/post", nil),
			expStatusCode:  http.StatusOK,
		},
		{
			name:           "success DELETE",
			allowedMethods: []string{http.MethodPatch, http.MethodDelete},
			request:        httptest.NewRequest(http.MethodDelete, "/post", nil),
			expStatusCode:  http.StatusOK,
		},
		{
			name:           "success POST",
			allowedMethods: []string{http.MethodPost},
			request:        httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader([]byte("{}"))),
			expStatusCode:  http.StatusOK,
		},
		{
			name:           "not allowed POST",
			allowedMethods: []string{http.MethodGet, http.MethodDelete},
			request:        httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader([]byte("{}"))),
			expStatusCode:  http.StatusMethodNotAllowed,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockRepo := NewMockRedditPostsRepo(mockCtrl)

			switch test.name {
			case "success GET":
				mockRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, true, nil)
			case "success DELETE":
				mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(true, nil)
			case "success POST":
				mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			}

			h := handler.New(mockRepo)
			rr := httptest.NewRecorder()
			actualHandlerFunc := h.Route(test.allowedMethods...)
			actualHandlerFunc(rr, test.request)

			assert.Equal(t, test.expStatusCode, rr.Result().StatusCode)
		})
	}
}
