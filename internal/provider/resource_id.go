// Copyright (c) The Nanoid Provider for Terraform Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gonanoid "github.com/matoous/go-nanoid"
)

const DEFAULT_ID_ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz-"
const DEFAULT_ID_LENGTH = 21

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IdResource{}
var _ resource.ResourceWithImportState = &IdResource{}

func NewIdResource() resource.Resource {
	return &IdResource{}
}

// IdResource defines the data source implementation.
type IdResource struct{}

// IdResourceModel describes the data source data model.
type IdResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Alphabet types.String `tfsdk:"alphabet"`
	Keepers  types.Map    `tfsdk:"keepers"`
	Length   types.Int64  `tfsdk:"length"`
}

func (d *IdResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_id"
}

func (d *IdResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The id resource generates random strings that are intended to be used as unique identifiers for other resources.\n\n" +
			"This resource can be used in conjunction with resources that have the `create_before_destroy` lifecycle flag set to avoid conflicts with " +
			"unique names during the brief period where both the old and new resources exist concurrently.",
		Attributes: map[string]schema.Attribute{
			"alphabet": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("Supply your own list of characters to use for id generation.\n"+
					"Should be between 1 and 255 characters long.\n"+
					"The default value is `\"%q\"`.", DEFAULT_ID_ALPHABET),
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(DEFAULT_ID_ALPHABET),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},

			"length": schema.Int64Attribute{
				MarkdownDescription: fmt.Sprintf("The length of the desired nanoid.\nShould be between 1 and 64.\nThe default value is %d.", DEFAULT_ID_LENGTH),
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(DEFAULT_ID_LENGTH),
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

func (d *IdResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IdResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IdResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	alphabet := data.Alphabet.ValueString()
	if data.Alphabet.IsNull() {
		alphabet = DEFAULT_ID_ALPHABET
	}

	length := data.Length.ValueInt64()
	if data.Length.IsNull() {
		length = DEFAULT_ID_LENGTH
	}

	id, err := gonanoid.Generate(alphabet, int(length))
	if err != nil {
		resp.Diagnostics.AddError("Failed to generate id", fmt.Sprintf("Failed to generate id: %s.", err))
		return
	}

	data.Id = types.StringValue(id)
	data.Alphabet = types.StringValue(alphabet)
	data.Length = types.Int64Value(length)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *IdResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IdResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IdResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IdResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IdResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *IdResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	length := len(id)
	if length > 64 {
		resp.Diagnostics.AddError("Invalid id", "The id must be at most 64 characters long.")
		return
	}

	state := &IdResourceModel{
		Id:       types.StringValue(id),
		Length:   types.Int64Value(int64(length)),
		Keepers:  types.MapNull(types.StringType),
		Alphabet: types.StringValue(DEFAULT_ID_ALPHABET),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
