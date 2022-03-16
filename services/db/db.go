package db

type DBClient struct {
}

func NewDB() (*DBClient, error) {

	return &DBClient{}, nil
}
