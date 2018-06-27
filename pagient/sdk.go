package pagient

import (
	"net/http"
	"time"
	"fmt"
)

const (
	pathClients     = "%s/clients"
	pathPagers      = "%s/pagers"
	pathPatients    = "%s/patients"
	pathPatient     = "%s/patients/%v"
)

// ClientAPI describes a client API
type ClientAPI interface {
	// ClientList returns a list of all clients
	ClientList() ([]*Client, error)

	// PagerList returns a list of all pagers
	PagerList() ([]*Pager, error)

	// PatientList returns a list of all patients
	PatientList() ([]*Patient, error)

	// PatientGet returns a patient
	PatientGet(string) (*Patient, error)

	// PatientAdd adds a patient
	PatientAdd(*Patient) (*Patient, error)
}

// Default implements the client interface
type Default struct {
	client      *http.Client
	base        string
	username    string
	password    string
}

func NewClient(uri string, username, password string) ClientAPI {
	return &Default{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		base:     uri,
		username: username,
		password: password,
	}
}

func (c *Default) ClientList() ([]*Client, error) {
	var out []*Client

	uri := fmt.Sprintf(pathClients, c.base)
	err := c.get(uri, &out)

	return out, err
}

func (c *Default) PagerList() ([]*Pager, error) {
	var out []*Pager

	uri := fmt.Sprintf(pathPagers, c.base)
	err := c.get(uri, &out)

	return out, err
}

func (c *Default) PatientList() ([]*Patient, error) {
	var out []*Patient

	uri := fmt.Sprintf(pathPatient, c.base)
	err := c.get(uri, &out)

	return out, err
}

func (c *Default) PatientGet(id string) (*Patient, error) {
	out := &Patient{}

	uri := fmt.Sprintf(pathPatient, c.base, id)
	err := c.get(uri, out)

	return out, err
}

func (c *Default) PatientAdd(patient *Patient) (*Patient, error) {
	out := &Patient{}

	uri := fmt.Sprintf(pathPatients, c.base)
	err := c.post(uri, patient, out)

	return out, err
}
