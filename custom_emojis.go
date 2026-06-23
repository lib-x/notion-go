package notion

import (
	"context"
	"net/url"
)

// CustomEmojisService implements Notion custom emoji endpoints.
type CustomEmojisService struct {
	client *Client
}

// ListCustomEmojisParams configures custom emoji listing.
type ListCustomEmojisParams struct {
	Name string
	PaginationParams
}

// List lists custom emojis.
func (s *CustomEmojisService) List(ctx context.Context, params *ListCustomEmojisParams) (*ListResponse, error) {
	q := make(url.Values)
	if params != nil {
		addString(q, "name", params.Name)
		for key, values := range paginationValues(&params.PaginationParams) {
			q[key] = values
		}
	}
	var out ListResponse
	err := s.client.get(ctx, apiPath("v1", "custom_emojis"), q, &out)
	return &out, err
}
