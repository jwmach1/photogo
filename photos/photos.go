package photos

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"velocitizer.com/photogo/data"
)

type MediaService interface {
	List(ctx context.Context, nextPageToken string) (*data.MediaResponse, error)
	Get(ctx context.Context, mediaItem data.MediaItem) ([]byte, error)
}

func Extract(ctx context.Context, client MediaService, outputDir string, workerCount int, readOnly bool) error {
	var total int64
	var nextPageToken string
	for {
		medias, err := client.List(ctx, nextPageToken)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return fmt.Errorf("failed to get mediaitems: %s", err)
		}
		total += int64(len(medias.MediaItems))
		fmt.Printf("%d items, has more %t\n", len(medias.MediaItems), len(medias.NextPageToken) > 0)
		eg, ctx := errgroup.WithContext(ctx)
		eg.SetLimit(workerCount)
		for _, media := range medias.MediaItems {
			media := *media
			eg.Go(func() error {
				if readOnly {
					fmt.Printf("%s/%s\n", buildPath(outputDir, media.Metadata.CreationTime), nameCleaner(media.Filename))
					return nil
				}
				return saveMedia(ctx, client, outputDir, media)
			})
		}
		err = eg.Wait()
		if err != nil {
			return err
		}
		nextPageToken = medias.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	p := message.NewPrinter(language.English)
	p.Printf("%d media processed\n", total)
	return nil
}

func saveMedia(ctx context.Context, client MediaService, outputDir string, mediaItem data.MediaItem) error {
	f, closer, err := openFile(outputDir, mediaItem)
	if err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}
	defer closer()

	imgBytes, err := client.Get(ctx, mediaItem)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", mediaItem.Filename, err)
	}
	count, err := f.Write(imgBytes)
	if err != nil {
		return fmt.Errorf("failed to write %s: %v", mediaItem.Filename, err)
	}
	fmt.Printf("wrote %s (%s) of %d\n", mediaItem.Filename, mediaItem.MimeType, count)

	return nil
}

func openFile(outputDir string, mediaItem data.MediaItem) (*os.File, func(), error) {
	emptyCloser := func() {}
	filePath := buildPath(outputDir, mediaItem.Metadata.CreationTime)
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return nil, emptyCloser, fmt.Errorf("failed to create directory structure for media: %+v", err)
	}
	cleanName := nameCleaner(mediaItem.Filename)
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", filePath, cleanName), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			info, err := os.Stat(fmt.Sprintf("%s/%s", filePath, cleanName))
			if err != nil {
				return nil, emptyCloser, fmt.Errorf("failed to determine if existing file was empty (%s): %+v", mediaItem.Filename, err)
			}
			if info.Size() > 0 {
				//if file has contents, assume it's complete. future support could confirm mimetype == mediaItem
				return nil, emptyCloser, os.ErrExist
			}
			f, _ = os.OpenFile(fmt.Sprintf("%s/%s", filePath, mediaItem.Filename), os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			return nil, emptyCloser, fmt.Errorf("failed to open file for media item %+v: %+v", mediaItem, err)
		}
	}

	return f, func() {
		f.Close()
		os.Chtimes(fmt.Sprintf("%s/%s", filePath, cleanName), mediaItem.Metadata.CreationTime, mediaItem.Metadata.CreationTime)
	}, nil
}
func nameCleaner(input string) string {
	output := strings.ReplaceAll(input, " ", "_")
	output = strings.ReplaceAll(input, "/", "_")
	return output
}

// buildPath returns the path where the media should be written
func buildPath(outputDir string, createdAt time.Time) string {
	year := createdAt.Year()
	month := createdAt.Month()
	return fmt.Sprintf("%s/%d/%02d", outputDir, year, month)
}
