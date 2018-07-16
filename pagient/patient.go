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
	// PatientStateFinished is for when the Patient is Finished with his medical examination
	PatientStateFinished PatientState = "finished"
)

// Patient API model
type Patient struct {
	ID       int          `json:"id"`
	Ssn      string       `json:"ssn"`
	Name     string       `json:"name"`
	PagerID  int          `json:"pagerId,omitempty"`
	ClientID int          `json:"clientId,omitempty"`
	Status   PatientState `json:"status"`
	Active   bool         `json:"active"`
}

func (p *Patient) String() string {
	return p.Name
}
