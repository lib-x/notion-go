package notion

import (
	"context"
	"net/url"
)

// DataSourcesService implements Notion data source endpoints.
type DataSourcesService struct {
	client *Client
}

// QueryDataSourceParams configures a data source query.
type QueryDataSourceParams struct {
	FilterProperties []string
}

// ListDataSourceTemplatesParams configures data source template listing.
type ListDataSourceTemplatesParams struct {
	Name string
	PaginationParams
}

// Get retrieves a data source.
func (s *DataSourcesService) Get(ctx context.Context, dataSourceID string) (DataSource, error) {
	out, err := s.client.getObject(ctx, apiPath("v1", "data_sources", dataSourceID), nil)
	return DataSource(out), err
}

// Update updates a data source.
func (s *DataSourcesService) Update(ctx context.Context, dataSourceID string, body any) (DataSource, error) {
	out, err := s.client.patchObject(ctx, apiPath("v1", "data_sources", dataSourceID), body)
	return DataSource(out), err
}

// Query queries a data source.
func (s *DataSourcesService) Query(ctx context.Context, dataSourceID string, body any, params *QueryDataSourceParams) (*ListResponse, error) {
	q := make(url.Values)
	if params != nil {
		addStrings(q, "filter_properties", params.FilterProperties)
	}
	return s.client.postList(ctx, apiPath("v1", "data_sources", dataSourceID, "query"), q, body)
}

// Create creates a data source.
func (s *DataSourcesService) Create(ctx context.Context, body any) (DataSource, error) {
	out, err := s.client.postObject(ctx, apiPath("v1", "data_sources"), body)
	return DataSource(out), err
}

// ListTemplates lists templates in a data source.
func (s *DataSourcesService) ListTemplates(ctx context.Context, dataSourceID string, params *ListDataSourceTemplatesParams) (*ListResponse, error) {
	q := make(url.Values)
	if params != nil {
		addString(q, "name", params.Name)
		for key, values := range paginationValues(&params.PaginationParams) {
			q[key] = values
		}
	}
	return s.client.getList(ctx, apiPath("v1", "data_sources", dataSourceID, "templates"), q)
}
