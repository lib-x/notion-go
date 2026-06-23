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
func (s *CommentsService) Create(ctx context.Context, body any) (Comment, error) {
	out, err := s.client.postObject(ctx, apiPath("v1", "comments"), body)
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
	return s.client.getList(ctx, apiPath("v1", "comments"), q)
}

// Get retrieves a comment.
func (s *CommentsService) Get(ctx context.Context, commentID string) (Comment, error) {
	out, err := s.client.getObject(ctx, apiPath("v1", "comments", commentID), nil)
	return Comment(out), err
}

// Update updates a comment.
func (s *CommentsService) Update(ctx context.Context, commentID string, body any) (Comment, error) {
	out, err := s.client.patchObject(ctx, apiPath("v1", "comments", commentID), body)
	return Comment(out), err
}

// Delete deletes a comment.
func (s *CommentsService) Delete(ctx context.Context, commentID string) (Comment, error) {
	out, err := s.client.deleteObject(ctx, apiPath("v1", "comments", commentID))
	return Comment(out), err
}
