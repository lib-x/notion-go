package notion

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestEndpointCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	calls := []struct {
		name       string
		method     string
		path       string
		query      url.Values
		response   string
		invoke     func(*Client) error
		assertBody func(*testing.T, *http.Request)
	}{
		{
			name:     "users me",
			method:   http.MethodGet,
			path:     "/v1/users/me",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Users.Me(ctx)
				return err
			},
		},
		{
			name:     "users get",
			method:   http.MethodGet,
			path:     "/v1/users/user-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Users.Get(ctx, "user-id")
				return err
			},
		},
		{
			name:   "users list",
			method: http.MethodGet,
			path:   "/v1/users",
			query: url.Values{
				"start_cursor": {"cursor"},
				"page_size":    {"25"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.Users.List(ctx, &PaginationParams{StartCursor: "cursor", PageSize: 25})
				return err
			},
		},
		{
			name:     "pages create",
			method:   http.MethodPost,
			path:     "/v1/pages",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Pages.Create(ctx, Object{"parent": Object{"page_id": "parent-id"}})
				return err
			},
		},
		{
			name:   "pages get",
			method: http.MethodGet,
			path:   "/v1/pages/page-id",
			query: url.Values{
				"filter_properties": {"title", "status"},
			},
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Pages.Get(ctx, "page-id", &GetPageParams{FilterProperties: []string{"title", "status"}})
				return err
			},
		},
		{
			name:     "pages update",
			method:   http.MethodPatch,
			path:     "/v1/pages/page-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Pages.Update(ctx, "page-id", Object{"archived": true})
				return err
			},
		},
		{
			name:     "pages move",
			method:   http.MethodPost,
			path:     "/v1/pages/page-id/move",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Pages.Move(ctx, "page-id", Object{"parent": Object{"page_id": "new-parent"}})
				return err
			},
		},
		{
			name:   "pages property",
			method: http.MethodGet,
			path:   "/v1/pages/page-id/properties/title",
			query: url.Values{
				"start_cursor": {"cursor"},
				"page_size":    {"10"},
			},
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Pages.GetProperty(ctx, "page-id", "title", &PaginationParams{StartCursor: "cursor", PageSize: 10})
				return err
			},
		},
		{
			name:   "pages retrieve markdown",
			method: http.MethodGet,
			path:   "/v1/pages/page-id/markdown",
			query: url.Values{
				"include_transcript": {"true"},
			},
			response: `{"markdown":"# Title"}`,
			invoke: func(c *Client) error {
				include := true
				_, err := c.Pages.RetrieveMarkdown(ctx, "page-id", &RetrieveMarkdownParams{IncludeTranscript: &include})
				return err
			},
		},
		{
			name:     "pages update markdown",
			method:   http.MethodPatch,
			path:     "/v1/pages/page-id/markdown",
			response: `{"markdown":"# Title"}`,
			invoke: func(c *Client) error {
				_, err := c.Pages.UpdateMarkdown(ctx, "page-id", Object{"type": "replace_content"})
				return err
			},
		},
		{
			name:     "blocks get",
			method:   http.MethodGet,
			path:     "/v1/blocks/block-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Blocks.Get(ctx, "block-id")
				return err
			},
		},
		{
			name:     "blocks update",
			method:   http.MethodPatch,
			path:     "/v1/blocks/block-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Blocks.Update(ctx, "block-id", Object{"paragraph": Object{}})
				return err
			},
		},
		{
			name:     "blocks delete",
			method:   http.MethodDelete,
			path:     "/v1/blocks/block-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Blocks.Delete(ctx, "block-id")
				return err
			},
		},
		{
			name:   "blocks children",
			method: http.MethodGet,
			path:   "/v1/blocks/block-id/children",
			query: url.Values{
				"page_size": {"50"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.Blocks.Children(ctx, "block-id", &PaginationParams{PageSize: 50})
				return err
			},
		},
		{
			name:     "blocks append children",
			method:   http.MethodPatch,
			path:     "/v1/blocks/block-id/children",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Blocks.AppendChildren(ctx, "block-id", Object{"children": []Object{}})
				return err
			},
		},
		{
			name:     "data sources get",
			method:   http.MethodGet,
			path:     "/v1/data_sources/data-source-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.DataSources.Get(ctx, "data-source-id")
				return err
			},
		},
		{
			name:     "data sources update",
			method:   http.MethodPatch,
			path:     "/v1/data_sources/data-source-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.DataSources.Update(ctx, "data-source-id", Object{"title": []Object{}})
				return err
			},
		},
		{
			name:   "data sources query",
			method: http.MethodPost,
			path:   "/v1/data_sources/data-source-id/query",
			query: url.Values{
				"filter_properties": {"title"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.DataSources.Query(ctx, "data-source-id", QueryRequest{"page_size": 10}, &QueryDataSourceParams{FilterProperties: []string{"title"}})
				return err
			},
		},
		{
			name:     "data sources create",
			method:   http.MethodPost,
			path:     "/v1/data_sources",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.DataSources.Create(ctx, Object{"parent": Object{}})
				return err
			},
		},
		{
			name:   "data sources templates",
			method: http.MethodGet,
			path:   "/v1/data_sources/data-source-id/templates",
			query: url.Values{
				"name":      {"Default"},
				"page_size": {"10"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.DataSources.ListTemplates(ctx, "data-source-id", &ListDataSourceTemplatesParams{
					Name:             "Default",
					PaginationParams: PaginationParams{PageSize: 10},
				})
				return err
			},
		},
		{
			name:     "databases get",
			method:   http.MethodGet,
			path:     "/v1/databases/database-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Databases.Get(ctx, "database-id")
				return err
			},
		},
		{
			name:     "databases update",
			method:   http.MethodPatch,
			path:     "/v1/databases/database-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Databases.Update(ctx, "database-id", Object{"title": []Object{}})
				return err
			},
		},
		{
			name:     "databases create",
			method:   http.MethodPost,
			path:     "/v1/databases",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Databases.Create(ctx, Object{"parent": Object{}})
				return err
			},
		},
		{
			name:     "search",
			method:   http.MethodPost,
			path:     "/v1/search",
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.Search.Do(ctx, SearchRequest{"query": "roadmap"})
				return err
			},
		},
		{
			name:     "comments create",
			method:   http.MethodPost,
			path:     "/v1/comments",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Comments.Create(ctx, Object{"rich_text": []Object{}})
				return err
			},
		},
		{
			name:   "comments list",
			method: http.MethodGet,
			path:   "/v1/comments",
			query: url.Values{
				"block_id":     {"block-id"},
				"start_cursor": {"cursor"},
				"page_size":    {"10"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.Comments.List(ctx, &ListCommentsParams{
					BlockID:          "block-id",
					PaginationParams: PaginationParams{StartCursor: "cursor", PageSize: 10},
				})
				return err
			},
		},
		{
			name:     "comments get",
			method:   http.MethodGet,
			path:     "/v1/comments/comment-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Comments.Get(ctx, "comment-id")
				return err
			},
		},
		{
			name:     "comments update",
			method:   http.MethodPatch,
			path:     "/v1/comments/comment-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Comments.Update(ctx, "comment-id", Object{"resolved": true})
				return err
			},
		},
		{
			name:     "comments delete",
			method:   http.MethodDelete,
			path:     "/v1/comments/comment-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Comments.Delete(ctx, "comment-id")
				return err
			},
		},
		{
			name:     "file uploads create",
			method:   http.MethodPost,
			path:     "/v1/file_uploads",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.FileUploads.Create(ctx, CreateFileUploadRequest{Mode: "single_part", Filename: "doc.pdf"})
				return err
			},
		},
		{
			name:   "file uploads list",
			method: http.MethodGet,
			path:   "/v1/file_uploads",
			query: url.Values{
				"status":    {"uploaded"},
				"page_size": {"10"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.FileUploads.List(ctx, &ListFileUploadsParams{
					Status:           "uploaded",
					PaginationParams: PaginationParams{PageSize: 10},
				})
				return err
			},
		},
		{
			name:     "file uploads send",
			method:   http.MethodPost,
			path:     "/v1/file_uploads/upload-id/send",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.FileUploads.Send(ctx, "upload-id", UploadFileRequest{
					Filename:   "doc.txt",
					Reader:     strings.NewReader("hello"),
					PartNumber: "2",
				})
				return err
			},
			assertBody: assertMultipartUpload,
		},
		{
			name:     "file uploads complete",
			method:   http.MethodPost,
			path:     "/v1/file_uploads/upload-id/complete",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.FileUploads.Complete(ctx, "upload-id")
				return err
			},
		},
		{
			name:     "file uploads get",
			method:   http.MethodGet,
			path:     "/v1/file_uploads/upload-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.FileUploads.Get(ctx, "upload-id")
				return err
			},
		},
		{
			name:   "custom emojis list",
			method: http.MethodGet,
			path:   "/v1/custom_emojis",
			query: url.Values{
				"name":      {"ship"},
				"page_size": {"5"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.CustomEmojis.List(ctx, &ListCustomEmojisParams{
					Name:             "ship",
					PaginationParams: PaginationParams{PageSize: 5},
				})
				return err
			},
		},
		{
			name:   "views list",
			method: http.MethodGet,
			path:   "/v1/views",
			query: url.Values{
				"database_id":    {"database-id"},
				"data_source_id": {"data-source-id"},
				"page_size":      {"10"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.Views.List(ctx, &ListViewsParams{
					DatabaseID:       "database-id",
					DataSourceID:     "data-source-id",
					PaginationParams: PaginationParams{PageSize: 10},
				})
				return err
			},
		},
		{
			name:     "views create",
			method:   http.MethodPost,
			path:     "/v1/views",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Views.Create(ctx, Object{"type": "table"})
				return err
			},
		},
		{
			name:     "views get",
			method:   http.MethodGet,
			path:     "/v1/views/view-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Views.Get(ctx, "view-id")
				return err
			},
		},
		{
			name:     "views update",
			method:   http.MethodPatch,
			path:     "/v1/views/view-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Views.Update(ctx, "view-id", Object{"name": "Updated"})
				return err
			},
		},
		{
			name:     "views delete",
			method:   http.MethodDelete,
			path:     "/v1/views/view-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Views.Delete(ctx, "view-id")
				return err
			},
		},
		{
			name:     "views create query",
			method:   http.MethodPost,
			path:     "/v1/views/view-id/queries",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Views.CreateQuery(ctx, "view-id", Object{})
				return err
			},
		},
		{
			name:   "views query results",
			method: http.MethodGet,
			path:   "/v1/views/view-id/queries/query-id",
			query: url.Values{
				"page_size": {"10"},
			},
			response: listResponse,
			invoke: func(c *Client) error {
				_, err := c.Views.GetQueryResults(ctx, "view-id", "query-id", &PaginationParams{PageSize: 10})
				return err
			},
		},
		{
			name:     "views delete query",
			method:   http.MethodDelete,
			path:     "/v1/views/view-id/queries/query-id",
			response: objectResponse,
			invoke: func(c *Client) error {
				_, err := c.Views.DeleteQuery(ctx, "view-id", "query-id")
				return err
			},
		},
		{
			name:     "meeting notes query",
			method:   http.MethodPost,
			path:     "/v1/blocks/meeting_notes/query",
			response: `{"results":[],"has_more":false}`,
			invoke: func(c *Client) error {
				_, err := c.MeetingNotes.Query(ctx, Object{"limit": 10})
				return err
			},
		},
		{
			name:     "oauth token",
			method:   http.MethodPost,
			path:     "/v1/oauth/token",
			response: oauthTokenResponse,
			invoke: func(c *Client) error {
				_, err := c.OAuth.ExchangeAuthorizationCode(ctx, "client-id", "client-secret", "code", "https://example.test/callback")
				return err
			},
			assertBody: assertOAuthBasicAuth,
		},
		{
			name:     "oauth revoke",
			method:   http.MethodPost,
			path:     "/v1/oauth/revoke",
			response: `{"request_id":"request-id"}`,
			invoke: func(c *Client) error {
				_, err := c.OAuth.Revoke(ctx, "client-id", "client-secret", "token")
				return err
			},
			assertBody: assertOAuthBasicAuth,
		},
		{
			name:     "oauth introspect",
			method:   http.MethodPost,
			path:     "/v1/oauth/introspect",
			response: `{"active":true,"scope":"read_content","iat":1}`,
			invoke: func(c *Client) error {
				_, err := c.OAuth.Introspect(ctx, "client-id", "client-secret", "token")
				return err
			},
			assertBody: assertOAuthBasicAuth,
		},
	}

	for _, tc := range calls {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assertRequest(t, r, tc.method, tc.path, tc.query)
				if tc.assertBody != nil {
					tc.assertBody(t, r)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tc.response))
			}))
			t.Cleanup(server.Close)

			client := MustNewClient("secret", WithBaseURL(server.URL), WithUserAgent("test-agent"))
			if err := tc.invoke(client); err != nil {
				t.Fatalf("invoke returned error: %v", err)
			}
		})
	}
}

