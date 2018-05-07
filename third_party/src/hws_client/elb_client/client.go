package elb_client

type Elbclient struct {
	endpoint string
}

func NewElbClient(endpoint string) *Elbclient {
	return &Elbclient{
		endpoint: endpoint,
	}
}
