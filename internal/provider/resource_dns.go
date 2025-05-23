// Copyright (c) The Nanoid Provider for Terraform Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gonanoid "github.com/matoous/go-nanoid"
)

const DEFAULT_DNS_ALPHABET = "0123456789abcdefghijklmnopqrstuvwxyz"
const DEFAULT_DNS_LENGTH = 10

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DnsResource{}
var _ resource.ResourceWithImportState = &DnsResource{}

func NewDnsResource() resource.Resource {
	return &DnsResource{}
}

// DnsResource defines the data source implementation.
type DnsResource struct{}

// DnsResourceModel describes the data source data model.
type DnsResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Keepers types.Map    `tfsdk:"keepers"`
	Length  types.Int64  `tfsdk:"length"`
}

func (d *DnsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns"
}

func (d *DnsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("The dns resource generates hostname/dns friendly random strings that are intended to be used as unique identifiers for other resources.\n\n"+
			"The alphabet used is `\"%q\"`\n\n"+
			"This resource can be used in conjunction with resources that have the `create_before_destroy` lifecycle flag set to avoid conflicts with "+
			"unique names during the brief period where both the old and new resources exist concurrently.", DEFAULT_DNS_ALPHABET),
		Attributes: map[string]schema.Attribute{
			"length": schema.Int64Attribute{
				MarkdownDescription: fmt.Sprintf("The length of the desired nanoid.\nShould be between 1 and 64.\nThe default value is %d.", DEFAULT_ID_LENGTH),
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(DEFAULT_DNS_LENGTH),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(1, 64),
				},
			},

			"keepers": schema.MapAttribute{
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplaceIfConfigured(),
				},
			},

			"id": schema.StringAttribute{
				MarkdownDescription: "The generated random string.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (d *DnsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	_, ok := req.ProviderData.(*NanoidProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *provider.NanoidProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
}

func (r *DnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DnsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alphabet := DEFAULT_DNS_ALPHABET
	length := data.Length.ValueInt64()
	if data.Length.IsNull() {
		length = DEFAULT_DNS_LENGTH
	}

	id, err := gonanoid.Generate(alphabet, int(length))
	if err != nil {
		resp.Diagnostics.AddError("Failed to generate id", fmt.Sprintf("Failed to generate id: %s.", err))
		return
	}

	data.Id = types.StringValue(id)
	data.Length = types.Int64Value(length)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DnsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DnsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DnsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DnsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DnsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DnsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	length := len(id)
	if length > 64 {
		resp.Diagnostics.AddError("Invalid id", "The id must be at most 64 characters long.")
		return
	}

	state := &DnsResourceModel{
		Id:      types.StringValue(id),
		Length:  types.Int64Value(int64(length)),
		Keepers: types.MapNull(types.StringType),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
