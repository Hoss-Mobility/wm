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

func TestToDbAdmin(t *testing.T) {
	role := "admin"

	webModel, err := ToDb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating dbModel: %v", err)
	}

	// admin is allowed to write to all fields
	assert.Equal(t, secretItem.Name, webModel.Name)
	assert.Equal(t, secretItem.Comment, webModel.Comment)
	assert.Equal(t, secretItem.SecretInfo, webModel.SecretInfo)
	assert.Equal(t, secretItem.TopSecret, webModel.TopSecret)
	assert.Equal(t, secretItem.CanOnlyBeWrittenTo, webModel.CanOnlyBeWrittenTo)
}

func TestToDbDeveloper(t *testing.T) {
	role := "developer"

	webModel, err := ToDb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating dbModel: %v", err)
	}

	// developer is not allowed to write to all fields
	assert.Equal(t, secretItem.Name, webModel.Name)
	assert.Equal(t, secretItem.Comment, webModel.Comment)
	assert.Empty(t, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Equal(t, secretItem.CanOnlyBeWrittenTo, webModel.CanOnlyBeWrittenTo)
}

func TestToDbStaff(t *testing.T) {
	role := "staff"

	webModel, err := ToDb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating dbModel: %v", err)
	}

	// staff is not allowed to write to all fields
	assert.Empty(t, webModel.Name)
	assert.Equal(t, secretItem.Comment, webModel.Comment)
	assert.Empty(t, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Equal(t, secretItem.CanOnlyBeWrittenTo, webModel.CanOnlyBeWrittenTo)
}

func TestToDbUnauthorized(t *testing.T) {
	role := "unauthorized"

	webModel, err := ToDb(secretItem, role)
	if err != nil {
		t.Fatalf("Error creating dbModel: %v", err)
	}

	// unauthorized is not allowed to write to any fields
	assert.Empty(t, webModel.Name)
	assert.Empty(t, webModel.Comment)
	assert.Empty(t, webModel.SecretInfo)
	assert.Empty(t, webModel.TopSecret)
	assert.Empty(t, webModel.CanOnlyBeWrittenTo)
}

func TestApplyUpdateAdmin(t *testing.T) {
	role := "admin"
	newSecretItem := internal.SecretItem{
		Name:               "Updated",
		Comment:            "Updated",
		SecretInfo:         "Updated",
		TopSecret:          "Updated",
		CanOnlyBeWrittenTo: "Updated",
	}

	updated, err := ApplyUpdate(secretItem, newSecretItem, role)
	if err != nil {
		t.Fatalf("Error creating updated: %v", err)
	}

	// admin is allowed to write to all fields
	assert.Equal(t, newSecretItem.Name, updated.Name)
	assert.Equal(t, newSecretItem.Comment, updated.Comment)
	assert.Equal(t, newSecretItem.SecretInfo, updated.SecretInfo)
	assert.Equal(t, newSecretItem.TopSecret, updated.TopSecret)
	assert.Equal(t, newSecretItem.CanOnlyBeWrittenTo, updated.CanOnlyBeWrittenTo)
}

func TestApplyUpdateDeveloper(t *testing.T) {
	role := "developer"
	newSecretItem := internal.SecretItem{
		Name:               "Updated",
		Comment:            "Updated",
		SecretInfo:         "Updated",
		TopSecret:          "Updated",
		CanOnlyBeWrittenTo: "Updated",
	}

	updated, err := ApplyUpdate(secretItem, newSecretItem, role)
	if err != nil {
		t.Fatalf("Error creating updated: %v", err)
	}

	// developer is not allowed to write to all fields
	assert.Equal(t, newSecretItem.Name, updated.Name)
	assert.Equal(t, newSecretItem.Comment, updated.Comment)
	assert.Equal(t, secretItem.SecretInfo, updated.SecretInfo)
	assert.Equal(t, secretItem.TopSecret, updated.TopSecret)
	assert.Equal(t, newSecretItem.CanOnlyBeWrittenTo, updated.CanOnlyBeWrittenTo)
}

func TestApplyUpdateStaff(t *testing.T) {
	role := "staff"
	newSecretItem := internal.SecretItem{
		Name:               "Updated",
		Comment:            "Updated",
		SecretInfo:         "Updated",
		TopSecret:          "Updated",
		CanOnlyBeWrittenTo: "Updated",
	}

	updated, err := ApplyUpdate(secretItem, newSecretItem, role)
	if err != nil {
		t.Fatalf("Error creating updated: %v", err)
	}

	// staff is not allowed to write to all fields
	assert.Equal(t, secretItem.Name, updated.Name)
	assert.Equal(t, newSecretItem.Comment, updated.Comment)
	assert.Equal(t, secretItem.SecretInfo, updated.SecretInfo)
	assert.Equal(t, secretItem.TopSecret, updated.TopSecret)
	assert.Equal(t, newSecretItem.CanOnlyBeWrittenTo, updated.CanOnlyBeWrittenTo)
}

func TestApplyUpdateUnauthorized(t *testing.T) {
	role := "unauthorized"
	newSecretItem := internal.SecretItem{
		Name:               "Updated",
		Comment:            "Updated",
		SecretInfo:         "Updated",
		TopSecret:          "Updated",
		CanOnlyBeWrittenTo: "Updated",
	}

	updated, err := ApplyUpdate(secretItem, newSecretItem, role)
	if err != nil {
		t.Fatalf("Error creating updated: %v", err)
	}

	// staff is not allowed to write to all fields
	assert.Equal(t, secretItem.Name, updated.Name)
	assert.Equal(t, secretItem.Comment, updated.Comment)
	assert.Equal(t, secretItem.SecretInfo, updated.SecretInfo)
	assert.Equal(t, secretItem.TopSecret, updated.TopSecret)
	assert.Equal(t, secretItem.CanOnlyBeWrittenTo, updated.CanOnlyBeWrittenTo)
}
