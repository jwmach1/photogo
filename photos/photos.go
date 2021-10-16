package photos

import (
	"fmt"
	"os"
	"time"

	"velocitizer.com/photogo/data"
)

type MediaService interface {
	List(nextPageToken string) (*data.MediaResponse, error)
	Get(mediaItem data.MediaItem) ([]byte, error)
}

func Extract(client MediaService, outputDir string) error {
	var nextPageToken string
	for {
		medias, err := client.List(nextPageToken)
		if err != nil {
			return fmt.Errorf("failed to get mediaitems: %s", err)
		}
		fmt.Printf("%d items, has more %t\n", len(medias.MediaItems), len(medias.NextPageToken) > 0)

		for _, media := range medias.MediaItems {
			if err := saveMedia(client, outputDir, *media); err != nil {
				return err
			}
		}
		nextPageToken = medias.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return nil
}

func saveMedia(client MediaService, outputDir string, mediaItem data.MediaItem) error {
	f, closer, err := openFile(outputDir, mediaItem)
	if err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}
	defer closer()

	imgBytes, err := client.Get(mediaItem)
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
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", filePath, mediaItem.Filename), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			info, err := os.Stat(fmt.Sprintf("%s/%s", filePath, mediaItem.Filename))
			if err != nil {
				return nil, emptyCloser, fmt.Errorf("failed to determine if existing file was empty (%s): %+v", mediaItem.Filename, err)
			}
			if info.Size() > 0 {
				//if file as contents, assume it's complete. future support could confirm mimetype == mediaItem
				return nil, emptyCloser, os.ErrExist
			}
			f, _ = os.OpenFile(fmt.Sprintf("%s/%s", filePath, mediaItem.Filename), os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			return nil, emptyCloser, fmt.Errorf("failed to open file: %+v", err)
		}
	}

	return f, func() {
		f.Close()
		os.Chtimes(fmt.Sprintf("%s/%s", filePath, mediaItem.Filename), mediaItem.Metadata.CreationTime, mediaItem.Metadata.CreationTime)
	}, nil
}

// buildPath returns the path where the media should be written
func buildPath(outputDir string, createdAt time.Time) string {
	year := createdAt.Year()
	month := createdAt.Month()
	return fmt.Sprintf("%s/%d/%02d", outputDir, year, month)
}
