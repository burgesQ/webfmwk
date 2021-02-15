package webfmwk

type RequestID interface {
	// GetRequest return the current request ID
	GetRequestID() string

	// SetRequest set the id of the current request
	SetRequestID(id string) Context
}

// GetRequestID implement RequestID and Context by returning
// the context unique request identifier.
func (c *icontext) GetRequestID() string {
	return c.uid
}

// SetRequestID implement RequestID and Context by persisting
// the context unique request identifier.
func (c *icontext) SetRequestID(id string) Context {
	c.uid = id
	return c
}
