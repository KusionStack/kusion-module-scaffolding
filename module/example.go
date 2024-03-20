package module

import (
	"context"

	"kusionstack.io/kusion-module-framework/pkg/module"
)

type ExampleNetworkModule struct{}

func (o *ExampleNetworkModule) Generate(context context.Context, request *module.GeneratorRequest) (*module.GeneratorResponse, error) {
	return nil, nil
}
