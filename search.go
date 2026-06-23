package notion

import "context"

// SearchService implements the Notion search endpoint.
type SearchService struct {
	client *Client
}

// Do searches pages and databases by title.
func (s *SearchService) Do(ctx context.Context, body SearchRequest) (*ListResponse, error) {
	return s.client.postList(ctx, apiPath("v1", "search"), nil, Object(body))
}
