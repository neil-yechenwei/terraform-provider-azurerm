// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/billingprofile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BillingProfileModel struct {
	Name                     string             `tfschema:"name"`
	BillingAccountName       string             `tfschema:"billing_account_name"`
	BillTo                   []BillTo           `tfschema:"bill_to"`
	DisplayName              string             `tfschema:"display_name"`
	EnabledAzurePlans        []EnabledAzurePlan `tfschema:"enabled_azure_plan"`
	InvoiceEmailOptInEnabled bool               `tfschema:"invoice_email_opt_in_enabled"`
	PoNumber                 string             `tfschema:"po_number"`
	Tags                     map[string]string  `tfschema:"tags"`
}

type BillTo struct {
	AddressLine1 string `tfschema:"address_line_1"`
	AddressLine2 string `tfschema:"address_line_2"`
	AddressLine3 string `tfschema:"address_line_3"`
	City         string `tfschema:"city"`
	CompanyName  string `tfschema:"company_name"`
	Country      string `tfschema:"country"`
	District     string `tfschema:"district"`
	Email        string `tfschema:"email"`
	FirstName    string `tfschema:"first_name"`
	LastName     string `tfschema:"last_name"`
	MiddleName   string `tfschema:"middle_name"`
	PhoneNumber  string `tfschema:"phone_number"`
	PostalCode   string `tfschema:"postal_code"`
	Region       string `tfschema:"region"`
}

type EnabledAzurePlan struct {
	SkuDescription string `tfschema:"sku_description"`
	SkuId          string `tfschema:"sku_id"`
}

var (
	_ sdk.Resource           = BillingProfileResource{}
	_ sdk.ResourceWithUpdate = BillingProfileResource{}
)

type BillingProfileResource struct{}

func (r BillingProfileResource) ModelObject() interface{} {
	return &BillingProfileModel{}
}

func (r BillingProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return billingprofile.ValidateBillingProfileID
}

func (r BillingProfileResource) ResourceType() string {
	return "azurerm_billing_profile"
}

