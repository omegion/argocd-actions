package argocd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	log "github.com/sirupsen/logrus"

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

func (a API) Refresh(appName string) error {
    refreshType := "normal" // or "hard" based on your preference

    // Get the current application state with refresh
    getRequest := applicationpkg.ApplicationQuery{
        Name:    &appName,
        Refresh: &refreshType,
    }
    _, err := a.client.Get(context.Background(), &getRequest)
    return err
}

// Sync syncs given application.
func (a API) Sync(appName string) error {
    maxRetries := 5

    for i := 0; i < maxRetries; i++ {
        // Refresh the application to detect latest changes
        err := a.Refresh(appName)
        if err != nil {
            return err
        }

        // Sync the application
        request := applicationpkg.ApplicationSyncRequest{
            Name:  &appName,
            Prune: true,
        }

        _, err = a.client.Sync(context.Background(), &request)
        if err != nil {
            return err
        }

        // Check if there's any diff
        app, err := a.client.Get(context.Background(), &applicationpkg.ApplicationQuery{Name: &appName})
        if err != nil {
            return err
        }

        // If the application is synced (no diff), break out of the loop
        if app.Status.Sync.Status == v1alpha1.SyncStatusCodeSynced {
            break
        }

        // If not, wait for a short duration before retrying
        time.Sleep(10 * time.Second)
    }

    return nil
}

// SyncWithLabels syncs applications based on provided labels.
func (a API) SyncWithLabels(labels string) ([]*v1alpha1.Application, error) {
    // 1. Fetch all applications
    listResponse, err := a.client.List(context.Background(), &applicationpkg.ApplicationQuery{})
    if err != nil {
        argoio.Close(a.connection)  // Close the connection here if there's an error
        return nil, err
    }

    var syncedApps []*v1alpha1.Application
    var syncErrors []string

    // 2. Iterate through each application, check labels, and sync if it matches
    for _, app := range listResponse.Items {
        log.Infof("Fetched app: %s with labels: %v", app.Name, app.ObjectMeta.Labels)

        if matchesLabels(&app, labels) {
            // Refresh and retry sync logic
            maxRetries := 5
            for i := 0; i < maxRetries; i++ {
                err := a.Refresh(app.Name)
                if err != nil {
                    syncErrors = append(syncErrors, fmt.Sprintf("Error refreshing %s: %v", app.Name, err))
                    break
                }

                // Check for differences (you might need to implement this)
                hasDifferences, err := a.HasDifferences(app.Name)
                if err != nil {
                    syncErrors = append(syncErrors, fmt.Sprintf("Error checking differences for %s: %v", app.Name, err))
                    break
                }

                if !hasDifferences {
                    break
                }

                err = a.Sync(app.Name)
                if err != nil {
                    syncErrors = append(syncErrors, fmt.Sprintf("Error syncing %s: %v", app.Name, err))
                    continue
                }
                log.Infof("Synced app %s based on labels", app.Name)
            }
            syncedApps = append(syncedApps, &app)
        }
    }

    // Close the gRPC connection after all sync operations are complete
    defer argoio.Close(a.connection)

    // Check if no applications were synced based on labels
    if len(syncedApps) == 0 {
        return nil, fmt.Errorf("No applications found with matching labels: %s", labels)
    }

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

// HasDifferences checks if the given application has differences between the desired and live state.
func (a API) HasDifferences(appName string) (bool, error) {
    // Get the application details
    appResponse, err := a.client.Get(context.Background(), &applicationpkg.ApplicationQuery{
        Name: &appName,
    })
    if err != nil {
        return false, fmt.Errorf("Error fetching application %s: %v", appName, err)
    }

    // Check the application's sync status
    if appResponse.Status.Sync.Status == v1alpha1.SyncStatusCodeOutOfSync {
        return true, nil
    }

    return false, nil
}
