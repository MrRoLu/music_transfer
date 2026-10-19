package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
)

type requestConfig struct {
	body        io.Reader
	contentType string
	query       map[string]string
	headers     map[string]string
}

type RequestOption func(*requestConfig) error

func WithJSON(body interface{}) RequestOption {
	return func(rc *requestConfig) error {
		data, err := json.Marshal(body)

		if err != nil {
			return err
		}

		rc.body = bytes.NewBuffer(data)
		rc.contentType = "application/json"

		return nil
	}
}

func WithForm(form url.Values) RequestOption { //или получать map[string][]string
	return func(rc *requestConfig) error {

		rc.body = bytes.NewBufferString(form.Encode())
		rc.contentType = "application/x-www-form-urlencoded"

		return nil
	}
}

func WithMultipart(fields map[string][]string) RequestOption {
	return func(rc *requestConfig) error {
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)

		for k, value := range fields {
			for _, v := range value {
				err := writer.WriteField(k, v)

				if err != nil {
					return err
				}
			}

		}

		if err := writer.Close(); err != nil {
			return err
		}

		rc.body = buf
		rc.contentType = writer.FormDataContentType()

		return nil
	}
}

func WithQuery(query map[string]string) RequestOption {
	return func(rc *requestConfig) error {
		if rc.query == nil {
			rc.query = make(map[string]string)
		}

		for k, v := range query {
			rc.query[k] = v
		}

		return nil
	}
}

func WithReqHeaders(key, value string) RequestOption {
	return func(rc *requestConfig) error {
		if rc.headers == nil {
			rc.headers = make(map[string]string)
		}

		rc.headers[key] = value

		return nil
	}
}
