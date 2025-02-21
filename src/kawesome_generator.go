package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/log"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
)

func main() {
	server.Start(&KawesomeModGenerator{})
}

// KawesomeModGenerator implements the Kusion Module generator interface.
type KawesomeModGenerator struct{}

// KawesomeCfgContext records the configuration context of Kawesome Module, which can be used for
// the unmarshalling of `devConfig` items and `platformConfig` items.
//
// Note: separating the definitions of `KawesomeModGenerator` and `KawesomeCfgContext` can ensure that
// the call of module generator is stateless, avoiding the issues concurrency.
//
// As an example of a Kusion Module, Kawesome consists of two components, one of which
// is a 'Service', which is used to generate a Kubernetes Service resource, and the other is a
// 'RandomePassword', which is used to generate a Terraform random_password resource.
//
// Typically, these two resources are not particularly related, but here they are combined to primarily
// illustrate how to develop a Kusion Module.
type KawesomeCfgContext struct {
	// Service is for service configs of kawesome module.
	Service Service `yaml:"service,omitempty" json:"service,omitempty"`

	// RandomPassword is for random password configs of kawesome module.
	RandomPassword RandomPassword `yaml:"randomPassword,omitempty" json:"randomPassword,omitempty"`
}

type Service struct {
	// Port is the exposed port of the workload.
	Port int `yaml:"port,omitempty" json:"port,omitempty"`

	// TargetPort is the backend container.Container port.
	TargetPort int `yaml:"targetPort,omitempty" json:"targetPort,omitempty"`

	// Protocol is protocol used to expose the port, support ProtocolTCP and ProtocolUDP.
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`

	// Labels are the attached labels of the port, works only when the Public is true.
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`

	// Annotations are the attached annotations of the port, works only when the Public is true.
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

type RandomPassword struct {
	// Length is the specified length of the random password string.
	Length int `yaml:"length,omitempty" json:"length,omitempty"`
}

// Generate implements the generation logic of kawesome module, including a Kubernetes Service and
// a Terraform random_password resource.
func (*KawesomeModGenerator) Generate(ctx context.Context, request *module.GeneratorRequest) (response *module.GeneratorResponse, err error) {
	// Get the module logger with the generator context.
	logger := log.GetModuleLogger(ctx)
	logger.Info("Generating resources...")

	defer func() {
		if r := recover(); r != nil {
			logger.Debug("failed to generate kawesome module: %v", r)
			response = nil
			rawRequest, _ := json.Marshal(request)
			err = fmt.Errorf("panic in kawesome module generator but recovered with error: [%v] and stack %v and request %v",
				r, string(debug.Stack()), string(rawRequest))
		}
	}()

	// Kawesome module does not exist in AppConfiguration configs.
	if request.DevConfig == nil {
		logger.Info("Kawesome module does not exist in AppConfiguration configs")
	}

	// Port should be binded to a Service Workload.
	if request.Workload == nil {
		return nil, errors.New("empty workload")
	} else if workloadType, ok := request.Workload["_type"]; !ok ||
		!strings.Contains(workloadType.(string), ".Service") {
		return nil, errors.New("port should be binded to a service workload")
	}

	// Initiate a new KawesomeCfgContext instance.
	k := &KawesomeCfgContext{}

	// Get the complete kawesome module configs.
	if err := k.CompleteConfig(request.DevConfig, request.PlatformConfig); err != nil {
		logger.Debug("failed to get complete kawesome module configs: %v", err)
		return nil, err
	}

	// Validate the completed kawesome module configs.
	if err := k.ValidateConfig(); err != nil {
		logger.Debug("failed to validate the kawesome module configs: %v", err)
		return nil, err
	}

	var resources []kusionapiv1.Resource
	var patcher *kusionapiv1.Patcher

	// Generate the Kubernetes Service related resource.
	resource, err := k.GenerateServiceResource(request)
	if err != nil {
		return nil, err
	}
	resources = append(resources, *resource)

	// Generate the Terraform random_password related resource and patcher.
	resource, patcher, err = k.GenerateRandomPasswordResource(request)
	if err != nil {
		return nil, err
	}
	resources = append(resources, *resource)

	// Return the Kusion generator response.
	return &module.GeneratorResponse{
		Resources: resources,
		Patcher:   patcher,
	}, nil
}

// CompleteConfig completes the kawesome module configs with both devModuleConfig and platformModuleConfig.
func (k *KawesomeCfgContext) CompleteConfig(devConfig kusionapiv1.Accessory, platformConfig kusionapiv1.GenericConfig) error {
	// Retrieve the config items the developers are concerned about.
	if devConfig != nil {
		devCfgYamlStr, err := yaml.Marshal(devConfig)
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal(devCfgYamlStr, k); err != nil {
			return err
		}
	}

	if k.Service.TargetPort == 0 {
		k.Service.TargetPort = k.Service.Port
	}

	// Retrieve the config items the platform engineers care about.
	if platformConfig != nil {
		platformCfgYamlStr, err := yaml.Marshal(platformConfig)
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal(platformCfgYamlStr, k); err != nil {
			return err
		}
	}

	return nil
}

// ValidateConfig validates the completed kawesome configs are valid or not.
func (k *KawesomeCfgContext) ValidateConfig() error {
	if k.Service.Port < 1 || k.Service.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}

	if k.Service.TargetPort < 1 || k.Service.TargetPort > 65535 {
		return errors.New("targetPort must be between 1 and 65535 if exist")
	}

	if k.Service.Protocol != "TCP" && k.Service.Protocol != "UDP" {
		return errors.New("protocol must be TCP or UDP")
	}

	if k.RandomPassword.Length < 1 {
		return errors.New("random password length must be a positive integer")
	}

	return nil
}

