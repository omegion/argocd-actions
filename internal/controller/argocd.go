package controller

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/cheelim1/argocd-actions/internal/argocd"
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
func (c Controller) SyncWithLabels(labels string) ([]*v1alpha1.Application, error) {
    return c.API.SyncWithLabels(labels)
}