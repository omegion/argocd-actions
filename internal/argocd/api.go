package argocd

import (
	"context"
	"fmt"
	"io"
	"strings"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"google.golang.org/appengine/log"

	argoio "github.com/argoproj/argo-cd/v2/util/io"
)

// Interface is an interface for API.
type Interface interface {
    Sync(appName string) error
    SyncWithLabels(labels string) ([]*v1alpha1.Application, error)
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
func (a API) SyncWithLabels(labels string) ([]*v1alpha1.Application, error) {
    // 1. Fetch all applications
    listResponse, err := a.client.List(context.Background(), &applicationpkg.ApplicationQuery{})
    if err != nil {
        return nil, err
    }

    var syncedApps []*v1alpha1.Application
    var syncErrors []string

    // 2. Iterate through each application, check labels, and sync if it matches
    for _, app := range listResponse.Items {
        ctx := context.Background()
        log.Infof(ctx, "Fetched app: %s with labels: %v", app.Name, app.ObjectMeta.Labels)

        if matchesLabels(&app, labels) {
            err := a.Sync(app.Name)
            if err != nil {
                syncErrors = append(syncErrors, fmt.Sprintf("Error syncing %s: %v", app.Name, err))
                continue
            }
            log.Infof(ctx, "Synced app %s based on labels", app.Name)
            syncedApps = append(syncedApps, &app)
        }
    }

    // Check if no applications were synced based on labels
    if len(syncedApps) == 0 {
        return nil, fmt.Errorf("No applications found with matching labels: %s", labels)
    }

    // Close connection
    defer argoio.Close(a.connection)

    // Return errors if any
    if len(syncErrors) > 0 {
        return syncedApps, fmt.Errorf(strings.Join(syncErrors, "; "))
    }

    return syncedApps, nil
}

func matchesLabels(app *v1alpha1.Application, labelsStr string) bool {
    pairs := strings.Split(labelsStr, ",")
    appLabels := app.ObjectMeta.Labels

    for _, pair := range pairs {
        // Handle negative matches
        if strings.Contains(pair, "!=") {
            keyValue := strings.Split(pair, "!=")
            if len(keyValue) != 2 {
                // Malformed label string
                continue
            }
            key, value := keyValue[0], keyValue[1]
            if appLabels[key] == value {
                return false
            }
        } else if strings.Contains(pair, "=") {
            keyValue := strings.Split(pair, "=")
            if len(keyValue) != 2 {
                // Malformed label string
                continue
            }
            key, value := keyValue[0], keyValue[1]
            if appLabels[key] != value {
                return false
            }
        } else if strings.Contains(pair, "notin") {
            parts := strings.Split(pair, "notin")
            if len(parts) != 2 {
                continue // or handle error
            }
            
            key := strings.TrimSpace(parts[0])
            valueStr := strings.TrimSpace(parts[1])
            
            // Trim brackets and split by comma
            values := strings.Split(strings.Trim(valueStr, "()"), ",")
            
            for _, v := range values {
                if appLabels[key] == strings.TrimSpace(v) {
                    return false
                }
            }
        } else if strings.HasPrefix(pair, "!") {
            key := strings.TrimPrefix(pair, "!")
            if _, exists := appLabels[key]; exists {
                return false
            }
        } else {
            // Existence checks
            if _, exists := appLabels[pair]; !exists {
                return false
            }
        }
    }
    return true
}

