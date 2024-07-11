package main

import (
	"fmt"
	"github.com/Hoss-Mobility/wm"
	"github.com/Hoss-Mobility/wm/internal"
)

func main() {
	secretItem := internal.SecretItem{
		Name:               "Crab Burger Recipe",
		Comment:            "The recipe of the famous crusty crab burger",
		SecretInfo:         "Bun, Pickle, Patty, Lettuce",
		TopSecret:          "Do not forget the tomato",
		CanOnlyBeWrittenTo: "Hecho en Crustáceo Cascarudo",
	}

	roles := []string{"staff", "developer", "admin", "unauthorized"}

	fmt.Printf("ToWeb()\n---------------\n")
	for _, role := range roles {
		webModel, err := wm.ToWeb(secretItem, role)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s sees: \n%#v\n\n", role, webModel)
	}

	fmt.Printf("\nToDb()\n---------------\n")
	for _, role := range roles {
		dbModel, err := wm.ToDb(secretItem, role)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s can set: \n%#v\n\n", role, dbModel)
	}

	update := internal.SecretItem{
		Name:               "Updated",
		Comment:            "Updated",
		SecretInfo:         "Updated",
		TopSecret:          "Updated",
		CanOnlyBeWrittenTo: "Updated",
	}

	fmt.Printf("\nApplyUpdate()\n---------------\n")
	for _, role := range roles {
		updatedModel, err := wm.ApplyUpdate(secretItem, update, role)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s can set: \n%#v\n\n", role, updatedModel)
	}
}
