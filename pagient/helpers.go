package pagient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	// UserAgent defines the user agent header sent with each request.
	UserAgent = "PagientGo"
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
	body, err := c.stream(
		rawurl,
		method,
		in,
		out,
	)

	if err != nil {
		return err
	}

	defer body.Close()

	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}

	return nil
}

// Helper function to stream an HTTP request
func (c *Default) stream(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)

	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter

	if in != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(in)

		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, uri.String(), buf)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		UserAgent,
	)

	if in != nil {
		req.Header.Set(
			"Content-Type",
			"application/json",
		)
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode > http.StatusPartialContent {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)

		msg := &Message{}
		parse := json.Unmarshal(out, msg)

		if parse != nil {
			return nil, fmt.Errorf(string(out))
		}

		return nil, fmt.Errorf(msg.Message)
	}

	return resp.Body, nil
}

// IsBadRequestErr returns whether it's a 400 or not
func IsBadRequestErr(err error) bool {
	return strings.TrimSpace(err.Error()) == http.StatusText(http.StatusBadRequest)
}

// IsUnauthorizedErr returns whether it's a 401 or not
func IsUnauthorizedErr(err error) bool {
	return strings.TrimSpace(err.Error()) == http.StatusText(http.StatusUnauthorized)
}

// IsNotFoundErr returns whether it's a 404 or not
func IsNotFoundErr(err error) bool {
	return strings.TrimSpace(err.Error()) == http.StatusText(http.StatusNotFound)
}

// IsConflictErr returns whether it's a 409 or not
func IsConflictErr(err error) bool {
	return strings.TrimSpace(err.Error()) == http.StatusText(http.StatusConflict)
}

// IsUnprocessableEntityErr returns whether it's a 422 or not
func IsUnprocessableEntityErr(err error) bool {
	return strings.TrimSpace(err.Error()) == http.StatusText(http.StatusUnprocessableEntity)
}

// IsInternalServerErr returns whether it's a 500 or not
func IsInternalServerErr(err error) bool {
	return strings.TrimSpace(err.Error()) == http.StatusText(http.StatusInternalServerError)
}

// IsGatewayTimeoutErr returns whether it's a 504 or not
func IsGatewayTimeoutErr(err error) bool {
	return strings.TrimSpace(err.Error()) == http.StatusText(http.StatusGatewayTimeout)
}
