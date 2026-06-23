package notion

import "context"

// DatabasesService implements Notion database endpoints.
//
// Data sources are the current API surface for querying database-like content;
// these database endpoints remain available for compatibility.
type DatabasesService struct {
	client *Client
}

// Get retrieves a database.
func (s *DatabasesService) Get(ctx context.Context, databaseID string) (Database, error) {
	out, err := s.client.getObject(ctx, apiPath("v1", "databases", databaseID), nil)
	return Database(out), err
}

// Update updates a database.
func (s *DatabasesService) Update(ctx context.Context, databaseID string, body any) (Database, error) {
	out, err := s.client.patchObject(ctx, apiPath("v1", "databases", databaseID), body)
	return Database(out), err
}

// Create creates a database.
func (s *DatabasesService) Create(ctx context.Context, body any) (Database, error) {
	out, err := s.client.postObject(ctx, apiPath("v1", "databases"), body)
	return Database(out), err
}
