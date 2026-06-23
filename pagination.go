package notion

import (
	"context"
	"encoding/json"
	"fmt"
)

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

// ListResponseOf is a typed view of a Notion cursor-paginated list envelope.
type ListResponseOf[T any] struct {
	Object        string        `json:"object,omitempty"`
	Results       []T           `json:"results,omitempty"`
	NextCursor    *string       `json:"next_cursor"`
	HasMore       bool          `json:"has_more"`
	Type          string        `json:"type,omitempty"`
	RequestStatus RequestStatus `json:"request_status,omitempty"`
	Raw           Object        `json:"-"`
}

// DecodeList decodes a generic ListResponse into typed result values.
func DecodeList[T any](list *ListResponse) (*ListResponseOf[T], error) {
	if list == nil {
		return nil, fmt.Errorf("notion: decode nil list response")
	}
	var results []T
	raw := list.RawResults()
	if len(raw) == 0 && len(list.Results) > 0 {
		var err error
		raw, err = json.Marshal(list.Results)
		if err != nil {
			return nil, err
		}
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &results); err != nil {
			return nil, err
		}
	}
	return &ListResponseOf[T]{
		Object:        list.Object,
		Results:       results,
		NextCursor:    list.NextCursor,
		HasMore:       list.HasMore,
		Type:          list.Type,
		RequestStatus: list.RequestStatus,
		Raw:           list.Raw,
	}, nil
}

// CollectPaginatedDecode fetches all pages from a generic ListResponse
// endpoint and decodes each result into T.
func CollectPaginatedDecode[T any](ctx context.Context, fetch PageFetcher) ([]T, error) {
	return CollectPaginatedAs(ctx, func(ctx context.Context, cursor string) (*ListResponseOf[T], error) {
		page, err := fetch(ctx, cursor)
		if err != nil {
			return nil, err
		}
		return DecodeList[T](page)
	})
}

// PageFetcherOf fetches one generic list page.
type PageFetcherOf[T any] func(ctx context.Context, cursor string) (*ListResponseOf[T], error)

// CollectPaginatedAs fetches all typed results from a cursor-paginated endpoint.
func CollectPaginatedAs[T any](ctx context.Context, fetch PageFetcherOf[T]) ([]T, error) {
	var all []T
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
