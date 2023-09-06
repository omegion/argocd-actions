package argocd

import (
	"context"
	"io"
	"strings"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"

	argoio "github.com/argoproj/argo-cd/v2/util/io"
)

// Interface is an interface for API.
type Interface interface {
	Sync(appName string) error
	SyncWithLabels(labels string) error
}

// API is struct for ArgoCD api.
type API struct {
	client     applicationpkg.ApplicationServiceClient
	connection io.Closer
}

// APIOptions is options for API.
type APIOptions struct {
	Address string
	Token   string
}

// NewAPI creates new API.
func NewAPI(options *APIOptions) API {
	clientOptions := argocdclient.ClientOptions{
		ServerAddr: options.Address,
		AuthToken:  options.Token,
		GRPCWeb:    true,
	}

	connection, client := argocdclient.NewClientOrDie(&clientOptions).NewApplicationClientOrDie()

	return API{client: client, connection: connection}
}

// Sync syncs given application.
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

// SyncWithLabels syncs applications based on provided labels.
func (a API) SyncWithLabels(labels string) error {
	// 1. Fetch all applications
	listResponse, err := a.client.List(context.Background(), &applicationpkg.ApplicationQuery{})
	if err != nil {
		return err
	}

	// 2. Iterate through each application, check labels, and sync if it matches
	for _, app := range listResponse.Items {
		if matchesLabels(&app, labels) {
			err := a.Sync(app.Name)
			if err != nil {
				// Decide how you want to handle individual sync errors
				// E.g., log the error and continue, or return
			}
		}
	}

	defer argoio.Close(a.connection)
	return nil
}

// matchesLabels checks if an application has the specified labels.
func matchesLabels(app *v1alpha1.Application, labelsStr string) bool {
    // Split the string into individual label key-value pairs.
    pairs := strings.Split(labelsStr, ",")

    appLabels := app.ObjectMeta.Labels

    for _, pair := range pairs {
        keyValue := strings.Split(pair, "=")
        if len(keyValue) != 2 {
          	// Malformed label string; skip this pair and continue.
            continue
        }
        key, value := keyValue[0], keyValue[1]

        appValue, exists := appLabels[key]
        if !exists || appValue != value {
            return false
        }
    }
    return true
}
