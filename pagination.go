package notion

import "context"

// PageFetcher fetches one page of results. The cursor is empty for the first
// page and should be passed to the underlying endpoint as start_cursor.
type PageFetcher func(ctx context.Context, cursor string) (*ListResponse, error)

// CollectPaginated fetches all pages from a cursor-paginated endpoint.
func CollectPaginated(ctx context.Context, fetch PageFetcher) ([]Object, error) {
	var all []Object
	var cursor string
	for {
		page, err := fetch(ctx, cursor)
		if err != nil {
			return nil, err
		}
		all = append(all, page.Results...)
		if !page.HasMore || page.NextCursor == nil || *page.NextCursor == "" {
			return all, nil
		}
		cursor = *page.NextCursor
	}
}
