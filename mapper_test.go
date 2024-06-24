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
	assert.Equal(t, webModel.Name, secretItem.Name)
	assert.Equal(t, webModel.Comment, secretItem.Comment)
	assert.Equal(t, webModel.SecretInfo, secretItem.SecretInfo)
	assert.Equal(t, webModel.TopSecret, secretItem.TopSecret)
	assert.Equal(t, webModel.CanOnlyBeWrittenTo, secretItem.CanOnlyBeWrittenTo)
}

func TestToWebDeveloper(t *testing.T) {
	role := "developer"

	webModel, err := ToWeb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// developer is not allowed to see all fields
	assert.Equal(t, webModel.Name, secretItem.Name)
	assert.Equal(t, webModel.Comment, secretItem.Comment)
	assert.Equal(t, webModel.SecretInfo, secretItem.SecretInfo)
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
	assert.Equal(t, webModel.Name, secretItem.Name)
	assert.Equal(t, webModel.Comment, secretItem.Comment)
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

func TestToDbAdmin(t *testing.T) {
	role := "admin"

	webModel, err := ToDb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// admin is allowed to write to all fields
	assert.Equal(t, webModel.Name, secretItem.Name)
	assert.Equal(t, webModel.Comment, secretItem.Comment)
	assert.Equal(t, webModel.SecretInfo, secretItem.SecretInfo)
	assert.Equal(t, webModel.TopSecret, secretItem.TopSecret)
	assert.Equal(t, webModel.CanOnlyBeWrittenTo, secretItem.CanOnlyBeWrittenTo)
}

func TestToDbDeveloper(t *testing.T) {
	role := "developer"

	webModel, err := ToDb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// developer is not allowed to write to all fields
	assert.Equal(t, webModel.Name, secretItem.Name)
	assert.Equal(t, webModel.Comment, secretItem.Comment)
	assert.Empty(t, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Equal(t, webModel.CanOnlyBeWrittenTo, secretItem.CanOnlyBeWrittenTo)
}

func TestToDbStaff(t *testing.T) {
	role := "staff"

	webModel, err := ToDb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// staff is not allowed to write to all fields
	assert.Empty(t, webModel.Name)
	assert.Equal(t, webModel.Comment, secretItem.Comment)
	assert.Empty(t, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Equal(t, webModel.CanOnlyBeWrittenTo, secretItem.CanOnlyBeWrittenTo)
}

func TestToDbUnauthorized(t *testing.T) {
	role := "unauthorized"

	webModel, err := ToDb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	// unauthorized is not allowed to write to any fields
	assert.Empty(t, webModel.Name)
	assert.Empty(t, webModel.Comment)
	assert.Empty(t, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Empty(t, webModel.CanOnlyBeWrittenTo)
}
