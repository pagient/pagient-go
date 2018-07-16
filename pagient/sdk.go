package pagient

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jackspirou/syscerts"
	"golang.org/x/oauth2"
)

const (
	pathAuthLogin    = "%s/oauth/token"
	pathAuthSessions = "%s/oauth/sessions"
	pathClients      = "%s/api/clients"
	pathPagers       = "%s/api/pagers"
	pathPatients     = "%s/api/patients"
	pathPatient      = "%s/api/patients/%v"
)

// ClientAPI describes a client API
type ClientAPI interface {
	// SetClient sets the default http client. This should
	// be used in conjunction with golang.org/x/oauth2 to
	// authenticate requests to the API.
	SetClient(client *http.Client)

	// IsAuthenticated checks if we already provided an authentication
	// token for our client requests. If it returns false you can update
	// the client after fetching a valid token.
	IsAuthenticated() bool

	// AuthLogin signs in based on credentials and returns a token.
	AuthLogin(string, string) (*Token, error)

	// ClientList returns a list of all clients
	ClientList() ([]*Client, error)

	// PagerList returns a list of all pagers
	PagerList() ([]*Pager, error)

	// PatientList returns a list of all patients
	PatientList() ([]*Patient, error)

	// PatientGet returns a patient
	PatientGet(int) (*Patient, error)

	// PatientAdd adds a patient
	PatientAdd(*Patient) error

	// PatientUpdate updates a patient
	PatientUpdate(*Patient) error

	// PatientRemove removes a patient
	PatientRemove(int) error
}

// Default implements the client interface
type Default struct {
	client *http.Client
	base   string
	token  string
}

// NewClient returns a default client
func NewClient(uri string) ClientAPI {
	return &Default{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		base: uri,
	}
}

// NewTokenClient returns a client that authenticates
// all outbound requests with the given token.
func NewTokenClient(uri, token string) ClientAPI {
	config := oauth2.Config{}

	client := config.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: token,
		},
	)

	if trans, ok := client.Transport.(*oauth2.Transport); ok {
		trans.Base = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				RootCAs: syscerts.SystemRootsPool(),
			},
		}
	}

	return &Default{
		client: client,
		base:   uri,
		token:  token,
	}
}

// IsAuthenticated checks if we already provided an authentication
// token for our client requests. If it returns false you can update
// the client after fetching a valid token.
func (c *Default) IsAuthenticated() bool {
	if c.token == "" {
		return false
	}

	uri, err := url.Parse(fmt.Sprintf(pathAuthSessions, c.base))
	if err != nil {
		return false
	}

	req, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return false
	}

	req.Header.Set(
		"User-Agent",
		UserAgent,
	)

	resp, err := c.client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// SetClient sets the default http client. This should
// be used in conjunction with golang.org/x/oauth2 to
// authenticate requests to the API.
func (c *Default) SetClient(client *http.Client) {
	c.client = client
}

// AuthLogin signs in based on credentials and returns a token.
func (c *Default) AuthLogin(username, password string) (*Token, error) {
	out := &Token{}

	in := struct {
		Username string
		Password string
	}{
		username,
		password,
	}

	uri := fmt.Sprintf(pathAuthLogin, c.base)
	err := c.post(uri, in, out)

	return out, err
}

// ClientList returns a list of clients
func (c *Default) ClientList() ([]*Client, error) {
	var out []*Client

	uri := fmt.Sprintf(pathClients, c.base)
	err := c.get(uri, &out)

	return out, err
}

// PagerList returns a list of pagers
func (c *Default) PagerList() ([]*Pager, error) {
	var out []*Pager

	uri := fmt.Sprintf(pathPagers, c.base)
	err := c.get(uri, &out)

	return out, err
}

// PatientList returns a list of patients
func (c *Default) PatientList() ([]*Patient, error) {
	var out []*Patient

	uri := fmt.Sprintf(pathPatients, c.base)
	err := c.get(uri, &out)

	return out, err
}

// PatientGet returns a patient by ID
func (c *Default) PatientGet(id int) (*Patient, error) {
	out := &Patient{}

	uri := fmt.Sprintf(pathPatient, c.base, id)
	err := c.get(uri, out)

	return out, err
}

// PatientAdd adds a patient
func (c *Default) PatientAdd(patient *Patient) error {
	uri := fmt.Sprintf(pathPatients, c.base)
	err := c.post(uri, patient, patient)

	return err
}

// PatientUpdate updates a patient
func (c *Default) PatientUpdate(patient *Patient) error {
	uri := fmt.Sprintf(pathPatient, c.base, patient.ID)
	err := c.post(uri, patient, patient)

	return err
}

// PatientRemove removes a patient by ID
func (c *Default) PatientRemove(id int) error {
	uri := fmt.Sprintf(pathPatient, c.base, id)
	err := c.delete(uri, nil)

	return err
}
