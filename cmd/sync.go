package cmd

import (
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
			labels, _ := cmd.Flags().GetString("labels")  // Capture the labels from the flag

			api := argocd.NewAPI(&argocd.APIOptions{
				Address: address,
				Token:   token,
			})

			controller := ctrl.NewController(api)

			// Conditionally sync based on labels or application name
			var err error
			if labels != "" {
				err = controller.SyncWithLabels(labels)  // This function needs to be implemented
			} else {
				err = controller.Sync(application)
			}

			if err != nil {
				return err
			}

			log.Infof("Application %s synced", application)

			return nil
		},
	}

	cmd.Flags().String("application", "", "ArgoCD application name")
	cmd.Flags().String("labels", "", "Labels to sync the ArgoCD app with")  // Add the new flag for labels

	if err := cmd.MarkFlagRequired("application"); err != nil {
		log.Fatalf("Lethal damage: %s\n\n", err)
	}

	return cmd
}
