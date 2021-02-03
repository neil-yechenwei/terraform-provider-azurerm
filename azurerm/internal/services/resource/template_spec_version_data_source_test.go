package resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type TemplateSpecVersionDataSource struct {
}

func TestAccTemplateSpecVersionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_template_spec_version", "test")
	r := TemplateSpecVersionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("artifact").Exists(),
			),
		},
	})
}

func (TemplateSpecVersionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_template_spec" "test" {
  name                = azurerm_template_spec_version.test.name
  resource_group_name = azurerm_resource_group.test.name
  template_spec_name  = azurerm_template_spec.test.name
}
`, TemplateSpecVersionResource{}.complete(data))
}
