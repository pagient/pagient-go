package pagient

// Client API model
type Client struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) String() string {
	return c.Name
}
