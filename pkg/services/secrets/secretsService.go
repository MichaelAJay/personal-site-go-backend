package secrets

import "os"

// SecretManager defines the interface for fetching secrets
type SecretManager interface {
	GetSecret(secretName string) (string, error)
}

type SecretManagerService struct {
	secretManager SecretManager
	local         bool
	prefix        string
}

func NewSecretManagerService() (*SecretManagerService, error) {
	var secretManager SecretManager
	var err error
	env := os.Getenv("ENV")
	local := env == "local"
	prefix := ""

	if !local {
		secretManager, err = newGoogleSecretManagerService()
		if err != nil {
			return nil, err
		}

		switch env {
		case "development":
			prefix = "GO_PERSONAL_DEV_"
		case "production":
			prefix = "GO_PERSONAL_PROD_"
		}

	} else {
		secretManager, err = newLocalSecretManagerService()
		if err != nil {
			return nil, err
		}
	}
	return &SecretManagerService{
		secretManager: secretManager,
		local:         local,
		prefix:        prefix,
	}, nil
}

func (s *SecretManagerService) GetSecret(secretName string) (string, error) {
	return s.secretManager.GetSecret(s.prefix + secretName)
}
