package module

import (
	"context"
	"reflect"
	"testing"

	"kusionstack.io/kusion-module-framework/pkg/module"
)

func TestExampleNetworkModule_Generate(t *testing.T) {
	type args struct {
		context context.Context
		request *module.GeneratorRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *module.GeneratorResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ExampleNetworkModule{}
			got, err := o.Generate(tt.args.context, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
