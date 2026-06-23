package notion

import (
	"context"
	"net/url"
)

// ViewsService implements Notion view endpoints.
type ViewsService struct {
	client *Client
}

// ListViewsParams configures view listing.
type ListViewsParams struct {
	DatabaseID   string
	DataSourceID string
	PaginationParams
}

// List lists views for a database or data source.
func (s *ViewsService) List(ctx context.Context, params *ListViewsParams) (*ListResponse, error) {
	q := make(url.Values)
	if params != nil {
		addString(q, "database_id", params.DatabaseID)
		addString(q, "data_source_id", params.DataSourceID)
		for key, values := range paginationValues(&params.PaginationParams) {
			q[key] = values
		}
	}
	var out ListResponse
	err := s.client.get(ctx, apiPath("v1", "views"), q, &out)
	return &out, err
}

// Create creates a view.
func (s *ViewsService) Create(ctx context.Context, body Object) (View, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "views"), body, &out)
	return View(out), err
}

// Get retrieves a view.
func (s *ViewsService) Get(ctx context.Context, viewID string) (View, error) {
	var out Object
	err := s.client.get(ctx, apiPath("v1", "views", viewID), nil, &out)
	return View(out), err
}

// Update updates a view.
func (s *ViewsService) Update(ctx context.Context, viewID string, body Object) (View, error) {
	var out Object
	err := s.client.patch(ctx, apiPath("v1", "views", viewID), body, &out)
	return View(out), err
}

// Delete deletes a view.
func (s *ViewsService) Delete(ctx context.Context, viewID string) (View, error) {
	var out Object
	err := s.client.delete(ctx, apiPath("v1", "views", viewID), &out)
	return View(out), err
}

// CreateQuery creates a view query.
func (s *ViewsService) CreateQuery(ctx context.Context, viewID string, body Object) (ViewQuery, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "views", viewID, "queries"), body, &out)
	return ViewQuery(out), err
}

// GetQueryResults returns paginated results for a view query.
func (s *ViewsService) GetQueryResults(ctx context.Context, viewID, queryID string, params *PaginationParams) (*ListResponse, error) {
	var out ListResponse
	err := s.client.get(ctx, apiPath("v1", "views", viewID, "queries", queryID), paginationValues(params), &out)
	return &out, err
}

// DeleteQuery deletes a view query.
func (s *ViewsService) DeleteQuery(ctx context.Context, viewID, queryID string) (ViewQuery, error) {
	var out Object
	err := s.client.delete(ctx, apiPath("v1", "views", viewID, "queries", queryID), &out)
	return ViewQuery(out), err
}
