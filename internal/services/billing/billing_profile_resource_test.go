// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/billingprofile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BillingProfileResource struct{}

func TestAccBillingProfile_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_BILLING_ACCOUNT_NAME") == "" {
		t.Skip("Skipping as `ARM_TEST_BILLING_ACCOUNT_NAME` is not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_billing_profile", "test")
	r := BillingProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBillingProfile_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_BILLING_ACCOUNT_NAME") == "" {
		t.Skip("Skipping as `ARM_TEST_BILLING_ACCOUNT_NAME` is not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_billing_profile", "test")
	r := BillingProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccBillingProfile_complete(t *testing.T) {
	if os.Getenv("ARM_TEST_BILLING_ACCOUNT_NAME") == "" {
		t.Skip("Skipping as `ARM_TEST_BILLING_ACCOUNT_NAME` is not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_billing_profile", "test")
	r := BillingProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBillingProfile_update(t *testing.T) {
	if os.Getenv("ARM_TEST_BILLING_ACCOUNT_NAME") == "" {
		t.Skip("Skipping as `ARM_TEST_BILLING_ACCOUNT_NAME` is not specified")
	}
	data := acceptance.BuildTestData(t, "azurerm_billing_profile", "test2")
	r := BillingProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r BillingProfileResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := billingprofile.ParseBillingProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Billing.BillingProfile.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r BillingProfileResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_billing_profile" "test" {
  name                 = "acctest-bp-%d"
  billing_account_name = "%s"
}
`, data.RandomInteger, os.Getenv("ARM_TEST_BILLING_ACCOUNT_NAME"))
}

func (r BillingProfileResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_billing_profile" "import" {
  name                 = azurerm_billing_profile.test.name
  billing_account_name = azurerm_billing_profile.test.billing_account_name
}
`, r.basic(data))
}

func (r BillingProfileResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_billing_profile" "test" {
  name                 = "acctest-bp-%d"
  billing_account_name = "%s"

  tags = {
    Env = "Test"
  }
}
`, data.RandomInteger, os.Getenv("ARM_TEST_BILLING_ACCOUNT_NAME"))
}

func (r BillingProfileResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_billing_profile" "test" {
  name                 = "acctest-bp-%d"
  billing_account_name = "%s"

  tags = {
    Env = "Test2"
  }
}
`, data.RandomInteger, os.Getenv("ARM_TEST_BILLING_ACCOUNT_NAME"))
}
