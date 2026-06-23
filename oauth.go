package notion

import "context"

// OAuthService implements Notion OAuth endpoints.
type OAuthService struct {
	client *Client
}

// Token exchanges an authorization code or refresh token using OAuth Basic auth.
func (s *OAuthService) Token(ctx context.Context, clientID, clientSecret string, body CreateTokenRequest) (*OAuthTokenResponse, error) {
	var out OAuthTokenResponse
	err := s.client.postBasic(ctx, apiPath("v1", "oauth", "token"), clientID, clientSecret, body, &out)
	return &out, err
}

// ExchangeAuthorizationCode exchanges an authorization code for tokens.
func (s *OAuthService) ExchangeAuthorizationCode(ctx context.Context, clientID, clientSecret, code, redirectURI string) (*OAuthTokenResponse, error) {
	return s.Token(ctx, clientID, clientSecret, CreateTokenRequest{
		GrantType:   "authorization_code",
		Code:        code,
		RedirectURI: redirectURI,
	})
}

// RefreshToken refreshes an OAuth access token.
func (s *OAuthService) RefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*OAuthTokenResponse, error) {
	return s.Token(ctx, clientID, clientSecret, CreateTokenRequest{
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
	})
}

// Revoke revokes an OAuth token.
func (s *OAuthService) Revoke(ctx context.Context, clientID, clientSecret, token string) (*OAuthRequestResponse, error) {
	var out OAuthRequestResponse
	err := s.client.postBasic(ctx, apiPath("v1", "oauth", "revoke"), clientID, clientSecret, Object{"token": token}, &out)
	return &out, err
}

// Introspect introspects an OAuth token.
func (s *OAuthService) Introspect(ctx context.Context, clientID, clientSecret, token string) (*OAuthIntrospectionResponse, error) {
	var out OAuthIntrospectionResponse
	err := s.client.postBasic(ctx, apiPath("v1", "oauth", "introspect"), clientID, clientSecret, Object{"token": token}, &out)
	return &out, err
}
