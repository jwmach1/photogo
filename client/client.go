package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"velocitizer.com/photogo/data"
)

type Getter func(*http.Request) (resp *http.Response, err error)

type Client struct {
	getter Getter
}

const pageSize = "25"

func New(getter Getter) *Client {
	return &Client{getter: getter}
}

func (c Client) List(nextPageToken string) (*data.MediaResponse, error) {
	values := url.Values{}
	values.Set("pageSize", pageSize)
	if nextPageToken != "" {
		values.Set("pageToken", nextPageToken)
	}
	url := fmt.Sprintf("%s?%s", "https://photoslibrary.googleapis.com/v1/mediaItems", values.Encode())
	get, _ := http.NewRequest("GET", url, nil)
	response, err := c.getter(get)
	if err != nil {
		if response != nil && response.Body != nil {
			defer response.Body.Close()
			b, _ := ioutil.ReadAll(response.Body)
			fmt.Println("body from error:", string(b))
		}
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		if response.Body != nil {
			defer response.Body.Close()
			b, _ := ioutil.ReadAll(response.Body)
			fmt.Println("body from error:", string(b))
		}
		return nil, fmt.Errorf("list call returned: %d:%s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	defer response.Body.Close()
	var medias data.MediaResponse
	if err := json.NewDecoder(response.Body).Decode(&medias); err != nil {
		return nil, err
	}
	return &medias, nil
}

func (c Client) Get(mediaItem data.MediaItem) ([]byte, error) {
	get, _ := http.NewRequest("GET", buildURL(mediaItem.MimeType, mediaItem.BaseUrl), nil)
	imgResponse, err := c.getter(get)
	if err != nil {
		return nil, fmt.Errorf("failed to get (%s): %v", mediaItem.ID, err)
	}
	if imgResponse.StatusCode != http.StatusOK {
		if imgResponse.Body != nil {
			defer imgResponse.Body.Close()
			b, _ := ioutil.ReadAll(imgResponse.Body)
			fmt.Println("body from error:", string(b))
		}
		return nil, fmt.Errorf("list call returned: %d:%s", imgResponse.StatusCode, http.StatusText(imgResponse.StatusCode))
	}
	defer imgResponse.Body.Close()

	return ioutil.ReadAll(imgResponse.Body)
}

// buildURL based on details from https://developers.google.com/photos/library/guides/access-media-items#base-urls
func buildURL(mimeType, baseURL string) string {
	if strings.Contains(mimeType, "video") {
		return fmt.Sprintf("%s=dv", baseURL)
	}

	return fmt.Sprintf("%s=d", baseURL)
}
