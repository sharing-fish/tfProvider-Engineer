package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewEngineerResource() resource.Resource {
	return &EngineerResource{}
}

var _ resource.Resource = &EngineerResource{}

type EngineerResource struct {
	client *Client
}

type EngineerResourceModel struct {
	Name  types.String `tfsdk:"name"`
	Id    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
}

func (r *EngineerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer-resource"
}

func (r *EngineerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *EngineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EngineerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create engineer via API
	engineer, err := r.client.CreateEngineer(data.Name.ValueString(), data.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create engineer",
			"An error occurred while creating the engineer: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(engineer.Id)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EngineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EngineerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch engineer from the API using GetEngineerById
	engineer, err := r.client.GetEngineerById(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch engineer",
			"An error occurred while fetching the engineer: "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(engineer.Name)
	data.Email = types.StringValue(engineer.Email)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EngineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EngineerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update engineer via API
	engineer, err := r.client.UpdateEngineer(data.Id.ValueString(), data.Name.ValueString(), data.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update engineer",
			"An error occurred while updating the engineer: "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(engineer.Name)
	data.Email = types.StringValue(engineer.Email)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EngineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EngineerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete engineer via API
	err := r.client.DeleteEngineer(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete engineer",
			"An error occurred while deleting the engineer: "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// Configure adds the provider configured client to the resource.
func (r *EngineerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
