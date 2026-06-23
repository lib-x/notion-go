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
	var out Object
	err := s.client.get(ctx, apiPath("v1", "data_sources", dataSourceID), nil, &out)
	return DataSource(out), err
}

// Update updates a data source.
func (s *DataSourcesService) Update(ctx context.Context, dataSourceID string, body Object) (DataSource, error) {
	var out Object
	err := s.client.patch(ctx, apiPath("v1", "data_sources", dataSourceID), body, &out)
	return DataSource(out), err
}

// Query queries a data source.
func (s *DataSourcesService) Query(ctx context.Context, dataSourceID string, body QueryRequest, params *QueryDataSourceParams) (*ListResponse, error) {
	q := make(url.Values)
	if params != nil {
		addStrings(q, "filter_properties", params.FilterProperties)
	}
	var out ListResponse
	err := s.client.sendJSON(ctx, "POST", apiPath("v1", "data_sources", dataSourceID, "query"), q, Object(body), &out)
	return &out, err
}

// Create creates a data source.
func (s *DataSourcesService) Create(ctx context.Context, body Object) (DataSource, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "data_sources"), body, &out)
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
	var out ListResponse
	err := s.client.get(ctx, apiPath("v1", "data_sources", dataSourceID, "templates"), q, &out)
	return &out, err
}
