package cmd

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cheelim1/argocd-actions/internal/argocd"
	ctrl "github.com/cheelim1/argocd-actions/internal/controller"
)

// Sync syncs given application.
func Sync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync ArgoCD application.",
		RunE: func(cmd *cobra.Command, args []string) error {
			address, _ := cmd.Flags().GetString("address")
			token, _ := cmd.Flags().GetString("token")
			application, _ := cmd.Flags().GetString("application")
			labels, _ := cmd.Flags().GetString("labels")
		
			// Validation logic
			if (application == "" && labels == "") || (application != "" && labels != "") {
				return errors.New("You must specify either 'application' or 'labels', but not both")
			}
		
			api := argocd.NewAPI(&argocd.APIOptions{
				Address: address,
				Token:   token,
			})
		
			controller := ctrl.NewController(api)
		
			if application != "" {
				err := controller.Sync(application)
				if err != nil {
					return err
				}
				log.Infof("Application %s synced", application)
			} else if labels != "" {
				log.Infof("To sync app %s based on labels %s", application, labels)
				matchedApps, err := controller.SyncWithLabels(labels)
				if err != nil {
					return err
				}

				for _, app := range matchedApps {
					log.Infof("Application %s synced using labels", app.Name)
				}
			}
			log.Infof("labels passed in: %s", labels)
			log.Infof("ArgoCD did not trigger sync! %s", application)
			return nil
		},
	}

	cmd.Flags().String("application", "", "ArgoCD application name")
	cmd.Flags().String("labels", "", "Labels to sync the ArgoCD app with")  // Add the new flag for labels
	return cmd
}
