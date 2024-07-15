package wm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
