# `wm` - A role based **W**ebmodel **M**apper
`wm` solves a problem as old as multiuser data access:
**How to define which fields certain users are allowed to access.**

Typically users accessing REST routes are authenticated via some form of identity provider. 
`wm` bridges the gap between the persistence layer, User Auth and exposing models via REST (or other technologies).
With `wm` it is possible to convert your DB model to a web model, **without** manually creating or generating web models.
To define which fields should be visible to what roles, just add struct tags and call `ToWeb()`.

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
}
```

The actual mapping is done via two functions: `ToWeb` and `ToDb`.

`ToWeb` creates a web model. A web model most often only contains subset of 
data in a model (i.e. don't expose secrets).
This is achieved by creating a copy of the source struct and
only setting values the specified user is allowed to read.

```go
source := SmallExample{Name: "Linus the cat"}
role := "staff"
webModel, err := wm.ToWeb(source, role)
// dbModel.Name is populated because staff is allowed to read from the field `Name`
```

`ToDb` creates a copy of the source struct to be used in your datalayer, but only
sets fields that the user is allowed to write to.

```go
source := SmallExample{Name: "Linus the cat"}
role := "staff"
dbModel, err := wm.ToDb(source, role)
// dbModel.Name is empty because staff is not allowed to write to the field `Name`
```

Please do not rely solely on `wm` to "sanitize" your models before storing it in the database. 
Make sure to check for SQL injections and other malicious techniques.

## Real World Examples
The following pseudo API provides an endpoint to `GET` and `POST` recipes. 
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
internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:""}

developer sees: 
internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"", CanOnlyBeWrittenTo:""}

admin sees: 
internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"Do not forget the tomato", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}

unauthorized sees: 
internal.SecretItem{Name:"", Comment:"", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:""}
```

```go
ToDb()
---------------
staff can set: 
internal.SecretItem{Name:"", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}

developer can set: 
internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}

admin can set: 
internal.SecretItem{Name:"Crab Burger Recipe", Comment:"The recipe of the famous crusty crab burger", SecretInfo:"Bun, Pickle, Patty, Lettuce", TopSecret:"Do not forget the tomato", CanOnlyBeWrittenTo:"Hecho en Crustáceo Cascarudo"}

unauthorized can set: 
internal.SecretItem{Name:"", Comment:"", SecretInfo:"", TopSecret:"", CanOnlyBeWrittenTo:""}
```

