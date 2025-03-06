// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/billingprofile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BillingProfileModel struct {
	Name                     string                     `tfschema:"name"`
	BillingAccountName       string                     `tfschema:"billing_account_name"`
	BillTo                   []AddressDetails           `tfschema:"bill_to"`
	CurrentPaymentTerm       []CurrentPaymentTerm       `tfschema:"current_payment_term"`
	DisplayName              string                     `tfschema:"display_name"`
	EnabledAzurePlans        []EnabledAzurePlan         `tfschema:"enabled_azure_plan"`
	IndirectRelationshipInfo []IndirectRelationshipInfo `tfschema:"indirect_relationship_info"`
	InvoiceEmailOptInEnabled bool                       `tfschema:"invoice_email_opt_in_enabled"`
	InvoiceRecipients        []string                   `tfschema:"invoice_recipients"`
	PoNumber                 string                     `tfschema:"po_number"`
	ShipTo                   []AddressDetails           `tfschema:"ship_to"`
	SoldTo                   []AddressDetails           `tfschema:"sold_to"`
	Tags                     map[string]string          `tfschema:"tags"`
}

type AddressDetails struct {
	AddressLine1        string `tfschema:"address_line_1"`
	AddressLine2        string `tfschema:"address_line_2"`
	AddressLine3        string `tfschema:"address_line_3"`
	City                string `tfschema:"city"`
	CompanyName         string `tfschema:"company_name"`
	Country             string `tfschema:"country"`
	District            string `tfschema:"district"`
	Email               string `tfschema:"email"`
	FirstName           string `tfschema:"first_name"`
	LastName            string `tfschema:"last_name"`
	MiddleName          string `tfschema:"middle_name"`
	PhoneNumber         string `tfschema:"phone_number"`
	PostalCode          string `tfschema:"postal_code"`
	Region              string `tfschema:"region"`
	ValidAddressEnabled bool   `tfschema:"valid_address_enabled"`
}

type CurrentPaymentTerm struct {
	EndDate   string `tfschema:"end_date"`
	StartDate string `tfschema:"start_date"`
	Term      string `tfschema:"term"`
}

type EnabledAzurePlan struct {
	ProductId      string `tfschema:"product_id"`
	SkuDescription string `tfschema:"sku_description"`
	SkuId          string `tfschema:"sku_id"`
}

