package notion

import (
	"context"
	"net/url"
)

// PagesService implements Notion page endpoints.
type PagesService struct {
	client *Client
}

// GetPageParams configures page retrieval.
type GetPageParams struct {
	FilterProperties []string
}

// RetrieveMarkdownParams configures page markdown retrieval.
type RetrieveMarkdownParams struct {
	IncludeTranscript *bool
}

// Create creates a page.
func (s *PagesService) Create(ctx context.Context, body Object) (Page, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "pages"), body, &out)
	return Page(out), err
}

// Get retrieves a page.
func (s *PagesService) Get(ctx context.Context, pageID string, params *GetPageParams) (Page, error) {
	q := make(url.Values)
	if params != nil {
		addStrings(q, "filter_properties", params.FilterProperties)
	}
	var out Object
	err := s.client.get(ctx, apiPath("v1", "pages", pageID), q, &out)
	return Page(out), err
}

// Update updates a page.
func (s *PagesService) Update(ctx context.Context, pageID string, body Object) (Page, error) {
	var out Object
	err := s.client.patch(ctx, apiPath("v1", "pages", pageID), body, &out)
	return Page(out), err
}

// Move moves a page.
func (s *PagesService) Move(ctx context.Context, pageID string, body Object) (Page, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "pages", pageID, "move"), body, &out)
	return Page(out), err
}

// GetProperty retrieves a page property item.
func (s *PagesService) GetProperty(ctx context.Context, pageID, propertyID string, params *PaginationParams) (PageProperty, error) {
	var out Object
	err := s.client.get(ctx, apiPath("v1", "pages", pageID, "properties", propertyID), paginationValues(params), &out)
	return PageProperty(out), err
}

// RetrieveMarkdown retrieves a page or block subtree as markdown.
func (s *PagesService) RetrieveMarkdown(ctx context.Context, pageID string, params *RetrieveMarkdownParams) (*MarkdownResponse, error) {
	q := make(url.Values)
	if params != nil {
		addBool(q, "include_transcript", params.IncludeTranscript)
	}
	var out MarkdownResponse
	err := s.client.get(ctx, apiPath("v1", "pages", pageID, "markdown"), q, &out)
	return &out, err
}

// UpdateMarkdown updates a page's content using Notion's markdown mutation API.
func (s *PagesService) UpdateMarkdown(ctx context.Context, pageID string, body Object) (*MarkdownResponse, error) {
	var out MarkdownResponse
	err := s.client.patch(ctx, apiPath("v1", "pages", pageID, "markdown"), body, &out)
	return &out, err
}
