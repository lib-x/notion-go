# notion-go

`notion-go` is a Go SDK for the Notion API. It targets the current Notion API version `2026-03-11` and covers every endpoint in the official OpenAPI document.

Official update sources:

- API reference: https://developers.notion.com/reference/intro
- OpenAPI document: https://developers.notion.com/openapi.json
- LLM documentation index: https://developers.notion.com/llms.txt

## Install

```sh
go get github.com/lib-x/notion-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	notion "github.com/lib-x/notion-go"
)

func main() {
	client := notion.MustNewClient(os.Getenv("NOTION_API_KEY"))

	page, err := client.Pages.Get(context.Background(), "page-id", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(page["id"])
}
```

## API Design

The SDK exposes stable service methods for every Notion REST endpoint. Common request structures such as parents, rich text, page properties, blocks, data source filters, and sorts have typed builders. The lower-level `notion.Object` (`map[string]any`) remains available as an escape hatch for newly released Notion fields or less common object shapes.

List responses can also be decoded into caller-owned generic result types with `DecodeList[T]`, `CollectPaginatedDecode[T]`, or `CollectPaginatedAs[T]`.

Use `client.Do` when Notion adds a new endpoint before this SDK has a convenience method.

```go
var out notion.Object
err := client.Do(ctx, "POST", "/v1/new_endpoint", &notion.RequestOptions{
	Body: notion.Object{"key": "value"},
}, &out)
```

## Endpoint Coverage

- Users: retrieve bot user, retrieve user, list users
- Pages: create, retrieve, update, move, retrieve property, retrieve markdown, update markdown
- Blocks: retrieve, update, delete, list children, append children
- Data sources: create, retrieve, update, query, list templates
- Databases: create, retrieve, update
- Search: search by title
- Comments: create, list, retrieve, update, delete
- File uploads: create, list, send, complete multipart upload, retrieve
- Custom emojis: list
- Views: list, create, retrieve, update, delete, create query, get query results, delete query
- Meeting notes: query
- OAuth: token exchange, refresh, revoke, introspect

## Pagination

List endpoints return `*notion.ListResponse`. Use `notion.CollectPaginated` for endpoints that accept `start_cursor`.

```go
allUsers, err := notion.CollectPaginated(ctx, func(ctx context.Context, cursor string) (*notion.ListResponse, error) {
	return client.Users.List(ctx, &notion.PaginationParams{
		StartCursor: cursor,
		PageSize:    100,
	})
})
```

For endpoints with extra filters, keep those filters in the closure and only vary `StartCursor`.

## Create a Page

```go
page, err := client.Pages.Create(ctx, notion.CreatePageRequest{
	Parent: notion.PageParent("parent-page-id"),
	Properties: notion.PageProperties{
		"Name":   notion.TitleProperty(notion.Text("Roadmap")),
		"Status": notion.StatusPropertyName("In progress"),
		"Due":    notion.DateProperty(notion.Date{Start: "2026-06-23"}),
	},
	Children: []notion.BlockRequest{
		notion.Heading1(notion.Text("Plan")),
		notion.Paragraph(notion.Text("Hello from Go")),
		notion.ToDo(false, notion.Text("Ship typed helpers")),
	},
})
```

## Query a Data Source

```go
results, err := client.DataSources.Query(ctx, "data-source-id", notion.DataSourceQueryRequest{
	Filter: notion.And(
		notion.TitleContains("Name", "roadmap"),
		notion.StatusEquals("Status", "In progress"),
	),
	Sorts: []notion.Sort{
		notion.PropertySort("Due", notion.Descending),
	},
	PageSize: 25,
}, nil)
```

## Typed List Results

```go
type PageSummary struct {
	ID     string `json:"id"`
	Object string `json:"object"`
}

list, err := client.DataSources.Query(ctx, "data-source-id", notion.DataSourceQueryRequest{
	PageSize: 50,
}, nil)
if err != nil {
	return err
}

typed, err := notion.DecodeList[PageSummary](list)
if err != nil {
	return err
}

for _, page := range typed.Results {
	fmt.Println(page.ID)
}
```

For paginated endpoints:

```go
pages, err := notion.CollectPaginatedDecode[PageSummary](ctx, func(ctx context.Context, cursor string) (*notion.ListResponse, error) {
	return client.DataSources.Query(ctx, "data-source-id", notion.DataSourceQueryRequest{
		StartCursor: cursor,
		PageSize:    100,
	}, nil)
})
```

The typed builders intentionally do not try to cover every Notion union branch on day one. Use `notion.RawProperty`, `notion.RawBlock`, `notion.RawFilter`, or plain `notion.Object` request bodies when you need an API shape that is not modeled yet.

## Markdown Content

```go
md, err := client.Pages.RetrieveMarkdown(ctx, "page-id", nil)
if err == nil {
	fmt.Println(md.Markdown)
}

updated, err := client.Pages.UpdateMarkdown(ctx, "page-id", notion.Object{
	"type": "update_content",
	"update_content": notion.Object{
		"content_updates": []notion.Object{
			{"old_str": "old text", "new_str": "new text"},
		},
	},
})
_ = updated
```

## File Uploads

```go
upload, err := client.FileUploads.Create(ctx, notion.CreateFileUploadRequest{
	Mode:        "single_part",
	Filename:    "document.txt",
	ContentType: "text/plain",
})
if err != nil {
	return err
}

fileID, _ := upload["id"].(string)
_, err = client.FileUploads.Send(ctx, fileID, notion.UploadFileRequest{
	Filename: "document.txt",
	Reader:   strings.NewReader("hello"),
})
```

## OAuth

```go
token, err := client.OAuth.ExchangeAuthorizationCode(
	ctx,
	os.Getenv("NOTION_OAUTH_CLIENT_ID"),
	os.Getenv("NOTION_OAUTH_CLIENT_SECRET"),
	authorizationCode,
	"https://example.com/callback",
)
```

OAuth methods use HTTP Basic authentication as required by Notion. The client token passed to `NewClient` may be empty if you only call OAuth methods.

## Error Handling

Non-2xx responses return `*notion.APIError`.

```go
_, err := client.Users.Me(ctx)
if err != nil {
	var apiErr *notion.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode, apiErr.Code, apiErr.Message, apiErr.RequestID)
	}
}
```

## Configuration

```go
client, err := notion.NewClient(
	os.Getenv("NOTION_API_KEY"),
	notion.WithNotionVersion(notion.LatestVersion),
	notion.WithUserAgent("my-app/1.0"),
)
```

`WithBaseURL` is available for tests or proxies.

## Updating From Notion OpenAPI

When Notion publishes API changes:

```sh
curl -fsSL https://developers.notion.com/openapi.json -o .firecrawl/notion-openapi.json
jq -r '.components.parameters.notionVersion.schema.enum[]' .firecrawl/notion-openapi.json
jq -r '.paths | to_entries[] | .key as $path | .value | to_entries[] | "\(.value.operationId)\t\(.key | ascii_upcase)\t\($path)"' .firecrawl/notion-openapi.json
```

Then update `LatestVersion`, add or adjust service methods, and run:

```sh
go test ./...
```
