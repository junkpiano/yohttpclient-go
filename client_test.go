package yohttpclient

import (
	"fmt"
	"testing"
)

type Response struct {
	URL string `json:"url"`
}

func TestGet(t *testing.T) {
	url := "https://httpbin.org"

	client, err := NewClient(url, nil)

	if err != nil {
		t.Errorf("FAILED to create client!!!")
	}

	var res Response
	if err := client.Get("/get", &res); err != nil {
		t.Errorf("FAILED!!!")
	}

	fmt.Println(res.URL)
}
