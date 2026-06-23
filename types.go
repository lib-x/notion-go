package notion

import (
	"encoding/json"
	"fmt"
)

const (
	// DefaultBaseURL is the production Notion API root.
	DefaultBaseURL = "https://api.notion.com"

	// LatestVersion is the latest Notion-Version advertised by the official
	// reference documentation used to build this SDK.
	LatestVersion = "2026-03-11"
)

// Object is a JSON object returned by Notion.
//
// Notion's API has many polymorphic resources whose shape changes by type
// field. Keeping resources as Object preserves forward compatibility while the
// SDK provides stable endpoint methods, pagination, upload, OAuth, and error
// behavior.
type Object map[string]any

type (
	User          = Object
	Page          = Object
	PageProperty  = Object
	Block         = Object
	DataSource    = Object
	Database      = Object
	Comment       = Object
	FileUpload    = Object
	CustomEmoji   = Object
	View          = Object
	ViewQuery     = Object
	MeetingNote   = Object
	OAuthOwner    = Object
	RequestStatus = Object
)

// ListResponse is Notion's cursor-paginated list envelope.
type ListResponse struct {
	Object        string          `json:"object,omitempty"`
	Results       []Object        `json:"results,omitempty"`
	NextCursor    *string         `json:"next_cursor"`
	HasMore       bool            `json:"has_more"`
	Type          string          `json:"type,omitempty"`
	RequestStatus RequestStatus   `json:"request_status,omitempty"`
	Raw           Object          `json:"-"`
	rawResults    json.RawMessage `json:"-"`
}

// UnmarshalJSON keeps the full list object while decoding results into
// generic Notion objects.
func (l *ListResponse) UnmarshalJSON(data []byte) error {
	type alias ListResponse
	var aux struct {
		*alias
		Results json.RawMessage `json:"results"`
	}
	aux.alias = (*alias)(l)
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if len(aux.Results) > 0 && string(aux.Results) != "null" {
		if err := json.Unmarshal(aux.Results, &l.Results); err != nil {
			return fmt.Errorf("decode list results: %w", err)
		}
		l.rawResults = aux.Results
	}
	_ = json.Unmarshal(data, &l.Raw)
	return nil
}

// RawResults returns the original JSON array for callers that want to decode a
// list into their own concrete structs.
func (l *ListResponse) RawResults() json.RawMessage {
	if len(l.rawResults) == 0 {
		return nil
	}
	out := make([]byte, len(l.rawResults))
	copy(out, l.rawResults)
	return out
}

// MarkdownResponse is returned by page markdown endpoints.
type MarkdownResponse struct {
	Markdown string `json:"markdown,omitempty"`
	Raw      Object `json:"-"`
}

func (r *MarkdownResponse) UnmarshalJSON(data []byte) error {
	type alias MarkdownResponse
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	_ = json.Unmarshal(data, &r.Raw)
	return nil
}

// OAuthTokenResponse is returned by /v1/oauth/token.
type OAuthTokenResponse struct {
	AccessToken          string     `json:"access_token,omitempty"`
	TokenType            string     `json:"token_type,omitempty"`
	RefreshToken         *string    `json:"refresh_token"`
	BotID                string     `json:"bot_id,omitempty"`
	WorkspaceIcon        *string    `json:"workspace_icon"`
	WorkspaceName        *string    `json:"workspace_name"`
	WorkspaceID          string     `json:"workspace_id,omitempty"`
	Owner                OAuthOwner `json:"owner,omitempty"`
	DuplicatedTemplateID *string    `json:"duplicated_template_id"`
	RequestID            string     `json:"request_id,omitempty"`
	Raw                  Object     `json:"-"`
}

func (r *OAuthTokenResponse) UnmarshalJSON(data []byte) error {
	type alias OAuthTokenResponse
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	_ = json.Unmarshal(data, &r.Raw)
	return nil
}

// OAuthIntrospectionResponse is returned by /v1/oauth/introspect.
type OAuthIntrospectionResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Raw       Object `json:"-"`
}

func (r *OAuthIntrospectionResponse) UnmarshalJSON(data []byte) error {
	type alias OAuthIntrospectionResponse
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	_ = json.Unmarshal(data, &r.Raw)
	return nil
}

// OAuthRequestResponse is a small success response used by token revocation.
type OAuthRequestResponse struct {
	RequestID string `json:"request_id,omitempty"`
	Raw       Object `json:"-"`
}

func (r *OAuthRequestResponse) UnmarshalJSON(data []byte) error {
	type alias OAuthRequestResponse
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	_ = json.Unmarshal(data, &r.Raw)
	return nil
}

// MeetingNotesResponse is returned by the meeting notes query endpoint.
type MeetingNotesResponse struct {
	Results []Object `json:"results,omitempty"`
	HasMore bool     `json:"has_more"`
	Raw     Object   `json:"-"`
}

func (r *MeetingNotesResponse) UnmarshalJSON(data []byte) error {
	type alias MeetingNotesResponse
	if err := json.Unmarshal(data, (*alias)(r)); err != nil {
		return err
	}
	_ = json.Unmarshal(data, &r.Raw)
	return nil
}

// PaginationParams are used by GET endpoints that accept start_cursor/page_size.
type PaginationParams struct {
	StartCursor string
	PageSize    int
}

// SearchRequest is the request body for /v1/search.
type SearchRequest Object

// QueryRequest is the request body for /query endpoints.
type QueryRequest Object

// CreateTokenRequest is the request body for exchanging an authorization code
// or refresh token.
type CreateTokenRequest struct {
	GrantType       string `json:"grant_type"`
	Code            string `json:"code,omitempty"`
	RedirectURI     string `json:"redirect_uri,omitempty"`
	RefreshToken    string `json:"refresh_token,omitempty"`
	ExternalAccount Object `json:"external_account,omitempty"`
}
