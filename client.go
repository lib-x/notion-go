package notion

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// Client is a Notion API client.
type Client struct {
	httpClient    *http.Client
	baseURL       *url.URL
	token         string
	notionVersion string
	userAgent     string

	Users        *UsersService
	Pages        *PagesService
	Blocks       *BlocksService
	DataSources  *DataSourcesService
	Databases    *DatabasesService
	Search       *SearchService
	Comments     *CommentsService
	FileUploads  *FileUploadsService
	CustomEmojis *CustomEmojisService
	Views        *ViewsService
	MeetingNotes *MeetingNotesService
	OAuth        *OAuthService
}

// Option configures a Client.
type Option func(*Client) error

// WithHTTPClient sets the HTTP client. A nil client is rejected.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("notion: nil http client")
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithBaseURL overrides the API base URL. It is mainly useful for tests.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) error {
		u, err := url.Parse(strings.TrimRight(rawURL, "/"))
		if err != nil {
			return fmt.Errorf("notion: parse base URL: %w", err)
		}
		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("notion: base URL must be absolute")
		}
		c.baseURL = u
		return nil
	}
}

// WithNotionVersion overrides the Notion-Version header.
func WithNotionVersion(version string) Option {
	return func(c *Client) error {
		if strings.TrimSpace(version) == "" {
			return fmt.Errorf("notion: empty notion version")
		}
		c.notionVersion = version
		return nil
	}
}

// WithUserAgent overrides the default User-Agent.
func WithUserAgent(userAgent string) Option {
	return func(c *Client) error {
		if strings.TrimSpace(userAgent) == "" {
			return fmt.Errorf("notion: empty user agent")
		}
		c.userAgent = userAgent
		return nil
	}
}

// NewClient creates a Notion API client. token may be empty for OAuth methods
// that use HTTP Basic authentication.
func NewClient(token string, opts ...Option) (*Client, error) {
	baseURL, err := url.Parse(DefaultBaseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		httpClient:    http.DefaultClient,
		baseURL:       baseURL,
		token:         token,
		notionVersion: LatestVersion,
		userAgent:     "notion-go",
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	c.Users = &UsersService{client: c}
	c.Pages = &PagesService{client: c}
	c.Blocks = &BlocksService{client: c}
	c.DataSources = &DataSourcesService{client: c}
	c.Databases = &DatabasesService{client: c}
	c.Search = &SearchService{client: c}
	c.Comments = &CommentsService{client: c}
	c.FileUploads = &FileUploadsService{client: c}
	c.CustomEmojis = &CustomEmojisService{client: c}
	c.Views = &ViewsService{client: c}
	c.MeetingNotes = &MeetingNotesService{client: c}
	c.OAuth = &OAuthService{client: c}
	return c, nil
}

// MustNewClient is like NewClient but panics on configuration errors.
func MustNewClient(token string, opts ...Option) *Client {
	c, err := NewClient(token, opts...)
	if err != nil {
		panic(err)
	}
	return c
}

// RequestOptions configures a raw request made with Do.
type RequestOptions struct {
	Query     url.Values
	Body      any
	Headers   http.Header
	BasicAuth *BasicAuth
}

// BasicAuth configures HTTP Basic authentication for OAuth endpoints.
type BasicAuth struct {
	Username string
	Password string
}

// Do performs a raw JSON request against the Notion API.
func (c *Client) Do(ctx context.Context, method, path string, opts *RequestOptions, out any) error {
	if opts == nil {
		opts = &RequestOptions{}
	}
	req, err := c.newJSONRequest(ctx, method, path, opts.Query, opts.Body)
	if err != nil {
		return err
	}
	if opts.BasicAuth != nil {
		setBasicAuth(req, opts.BasicAuth.Username, opts.BasicAuth.Password)
	}
	for key, values := range opts.Headers {
		req.Header.Del(key)
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return c.do(req, out)
}

func (c *Client) get(ctx context.Context, path string, query url.Values, out any) error {
	req, err := c.newJSONRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) post(ctx context.Context, path string, body any, out any) error {
	return c.sendJSON(ctx, http.MethodPost, path, nil, body, out)
}

func (c *Client) patch(ctx context.Context, path string, body any, out any) error {
	return c.sendJSON(ctx, http.MethodPatch, path, nil, body, out)
}

func (c *Client) delete(ctx context.Context, path string, out any) error {
	return c.sendJSON(ctx, http.MethodDelete, path, nil, nil, out)
}

func (c *Client) postBasic(ctx context.Context, path, clientID, clientSecret string, body any, out any) error {
	req, err := c.newJSONRequest(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return err
	}
	setBasicAuth(req, clientID, clientSecret)
	return c.do(req, out)
}

func (c *Client) sendJSON(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	req, err := c.newJSONRequest(ctx, method, path, query, body)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) newJSONRequest(ctx context.Context, method, path string, query url.Values, body any) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("notion: encode request body: %w", err)
		}
		r = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.resolve(path, query), r)
	if err != nil {
		return nil, err
	}
	c.setDefaultHeaders(req)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) doMultipart(ctx context.Context, path string, body *bytes.Buffer, contentType string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.resolve(path, nil), body)
	if err != nil {
		return err
	}
	c.setDefaultHeaders(req)
	req.Header.Set("Content-Type", contentType)
	return c.do(req, out)
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("notion: read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return decodeAPIError(resp, data)
	}
	if out == nil || len(data) == 0 {
		return nil
	}
	if raw, ok := out.(*[]byte); ok {
		*raw = append((*raw)[:0], data...)
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("notion: decode response body: %w", err)
	}
	return nil
}

func (c *Client) resolve(path string, query url.Values) string {
	u := *c.baseURL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		parsed, _ := url.Parse(path)
		u = *parsed
	} else {
		u.Path = strings.TrimRight(c.baseURL.Path, "/") + "/" + strings.TrimLeft(path, "/")
	}
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String()
}

func (c *Client) setDefaultHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Notion-Version", c.notionVersion)
	req.Header.Set("User-Agent", c.userAgent)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}

func setBasicAuth(req *http.Request, username, password string) {
	token := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+token)
}

func paginationValues(params *PaginationParams) url.Values {
	q := make(url.Values)
	if params == nil {
		return q
	}
	if params.StartCursor != "" {
		q.Set("start_cursor", params.StartCursor)
	}
	if params.PageSize > 0 {
		q.Set("page_size", fmt.Sprintf("%d", params.PageSize))
	}
	return q
}

func addString(q url.Values, key, value string) {
	if value != "" {
		q.Set(key, value)
	}
}

func addBool(q url.Values, key string, value *bool) {
	if value != nil {
		q.Set(key, fmt.Sprintf("%t", *value))
	}
}

func addStrings(q url.Values, key string, values []string) {
	for _, value := range values {
		if value != "" {
			q.Add(key, value)
		}
	}
}

func apiPath(parts ...string) string {
	if len(parts) == 0 {
		return "/"
	}
	escaped := make([]string, 0, len(parts))
	for _, part := range parts {
		escaped = append(escaped, url.PathEscape(part))
	}
	return "/" + strings.Join(escaped, "/")
}

func writeMultipartFile(writer *multipart.Writer, fieldName, filename string, r io.Reader) error {
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, r)
	return err
}
