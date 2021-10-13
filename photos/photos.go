package photos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

func Extract(client Getter, outputDir string) {
	response, err := client.Get("https://photoslibrary.googleapis.com/v1/mediaItems")
	if err != nil {
		log.Fatalf("failed to get mediaitems: %s", err)
	}
	defer response.Body.Close()
	var medias MediaResponse
	if err := json.NewDecoder(response.Body).Decode(&medias); err != nil {
		log.Fatalf("failed to read body of mediaitems: %s", err)
	}

	fmt.Printf("%d items\n", len(medias.MediaItems))

	for _, media := range medias.MediaItems {
		if err := saveMedia(client, outputDir, media); err != nil {
			log.Fatal(err)
		}
	}
}

func saveMedia(client Getter, outputDir string, mediaItem *MediaItem) error {
	f, closer, err := openFile(outputDir, *mediaItem)
	if err != nil {
		if os.IsExist(err) {
			return nil
		} else {
			return err
		}
	}
	defer closer()

	imgResponse, err := client.Get(buildURL(mediaItem.MimeType, mediaItem.BaseUrl))
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", mediaItem.Filename, err)
	}
	defer imgResponse.Body.Close()

	imgBytes, err := ioutil.ReadAll(imgResponse.Body)
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

func openFile(outputDir string, mediaItem MediaItem) (*os.File, func(), error) {
	filePath := buildPath(outputDir, mediaItem.Metadata.CreationTime)
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create directory structure for media: %+v", err)
	}
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", filePath, mediaItem.Filename), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			info, err := os.Stat(fmt.Sprintf("%s/%s", filePath, mediaItem.Filename))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to determine if existing file was empty (%s): %+v", mediaItem.Filename, err)
			}
			if info.Size() > 0 {
				//if file as contents, assume it's complete. future support could confirm mimetype == mediaItem
				return nil, nil, err
			}
			f, _ = os.OpenFile(fmt.Sprintf("%s/%s", filePath, mediaItem.Filename), os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			return nil, nil, fmt.Errorf("failed to open file: %+v", err)
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

// buildURL based on details from https://developers.google.com/photos/library/guides/access-media-items#base-urls
func buildURL(mimeType, baseURL string) string {
	if strings.Contains(mimeType, "video") {
		return fmt.Sprintf("%s=dv", baseURL)
	}

	return fmt.Sprintf("%s=d", baseURL)
}