func TestAPIErrorDecoding(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "request-id")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{
			"object":"error",
			"status":429,
			"code":"rate_limited",
			"message":"slow down",
			"additional_data":{"retry_after":["1"]}
		}`))
	}))
	t.Cleanup(server.Close)

	client := MustNewClient("secret", WithBaseURL(server.URL))
	_, err := client.Users.Me(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", apiErr.StatusCode, http.StatusTooManyRequests)
	}
	if apiErr.Code != "rate_limited" || apiErr.Message != "slow down" {
		t.Fatalf("decoded error = %#v", apiErr)
	}
	if apiErr.RequestID != "request-id" {
		t.Fatalf("request id = %q, want request-id", apiErr.RequestID)
	}
}

func TestListResponsePreservesRawResults(t *testing.T) {
	t.Parallel()

	var list ListResponse
	if err := json.Unmarshal([]byte(`{
		"object":"list",
		"results":[{"id":"a"},{"id":"b"}],
		"has_more":true,
		"next_cursor":"next",
		"type":"page"
	}`), &list); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}

	if len(list.Results) != 2 {
		t.Fatalf("results length = %d, want 2", len(list.Results))
	}
	if list.Raw["type"] != "page" {
		t.Fatalf("raw type = %v, want page", list.Raw["type"])
	}
	if got := string(list.RawResults()); !strings.Contains(got, `"id":"a"`) {
		t.Fatalf("raw results = %s", got)
	}
}

func TestAPIPathEscapesSegments(t *testing.T) {
	t.Parallel()

	got := apiPath("v1", "views", "id/with/slash", "queries", "query id")
	want := "/v1/views/id%2Fwith%2Fslash/queries/query%20id"
	if got != want {
		t.Fatalf("apiPath = %q, want %q", got, want)
	}

	got = apiPath("v1", "pages", "")
	want = "/v1/pages/"
	if got != want {
		t.Fatalf("apiPath empty segment = %q, want %q", got, want)
	}
}

func TestCollectPaginated(t *testing.T) {
	t.Parallel()

	next := "next"
	calls := 0
	results, err := CollectPaginated(context.Background(), func(ctx context.Context, cursor string) (*ListResponse, error) {
		calls++
		switch calls {
		case 1:
			if cursor != "" {
				t.Fatalf("first cursor = %q, want empty", cursor)
			}
			return &ListResponse{
				Results:    []Object{{"id": "a"}},
				HasMore:    true,
				NextCursor: &next,
			}, nil
		case 2:
			if cursor != "next" {
				t.Fatalf("second cursor = %q, want next", cursor)
			}
			return &ListResponse{
				Results: []Object{{"id": "b"}},
				HasMore: false,
			}, nil
		default:
			t.Fatalf("unexpected call %d", calls)
			return nil, nil
		}
	})
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("results length = %d, want 2", len(results))
	}
}

func assertRequest(t *testing.T, r *http.Request, method, path string, query url.Values) {
	t.Helper()

	if r.Method != method {
		t.Fatalf("method = %s, want %s", r.Method, method)
	}
	if r.URL.Path != path {
		t.Fatalf("path = %s, want %s", r.URL.Path, path)
	}
	if got := r.URL.Query(); !equalQuery(got, query) {
		t.Fatalf("query = %v, want %v", got, query)
	}
	if r.Header.Get("Notion-Version") != LatestVersion {
		t.Fatalf("Notion-Version = %q, want %q", r.Header.Get("Notion-Version"), LatestVersion)
	}
	if r.Header.Get("User-Agent") != "test-agent" {
		t.Fatalf("User-Agent = %q, want test-agent", r.Header.Get("User-Agent"))
	}
}

func assertMultipartUpload(t *testing.T, r *http.Request) {
	t.Helper()

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		t.Fatalf("parse content type: %v", err)
	}
	if mediaType != "multipart/form-data" {
		t.Fatalf("media type = %s, want multipart/form-data", mediaType)
	}
	reader := multipart.NewReader(r.Body, params["boundary"])
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	t.Cleanup(func() { _ = form.RemoveAll() })

	if got := form.Value["part_number"]; len(got) != 1 || got[0] != "2" {
		t.Fatalf("part_number = %v, want [2]", got)
	}
	files := form.File["file"]
	if len(files) != 1 {
		t.Fatalf("file count = %d, want 1", len(files))
	}
	if files[0].Filename != "doc.txt" {
		t.Fatalf("filename = %s, want doc.txt", files[0].Filename)
	}
}

func assertOAuthBasicAuth(t *testing.T, r *http.Request) {
	t.Helper()

	want := "Basic " + base64.StdEncoding.EncodeToString([]byte("client-id:client-secret"))
	if got := r.Header.Get("Authorization"); got != want {
		t.Fatalf("Authorization = %q, want %q", got, want)
	}
}

func equalQuery(got, want url.Values) bool {
	if len(got) != len(want) {
		return false
	}
	for key, wantValues := range want {
		gotValues, ok := got[key]
		if !ok || len(gotValues) != len(wantValues) {
			return false
		}
		for i := range wantValues {
			if gotValues[i] != wantValues[i] {
				return false
			}
		}
	}
	return true
}

const objectResponse = `{"object":"test","id":"id"}`

const listResponse = `{
	"object":"list",
	"results":[{"object":"page","id":"page-id"}],
	"has_more":false,
	"next_cursor":null,
	"type":"page"
}`

const oauthTokenResponse = `{
	"access_token":"access",
	"token_type":"bearer",
	"refresh_token":"refresh",
	"bot_id":"00000000-0000-0000-0000-000000000000",
	"workspace_icon":null,
	"workspace_name":"Workspace",
	"workspace_id":"11111111-1111-1111-1111-111111111111",
	"owner":{"type":"workspace","workspace":true},
	"duplicated_template_id":null,
	"request_id":"request-id"
}`
