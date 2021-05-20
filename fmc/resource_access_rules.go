package fmc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var access_policies_type string = "AccessRule"

func resourceAccessRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessRulesCreate,
		ReadContext:   resourceAccessRulesRead,
		UpdateContext: resourceAccessRulesUpdate,
		DeleteContext: resourceAccessRulesDelete,
		Schema: map[string]*schema.Schema{
			"acp": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := strings.ToUpper(val.(string))
					allowedValues := []string{"ALLOW", "TRUST", "BLOCK", "MONITOR", "BLOCK_RESET", "BLOCK_INTERACTIVE", "BLOCK_RESET_INTERACTIVE"}
					for _, allowed := range allowedValues {
						if v == allowed {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be in %v, got: %q", key, allowedValues, v))
					return
				},
			},
			"syslog_severity": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := strings.ToUpper(val.(string))
					allowedValues := []string{"ALERT", "CRIT", "DEBUG", "EMERG", "ERR", "INFO", "NOTICE", "WARNING"}
					for _, allowed := range allowedValues {
						if v == allowed {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be in %v, got: %q", key, allowedValues, v))
					return
				},
			},
			"enable_syslog": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"send_events_to_fmc": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"log_files": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"log_begin": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"log_end": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_networks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_networks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_ports": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_ports": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"urls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ips_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"syslog_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comments": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAccessRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics
	var diags diag.Diagnostics

	var sourceZones, destinationZones, sourceNetworks, destinationNetworks, sourcePorts, destinationPorts, urls []AccessRuleSubConfig
	dynamicObjects := []*[]AccessRuleSubConfig{
		&sourceZones, &destinationZones, &sourceNetworks, &destinationNetworks, &sourcePorts, &destinationPorts, &urls,
	}
	for i, objType := range []string{"source_zones", "destination_zones", "source_networks", "destination_networks", "source_ports", "destination_ports", "urls"} {
		if inputEntries, ok := d.GetOk(objType); ok {
			for _, ent := range inputEntries.([]interface{}) {
				*dynamicObjects[i] = append(*dynamicObjects[i], AccessRuleSubConfig{
					ID: ent.(string),
				})
			}
		}
	}
	comments := []string{}
	for _, comment := range d.Get("comments").([]interface{}) {
		comments = append(comments, comment.(string))
	}
	res, err := c.CreateAccessRule(ctx, d.Get("acp").(string), &AccessRule{
		Name:            d.Get("name").(string),
		Type:            access_policies_type,
		Action:          strings.ToUpper(d.Get("action").(string)),
		Syslogseverity:  strings.ToUpper(d.Get("syslog_severity").(string)),
		Enablesyslog:    d.Get("enable_syslog").(bool),
		Enabled:         d.Get("enabled").(bool),
		Sendeventstofmc: d.Get("send_events_to_fmc").(bool),
		Logfiles:        d.Get("log_files").(bool),
		Logbegin:        d.Get("log_begin").(bool),
		Logend:          d.Get("log_end").(bool),
		Sourcezones: AccessRuleSubConfigs{
			Objects: sourceZones,
		},
		Destinationzones: AccessRuleSubConfigs{
			Objects: destinationZones,
		},
		Sourcenetworks: AccessRuleSubConfigs{
			Objects: sourceNetworks,
		},
		Destinationnetworks: AccessRuleSubConfigs{
			Objects: destinationNetworks,
		},
		Sourceports: AccessRuleSubConfigs{
			Objects: sourcePorts,
		},
		Destinationports: AccessRuleSubConfigs{
			Objects: destinationPorts,
		},
		Urls: AccessRuleSubConfigs{
			Objects: urls,
		},
		Ipspolicy: AccessRuleSubConfig{
			ID: d.Get("ips_policy").(string),
		},
		Filepolicy: AccessRuleSubConfig{
			ID: d.Get("file_policy").(string),
		},
		Syslogconfig: AccessRuleSubConfig{
			ID: d.Get("syslog_config").(string),
		},
		Newcomments: comments,
	})
	if err != nil {
		return returnWithDiag(diags, err)
	}
	d.SetId(res.ID)
	return resourceAccessRulesRead(ctx, d, m)
}

func returnWithDiag(diags diag.Diagnostics, err error) diag.Diagnostics {
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Error in access rule",
		Detail:   err.Error(),
	})
	return diags
}

func resourceAccessRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	item, err := c.GetAccessRule(ctx, d.Get("acp").(string), d.Id())
	if err != nil {
		return returnWithDiag(diags, err)
	}
	if err := d.Set("name", item.Name); err != nil {
		return returnWithDiag(diags, err)
	}

	if err := d.Set("type", item.Type); err != nil {
		return returnWithDiag(diags, err)
	}

	if err := d.Set("action", item.Action); err != nil {
		return returnWithDiag(diags, err)
	}
	if err := d.Set("syslog_severity", item.Syslogseverity); err != nil {
		return returnWithDiag(diags, err)
	}
	if err := d.Set("enable_syslog", item.Enablesyslog); err != nil {
		return returnWithDiag(diags, err)
	}
	if err := d.Set("enabled", item.Enabled); err != nil {
		return returnWithDiag(diags, err)
	}
	if err := d.Set("send_events_to_fmc", item.Sendeventstofmc); err != nil {
		return returnWithDiag(diags, err)
	}
	if err := d.Set("log_files", item.Logfiles); err != nil {
		return returnWithDiag(diags, err)
	}
	if err := d.Set("log_begin", item.Logbegin); err != nil {
		return returnWithDiag(diags, err)
	}
	if err := d.Set("log_end", item.Logend); err != nil {
		return returnWithDiag(diags, err)
	}

	dynamicObjects := []*[]AccessRuleResponseObject{
		&item.Sourcezones.Objects,
		&item.Destinationzones.Objects,
		&item.Sourcenetworks.Objects,
		&item.Destinationnetworks.Objects,
		&item.Sourceports.Objects,
		&item.Destinationports.Objects,
		&item.Urls.Objects,
	}

	dynamicObjectNames := []string{"source_zones", "destination_zones", "source_networks", "destination_networks", "source_ports", "destination_ports", "urls"}

	for i, objs := range dynamicObjects {
		response := make([]string, 0)
		for _, obj := range *objs {
			response = append(response, obj.ID)
		}
		if err := d.Set(dynamicObjectNames[i], response); err != nil {
			return returnWithDiag(diags, err)
		}
	}

	return diags
}

func resourceAccessRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//TODO: Need to write update rule
	return diag.Diagnostics{diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Error in access rule",
		Detail:   fmt.Errorf("update not implemented").Error(),
	}}
}

func resourceAccessRulesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := c.DeleteAccessRule(ctx, d.Get("acp").(string), d.Id())
	if err != nil {
		return returnWithDiag(diags, err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
