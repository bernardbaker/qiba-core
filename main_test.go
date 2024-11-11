package main

import (
	"testing"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/infrastructure"
)

func TestServerComponents(t *testing.T) {
	// Test repository initialization
	repo := infrastructure.NewInMemoryGameRepository()
	if repo == nil {
		t.Error("Failed to create repository")
	}

	// Test encrypter initialization
	encrypter := infrastructure.NewEncrypter([]byte("testkey"))
	if encrypter == nil {
		t.Error("Failed to create encrypter")
	}

	// Test service initialization
	service := app.NewGameService(repo, encrypter)
	if service == nil {
		t.Error("Failed to create game service")
	}
}
