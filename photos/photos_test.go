package photos_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"velocitizer.com/photogo/data"
	"velocitizer.com/photogo/photos"
	"velocitizer.com/photogo/photos/mocks"
)

//go:generate mockery --name=MediaService
func Test_Extract(t *testing.T) {
	t.Run("empty response exists", func(t *testing.T) {
		service := new(mocks.MediaService)

		service.On("List", "").Return(&data.MediaResponse{}, nil)

		err := photos.Extract(service, "testdata")

		assert.NoError(t, err)
	})

	t.Run("empty response exists", func(t *testing.T) {
		service := new(mocks.MediaService)

		service.On("List", "").Return(nil, errors.New("list fails"))

		err := photos.Extract(service, "testdata")
		assert.Error(t, err)
	})

	t.Run("media is passed to save", func(t *testing.T) {
		service := new(mocks.MediaService)

		item := &data.MediaItem{
			ID:       "doesn't matter",
			Filename: "foomedia.jpg",
			BaseUrl:  "http://localhost/bar",
			MimeType: "image/jpbeg",
			Metadata: data.MediaMetadata{
				CreationTime: time.Now(),
			},
		}
		service.On("List", "").Return(&data.MediaResponse{
			MediaItems: []*data.MediaItem{
				item,
			},
		}, nil)
		service.On("Get", *item).Return([]byte("foo"), nil)

		err := photos.Extract(service, "testdata")
		assert.NoError(t, err)
	})

	t.Run("loop twice and exit", func(t *testing.T) {
		require.NoError(t, os.RemoveAll("testdata/2021"))
		service := new(mocks.MediaService)
		service.Test(t)

		item := &data.MediaItem{
			ID:       "doesn't matter",
			Filename: "foomedia.jpg",
			BaseUrl:  "http://localhost/bar",
			MimeType: "image/jpbeg",
			Metadata: data.MediaMetadata{
				CreationTime: time.Now(),
			},
		}
		service.On("List", "").Return(&data.MediaResponse{
			MediaItems: []*data.MediaItem{
				item,
			},
			NextPageToken: "nextpage",
		}, nil).Once()
		service.On("List", "nextpage").Return(&data.MediaResponse{
			MediaItems:    []*data.MediaItem{},
			NextPageToken: "",
		}, nil).Once()
		service.On("Get", *item).Return([]byte("foo"), nil)

		err := photos.Extract(service, "testdata")
		assert.NoError(t, err)
		service.AssertExpectations(t)
	})
}
