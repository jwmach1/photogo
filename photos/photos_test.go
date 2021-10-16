package photos_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"velocitizer.com/photogo/data"
	"velocitizer.com/photogo/photos"
	"velocitizer.com/photogo/photos/mocks"
)

//go:generate mockery --name=MediaService
func Test_Extract(t *testing.T) {
	t.Run("empty response exists", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "verify the same ctx", "value")
		service := new(mocks.MediaService)

		service.On("List", ctx, "").Return(&data.MediaResponse{}, nil)

		err := photos.Extract(ctx, service, "testdata", 1)

		assert.NoError(t, err)
	})

	t.Run("empty response exists", func(t *testing.T) {
		service := new(mocks.MediaService)

		service.On("List", context.Background(), "").Return(nil, errors.New("list fails"))

		err := photos.Extract(context.Background(), service, "testdata", 4)
		assert.Error(t, err)
	})

	t.Run("existing media on disk is skipped", func(t *testing.T) {
		defer func() {
			require.NoError(t, os.RemoveAll("testdata/2021"))
		}()
		mediaTime, err := time.Parse(time.RFC3339, "2009-05-13T15:04:05Z")
		require.NoError(t, err)

		service := new(mocks.MediaService)

		item := &data.MediaItem{
			ID:       "doesn't matter",
			Filename: "sample.txt",
			BaseUrl:  "http://localhost/bar",
			MimeType: "image/jpeg",
			Metadata: data.MediaMetadata{
				CreationTime: mediaTime,
			},
		}
		service.On("List", context.Background(), "").Return(&data.MediaResponse{
			MediaItems: []*data.MediaItem{
				item,
			},
		}, nil)

		err = photos.Extract(context.Background(), service, "testdata", 2)
		assert.NoError(t, err)
	})

	t.Run("media is passed to save", func(t *testing.T) {
		defer func() {
			require.NoError(t, os.RemoveAll("testdata/2021"))
		}()
		mediaTime, err := time.Parse(time.RFC3339, "2021-09-13T15:04:05Z")
		require.NoError(t, err)

		service := new(mocks.MediaService)

		item := &data.MediaItem{
			ID:       "doesn't matter",
			Filename: "foomedia.jpg",
			BaseUrl:  "http://localhost/bar",
			MimeType: "image/jpeg",
			Metadata: data.MediaMetadata{
				CreationTime: mediaTime,
			},
		}
		service.On("List", context.Background(), "").Return(&data.MediaResponse{
			MediaItems: []*data.MediaItem{
				item,
			},
		}, nil)
		service.On("Get", mock.Anything, *item).Return([]byte("foo"), nil)

		err = photos.Extract(context.Background(), service, "testdata", 2)
		assert.NoError(t, err)
	})

	t.Run("media with weird filename passed to save", func(t *testing.T) {
		defer func() {
			require.NoError(t, os.RemoveAll("testdata/2014"))
		}()
		mediaTime, err := time.Parse(time.RFC3339, "2014-07-21T15:04:05Z")
		require.NoError(t, err)

		service := new(mocks.MediaService)

		item := &data.MediaItem{
			ID:       "doesn't matter",
			Filename: "7/21/14 - 1",
			BaseUrl:  "http://localhost/bar",
			MimeType: "image/jpeg",
			Metadata: data.MediaMetadata{
				CreationTime: mediaTime,
			},
		}
		service.On("List", context.Background(), "").Return(&data.MediaResponse{
			MediaItems: []*data.MediaItem{
				item,
			},
		}, nil)
		service.On("Get", mock.Anything, *item).Return([]byte("foo"), nil)

		err = photos.Extract(context.Background(), service, "testdata", 2)
		assert.NoError(t, err)
	})

	t.Run("loop twice and exit", func(t *testing.T) {
		defer func() {
			require.NoError(t, os.RemoveAll("testdata/2021"))
		}()
		service := new(mocks.MediaService)
		service.Test(t)

		mediaTime, err := time.Parse(time.RFC3339, "2021-10-13T15:04:05Z")
		require.NoError(t, err)
		item := &data.MediaItem{
			ID:       "doesn't matter",
			Filename: "foomedia.jpg",
			BaseUrl:  "http://localhost/bar",
			MimeType: "image/jpbeg",
			Metadata: data.MediaMetadata{
				CreationTime: mediaTime,
			},
		}
		service.On("List", context.Background(), "").Return(&data.MediaResponse{
			MediaItems: []*data.MediaItem{
				item,
			},
			NextPageToken: "nextpage",
		}, nil).Once()
		service.On("List", context.Background(), "nextpage").Return(&data.MediaResponse{
			MediaItems:    []*data.MediaItem{},
			NextPageToken: "",
		}, nil).Once()
		service.On("Get", mock.Anything, *item).Return([]byte("foo"), nil)

		err = photos.Extract(context.Background(), service, "testdata", 3)
		assert.NoError(t, err)
		service.AssertExpectations(t)
	})
}
