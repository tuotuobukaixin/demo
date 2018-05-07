package evs_client

type Evsclient struct {
	endpoint string
}

func NewEvsClient(endpoint string) *Evsclient {
	return &Evsclient{
		endpoint: endpoint,
	}
}
