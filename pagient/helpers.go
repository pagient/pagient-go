package pagient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	// UserAgent defines the user agent header sent with each request.
	userAgent = "PagientGo"
)

// Helper function for making an GET request.
func (c *Default) get(rawurl string, out interface{}) error {
	return c.do(rawurl, "GET", nil, out)
}

// Helper function for making an POST request.
func (c *Default) post(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "POST", in, out)
}

// Helper function for making an PUT request.
func (c *Default) put(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "PUT", in, out)
}

// Helper function for making an DELETE request.
func (c *Default) delete(rawurl string, in interface{}) error {
	return c.do(rawurl, "DELETE", in, nil)
}

// Helper function to make an HTTP request
func (c *Default) do(rawurl, method string, in, out interface{}) error {
	body, err := c.stream(rawurl, method, in)
	if err != nil {
		return errors.WithStack(err)
	}

	defer body.Close()

	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}

	return nil
}

// Helper function to stream an HTTP request
func (c *Default) stream(rawurl, method string, in interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var buf io.ReadWriter

	if in != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(in)
		if err != nil {
			return nil, errors.Wrap(err, "json encode")
		}
	}

	req, err := http.NewRequest(method, uri.String(), buf)
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}

	req.Header.Set(
		"User-Agent",
		userAgent,
	)

	if in != nil {
		req.Header.Set(
			"Content-Type",
			"application/json",
		)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, &httpTransportErr{err}
	}

	if resp.StatusCode > http.StatusPartialContent {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		msg := &Message{}
		err := json.Unmarshal(body, msg)
		if err != nil {
			return nil, errors.Wrap(err, "json unmarshal")
		}

		return nil, &httpResponseErr{
			msg:        msg.ErrorText,
			statusCode: msg.StatusCode,
		}
	}

	return resp.Body, nil
}
