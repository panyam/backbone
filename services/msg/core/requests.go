package core

/**
 * Base struct for all the requests that will be handled by all the operations
 * that are handled and serviced by our service!
 */
type Request struct {
	// Every request must have a requestor
	Requestor *User
}
