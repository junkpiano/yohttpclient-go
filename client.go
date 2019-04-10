package yohttpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

var version = "1.0.0-alpha"
var userAgent = fmt.Sprintf("YOHTTPClient/%s (%s)", version, runtime.Version())

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Logger     *log.Logger
}

func NewClient(urlStr string, logger *log.Logger) (*Client, error) {
	parsedURL, err := url.ParseRequestURI(urlStr)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", urlStr)
	}

	if logger == nil {
		var defaultLogger = log.New(ioutil.Discard, "", log.LstdFlags)
		logger = defaultLogger
	}

	httpClient := new(http.Client)

	client := &Client{URL: parsedURL, HTTPClient: httpClient, Logger: logger}

	return client, nil
}

func (c *Client) Get(path string, out interface{}) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := c.newRequest(ctx, "GET", path, nil)

	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	return decodeBody(resp, out)
}

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {

	u := c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
