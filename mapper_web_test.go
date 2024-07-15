package wm

import (
	"github.com/Hoss-Mobility/wm/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

var secretItem = internal.SecretItem{
	Name:               "Crab Burger Recipe",
	Comment:            "The recipe of the famous crusty crab burger",
	SecretInfo:         "Bun, Pickle, Patty, Lettuce",
	TopSecret:          "Do not forget the tomato",
	CanOnlyBeWrittenTo: "Hecho en Crustáceo Cascarudo",
}

func TestToWebAdmin(t *testing.T) {
	role := "admin"

	webModel, err := ToWeb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// admin is allowed to see all fields
	assert.Equal(t, secretItem.Name, webModel.Name)
	assert.Equal(t, secretItem.Comment, webModel.Comment)
	assert.Equal(t, secretItem.SecretInfo, webModel.SecretInfo)
	assert.Equal(t, secretItem.TopSecret, webModel.TopSecret)
	assert.Equal(t, secretItem.CanOnlyBeWrittenTo, webModel.CanOnlyBeWrittenTo)
}

func TestToWebDeveloper(t *testing.T) {
	role := "developer"

	webModel, err := ToWeb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// developer is not allowed to see all fields
	assert.Equal(t, secretItem.Name, webModel.Name)
	assert.Equal(t, secretItem.Comment, webModel.Comment)
	assert.Equal(t, secretItem.SecretInfo, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Empty(t, webModel.CanOnlyBeWrittenTo)
}

func TestToWebStaff(t *testing.T) {
	role := "staff"

	webModel, err := ToWeb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// staff is not allowed to see all fields
	assert.Equal(t, secretItem.Name, webModel.Name)
	assert.Equal(t, secretItem.Comment, webModel.Comment)
	assert.Empty(t, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Empty(t, webModel.CanOnlyBeWrittenTo)
}

func TestToWebUnauthorized(t *testing.T) {
	role := "unauthorized"

	webModel, err := ToWeb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// unauthorized is not allowed to see any fields
	assert.Empty(t, webModel.Name)
	assert.Empty(t, webModel.Comment)
	assert.Empty(t, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Empty(t, webModel.CanOnlyBeWrittenTo)
}
