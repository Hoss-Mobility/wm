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

func ToWeb[T any](dbModel T, role string) (webModel T, err error) {
	return doMapping(dbModel, role, toWebActions)
}

func ToDb[T any](webModel T, role string) (dbModel T, err error) {
	return doMapping(webModel, role, toDBActions)
}

func doMapping[T any](sourceModel T, role string, allowedActions []string) (targetModel T, err error) {
	// sourceSchema needed to get struct tags
	sourceSchema := reflect.TypeOf(sourceModel)
	// sourceValues is needed to get values from sourceModel
	sourceValues := reflect.ValueOf(&sourceModel).Elem()
	// targetValues is needed to access values of targetModel
	targetValues := reflect.ValueOf(&targetModel).Elem()

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
				return targetModel, ImproperFormatErr
			}
			// check if allowed to set field
			webModelField := targetValues.FieldByName(field.Name)
			if !webModelField.IsValid() {
				return targetModel, FieldNotAddressableErr
			}
			if permMapping[0] == role {
				if slices.Contains(allowedActions, permMapping[1]) {
					// only set value on targetModel if role is allowed to perform action
					test := sourceValues.FieldByName(field.Name)
					webModelField.Set(test)
				}
			}
		}
	}
	return targetModel, nil
}