type IndirectRelationshipInfo struct {
	BillingAccountName string `tfschema:"billing_account_name"`
	BillingProfileName string `tfschema:"billing_profile_name"`
	DisplayName        string `tfschema:"display_name"`
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
			ValidateFunc: validate.BillingProfileName,
		},

		"billing_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BillingAccountName,
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
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: location.DiffSuppressFunc,
					},

					"valid_address_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"current_payment_term": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"end_date": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"start_date": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"term": {
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
					"product_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

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

		"indirect_relationship_info": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"billing_account_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"billing_profile_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"display_name": {
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

		"invoice_recipients": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"po_number": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"ship_to": {
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

					"country": {
						Type:         pluginsdk.TypeString,
						Required:     true,
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
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: location.DiffSuppressFunc,
					},

					"valid_address_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"sold_to": {
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

					"country": {
						Type:         pluginsdk.TypeString,
						Required:     true,
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
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: location.DiffSuppressFunc,
					},

					"valid_address_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
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
					BillTo:                   expandAddressDetails(model.BillTo),
					CurrentPaymentTerm:       expandCurrentPaymentTerm(model.CurrentPaymentTerm),
					EnabledAzurePlans:        expandEnabledAzurePlans(model.EnabledAzurePlans),
					IndirectRelationshipInfo: expandIndirectRelationshipInfo(model.IndirectRelationshipInfo),
					InvoiceEmailOptIn:        pointer.To(model.InvoiceEmailOptInEnabled),
					InvoiceRecipients:        pointer.To(model.InvoiceRecipients),
					ShipTo:                   expandAddressDetails(model.ShipTo),
					SoldTo:                   expandAddressDetails(model.SoldTo),
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
					state.BillTo = flattenAddressDetails(props.BillTo)
					state.CurrentPaymentTerm = flattenCurrentPaymentTerm(props.CurrentPaymentTerm)
					state.DisplayName = pointer.From(props.DisplayName)
					state.EnabledAzurePlans = flattenEnabledAzurePlans(props.EnabledAzurePlans)
					state.IndirectRelationshipInfo = flattenIndirectRelationshipInfo(props.IndirectRelationshipInfo)
					state.InvoiceEmailOptInEnabled = pointer.From(props.InvoiceEmailOptIn)
					state.InvoiceRecipients = pointer.From(props.InvoiceRecipients)
					state.PoNumber = pointer.From(props.PoNumber)
					state.ShipTo = flattenAddressDetails(props.ShipTo)
					state.SoldTo = flattenAddressDetails(props.SoldTo)
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
				parameters.Properties.BillTo = expandAddressDetails(model.BillTo)
			}

			if metadata.ResourceData.HasChange("current_payment_term") {
				parameters.Properties.CurrentPaymentTerm = expandCurrentPaymentTerm(model.CurrentPaymentTerm)
			}

			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if metadata.ResourceData.HasChange("enabled_azure_plan") {
				parameters.Properties.EnabledAzurePlans = expandEnabledAzurePlans(model.EnabledAzurePlans)
			}

			if metadata.ResourceData.HasChange("indirect_relationship_info") {
				parameters.Properties.IndirectRelationshipInfo = expandIndirectRelationshipInfo(model.IndirectRelationshipInfo)
			}

			if metadata.ResourceData.HasChange("invoice_email_opt_in_enabled") {
				parameters.Properties.InvoiceEmailOptIn = pointer.To(model.InvoiceEmailOptInEnabled)
			}

			if metadata.ResourceData.HasChange("invoice_recipients") {
				parameters.Properties.InvoiceRecipients = pointer.To(model.InvoiceRecipients)
			}

			if metadata.ResourceData.HasChange("po_number") {
				parameters.Properties.PoNumber = pointer.To(model.PoNumber)
			}

			if metadata.ResourceData.HasChange("ship_to") {
				parameters.Properties.ShipTo = expandAddressDetails(model.ShipTo)
			}

			if metadata.ResourceData.HasChange("sold_to") {
				parameters.Properties.SoldTo = expandAddressDetails(model.SoldTo)
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

func expandAddressDetails(input []AddressDetails) *billingprofile.AddressDetails {
	if len(input) == 0 {
		return nil
	}

	addressDetails := input[0]

	result := &billingprofile.AddressDetails{
		AddressLine1:   addressDetails.AddressLine1,
		Country:        addressDetails.Country,
		IsValidAddress: pointer.To(addressDetails.ValidAddressEnabled),
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

func flattenAddressDetails(input *billingprofile.AddressDetails) []AddressDetails {
	result := make([]AddressDetails, 0)
	if input == nil {
		return result
	}

	return append(result, AddressDetails{
		AddressLine1:        input.AddressLine1,
		AddressLine2:        pointer.From(input.AddressLine2),
		AddressLine3:        pointer.From(input.AddressLine3),
		City:                pointer.From(input.City),
		CompanyName:         pointer.From(input.CompanyName),
		Country:             input.Country,
		District:            pointer.From(input.District),
		Email:               pointer.From(input.Email),
		FirstName:           pointer.From(input.FirstName),
		LastName:            pointer.From(input.LastName),
		MiddleName:          pointer.From(input.MiddleName),
		PhoneNumber:         pointer.From(input.PhoneNumber),
		PostalCode:          pointer.From(input.PostalCode),
		Region:              pointer.From(input.Region),
		ValidAddressEnabled: pointer.From(input.IsValidAddress),
	})
}

func expandCurrentPaymentTerm(input []CurrentPaymentTerm) *billingprofile.PaymentTerm {
	if len(input) == 0 {
		return nil
	}

	paymentTerm := input[0]

	result := billingprofile.PaymentTerm{}

	if v := paymentTerm.EndDate; v != "" {
		result.EndDate = pointer.To(v)
	}

	if v := paymentTerm.StartDate; v != "" {
		result.StartDate = pointer.To(v)
	}

	if v := paymentTerm.Term; v != "" {
		result.Term = pointer.To(v)
	}

	return &result
}

func flattenCurrentPaymentTerm(input *billingprofile.PaymentTerm) []CurrentPaymentTerm {
	result := make([]CurrentPaymentTerm, 0)
	if input == nil {
		return result
	}

	return append(result, CurrentPaymentTerm{
		EndDate:   pointer.From(input.EndDate),
		StartDate: pointer.From(input.StartDate),
		Term:      pointer.From(input.Term),
	})
}

func expandEnabledAzurePlans(input []EnabledAzurePlan) *[]billingprofile.AzurePlan {
	result := make([]billingprofile.AzurePlan, 0)
	if len(input) == 0 {
		return &result
	}

	for _, item := range input {
		azurePlan := billingprofile.AzurePlan{}

		if v := item.ProductId; v != "" {
			azurePlan.ProductId = pointer.To(v)
		}

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
			ProductId:      pointer.From(item.ProductId),
			SkuDescription: pointer.From(item.SkuDescription),
			SkuId:          pointer.From(item.SkuId),
		})
	}

	return result
}

func expandIndirectRelationshipInfo(input []IndirectRelationshipInfo) *billingprofile.IndirectRelationshipInfo {
	if len(input) == 0 {
		return nil
	}

	indirectRelationshipInfo := input[0]

	result := billingprofile.IndirectRelationshipInfo{}

	if v := indirectRelationshipInfo.BillingAccountName; v != "" {
		result.BillingAccountName = pointer.To(v)
	}

	if v := indirectRelationshipInfo.BillingProfileName; v != "" {
		result.BillingProfileName = pointer.To(v)
	}

	if v := indirectRelationshipInfo.DisplayName; v != "" {
		result.DisplayName = pointer.To(v)
	}

	return &result
}

func flattenIndirectRelationshipInfo(input *billingprofile.IndirectRelationshipInfo) []IndirectRelationshipInfo {
	result := make([]IndirectRelationshipInfo, 0)
	if input == nil {
		return result
	}

	return append(result, IndirectRelationshipInfo{
		BillingAccountName: pointer.From(input.BillingAccountName),
		BillingProfileName: pointer.From(input.BillingProfileName),
		DisplayName:        pointer.From(input.DisplayName),
	})
}
