package argocd

import (
	"context"
	"io"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	argoio "github.com/argoproj/argo-cd/v2/util/io"
)

type API struct {
	client     applicationpkg.ApplicationServiceClient
	connection io.Closer
}

type APIOptions struct {
	Address string
	Token   string
}

func NewAPI(options APIOptions) API {
	clientOptions := argocdclient.ClientOptions{
		ServerAddr: options.Address,
		AuthToken:  options.Token,
		GRPCWeb:    true,
	}

	connection, client := argocdclient.NewClientOrDie(&clientOptions).NewApplicationClientOrDie()

	return API{client: client, connection: connection}
}

func (a API) Sync(appName string) error {
	request := applicationpkg.ApplicationSyncRequest{
		Name:  &appName,
		Prune: true,
	}

	_, err := a.client.Sync(context.Background(), &request)
	if err != nil {
		return err
	}

	defer argoio.Close(a.connection)

	return nil
}
