package cloud

type Session interface {
	Kind() string
}

type Client interface {
}

//type Cloud interface { // TODO: Implement an interface that can be an entry point to all the dependent cloud providers. See if that is possible. I don't think Go supports it as of yet.
//	Do() (string, error)
//	//BeginInteraction(ctx context.Context) (*Session, error) // TODO: when this method is defined under this interface I am not able to successfully initialize it in any of the cloud provider Compose consumer structs; the interface keeps rejecting the base cloud struct even though all of its methods are defined by the consumer struct.
//	// TODO: I think it has something to do with the fact that this method returns an interface Session.
//}

//type Config struct { // for S3 yet
//	provider string
//	credLoc  string
//	action   struct {
//		service string
//		verb    string
//		target  string
//	}
//}
//
//func (c *Config) NewCloud() (*Cloud, error) {
//	var cloud Cloud
//	switch c.provider {
//	case internal.AWS:
//		cloud = amazon.NewCompose(c.credLoc, c.action.service, c.action.verb, c.action.target)
//	}
//	return nil, nil
//}
