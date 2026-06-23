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
	var out Object
	err := s.client.get(ctx, apiPath("v1", "databases", databaseID), nil, &out)
	return Database(out), err
}

// Update updates a database.
func (s *DatabasesService) Update(ctx context.Context, databaseID string, body Object) (Database, error) {
	var out Object
	err := s.client.patch(ctx, apiPath("v1", "databases", databaseID), body, &out)
	return Database(out), err
}

// Create creates a database.
func (s *DatabasesService) Create(ctx context.Context, body Object) (Database, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "databases"), body, &out)
	return Database(out), err
}
