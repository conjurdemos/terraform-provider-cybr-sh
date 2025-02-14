package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aharriscybr/terraform-provider-cybr-sh/internal/provider"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAWSAccountResourceSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	// Instantiate the resource.Resource and call its Schema method
	provider.NewAWSAccountResource().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestAccAwsAccountResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccAWSAccountCreateData(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "name", os.Getenv("TF_AWS_NAME")),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "username", os.Getenv("TF_AWS_USERNAME")),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "platform", "AWSAccessKeys"),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "safe", "Testsafe"),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "secret", os.Getenv("TF_AWS_SECRET")),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "sm_manage", "false"),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "sm_manage_reason", "No CPM Associated with Safe."),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "aws_kid", os.Getenv("TF_AWS_KEY_ID")),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "aws_account_id", os.Getenv("TF_AWS_ACCOUNT_ID")),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "aws_alias", os.Getenv("TF_AWS_ALIAS")),
					resource.TestCheckResourceAttr("cybr-sh_aws_account.test", "aws_account_region", os.Getenv("TF_AWS_ACCOUNT_REGION")),
					resource.TestCheckResourceAttrSet("cybr-sh_aws_account.test", "id"),
					resource.TestCheckResourceAttrSet("cybr-sh_aws_account.test", "last_updated"),
				),
			},
			{
				Config: providerConfig + `
				  removed {
					from = cybr-sh_aws_account.test
					lifecycle {
						destroy = false
					}
				
				}`,
			},
		},
	})
}

func testAccAWSAccountCreateData() string {
	return fmt.Sprintf(`
	resource "cybr-sh_aws_account" "test" {
		name               = %[1]q
		username           = %[2]q
		platform           = "AWSAccessKeys"
		safe               = "Testsafe"
		secret             = %[3]q
		sm_manage          = false
		sm_manage_reason   = "No CPM Associated with Safe."
		aws_kid            = %[4]q
		aws_account_id     = %[5]q
		aws_alias          = %[6]q
		aws_account_region = %[7]q
}
	`, os.Getenv("TF_AWS_NAME"), os.Getenv("TF_AWS_USERNAME"), os.Getenv("TF_AWS_SECRET"),
		os.Getenv("TF_AWS_KEY_ID"), os.Getenv("TF_AWS_ACCOUNT_ID"), os.Getenv("TF_AWS_ALIAS"), os.Getenv("TF_AWS_ACCOUNT_REGION"))
}
