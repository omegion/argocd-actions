package controller

import (
	"github.com/omegion/argocd-actions/internal/argocd"
)

// Action is action type.
type Action int

// Controller is main struct of Vault.
type Controller struct {
	API argocd.Interface
}

// NewController is a factory to create a Controller.
func NewController(api argocd.Interface) *Controller {
	return &Controller{API: api}
}

// Sync syncs given application.
func (c Controller) Sync(appName string) error {
	return c.API.Sync(appName)
}

// SyncWithLabels syncs apps based on provided labels.
func (c Controller) SyncWithLabels(labels string) error {
	return c.API.SyncWithLabels(labels)  // We are assuming the API interface has a SyncWithLabels method
}