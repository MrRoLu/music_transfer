package httpclient

// НАПИСАТЬ УНИВЕРСАЛЬНЫЙ КЛИЕНТ

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string

	rateMu      sync.Mutex
	lastRequest time.Time
	// body any
}

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

func WithHeader(key, value string) Option {
	return func(c *Client) { c.headers[key] = value }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		headers:    make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

const minRequestInterval = 5 * time.Second
const maxRetries = 5

func (c *Client) waitForRateLimit() {
	c.rateMu.Lock()
	defer c.rateMu.Unlock()

	elapsed := time.Since(c.lastRequest)
	if elapsed < minRequestInterval {
		time.Sleep(minRequestInterval - elapsed)
	}
	c.lastRequest = time.Now()
}

func (c *Client) retryRequest(header string, attempt int) time.Duration {
	if header != "" {
		if timeout, err := strconv.Atoi(header); err == nil && timeout > 0 {
			return time.Duration(timeout) * time.Second
		}
	}

	backoff := time.Duration(1<<uint(attempt)) * time.Second

	if backoff > 60*time.Second {
		backoff = 60 * time.Second
	}

	return backoff
}

func (c *Client) DoRequest(method, path string, opts ...RequestOption) ([]byte, error) {
	rc := &requestConfig{}

	for _, opt := range opts {
		if err := opt(rc); err != nil {
			return nil, err
		}
	}

	fullURL := c.baseURL + path

	// fmt.Printf("FullURL: %v, baseURL: %v, path: %v\n", fullURL, c.baseURL, path)

	if len(rc.query) > 0 {
		u, err := url.Parse(fullURL)

		if err != nil {
			return nil, err
		}

		q := u.Query()

		for k, v := range rc.query {
			q.Add(k, v)
		}

		u.RawQuery = q.Encode()

		fullURL = u.String()
	}

	// fmt.Printf("fullURL with query: %v\n", fullURL)

	req, err := http.NewRequest(method, fullURL, rc.body)

	if err != nil {
		return nil, err
	}

	if rc.contentType != "" {
		req.Header.Set("Content-Type", rc.contentType)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	for k, v := range rc.headers {
		req.Header.Set(k, v)
	}

	for attempt := 0; ; attempt++ {
		c.waitForRateLimit()

		if attempt > 0 {
			seeker, ok := req.Body.(io.Seeker)

			if !ok {
				return nil, fmt.Errorf("Body is not seekeble")
			}

			seeker.Seek(0, io.SeekStart)
		}

		resp, err := c.httpClient.Do(req)

		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			if attempt >= maxRetries {
				return nil, fmt.Errorf("apple music api error: status=429 body=%s (max retries exceeded)", string(body))
			}

			wait := c.retryRequest(resp.Header.Get("Retry-After"), attempt)

			time.Sleep(wait)
			continue
		}

		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
		}

		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("apple music api error: status=%d body=%s", resp.StatusCode, body)
		}

		return body, nil
	}
}
