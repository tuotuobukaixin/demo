package iam_client

type Iamclient struct {
	endpoint string
}

func NewIamClient(endpoint string) *Iamclient {
	return &Iamclient{
		endpoint: endpoint,
	}
}
