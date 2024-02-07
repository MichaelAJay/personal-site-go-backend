package secrets

import (
	"context"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type googleSecretManagerService struct {
	client      *secretmanager.Client
	projectName string
}

func newGoogleSecretManagerService() (*googleSecretManagerService, error) {
	projectName := os.Getenv("GOOGLE_PROJECT")

	client, err := secretmanager.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	return &googleSecretManagerService{
		client:      client,
		projectName: projectName,
	}, nil
}

func (s *googleSecretManagerService) GetSecret(secretName string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/" + s.projectName + "/secrets/" + secretName + "/versions/latest",
	}

	result, err := s.client.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", err
	}

	return string(result.Payload.Data), nil
}
