package environment

import (
	"context"
	"os"
	"regexp"
	"strings"

	"encoding/base64"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVariables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVariablesRead,
		Schema: map[string]*schema.Schema{
			"sensitive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  ptrBool(false),
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"items": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVariablesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sensitive := d.Get("sensitive").(bool)
	filter := d.Get("filter").(string)
	err := d.Set("items", flattenVariables(os.Environ(), sensitive, filter))
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid)

	return nil
}

func flattenVariables(variables []string, sensitive bool, filter string) map[string]interface{} {
	filtering := len(filter) > 0
	re := regexp.MustCompile(filter)

	out := make(map[string]interface{})
	if variables != nil {
		for _, variable := range variables {
			fields := strings.SplitN(variable, "=", 2)
			name, value := fields[0], fields[1]

			if filtering && !re.MatchString(name) {
				continue
			}
			if sensitive {
				value = base64.StdEncoding.EncodeToString([]byte(value))
			}

			out[name] = value
		}
	}

	return out
}
