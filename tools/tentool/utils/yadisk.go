package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// https://tech.yandex.ru/disk/api/reference/public-docpage/#downloads

type YADiskItem struct {
	Type      string `json:"type"`
	File      string `json:"file"`
	Path      string `json:"path"`
	PublicURL string `json:"public_url"`
}

type YADiskAnswer struct {
	PublicURL string `json:"public_url"`
	Embedded  struct {
		Items []YADiskItem `json:"items"`
	} `json:"_embedded"`
}

func (i *YADiskItem) SetPublicURL(parent string) {
	if i.PublicURL != "" {
		return
	}
	i.PublicURL = parent + i.Path
}

func YaDiskResourcesLink(u string) string {
	vals := url.Values{
		"public_key": []string{u},
	}
	return "https://cloud-api.yandex.net/v1/disk/public/resources?" + vals.Encode()
}

func YaDiskParseItems(data string) ([]YADiskItem, error) {
	var result YADiskAnswer
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	items := result.Embedded.Items
	for k := range result.Embedded.Items {
		items[k].SetPublicURL(result.PublicURL)
	}
	return items, nil
}

func YaDiskDownload(u string, location string, interactive bool, onDownloaded func(string, string) error) error {
	cl := &http.Client{}
	dl := NewDownloader(Client(cl))
	ctx := context.TODO()
	resp, err := dl.DownloadString(ctx, YaDiskResourcesLink(u))
	if err != nil {
		return err
	}
	items, err := YaDiskParseItems(resp)
	if err != nil {
		return err
	}
	for _, v := range items {
		if v.Type != "file" {
			continue
		}
		path := path.Join(location, v.Path)
		exists, err := FileExists(path)
		if err != nil {
			return err
		}
		if exists {
			fmt.Printf("File %v already exists. Skip download.\n", path)
			continue
		}
		dl := NewDownloader(
			Client(cl),
			AddTracker(
				NewInteractiveTracker(v.PublicURL, path, interactive),
			),
		)
		err = dl.WriteFile(ctx, v.File, path)
		if err != nil {
			return err
		}
		err = onDownloaded(v.PublicURL, path)
		if err != nil {
			return err
		}
	}
	return nil
}
