package notion

import "context"

// BlocksService implements Notion block endpoints.
type BlocksService struct {
	client *Client
}

// Get retrieves a block.
func (s *BlocksService) Get(ctx context.Context, blockID string) (Block, error) {
	var out Object
	err := s.client.get(ctx, apiPath("v1", "blocks", blockID), nil, &out)
	return Block(out), err
}

// Update updates a block.
func (s *BlocksService) Update(ctx context.Context, blockID string, body Object) (Block, error) {
	var out Object
	err := s.client.patch(ctx, apiPath("v1", "blocks", blockID), body, &out)
	return Block(out), err
}

// Delete deletes a block.
func (s *BlocksService) Delete(ctx context.Context, blockID string) (Block, error) {
	var out Object
	err := s.client.delete(ctx, apiPath("v1", "blocks", blockID), &out)
	return Block(out), err
}

// Children lists a block's children.
func (s *BlocksService) Children(ctx context.Context, blockID string, params *PaginationParams) (*ListResponse, error) {
	var out ListResponse
	err := s.client.get(ctx, apiPath("v1", "blocks", blockID, "children"), paginationValues(params), &out)
	return &out, err
}

// AppendChildren appends child blocks.
func (s *BlocksService) AppendChildren(ctx context.Context, blockID string, body Object) (Block, error) {
	var out Object
	err := s.client.patch(ctx, apiPath("v1", "blocks", blockID, "children"), body, &out)
	return Block(out), err
}
