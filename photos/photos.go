package photos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	imgResponse, err := client.Get(mediaItem.BaseUrl)
	if err != nil {
		return fmt.Errorf("failed to get %s: %v", mediaItem.Filename, err)
	}
	defer imgResponse.Body.Close()

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
	fmt.Printf("wrote %s of %d", name, count)
	if f.Close() != nil {
		return fmt.Errorf("failed to close file %s: %s", name, err)
	}
	os.Chtimes(name, mediaItem.Metadata.CreationTime, mediaItem.Metadata.CreationTime)
	return nil
}
