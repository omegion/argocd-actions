package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// Root is root of cmd.
func Root() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "argocd-actions",
		Short:        "ArgoCD Actions.",
		Long:         "Operate your ArgoCD applications from GitHub.",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().String("address", "", "ArgoCD address")

	if err := cmd.MarkPersistentFlagRequired("address"); err != nil {
		log.Fatalf("Lethal damage: %s\n\n", err)
	}

	cmd.PersistentFlags().String("token", "", "ArgoCD token")

	if err := cmd.MarkPersistentFlagRequired("token"); err != nil {
		log.Fatalf("Lethal damage: %s\n\n", err)
	}

	return cmd
}
