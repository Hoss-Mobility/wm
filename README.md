# `wm` - A role based **W**ebmodel **M**apper
`wm` solves a problem as old as multiuser data access:
**How to define which fields certain users are allowed to access.**

Typically users accessing REST routes are authenticated via some form of identity provider. 
`wm` bridges the gap between the persistence layer, User Auth and exposing models via REST (or other technologies).
With `wm` it is possible to convert your DB model to a web model, **without** manually creating or generating web models.
To define which fields should be visible to what roles, just add struct tags.

<p align="center">
	<img src="https://dibiasi.dev/share/wm.jpg" width=30% height=30% class="center">
</p>


## Installation
This package requires Go 1.12 or newer. [Versioned releases](https://github.com/Hoss-Mobility/wm/releases) are available.

```sh
go get github.com/Hoss-Mobility/wm
```

## Getting started
`wm` uses [struct tags](https://go.dev/ref/spec#Tag) and predefined authorization labels.
Just add the name of a role to your struct and set one of the following tags.
- `r` - Read
- `w` - Write
- `rw` - Read & Write

Multiple roles can be used in one `wm` prefixed struct tag, to define fine granular access control:
```go
type SmallExample struct {
	// Staff is only allowed read the Name field
	// Developers and admins are also allowed to write / change it
	Name          string `wm:"staff:r;developer:rw;admin:rw"`
	// Staff, Developer and Admin are allowed to read and write to field Comment
	Comment       string `wm:"staff:rw;developer:rw;admin:rw"`
}
```

The actual mapping is done via three functions: `ToWeb`, `ToDb` and `ApplyUpdate`.

`ToWeb` creates a web model. A web model most often only contains a subset of 
data in a model (i.e. don't expose secrets).
This is achieved by creating a copy of the source struct and
only setting values the specified user is allowed to read.

```go
source := SmallExample{Name: "Linus the cat", Comment: "A fine boi"}
role := "staff"
webModel, err := wm.ToWeb(source, role)
// webModel.Name and webModel.Comment is populated because staff is allowed to read from the field `Name` and `Comment`
// {
//    "Name": "Linus the cat",
//    "Comment": "A fine boi"
// }
```

`ToDb` creates a new instance of the type of the source struct, but only
sets fields that the user is allowed to write to. 
This method can be useful to create and store a new model in the database.

```go
source := SmallExample{Name: "Linus the cat", Comment: "A heckin' chonker"}
role := "staff"
dbModel, err := wm.ToDb(source, role)
// dbModel.Name is empty because staff is not allowed to write to the field `Name`
// dbModel.Comment is populated because staff is allowed to write to `Comment`
// {
//    "Name": "",
//    "Comment": "A heckin' chonker"
// }
```

`ApplyUpdate` applies changes from a new model to an old model. 
Only sets fields from the new model on the old model, 
if the supplied role is allowed to write to (W or RW).
This method can be useful to update existing models in the database.

```go
old := SmallExample{Name: "Linus the cat", Comment: "A heckin' chonker"}
new := SmallExample{Name: "New name", Comment: "MEGACHONKER"}
role := "staff"
updatedModel, err := wm.ApplyUpdate(old, new, role)
// updatedModel.Name still is "Linus the cat" because staff is not allowed to change the name
// updatedModel.Comment is set to "MEGACHONKER"
// {
//    "Name": "Linus the cat",
//    "Comment": "MEGACHONKER"
// }
```

Please do not rely solely on `wm` to "sanitize" your models before storing it in the database. 
Make sure to check for SQL injections and other malicious techniques.

## Real World Examples
The following pseudo API provides an endpoint to `GET`, `POST` and `PUT` recipes. 
The `Recipe` struct is either served or consumed. `Recipe.SecretIngredients` must
not be exposed to Staff members. Staff members are only allowed to update the Details of 
a recipe, but are not allowed to write / change the name.

```go
type Recipe struct {
    Name               string `json:"name" wm:"staff:r;admin:rw" `
    Details            string `json:"details" wm:"staff:rw;admin:rw"`
    SecretIngredients  string `json:"secret_ingredients,omitempty" wm:"admin:rw"`
}

func GetRecipe(database db.YourDbHandler, manager *scs.SessionManager) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        // user role is set on session via middleware
        userRole := manager.GetString(r.Context(), "USER_ROLE_KEY")
        // get data from db
        dbRecipe, err := database.GetRecipe("crabburger")
        if err != nil {...}
        // convert to web model
        webRecipe, err := wm.ToWeb(dbRecipe, userRole)
        if err != nil {...}
        // render web model
        render.Status(r, http.StatusOK)
        render.JSON(w, r, webRecipe)
    }
}

func PostRecipe(database db.YourDbHandler, manager *scs.SessionManager) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        var webRecipe Recipe
        err := httptools.ParseBodyToStruct(r.Body, &webRecipe)
        if err != nil {...}
        // user role is set on session via middleware
        userRole := manager.GetString(r.Context(), "USER_ROLE_KEY")
        // convert to db model
        dbRecipe, err := wm.ToDb(webRecipe, userRole)
        if err != nil {...}
        // store in db
        dbRecipe, err := database.AddRecipe(dbRecipe)
        if err != nil {...}
        // set status to ok
        render.Status(r, http.StatusOK)
    }
}

func PutRecipe(database db.YourDbHandler, manager *scs.SessionManager) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        var webRecipe Recipe
        err := httptools.ParseBodyToStruct(r.Body, &webRecipe)
        if err != nil {...}
        // get already existing recipe from db
        dbRecipe, err := database.GetRecipeByName(webRecipe.Name)
        if err != nil {...}
        // user role is set on session via middleware
        userRole := manager.GetString(r.Context(), "USER_ROLE_KEY")
        // apply updates from webRecipe to dbRecipe
        updatedRecipe, err := wm.ApplyUpdate(dbRecipe, webRecipe, userRole)
        if err != nil {...}
        // store in db
        dbRecipe, err = database.AddRecipe(updatedRecipe)
        if err != nil {...}
        // set status to ok
        render.Status(r, http.StatusOK)
    }
}
```

## Showcase
The following code snippets showcase the example `main.go` included in the `wm` module.
Specific fields of `SecretItem` are only visible to specified roles.

```go
type SecretItem struct {
	Name               string `wm:"staff:r;developer:rw;admin:rw"`
	Comment            string `wm:"staff:rw;developer:rw;admin:rw"`
	SecretInfo         string `wm:"developer:r;admin:rw"`
	TopSecret          string `wm:"admin:rw"`
	CanOnlyBeWrittenTo string `wm:"staff:w;developer:w;admin:rw"`
}
```

The following snippet highlights what data each role sees:
```go
ToWeb()
---------------
staff sees:
&internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:""}

developer sees:
&internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"", CanOnlyBeWrittenTo:""}

admin sees:
&internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"Do not forget the tomato", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}

unauthorized sees:
&internal.SecretItem{Name:"", Comment:"", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:""}
```

```go
ToDb()
---------------
staff can set:
&internal.SecretItem{Name:"", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}

developer can set:
&internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}

admin can set:
&internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"Do not forget the tomato", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}

unauthorized can set:
&internal.SecretItem{Name:"", Comment:"", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:""}
```

Update all possible fields:
```go
ApplyUpdate()
---------------
staff can set:
&internal.SecretItem{Name:"Crab Burger Recipe", Comment:"Updated", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"Do not forget the tomato", CanOnlyBeWrittenTo:"Updated"}

developer can set:
&internal.SecretItem{Name:"Updated", Comment:"Updated", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"Do not forget the tomato", CanOnlyBeWrittenTo:"Updated"}

admin can set:
&internal.SecretItem{Name:"Updated", Comment:"Updated", SecretInfo:"Updated", TopSecret:"Updated", CanOnlyBeWrittenTo:"Updated"}

unauthorized can set:
&internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"Do not forget the tomato", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}
```

