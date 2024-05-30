package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
)

func TestKawesomeModule_Generate(t *testing.T) {
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
		{
			name: "Invalid workload type",
			args: args{
				request: &module.GeneratorRequest{
					Project: "kawesome-example",
					Stack:   "dev",
					App:     "kawesome",
					Workload: &v1.Workload{
						Header: v1.Header{
							Type: "Job",
						},
						Job: &v1.Job{},
					},
					DevConfig: v1.Accessory{
						"service": v1.Accessory{
							"port":       80,
							"targetPort": 8080,
							"protocol":   "TCP",
						},
						"randomPassword": v1.Accessory{
							"length": 10,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid service port number",
			args: args{
				request: &module.GeneratorRequest{
					Project: "kawesome-example",
					Stack:   "dev",
					App:     "kawesome",
					Workload: &v1.Workload{
						Header: v1.Header{
							Type: "Service",
						},
						Service: &v1.Service{},
					},
					DevConfig: v1.Accessory{
						"service": v1.Accessory{
							"port": 0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid service targetPort number",
			args: args{
				request: &module.GeneratorRequest{
					Project: "kawesome-example",
					Stack:   "dev",
					App:     "kawesome",
					Workload: &v1.Workload{
						Header: v1.Header{
							Type: "Service",
						},
						Service: &v1.Service{},
					},
					DevConfig: v1.Accessory{
						"service": v1.Accessory{
							"port":       80,
							"targetPort": 0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid service protocol",
			args: args{
				request: &module.GeneratorRequest{
					Project: "kawesome-example",
					Stack:   "dev",
					App:     "kawesome",
					Workload: &v1.Workload{
						Header: v1.Header{
							Type: "Service",
						},
						Service: &v1.Service{},
					},
					DevConfig: v1.Accessory{
						"service": v1.Accessory{
							"port":       80,
							"targetPort": 8080,
							"protocol":   "STCP",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid random password length",
			args: args{
				request: &module.GeneratorRequest{
					Project: "kawesome-example",
					Stack:   "dev",
					App:     "kawesome",
					Workload: &v1.Workload{
						Header: v1.Header{
							Type: "Service",
						},
						Service: &v1.Service{},
					},
					DevConfig: v1.Accessory{
						"service": v1.Accessory{
							"port":       80,
							"targetPort": 8080,
							"protocol":   "TCP",
						},
						"randomPassword": v1.Accessory{
							"length": 0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Valid kawesome module with labels and annotations",
			args: args{
				request: &module.GeneratorRequest{
					Project: "kawesome-example",
					Stack:   "dev",
					App:     "kawesome",
					Workload: &v1.Workload{
						Header: v1.Header{
							Type: "Service",
						},
						Service: &v1.Service{},
					},
					DevConfig: v1.Accessory{
						"service": v1.Accessory{
							"port":       80,
							"targetPort": 8080,
							"protocol":   "TCP",
						},
						"randomPassword": v1.Accessory{
							"length": 10,
						},
					},
					PlatformConfig: v1.GenericConfig{
						"service": v1.GenericConfig{
							"labels": v1.GenericConfig{
								"kusionstack.io/module-name": "kawesome",
							},
							"annotations": v1.GenericConfig{
								"kusionstack.io/module-version": "0.1.0",
							},
						},
					},
				},
			},
			want: &module.GeneratorResponse{
				Resources: []v1.Resource{
					{
						ID:   "v1:Service:kawesome-example:kawesome-example-dev-kawesome",
						Type: "Kubernetes",
						Attributes: map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "Service",
							"metadata": map[string]interface{}{
								"labels": map[string]interface{}{
									"kusionstack.io/module-name": "kawesome",
								},
								"annotations": map[string]interface{}{
									"kusionstack.io/module-version": "0.1.0",
								},
								"name":              "kawesome-example-dev-kawesome",
								"namespace":         "kawesome-example",
								"creationTimestamp": nil,
							},
							"spec": map[string]interface{}{
								"ports": []interface{}{
									map[string]interface{}{
										"name":       "kawesome-example-dev-kawesome-80-tcp",
										"port":       int64(80),
										"protocol":   "TCP",
										"targetPort": int64(8080),
									},
								},
								"selector": map[string]interface{}{
									"app.kubernetes.io/name":    "kawesome",
									"app.kubernetes.io/part-of": "kawesome-example",
								},
								"type": "ClusterIP",
							},
							"status": map[string]interface{}{
								"loadBalancer": map[string]interface{}{},
							},
						},
						DependsOn: nil,
						Extensions: map[string]interface{}{
							"GVK": "/v1, Kind=Service",
						},
					},
					{
						ID:   "hashicorp:random:random_password:kawesome-example-dev-kawesome",
						Type: v1.Terraform,
						Attributes: map[string]interface{}{
							"length":           10,
							"override_special": "_",
							"special":          true,
						},
						Extensions: map[string]interface{}{
							"provider":     "registry.terraform.io/hashicorp/random/3.6.0",
							"providerMeta": map[string]interface{}(nil),
							"resourceType": "random_password",
						},
					},
				},
				Patcher: &v1.Patcher{
					Environments: []corev1.EnvVar{
						{
							Name:  "KUSION_KAWESOME_RANDOM_PASSWORD",
							Value: "$kusion_path.hashicorp:random:random_password:kawesome-example-dev-kawesome.result",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Kawesome{}
			got, err := o.Generate(tt.args.context, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
