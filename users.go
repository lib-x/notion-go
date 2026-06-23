package notion

import "context"

// UsersService implements Notion user endpoints.
type UsersService struct {
	client *Client
}

// Me retrieves the bot user associated with the current token.
func (s *UsersService) Me(ctx context.Context) (User, error) {
	var out Object
	err := s.client.get(ctx, apiPath("v1", "users", "me"), nil, &out)
	return User(out), err
}

// Get retrieves a user by ID.
func (s *UsersService) Get(ctx context.Context, userID string) (User, error) {
	var out Object
	err := s.client.get(ctx, apiPath("v1", "users", userID), nil, &out)
	return User(out), err
}

// List lists users with cursor pagination.
func (s *UsersService) List(ctx context.Context, params *PaginationParams) (*ListResponse, error) {
	var out ListResponse
	err := s.client.get(ctx, apiPath("v1", "users"), paginationValues(params), &out)
	return &out, err
}
