package notion

import (
	"context"
	"net/url"
)

// CommentsService implements Notion comment endpoints.
type CommentsService struct {
	client *Client
}

// ListCommentsParams configures comment listing.
type ListCommentsParams struct {
	BlockID string
	PaginationParams
}

// Create creates a comment.
func (s *CommentsService) Create(ctx context.Context, body Object) (Comment, error) {
	var out Object
	err := s.client.post(ctx, apiPath("v1", "comments"), body, &out)
	return Comment(out), err
}

// List lists comments.
func (s *CommentsService) List(ctx context.Context, params *ListCommentsParams) (*ListResponse, error) {
	q := make(url.Values)
	if params != nil {
		addString(q, "block_id", params.BlockID)
		for key, values := range paginationValues(&params.PaginationParams) {
			q[key] = values
		}
	}
	var out ListResponse
	err := s.client.get(ctx, apiPath("v1", "comments"), q, &out)
	return &out, err
}

// Get retrieves a comment.
func (s *CommentsService) Get(ctx context.Context, commentID string) (Comment, error) {
	var out Object
	err := s.client.get(ctx, apiPath("v1", "comments", commentID), nil, &out)
	return Comment(out), err
}

// Update updates a comment.
func (s *CommentsService) Update(ctx context.Context, commentID string, body Object) (Comment, error) {
	var out Object
	err := s.client.patch(ctx, apiPath("v1", "comments", commentID), body, &out)
	return Comment(out), err
}

// Delete deletes a comment.
func (s *CommentsService) Delete(ctx context.Context, commentID string) (Comment, error) {
	var out Object
	err := s.client.delete(ctx, apiPath("v1", "comments", commentID), &out)
	return Comment(out), err
}
