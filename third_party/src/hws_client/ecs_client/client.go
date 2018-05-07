package ecs_client

type Ecsclient struct {
	endpoint string
}

func NewEcsClient(endpoint string) *Ecsclient {
	return &Ecsclient{
		endpoint: endpoint,
	}
}
