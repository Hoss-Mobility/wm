package wm

import (
	"errors"
	"reflect"
	"slices"
	"strings"
)

var (
	FieldNotAddressableErr = errors.New("cannot dynamically set data on field")
	ImproperFormatErr      = errors.New("improper WM struct field format, example: `staff:r;developer:rw`")
)

const (
	READ       string = "r"
	WRITE             = "w"
	READ_WRITE        = "rw"
)

var (
	toWebActions = []string{READ, READ_WRITE}
	toDBActions  = []string{WRITE, READ_WRITE}
)

// ToWeb converts dbModel to a webModel.
// Only sets fields on webModel the supplied role is allowed to read from (R or RW).
func ToWeb[T any](dbModel T, role string) (webModel *T, err error) {
	return doMapping(dbModel, nil, role, toWebActions)
}

// SliceToWeb converts a slice of dbModels to a slice of webModels.
// Only sets fields on webModels the supplied role is allowed to read from (R or RW).
func SliceToWeb[T any](dbSlice []T, role string) (*[]T, error) {
	webSlice := make([]T, len(dbSlice))
	for i, item := range dbSlice {
		m, err := doMapping(item, nil, role, toWebActions)
		if err != nil {
			return nil, err
		}
		webSlice[i] = *m
	}
	return &webSlice, nil
}

// ToDb converts webModel to a dbModel.
// Only sets fields on dbModel the supplied role is allowed to write to (W or RW).
func ToDb[T any](webModel T, role string) (dbModel *T, err error) {
	return doMapping(webModel, nil, role, toDBActions)
}

// SliceToDb converts a slice of webModels to a slice of dbModels.
// Only sets fields on dbModels the supplied role is allowed to write to (W or RW).
func SliceToDb[T any](webSlice []T, role string) (*[]T, error) {
	dbSlice := make([]T, len(webSlice))
	for i, item := range webSlice {
		m, err := doMapping(item, nil, role, toDBActions)
		if err != nil {
			return nil, err
		}
		dbSlice[i] = *m
	}
	return &dbSlice, nil
}

// ApplyUpdate applies changes from newModel to oldModel.
// Only sets fields from newModel on oldModel,
// if the supplied role is allowed to write to (W or RW).
func ApplyUpdate[T any](oldModel T, newModel T, role string) (diffModel *T, err error) {
	return doMapping(newModel, &oldModel, role, toDBActions)
}

// doMapping maps sets the fields of sourceModel on targetModel
// if the role is allowed to do so, according to allowedActions.
// If targetModel is nil, a new model of type T is created.
func doMapping[T any](sourceModel T, targetModel *T, role string, allowedActions []string) (*T, error) {
	if targetModel == nil {
		targetModel = new(T)
	}

	// sourceSchema needed to get struct tags
	sourceSchema := reflect.TypeOf(sourceModel)
	// sourceValues is needed to get values from sourceModel
	sourceValues := reflect.ValueOf(&sourceModel).Elem()
	// targetValues is needed to access values of targetModel
	targetValues := reflect.ValueOf(targetModel).Elem()

	for i := 0; i < sourceSchema.NumField(); i++ {
		field := sourceSchema.Field(i)
		aclRaw := field.Tag.Get("wm")
		if aclRaw == "" {
			// no wm fields have been found
			continue
		}
		acls := strings.Split(aclRaw, ";")

		// run through all permissions
		for _, acl := range acls {
			// permMapping == role:permission
			permMapping := strings.Split(acl, ":")
			if len(permMapping) != 2 {
				return nil, ImproperFormatErr
			}
			// check if allowed to set field
			webModelField := targetValues.FieldByName(field.Name)
			if !webModelField.IsValid() {
				return nil, FieldNotAddressableErr
			}
			if permMapping[0] == role {
				if slices.Contains(allowedActions, permMapping[1]) {
					// only set value on targetModel if role is allowed to perform action
					sourceVal := sourceValues.FieldByName(field.Name)
					webModelField.Set(sourceVal)
				}
			}
		}
	}
	return targetModel, nil
}
