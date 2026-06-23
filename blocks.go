package notion

import "context"

// BlocksService implements Notion block endpoints.
type BlocksService struct {
	client *Client
}

// Get retrieves a block.
func (s *BlocksService) Get(ctx context.Context, blockID string) (Block, error) {
	out, err := s.client.getObject(ctx, apiPath("v1", "blocks", blockID), nil)
	return Block(out), err
}

// Update updates a block.
func (s *BlocksService) Update(ctx context.Context, blockID string, body any) (Block, error) {
	out, err := s.client.patchObject(ctx, apiPath("v1", "blocks", blockID), body)
	return Block(out), err
}

// Delete deletes a block.
func (s *BlocksService) Delete(ctx context.Context, blockID string) (Block, error) {
	out, err := s.client.deleteObject(ctx, apiPath("v1", "blocks", blockID))
	return Block(out), err
}

// Children lists a block's children.
func (s *BlocksService) Children(ctx context.Context, blockID string, params *PaginationParams) (*ListResponse, error) {
	return s.client.getList(ctx, apiPath("v1", "blocks", blockID, "children"), paginationValues(params))
}

// AppendChildren appends child blocks.
func (s *BlocksService) AppendChildren(ctx context.Context, blockID string, body any) (Block, error) {
	out, err := s.client.patchObject(ctx, apiPath("v1", "blocks", blockID, "children"), body)
	return Block(out), err
}
