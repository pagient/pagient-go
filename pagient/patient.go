package pagient

// PatientState hold the state of the Patient
type PatientState string

// enumerates all states a patient can be in
const (
	// PatientStateNew is for when the Patient is Pending
	PatientStatePending PatientState = "pending"
	// PatientStateCall is for when the  Patient's pager gets called
	PatientStateCall PatientState = "call"
	// PatientStateCalled is for when the Patient's Pager has been called
	PatientStateCalled PatientState = "called"
)

// Patient API model
type Patient struct {
	ID       int          `json:"id"`
	Ssn      string       `json:"ssn"`
	Name     string       `json:"name"`
	PagerID  int          `json:"pager_id,omitempty"`
	ClientID int          `json:"client_id,omitempty"`
	Status   PatientState `json:"status"`
	Active   bool         `json:"active"`
}

func (p *Patient) String() string {
	return p.Name
}
