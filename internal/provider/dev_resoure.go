package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDevResource() resource.Resource {
	return &DevResource{}
}

var _ resource.Resource = &DevResource{}

type DevResource struct {
	client *Client
}

type DevResourceModel struct {
	Name      types.String    `tfsdk:"name"`
	Id        types.String    `tfsdk:"id"`
	Engineers []EngineerModel `tfsdk:"engineers"`
}

func (r *DevResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev_resource"
}

func (r *DevResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"engineers": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional: true,
						},
						"id": schema.StringAttribute{
							Optional: true,
						},
						"email": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (r *DevResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DevResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create Dev via API
	dev, err := r.client.CreateDev(data.Name.ValueString(), data.Engineers)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create dev",
			"An error occurred while creating the dev: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(dev.Id)
	data.Name = types.StringValue(dev.Name)
	data.Engineers = make([]EngineerModel, len(dev.Engineers))
	for i, engineer := range dev.Engineers {
		data.Engineers[i].Name = engineer.Name
		data.Engineers[i].Id = engineer.Id
		data.Engineers[i].Email = engineer.Email
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DevResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DevResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch dev from the API using GetDevById
	dev, err := r.client.GetDevById(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch dev",
			"An error occurred while fetching the dev: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(dev.Id)
	data.Name = types.StringValue(dev.Name)
	data.Engineers = make([]EngineerModel, len(dev.Engineers))
	for i, engineer := range dev.Engineers {
		data.Engineers[i].Name = engineer.Name
		data.Engineers[i].Id = engineer.Id
		data.Engineers[i].Email = engineer.Email
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DevResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DevResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update dev via API
	dev, err := r.client.UpdateDev(data.Id.ValueString(), data.Name.ValueString(), data.Engineers)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update dev",
			"An error occurred while updating the dev: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(dev.Id)
	data.Name = types.StringValue(dev.Name)
	data.Engineers = make([]EngineerModel, len(dev.Engineers))
	for i, engineer := range dev.Engineers {
		data.Engineers[i].Name = engineer.Name
		data.Engineers[i].Id = engineer.Id
		data.Engineers[i].Email = engineer.Email
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DevResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DevResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete dev via API
	err := r.client.DeleteDev(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete dev",
			"An error occurred while deleting the dev: "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// Configure adds the provider configured client to the resource.
func (r *DevResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DevResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
