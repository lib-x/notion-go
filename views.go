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
	return s.client.getList(ctx, apiPath("v1", "views"), q)
}

// Create creates a view.
func (s *ViewsService) Create(ctx context.Context, body any) (View, error) {
	out, err := s.client.postObject(ctx, apiPath("v1", "views"), body)
	return View(out), err
}

// Get retrieves a view.
func (s *ViewsService) Get(ctx context.Context, viewID string) (View, error) {
	out, err := s.client.getObject(ctx, apiPath("v1", "views", viewID), nil)
	return View(out), err
}

// Update updates a view.
func (s *ViewsService) Update(ctx context.Context, viewID string, body any) (View, error) {
	out, err := s.client.patchObject(ctx, apiPath("v1", "views", viewID), body)
	return View(out), err
}

// Delete deletes a view.
func (s *ViewsService) Delete(ctx context.Context, viewID string) (View, error) {
	out, err := s.client.deleteObject(ctx, apiPath("v1", "views", viewID))
	return View(out), err
}

// CreateQuery creates a view query.
func (s *ViewsService) CreateQuery(ctx context.Context, viewID string, body any) (ViewQuery, error) {
	out, err := s.client.postObject(ctx, apiPath("v1", "views", viewID, "queries"), body)
	return ViewQuery(out), err
}

// GetQueryResults returns paginated results for a view query.
func (s *ViewsService) GetQueryResults(ctx context.Context, viewID, queryID string, params *PaginationParams) (*ListResponse, error) {
	return s.client.getList(ctx, apiPath("v1", "views", viewID, "queries", queryID), paginationValues(params))
}

// DeleteQuery deletes a view query.
func (s *ViewsService) DeleteQuery(ctx context.Context, viewID, queryID string) (ViewQuery, error) {
	out, err := s.client.deleteObject(ctx, apiPath("v1", "views", viewID, "queries", queryID))
	return ViewQuery(out), err
}
