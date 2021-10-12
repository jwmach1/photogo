package photos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func Extract(client *http.Client, outputDir string) {
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

func saveMedia(client *http.Client, outputDir string, mediaItem *MediaItem) error {
	imgResponse, err := client.Get(buildURL(mediaItem.MimeType, mediaItem.BaseUrl))
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", mediaItem.Filename, err)
	}
	defer imgResponse.Body.Close()
	// for k, v := range imgResponse.Header {
	// 	fmt.Printf("%s:\t%s=%s\n", mediaItem.Filename, k, v)
	// }

	name := fmt.Sprintf("%s/%s", outputDir, mediaItem.Filename)
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	imgBytes, err := ioutil.ReadAll(imgResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", mediaItem.Filename, err)
	}
	count, err := f.Write(imgBytes)
	if err != nil {
		return fmt.Errorf("failed to write %s: %v", mediaItem.Filename, err)
	}
	fmt.Printf("wrote %s (%s) of %d\n", name, mediaItem.MimeType, count)
	if f.Close() != nil {
		return fmt.Errorf("failed to close file %s: %s", name, err)
	}
	os.Chtimes(name, mediaItem.Metadata.CreationTime, mediaItem.Metadata.CreationTime)
	return nil
}

// buildURL based on details from https://developers.google.com/photos/library/guides/access-media-items#base-urls
func buildURL(mimeType, baseURL string) string {
	if strings.Contains(mimeType, "video") {
		fmt.Println(baseURL)
		return fmt.Sprintf("%s=dv", baseURL)
	}

	return fmt.Sprintf("%s=d", baseURL)
}
