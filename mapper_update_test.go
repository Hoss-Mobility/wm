package wm

import (
	"github.com/Hoss-Mobility/wm/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
