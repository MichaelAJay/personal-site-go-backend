package secrets

import "os"

type localSecretManagerService struct {
}

func newLocalSecretManagerService() (*localSecretManagerService, error) {
	return &localSecretManagerService{}, nil
}

func (s *localSecretManagerService) GetSecret(secretName string) (string, error) {
	return os.Getenv(secretName), nil
}
