package db

type DBClient struct {
}

func GetDB() (*DBClient, error) {

	return &DBClient{}, nil
}
