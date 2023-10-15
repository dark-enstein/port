package auth

type Options interface {
	IsValid() bool
	//IsCustom() bool // TODO impl later
}

type Director interface {
	// PingDependencies pings of the relevant dependencies of the reference struct,
	// thereby confirming that the Director can proceed
	PingDependencies() (bool, error)
	// IsEmpty checks if the reference struct has been set up with the request variables or not
	IsEmpty() bool
}
