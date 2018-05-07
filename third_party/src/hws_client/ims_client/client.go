package ims_client

type Imsclient struct {
	endpoint string
}

func NewImsClient(endpoint string) *Imsclient {
	return &Imsclient{
		endpoint: endpoint,
	}
}
