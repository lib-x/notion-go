package notion

import "context"

// UsersService implements Notion user endpoints.
type UsersService struct {
	client *Client
}

// Me retrieves the bot user associated with the current token.
func (s *UsersService) Me(ctx context.Context) (User, error) {
	out, err := s.client.getObject(ctx, apiPath("v1", "users", "me"), nil)
	return User(out), err
}

// Get retrieves a user by ID.
func (s *UsersService) Get(ctx context.Context, userID string) (User, error) {
	out, err := s.client.getObject(ctx, apiPath("v1", "users", userID), nil)
	return User(out), err
}

// List lists users with cursor pagination.
func (s *UsersService) List(ctx context.Context, params *PaginationParams) (*ListResponse, error) {
	return s.client.getList(ctx, apiPath("v1", "users"), paginationValues(params))
}
