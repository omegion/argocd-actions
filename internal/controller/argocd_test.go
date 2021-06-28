package controller

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/omegion/argocd-actions/internal/argocd/mocks"
)

func TestSign(t *testing.T) {
	ctrl := gomock.NewController(t)
	api := mocks.NewMockInterface(ctrl)

	expectedAppName := "testApp"

	api.EXPECT().Sync(expectedAppName).Return(nil)

	controller := NewController(api)
	err := controller.Sync(expectedAppName)

	assert.NoError(t, err)
}
