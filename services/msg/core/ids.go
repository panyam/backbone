package core

type CreateDomainRequest struct {
	*Request

	Domain     string
	StartValue int64
	Increment  int
}

type NextIDRequest struct {
	*Request
	Domain string
}

type IIDService interface {
	/**
	 * Creates a domain over which IDs can be obtained.
	 */
	CreateDomain(request *CreateDomainRequest) error

	/**
	 * Gets the next ID in a given domain
	 */
	NextID(request *NextIDRequest) (int64, error)

	/**
	 * Removes all entries.
	 */
	RemoveAllDomains(request *Request)
}
