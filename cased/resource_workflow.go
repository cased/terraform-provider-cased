package cased

import (
	"context"
	"terraform-provider-cased/workflows"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowCreate,
		ReadContext:   resourceWorkflowRead,
		UpdateContext: resourceWorkflowUpdate,
		DeleteContext: resourceWorkflowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
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
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1000,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"eq", "in", "not", "endsWith", "startsWith"}, false),
						},
					},
				},
			},
			"controls": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"authentication": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"approval": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"self_approval": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"duration": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"responders": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"responder": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Required: true,
															},
															"required": {
																Type:     schema.TypeBool,
																Optional: true,
																Default:  false,
															},
														},
													},
												},
											},
										},
									},
									"sources": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"slack": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"channel": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
												"email": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceWorkflowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*workflows.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	workflow, err := c.CreateWorkflow(buildWorkflow(d, diags))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(workflow.ID)

	return diags
}

func resourceWorkflowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*workflows.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	workflow, err := c.GetWorkflow(d.Id())
	if err != nil {
		if err == workflows.ErrNotFound {
			d.SetId("")
			return diags
		}

		return diag.FromErr(err)
	}

	if err := d.Set("name", workflow.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("conditions", flattenWorkflowConditions(workflow.Conditions)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("controls", flattenWorkflowControls(workflow.Controls)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("updated_at", workflow.UpdatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", workflow.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWorkflowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*workflows.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	_, err := c.UpdateWorkflow(buildWorkflow(d, diags))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkflowRead(ctx, d, m)
}

func resourceWorkflowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*workflows.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	_, err := c.DeleteWorkflow(d.Id())
	if err != nil {
		if err == workflows.ErrNotFound {
			return diags
		}

		return diag.FromErr(err)
	}

	return diags
}

func flattenWorkflowConditions(conditions []workflows.Condition) []interface{} {
	cs := make([]interface{}, len(conditions), len(conditions))

	for i, condition := range conditions {
		c := make(map[string]interface{})

		c["field"] = condition.Field
		c["value"] = condition.Value
		c["operator"] = condition.Operator

		cs[i] = c
	}

	return cs
}

func flattenWorkflowControls(controls workflows.Controls) []interface{} {
	control := map[string]interface{}{}

	if controls.Authentication {
		control["authentication"] = true
	}

	if controls.Reason {
		control["reason"] = true
	}

	if controls.Approval != nil {
		var sources map[string]interface{}
		if controls.Approval.Sources != nil {
			sources = map[string]interface{}{
				"email": controls.Approval.Sources.Email,
			}

			if controls.Approval.Sources.Slack != nil {
				slack := map[string]string{
					"channel": controls.Approval.Sources.Slack.Channel,
				}

				sources["slack"] = []interface{}{slack}
			}
		}

		responders := []interface{}{}
		if controls.Approval.Responders != nil {
			for responder, required := range *controls.Approval.Responders {
				table := map[string]interface{}{
					"name":     responder,
					"required": required == "required",
				}

				responders = append(responders, map[string]interface{}{
					"responder": []interface{}{table},
				})
			}
		}

		control["approval"] = []interface{}{
			map[string]interface{}{
				"count":         controls.Approval.Count,
				"self_approval": controls.Approval.SelfApproval,
				"duration":      controls.Approval.Duration,
				"timeout":       controls.Approval.Timeout,
				"responders":    responders,
				"sources":       []interface{}{sources},
			},
		}
	}

	return []interface{}{control}
}

func buildWorkflow(d *schema.ResourceData, diags diag.Diagnostics) workflows.Workflow {
	conditionsConfig := d.Get("conditions").([]interface{})
	controlsConfig := d.Get("controls").([]interface{})
	name := d.Get("name").(string)
	conditions := []workflows.Condition{}
	controls := workflows.Controls{}

	for _, control := range controlsConfig {
		c := control.(map[string]interface{})

		if val, ok := c["authentication"].(bool); ok {
			controls.Authentication = val
		}

		if val, ok := c["reason"].(bool); ok {
			controls.Reason = val
		}

		if approvals, ok := c["approval"].([]interface{}); ok {
			if controls.Approval == nil {
				controls.Approval = &workflows.ApprovalControl{}
			}

			for _, approval := range approvals {
				a := approval.(map[string]interface{})

				if count, ok := a["reason"].(int); ok {
					controls.Approval.Count = count
				}

				if selfApproval, ok := a["self_approval"].(bool); ok {
					controls.Approval.SelfApproval = selfApproval
				}

				if count, ok := a["count"].(int); ok {
					controls.Approval.Count = count
				}

				if duration, ok := a["duration"].(int); ok {
					controls.Approval.Duration = duration
				}

				if timeout, ok := a["timeout"].(int); ok {
					controls.Approval.Timeout = timeout
				}

				if responders, ok := a["responders"].([]interface{}); ok {
					table := workflows.Responders{}

					for _, responder := range responders {
						r := responder.(map[string]interface{})

						if resp, ok := r["responder"].([]interface{}); ok {
							for _, res := range resp {
								rc := res.(map[string]interface{})

								required := rc["required"].(bool)

								if name, ok := rc["name"].(string); ok {
									if required {
										table[name] = "required"
									} else {
										table[name] = "optional"
									}
								}
							}
						}
					}

					controls.Approval.Responders = &table
				}

				if sources, ok := a["sources"].([]interface{}); ok {
					if controls.Approval.Sources == nil {
						controls.Approval.Sources = &workflows.ApprovalControlSources{}
					}

					for _, source := range sources {
						s := source.(map[string]interface{})

						if email, ok := s["email"].(bool); ok {
							controls.Approval.Sources.Email = email
						}

						if slacks, ok := s["slack"].([]interface{}); ok {
							if controls.Approval.Sources.Slack == nil {
								controls.Approval.Sources.Slack = &workflows.ApprovalControlSourceSlack{}
							}

							for _, slack := range slacks {
								sc := slack.(map[string]interface{})

								if channel, ok := sc["channel"].(string); ok {
									controls.Approval.Sources.Slack.Channel = channel
								}
							}
						}
					}
				}
			}
		}
	}

	for _, condition := range conditionsConfig {
		c := condition.(map[string]interface{})

		conditions = append(conditions, workflows.Condition{
			Field:    c["field"].(string),
			Value:    c["value"].(string),
			Operator: c["operator"].(string),
		})
	}

	return workflows.Workflow{
		ID:         d.Id(),
		Name:       name,
		Conditions: conditions,
		Controls:   controls,
	}
}
