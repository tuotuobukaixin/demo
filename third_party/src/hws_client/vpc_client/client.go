package vpc_client

type Vpcclient struct {
	endpoint string
}

func NewVpcClient(endpoint string) *Vpcclient {
	return &Vpcclient{
		endpoint: endpoint,
	}
}
