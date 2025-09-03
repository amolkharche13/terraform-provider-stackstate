package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"

	"github.com/harvester/terraform-provider-harvester/internal/config"
	"github.com/harvester/terraform-provider-harvester/internal/provider/bootstrap"
	"github.com/harvester/terraform-provider-harvester/pkg/constants"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			constants.FieldProviderBootstrap: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "bootstrap harvester server, it will write content to kubeconfig file",
			},
			constants.FieldProviderKubeConfig: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "kubeconfig file path or content of the kubeconfig file as base64 encoded string, users can use the KUBECONFIG environment variable instead.",
			},
			constants.FieldProviderKubeContext: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "name of the kubernetes context to use",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			constants.ResourceTypeDashboard: dashboard.ResourceDashboard(),
		},
		ResourcesMap: map[string]*schema.Resource{
			constants.ResourceTypeDashboard: dashboard.ResourceDashboard(),
		},
		ConfigureContextFunc: providerConfig,
	}
	return p
}

func providerConfig(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	bootstrap := d.Get(constants.FieldProviderBootstrap).(bool)
	kubeConfig := d.Get(constants.FieldProviderKubeConfig).(string)
	kubeContext := d.Get(constants.FieldProviderKubeContext).(string)
	if bootstrap {
		if kubeConfig != "" {
			return nil, diag.Errorf("kubeconfig is not allowed when bootstrap is true")
		}

		if kubeContext != "" {
			return nil, diag.Errorf("kubecontext is not allowed when bootstrap is true")
		}

		return &config.Config{
			Bootstrap: bootstrap,
		}, nil
	}

	kubeConfig, err := homedir.Expand(d.Get(constants.FieldProviderKubeConfig).(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return &config.Config{
		KubeConfig:  kubeConfig,
		KubeContext: kubeContext,
	}, nil
}
