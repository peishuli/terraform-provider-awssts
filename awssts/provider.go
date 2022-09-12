package awssts

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		// ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"aws_federation_token": dataSourceAWSSTS(),
		},
	}
}
