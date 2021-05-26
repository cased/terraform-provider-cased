package cased

import (
	"context"
	"time"

	"github.com/cased/cased-go"
	"github.com/cased/cased-go/webhooks/endpoint"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWebhooksEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebhooksEndpointCreate,
		ReadContext:   resourceWebhooksEndpointRead,
		UpdateContext: resourceWebhooksEndpointUpdate,
		DeleteContext: resourceWebhooksEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"event.created", "workflow.result.created", "workflow.result.updated"}, false),
				},
				Optional: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceWebhooksEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	e, err := endpoint.New(buildWebhooksEndpointParams(d, diags))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(e.ID)

	return diags
}

func resourceWebhooksEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	e, err := endpoint.Get(d.Id())
	if err != nil {
		if casedErr, ok := err.(*cased.Error); ok {
			if casedErr.Code == cased.ErrorCodeNotFound {
				d.SetId("")
				return diags
			}
		}

		return diag.FromErr(err)
	}

	if err := d.Set("url", e.URL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("secret", e.Secret); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_types", e.EventTypes); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("updated_at", e.UpdatedAt.Format(time.RFC3339Nano)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", e.CreatedAt.Format(time.RFC3339Nano)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWebhooksEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	_, err := endpoint.Update(d.Id(), buildWebhooksEndpointParams(d, diags))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWebhooksEndpointRead(ctx, d, m)
}

func resourceWebhooksEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	_, err := endpoint.Delete(d.Id())
	if err != nil {
		if casedErr, ok := err.(*cased.Error); ok {
			if casedErr.Code == cased.ErrorCodeNotFound {
				return diags
			}
		}

		return diag.FromErr(err)
	}

	return diags
}

func buildWebhooksEndpointParams(d *schema.ResourceData, diags diag.Diagnostics) *cased.WebhooksEndpointParams {
	params := &cased.WebhooksEndpointParams{}
	if ok := d.HasChange("url"); ok {
		params.URL = cased.String(d.Get("url").(string))
	}

	if ok := d.HasChange("event_types"); ok {
		eventTypesSet := d.Get("event_types").(*schema.Set)
		eventTypes := []string{}

		for _, v := range eventTypesSet.List() {
			eventTypes = append(eventTypes, v.(string))
		}

		params.EventTypes = &eventTypes
	}

	return params
}
