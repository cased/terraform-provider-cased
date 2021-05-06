package cased

import (
	"context"
	"terraform-provider-cased/workflows"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CASED_API_URL", "https://api.cased.com"),
			},
			"workflows_api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CASED_WORKFLOWS_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cased_workflow": resourceWorkflow(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiURL := d.Get("api_url").(string)
	workflowsApiKey := d.Get("workflows_api_key").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c, err := workflows.NewClient(apiURL, workflowsApiKey)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
