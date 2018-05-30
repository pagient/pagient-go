package pagient

// Pager API model
type Pager struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (p *Pager) String() string {
	return p.Name
}
