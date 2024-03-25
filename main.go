package main

import (
	module "kusion-modules/kawesome"

	"kusionstack.io/kusion-module-framework/pkg/server"
)

func main() {
	server.Start(&module.Kawesome{})
}
