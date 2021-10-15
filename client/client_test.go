package client_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"velocitizer.com/photogo/client"
	"velocitizer.com/photogo/client/mocks"
	"velocitizer.com/photogo/data"
)

//go:generate mockery --name=Getter
func TestClient_List(t *testing.T) {
	t.Run("given empty page token", func(t *testing.T) {
		response := httptest.NewRecorder()
		response.Body = bytes.NewBuffer([]byte(`{"nextPageToken":"foopagetoken"}`))
		getter := new(mocks.Getter)
		getter.Test(t)
		getter.On("Execute", mock.MatchedBy(func(r *http.Request) bool {
			return r.URL.String() == "https://photoslibrary.googleapis.com/v1/mediaItems"
		})).Return(response.Result(), nil)

		actual, err := client.New(getter.Execute).List("")
		assert.NoError(t, err)
		expected := &data.MediaResponse{NextPageToken: "foopagetoken"}
		assert.Equal(t, expected, actual)

		getter.AssertExpectations(t)
	})
	t.Run("given a page token", func(t *testing.T) {
		response := httptest.NewRecorder()
		response.Body = bytes.NewBuffer([]byte(`{}`))
		getter := new(mocks.Getter)
		getter.Test(t)
		body := make(map[string]interface{})
		getter.On("Execute", mock.MatchedBy(func(r *http.Request) bool {
			json.NewDecoder(r.Body).Decode(&body)
			return r.URL.String() == "https://photoslibrary.googleapis.com/v1/mediaItems"
		})).Return(response.Result(), nil)

		client.New(getter.Execute).List("foopagetoken")
		getter.AssertExpectations(t)
		assert.Equal(t, "foopagetoken", body["pageToken"])
		assert.Equal(t, "50", body["pageSize"])
	})
	t.Run("bad content from REST body", func(t *testing.T) {
		response := httptest.NewRecorder()
		response.Body = bytes.NewBuffer([]byte(`thi}s is not { jason }`))
		getter := new(mocks.Getter)
		getter.Test(t)
		getter.On("Execute", mock.Anything).Return(response.Result(), nil)

		_, err := client.New(getter.Execute).List("foopagetoken")
		assert.Error(t, err)
	})
	t.Run("REST call fail", func(t *testing.T) {
		getter := new(mocks.Getter)
		getter.Test(t)
		getter.On("Execute", mock.Anything).Return(nil, errors.New("some error"))

		_, err := client.New(getter.Execute).List("foopagetoken")

		assert.EqualError(t, err, "some error")
	})
	t.Run("REST call fail w/ body that describes error", func(t *testing.T) {
		getter := new(mocks.Getter)
		getter.Test(t)
		response := httptest.NewRecorder()
		response.Body = bytes.NewBuffer([]byte(`{}`))
		getter.On("Execute", mock.Anything).Return(response.Result(), errors.New("some error"))

		_, err := client.New(getter.Execute).List("foopagetoken")

		assert.EqualError(t, err, "some error")
	})
	t.Run("REST call fail w/ non-200 error response", func(t *testing.T) {
		getter := new(mocks.Getter)
		getter.Test(t)
		response := httptest.NewRecorder()
		response.Body = bytes.NewBuffer([]byte(`{}`))
		response.Result().StatusCode = http.StatusFailedDependency
		getter.On("Execute", mock.Anything).Return(response.Result(), nil)

		_, err := client.New(getter.Execute).List("foopagetoken")

		assert.EqualError(t, err, "list call returned: 424:"+http.StatusText(http.StatusFailedDependency))
	})
}

func TestClient_Get(t *testing.T) {
	t.Run("get jpeg image", func(t *testing.T) {
		response := httptest.NewRecorder()
		response.Body = bytes.NewBuffer([]byte(`contents of the file`))
		getter := new(mocks.Getter)
		getter.Test(t)
		getter.On("Execute", mock.MatchedBy(func(r *http.Request) bool {
			return r.URL.String() == "https://base/url=d"
		})).Return(response.Result(), nil)

		actual, err := client.New(getter.Execute).Get(data.MediaItem{MimeType: "image/jpeg", BaseUrl: "https://base/url"})
		assert.NoError(t, err)
		assert.Equal(t, "contents of the file", string(actual))

		getter.AssertExpectations(t)
	})
	t.Run("get video", func(t *testing.T) {
		response := httptest.NewRecorder()
		response.Body = bytes.NewBuffer([]byte(`contents of the file`))
		getter := new(mocks.Getter)
		getter.Test(t)
		getter.On("Execute", mock.MatchedBy(func(r *http.Request) bool {
			return r.URL.String() == "https://base/videourl=dv"
		})).Return(response.Result(), nil)

		actual, err := client.New(getter.Execute).Get(data.MediaItem{MimeType: "video/mpeg", BaseUrl: "https://base/videourl"})
		assert.NoError(t, err)
		assert.Equal(t, "contents of the file", string(actual))

		getter.AssertExpectations(t)
	})
	t.Run("get error", func(t *testing.T) {
		getter := new(mocks.Getter)
		getter.Test(t)
		getter.On("Execute", mock.Anything).Return(nil, errors.New("expected"))

		_, err := client.New(getter.Execute).Get(data.MediaItem{ID: "the_id", MimeType: "video/mpeg", BaseUrl: "https://base/videourl"})

		assert.EqualError(t, err, "failed to get (the_id): expected")
	})
}
