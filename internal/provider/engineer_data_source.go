package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NewEngineerDataSource() datasource.DataSource {
	return &EngineerDataSource{}
}

var _ datasource.DataSource = &EngineerDataSource{}

type EngineerDataSourceModel struct {
	Engineer []EngineerModel `tfsdk:"engineer"`
}

type EngineerModel struct {
	Name  string `tfsdk:"name"`
	Id    string `tfsdk:"id"`
	Email string `tfsdk:"email"`
}

type EngineerDataSource struct {
	client *http.Client
}

func (d *EngineerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

func (d *EngineerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"engineer": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *EngineerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EngineerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	for _, engineer := range data.Engineer {
		engineerState := EngineerModel{
			Name:  engineer.Name,
			Id:    engineer.Id,
			Email: engineer.Email,
		}
		data.Engineer = append(data.Engineer, engineerState)
	}

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