func (r BillingProfileResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"billing_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"bill_to": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_line_1": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"country": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_line_2": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_line_3": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"city": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"company_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"district": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"email": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"first_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"last_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"middle_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"phone_number": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"postal_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"region": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled_azure_plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sku_description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sku_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"invoice_email_opt_in_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"po_number": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r BillingProfileResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BillingProfileResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.BillingProfile

			var model BillingProfileModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := billingprofile.NewBillingProfileID(model.BillingAccountName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := billingprofile.BillingProfile{
				Properties: &billingprofile.BillingProfileProperties{
					BillTo:            expandBillTo(model.BillTo),
					EnabledAzurePlans: expandEnabledAzurePlans(model.EnabledAzurePlans),
					InvoiceEmailOptIn: pointer.To(model.InvoiceEmailOptInEnabled),
				},
				Tags: pointer.To(model.Tags),
			}

			if v := model.DisplayName; v != "" {
				parameters.Properties.DisplayName = pointer.To(v)
			}

			if v := model.PoNumber; v != "" {
				parameters.Properties.PoNumber = pointer.To(v)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BillingProfileResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.BillingProfile

			id, err := billingprofile.ParseBillingProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BillingProfileModel{}
			if model := resp.Model; model != nil {
				state.Name = id.BillingProfileName
				state.BillingAccountName = id.BillingAccountName
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.BillTo = flattenBillTo(props.BillTo)
					state.DisplayName = pointer.From(props.DisplayName)
					state.EnabledAzurePlans = flattenEnabledAzurePlans(props.EnabledAzurePlans)
					state.InvoiceEmailOptInEnabled = pointer.From(props.InvoiceEmailOptIn)
					state.PoNumber = pointer.From(props.PoNumber)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r BillingProfileResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.BillingProfile

			id, err := billingprofile.ParseBillingProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BillingProfileModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("bill_to") {
				parameters.Properties.BillTo = expandBillTo(model.BillTo)
			}

			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if metadata.ResourceData.HasChange("enabled_azure_plan") {
				parameters.Properties.EnabledAzurePlans = expandEnabledAzurePlans(model.EnabledAzurePlans)
			}

			if metadata.ResourceData.HasChange("invoice_email_opt_in_enabled") {
				parameters.Properties.InvoiceEmailOptIn = pointer.To(model.InvoiceEmailOptInEnabled)
			}

			if metadata.ResourceData.HasChange("po_number") {
				parameters.Properties.PoNumber = pointer.To(model.PoNumber)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r BillingProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.BillingProfile

			id, err := billingprofile.ParseBillingProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandEnabledAzurePlans(input []EnabledAzurePlan) *[]billingprofile.AzurePlan {
	result := make([]billingprofile.AzurePlan, 0)
	if len(input) == 0 {
		return &result
	}

	for _, item := range input {
		azurePlan := billingprofile.AzurePlan{}

		if v := item.SkuDescription; v != "" {
			azurePlan.SkuDescription = pointer.To(v)
		}

		if v := item.SkuId; v != "" {
			azurePlan.SkuId = pointer.To(v)
		}

		result = append(result, azurePlan)
	}

	return &result
}

func flattenEnabledAzurePlans(input *[]billingprofile.AzurePlan) []EnabledAzurePlan {
	result := make([]EnabledAzurePlan, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, EnabledAzurePlan{
			SkuDescription: pointer.From(item.SkuDescription),
			SkuId:          pointer.From(item.SkuId),
		})
	}

	return result
}

func expandBillTo(input []BillTo) *billingprofile.AddressDetails {
	if len(input) == 0 {
		return nil
	}

	addressDetails := input[0]

	result := &billingprofile.AddressDetails{
		AddressLine1: addressDetails.AddressLine1,
		Country:      addressDetails.Country,
	}

	if v := addressDetails.AddressLine2; v != "" {
		result.AddressLine2 = pointer.To(v)
	}

	if v := addressDetails.AddressLine3; v != "" {
		result.AddressLine3 = pointer.To(v)
	}

	if v := addressDetails.City; v != "" {
		result.City = pointer.To(v)
	}

	if v := addressDetails.CompanyName; v != "" {
		result.CompanyName = pointer.To(v)
	}

	if v := addressDetails.District; v != "" {
		result.District = pointer.To(v)
	}

	if v := addressDetails.Email; v != "" {
		result.Email = pointer.To(v)
	}

	if v := addressDetails.FirstName; v != "" {
		result.FirstName = pointer.To(v)
	}

	if v := addressDetails.LastName; v != "" {
		result.LastName = pointer.To(v)
	}

	if v := addressDetails.MiddleName; v != "" {
		result.MiddleName = pointer.To(v)
	}

	if v := addressDetails.PhoneNumber; v != "" {
		result.PhoneNumber = pointer.To(v)
	}

	if v := addressDetails.PostalCode; v != "" {
		result.PostalCode = pointer.To(v)
	}

	if v := addressDetails.Region; v != "" {
		result.Region = pointer.To(v)
	}

	return result
}

func flattenBillTo(input *billingprofile.AddressDetails) []BillTo {
	result := make([]BillTo, 0)
	if input == nil {
		return result
	}

	billTo := BillTo{
		AddressLine1: input.AddressLine1,
		AddressLine2: pointer.From(input.AddressLine2),
		AddressLine3: pointer.From(input.AddressLine3),
		City:         pointer.From(input.City),
		CompanyName:  pointer.From(input.CompanyName),
		Country:      input.Country,
		District:     pointer.From(input.District),
		Email:        pointer.From(input.Email),
		FirstName:    pointer.From(input.FirstName),
		LastName:     pointer.From(input.LastName),
		MiddleName:   pointer.From(input.MiddleName),
		PhoneNumber:  pointer.From(input.PhoneNumber),
		PostalCode:   pointer.From(input.PostalCode),
		Region:       pointer.From(input.Region),
	}

	return append(result, billTo)
}