// GenerateServiceResource generates the Kubernetes Service related to the kawesome module service.
//
// Note that we will use the SDK provided by the kusion module framework to wrap the Kubernetes resource
// into Kusion resource.
func (k *KawesomeCfgContext) GenerateServiceResource(request *module.GeneratorRequest) (*kusionapiv1.Resource, error) {
	// Generate the unique application name with project, stack and app name.
	appUniqueName := module.UniqueAppName(request.Project, request.Stack, request.App)
	svcType := v1.ServiceTypeClusterIP

	// Generate the selector for the Service workload with the unique app labels SDK
	// provided by Kusion.
	selector := module.UniqueAppLabels(request.Project, request.App)
	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1.SchemeGroupVersion.String(),
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      appUniqueName,
			Namespace: request.Project,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name: fmt.Sprintf("%s-%d-%s",
						appUniqueName, k.Service.Port, strings.ToLower(k.Service.Protocol)),
					Port:       int32(k.Service.Port),
					TargetPort: intstr.FromInt(k.Service.TargetPort),
					Protocol:   v1.Protocol(k.Service.Protocol),
				},
			},
			Selector: selector,
			Type:     svcType,
		},
	}

	// Add the labels and annotations in kawesome module to the Service.
	if len(svc.Labels) == 0 {
		svc.Labels = make(map[string]string)
	}
	if len(svc.Annotations) == 0 {
		svc.Annotations = make(map[string]string)
	}

	for k, v := range k.Service.Labels {
		svc.Labels[k] = v
	}
	for k, v := range k.Service.Annotations {
		svc.Annotations[k] = v
	}

	// Generate Kusion resource ID and wrap the Kubernetes Service into Kusion resource
	// with the SDK provided by kusion module framework.
	resourceID := module.KubernetesResourceID(svc.TypeMeta, svc.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, svc)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GenerateRandomPasswordResource generates the Terraform random_password related to the kawesome module randomPassword.
//
// Note that we will use the SDK provided by the kusion module framework to wrap the Terraform resource
// into Kusion resource.
func (k *KawesomeCfgContext) GenerateRandomPasswordResource(request *module.GeneratorRequest) (*kusionapiv1.Resource, *kusionapiv1.Patcher, error) {
	// Set the random_password provider config.
	randomPasswordPvdCfg := module.ProviderConfig{
		Source:  "hashicorp/random",
		Version: "3.6.0",
	}

	// Set the random_password attributes.
	attrs := map[string]any{
		"length":           k.RandomPassword.Length,
		"special":          true,
		"override_special": "_",
	}

	// Generate Kusion resource ID & extentions and wrap the Terraform random_password into Kusion resource
	// with the SDK provided by kusion module framework.
	appUniqueName := module.UniqueAppName(request.Project, request.Stack, request.App)
	resourceID, err := module.TerraformResourceID(randomPasswordPvdCfg, "random_password", appUniqueName)
	if err != nil {
		return nil, nil, err
	}

	// Wrap the Terraform resource to Kusion resource in Spec.
	resource, err := module.WrapTFResourceToKusionResource(randomPasswordPvdCfg, "random_password", resourceID, attrs, nil)
	if err != nil {
		return nil, nil, err
	}

	// Inject the random password string into the workload as the environment variables
	// with Kusion Patcher.
	// Note that we use the Kusion path dependency to reference the result of the random_password.
	envVars := []v1.EnvVar{
		{
			Name:  "KUSION_KAWESOME_RANDOM_PASSWORD",
			Value: module.KusionPathDependency(resourceID, "result"),
		},
	}
	patcher := &kusionapiv1.Patcher{
		Environments: envVars,
	}

	return resource, patcher, nil
}
