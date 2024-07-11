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

// ToDb converts webModel to a dbModel.
// Only sets fields on dbModel the supplied role is allowed to write to (W or RW).
func ToDb[T any](webModel T, role string) (dbModel *T, err error) {
	return doMapping(webModel, nil, role, toDBActions)
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
