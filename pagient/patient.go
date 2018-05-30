package pagient

import "github.com/satori/go.uuid"

// PatientState hold the state of the Patient
type PatientState string

// enumerates all states a patient can be in
const (
	// PatientStateNew is for when the Patient is Pending
	PatientStatePending PatientState = "pending"
	// PatientStateCalled is for when the Patient's Pager has been called
	PatientStateCalled PatientState = "called"
)

// Patient API model
type Patient struct {
	ID       uuid.UUID    `json:"id"`
	Name     string       `json:"name"`
	PagerID  int          `json:"pager_id,omitempty"`
	ClientID int          `json:"client_id,omitempty"`
	Status   PatientState `json:"status"`
}

func (p *Patient) String() string {
	return p.Name
}
