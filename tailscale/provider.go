// Package tailscale describes the resources and data sources provided by the terraform provider. Each resource
// or data source is described within its own file.
package tailscale

import (
	"context"
	"fmt"

	"github.com/davidsbond/terraform-provider-tailscale/internal/tailscale"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderOption func(p *schema.Provider)

// Provider returns the *schema.Provider instance that implements the terraform provider.
func Provider(options ...ProviderOption) *schema.Provider {
	provider := &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("TAILSCALE_API_KEY", nil),
				Required:    true,
			},
			"tailnet": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("TAILSCALE_TAILNET", nil),
				Required:    true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"tailscale_acl":                  resourceACL(),
			"tailscale_dns_nameservers":      resourceDNSNameservers(),
			"tailscale_dns_preferences":      resourceDNSPreferences(),
			"tailscale_dns_search_paths":     resourceDNSSearchPaths(),
			"tailscale_device_subnet_routes": resourceDeviceSubnetRoutes(),
			"tailscale_device_authorization": resourceDeviceAuthorization(),
			"tailscale_tailnet_key":          resourceTailnetKey(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"tailscale_device":  dataSourceDevice(),
			"tailscale_devices": dataSourceDevices(),
		},
	}

	for _, option := range options {
		option(provider)
	}

	return provider
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	tailnet := d.Get("tailnet").(string)

	client, err := tailscale.NewClient(apiKey, tailnet)
	if err != nil {
		return nil, diagnosticsError(err, "failed to initialise client")
	}

	return client, nil
}

func diagnosticsError(err error, message string, args ...interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		{
			Severity: diag.Error,
			Summary:  fmt.Sprintf(message, args...),
			Detail:   err.Error(),
		},
	}
}

func diagnosticsErrorWithPath(err error, message string, path cty.Path, args ...interface{}) diag.Diagnostics {
	d := diagnosticsError(err, message, args...)
	d[0].AttributePath = path
	return d
}

func createUUID() string {
	val, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return val
}
