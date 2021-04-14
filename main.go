package main

import (
	"log"
	"os"

	"github.com/omegion/argocd-actions/pkg/argocd"
)

func main() {
	options := argocd.APIOptions{
		Address: os.Getenv("INPUT_ADDRESS"),
		Token:   os.Getenv("INPUT_TOKEN"),
	}

	api := argocd.NewAPI(options)

	err := api.Sync(os.Getenv("INPUT_APPNAME"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
