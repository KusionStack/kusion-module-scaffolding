package kusion_module_scaffolding

import (
	"kusionstack.io/kusion-module-framework/pkg/server"

	"kusion-modules/module"
)

func main() {
	server.Start(&module.ExampleNetworkModule{})
}
