package notion

import "context"

// SearchService implements the Notion search endpoint.
type SearchService struct {
	client *Client
}

// Do searches pages and databases by title.
func (s *SearchService) Do(ctx context.Context, body SearchRequest) (*ListResponse, error) {
	var out ListResponse
	err := s.client.post(ctx, apiPath("v1", "search"), Object(body), &out)
	return &out, err
}
