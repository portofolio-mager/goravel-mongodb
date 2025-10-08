package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

var config = `map[string]any{
        "uri":      config.Env("MONGODB_URI", "mongodb://localhost:27017"),
        "database": config.Env("MONGODB_DATABASE", "goravel"),
        "username": config.Env("MONGODB_USERNAME", ""),
        "password": config.Env("MONGODB_PASSWORD", ""),
        "options": map[string]any{
            "max_pool_size": 100,
            "min_pool_size": 5,
        },
        "via": func() (driver.Driver, error) {
            return mongodbfacades.MongoDBDriver("mongodb")
        },
    }`

func main() {
	packages.Setup(os.Args).
		Install(
			modify.GoFile(path.Config("app.go")).
				Find(match.Imports()).Modify(modify.AddImport(packages.GetModulePath())).
				Find(match.Providers()).Modify(modify.Register("&mongodb.ServiceProvider{}", "&database.ServiceProvider{}")),
			modify.GoFile(path.Config("database.go")).
				Find(match.Imports()).Modify(modify.AddImport("github.com/goravel/framework/contracts/database/driver"), modify.AddImport("github.com/portofolio-mager/goravel-mongodb/facades", "mongodbfacades")).
				Find(match.Config("database.connections")).Modify(modify.AddConfig("mongodb", config)),
		).
		Uninstall(
			modify.GoFile(path.Config("app.go")).
				Find(match.Providers()).Modify(modify.Unregister("&mongodb.ServiceProvider{}")).
				Find(match.Imports()).Modify(modify.RemoveImport(packages.GetModulePath())),
			modify.GoFile(path.Config("database.go")).
				Find(match.Config("database.connections")).Modify(modify.RemoveConfig("mongodb")).
				Find(match.Imports()).Modify(modify.RemoveImport("github.com/goravel/framework/contracts/database/driver"), modify.RemoveImport("github.com/portofolio-mager/goravel-mongodb/facades", "mongodbfacades")),
		).
		Execute()
}
