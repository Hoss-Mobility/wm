package wm

import (
	"github.com/Hoss-Mobility/wm/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

var secretsSlice = []internal.SecretItem{secretItem, secretItem, secretItem}

func TestSliceToWebAdmin(t *testing.T) {
	role := "admin"

	webModelSlice, err := SliceToWeb(secretsSlice, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	if webModelSlice == nil {
		t.Fatalf("Error creating webModelSlice")
	}

	deref := *webModelSlice

	// admin is allowed to see all fields
	for i, item := range secretsSlice {
		assert.Equal(t, item.Name, deref[i].Name)
		assert.Equal(t, item.Comment, deref[i].Comment)
		assert.Equal(t, item.SecretInfo, deref[i].SecretInfo)
		assert.Equal(t, item.TopSecret, deref[i].TopSecret)
		assert.Equal(t, item.CanOnlyBeWrittenTo, deref[i].CanOnlyBeWrittenTo)
	}
}

func TestSliceToWebDeveloper(t *testing.T) {
	role := "developer"

	webModelSlice, err := SliceToWeb(secretsSlice, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	if webModelSlice == nil {
		t.Fatalf("Error creating webModelSlice")
	}

	deref := *webModelSlice

	// developer is not allowed to see all fields
	for i, item := range secretsSlice {
		assert.Equal(t, item.Name, deref[i].Name)
		assert.Equal(t, item.Comment, deref[i].Comment)
		assert.Equal(t, item.SecretInfo, deref[i].SecretInfo)
		assert.Empty(t, deref[i].TopSecret)
		assert.Empty(t, deref[i].CanOnlyBeWrittenTo)
	}
}

func TestSliceToWebStaff(t *testing.T) {
	role := "staff"

	webModelSlice, err := SliceToWeb(secretsSlice, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	if webModelSlice == nil {
		t.Fatalf("Error creating webModelSlice")
	}

	deref := *webModelSlice

	// staff is not allowed to see all fields
	for i, item := range secretsSlice {
		assert.Equal(t, item.Name, deref[i].Name)
		assert.Equal(t, item.Comment, deref[i].Comment)
		assert.Empty(t, deref[i].SecretInfo)
		assert.Empty(t, deref[i].TopSecret)
		assert.Empty(t, deref[i].CanOnlyBeWrittenTo)
	}
}

func TestSliceToWebUnauthorized(t *testing.T) {
	role := "unauthorized"

	webModelSlice, err := SliceToWeb(secretsSlice, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	if webModelSlice == nil {
		t.Fatalf("Error creating webModelSlice")
	}

	deref := *webModelSlice

	// unauthorized is not allowed to see any fields
	for i, _ := range secretsSlice {
		assert.Empty(t, deref[i].Name)
		assert.Empty(t, deref[i].Comment)
		assert.Empty(t, deref[i].SecretInfo)
		assert.Empty(t, deref[i].TopSecret)
		assert.Empty(t, deref[i].CanOnlyBeWrittenTo)
	}
}

// To DB
func TestSliceToDbAdmin(t *testing.T) {
	role := "admin"

	dbModelSlice, err := SliceToDb(secretsSlice, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	if dbModelSlice == nil {
		t.Fatalf("Error creating dbModelSlice")
	}

	deref := *dbModelSlice

	// admin is allowed to see all fields
	for i, item := range secretsSlice {
		assert.Equal(t, item.Name, deref[i].Name)
		assert.Equal(t, item.Comment, deref[i].Comment)
		assert.Equal(t, item.SecretInfo, deref[i].SecretInfo)
		assert.Equal(t, item.TopSecret, deref[i].TopSecret)
		assert.Equal(t, item.CanOnlyBeWrittenTo, deref[i].CanOnlyBeWrittenTo)
	}
}

func TestSliceToDbDeveloper(t *testing.T) {
	role := "developer"

	dbModelSlice, err := SliceToDb(secretsSlice, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	if dbModelSlice == nil {
		t.Fatalf("Error creating dbModelSlice")
	}

	deref := *dbModelSlice

	// developer is not allowed to see all fields
	for i, item := range secretsSlice {
		assert.Equal(t, item.Name, deref[i].Name)
		assert.Equal(t, item.Comment, deref[i].Comment)
		assert.Empty(t, deref[i].SecretInfo)
		assert.Empty(t, deref[i].TopSecret)
		assert.Equal(t, item.CanOnlyBeWrittenTo, deref[i].CanOnlyBeWrittenTo)
	}
}

func TestSliceToDbStaff(t *testing.T) {
	role := "staff"

	dbModelSlice, err := SliceToDb(secretsSlice, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	if dbModelSlice == nil {
		t.Fatalf("Error creating dbModelSlice")
	}

	deref := *dbModelSlice

	// staff is not allowed to see all fields
	for i, item := range secretsSlice {
		assert.Empty(t, deref[i].Name)
		assert.Equal(t, item.Comment, deref[i].Comment)
		assert.Empty(t, deref[i].SecretInfo)
		assert.Empty(t, deref[i].TopSecret)
		assert.Equal(t, item.CanOnlyBeWrittenTo, deref[i].CanOnlyBeWrittenTo)
	}
}

func TestSliceToDbUnauthorized(t *testing.T) {
	role := "unauthorized"

	dbModelSlice, err := SliceToDb(secretsSlice, role)
	if err != nil {
		t.Fatalf("Error creating webModel: %v", err)
	}

	if dbModelSlice == nil {
		t.Fatalf("Error creating dbModelSlice")
	}

	deref := *dbModelSlice

	// unauthorized is not allowed to see any fields
	for i, _ := range secretsSlice {
		assert.Empty(t, deref[i].Name)
		assert.Empty(t, deref[i].Comment)
		assert.Empty(t, deref[i].SecretInfo)
		assert.Empty(t, deref[i].TopSecret)
		assert.Empty(t, deref[i].CanOnlyBeWrittenTo)
	}
}
