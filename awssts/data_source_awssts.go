package awssts

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Credentials struct {
	AccessKeyId string
	SecretAccessKey string
	SessionToken string
	Expiration time.Time
}

type FederatedUser struct {
	FederatedUserId string
	Arn string
}
type FederationToken struct {
	Credentials Credentials
	FederatedUser FederatedUser
	PackedPolicySize int32
	// ResultMetadata interface{} `json:"-"`
}

func dataSourceAWSSTS() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		// Description: "Sample data source in the Terraform provider awssts.",

		ReadContext: dataSourceAWSSTSRead,

		Schema: map[string]*schema.Schema{
			"federated_user_id": {
				// This description is used by the documentation generator and the language server.
				Description: "The id of a federated user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"federation_token": {
				// This description is used by the documentation generator and the language server.
				Description: "A set of temporary security credentials (consisting of an access key ID, a secret access key, and a security token) for a federated user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAWSSTSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	//TODO: refeactor the aws sts sdk client code to apiClinet
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon sts service client
	client := sts.NewFromConfig((cfg))
	output, err := client.GetFederationToken(context.TODO(), &sts.GetFederationTokenInput{
		Name: aws.String("sam"),
	})

	federated_user_id := strings.Split(*output.FederatedUser.FederatedUserId, ":")[0]
	
	// remove extra ResultMetadata from the SDK json payload
	federationToken := &FederationToken {
		Credentials: Credentials {
			AccessKeyId: *output.Credentials.AccessKeyId,
			SecretAccessKey: *output.Credentials.SecretAccessKey,
			SessionToken: *output.Credentials.SessionToken,
			Expiration: *output.Credentials.Expiration,
		},
		FederatedUser: FederatedUser {
			FederatedUserId: *output.FederatedUser.FederatedUserId,
			Arn: *output.FederatedUser.Arn,
		},
		PackedPolicySize: *output.PackedPolicySize,

	}

	federation_token_bytes, err := json.Marshal(federationToken)
	if err != nil {
		log.Fatal(err)
	}
	federation_token := string(federation_token_bytes)

	d.Set("federated_user_id", federated_user_id)
	d.Set("federation_token", federation_token)
	idFromAPI := "my-id"
	d.SetId(idFromAPI)

	// return diag.Errorf("not implemented")
	return diags
}
