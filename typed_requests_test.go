package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePageRequestMarshalsTypedContent(t *testing.T) {
	t.Parallel()

	done := false
	body := CreatePageRequest{
		Parent: PageParent("parent-page-id"),
		Properties: PageProperties{
			"Name":   TitleProperty(Text("Roadmap")),
			"Status": StatusPropertyName("In progress"),
			"Due":    DateProperty(Date{Start: "2026-06-23"}),
			"Done":   CheckboxProperty(done),
		},
		Children: []BlockRequest{
			Paragraph(Text("Hello from Go")),
			Heading1(Text("Plan")),
			ToDo(false, Text("Ship typed helpers")),
			Divider(),
		},
	}

	assertJSONEqual(t, mustJSON(t, body), `{
		"parent":{"type":"page_id","page_id":"parent-page-id"},
		"properties":{
			"Name":{"type":"title","title":[{"type":"text","text":{"content":"Roadmap"}}]},
			"Status":{"type":"status","status":{"name":"In progress"}},
			"Due":{"type":"date","date":{"start":"2026-06-23"}},
			"Done":{"type":"checkbox","checkbox":false}
		},
		"children":[
			{"object":"block","type":"paragraph","paragraph":{"rich_text":[{"type":"text","text":{"content":"Hello from Go"}}]}},
			{"object":"block","type":"heading_1","heading_1":{"rich_text":[{"type":"text","text":{"content":"Plan"}}]}},
			{"object":"block","type":"to_do","to_do":{"rich_text":[{"type":"text","text":{"content":"Ship typed helpers"}}],"checked":false}},
			{"object":"block","type":"divider","divider":{}}
		]
	}`)
}

func TestCreatePageRequestOmitsMissingParent(t *testing.T) {
	t.Parallel()

	body := CreatePageRequest{
		Properties: PageProperties{"Name": TitleProperty(Text("Roadmap"))},
	}

	assertJSONEqual(t, mustJSON(t, body), `{
		"properties":{"Name":{"type":"title","title":[{"type":"text","text":{"content":"Roadmap"}}]}}
	}`)
}

func TestDataSourceQueryRequestMarshalsTypedFilterAndSort(t *testing.T) {
	t.Parallel()

	inTrash := false
	body := DataSourceQueryRequest{
		Filter: And(
			TitleContains("Name", "roadmap"),
			StatusEquals("Status", "In progress"),
			CheckboxEquals("Done", false),
		),
		Sorts: []Sort{
			PropertySort("Due", Descending),
			TimestampSort("created_time", Ascending),
		},
		PageSize:   25,
		ResultType: ResultTypePage,
		InTrash:    &inTrash,
	}

	assertJSONEqual(t, mustJSON(t, body), `{
		"sorts":[
			{"property":"Due","direction":"descending"},
			{"timestamp":"created_time","direction":"ascending"}
		],
		"filter":{
			"and":[
				{"property":"Name","title":{"contains":"roadmap"}},
				{"property":"Status","status":{"equals":"In progress"}},
				{"property":"Done","checkbox":{"equals":false}}
			]
		},
		"page_size":25,
		"in_trash":false,
		"result_type":"page"
	}`)
}

func TestServicesAcceptTypedRequests(t *testing.T) {
	t.Parallel()

	var seenPageBody []byte
	var seenQueryBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/pages":
			seenPageBody = readRequestBody(t, r)
			_, _ = w.Write([]byte(objectResponse))
		case "/v1/data_sources/data-source-id/query":
			seenQueryBody = readRequestBody(t, r)
			_, _ = w.Write([]byte(listResponse))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := MustNewClient("secret", WithBaseURL(server.URL))
	_, err := client.Pages.Create(context.Background(), CreatePageRequest{
		Parent:     PageParent("parent-page-id"),
		Properties: PageProperties{"Name": TitleProperty(Text("Roadmap"))},
	})
	if err != nil {
		t.Fatalf("create page: %v", err)
	}
	_, err = client.DataSources.Query(context.Background(), "data-source-id", DataSourceQueryRequest{
		Filter:   TitleContains("Name", "roadmap"),
		PageSize: 10,
	}, nil)
	if err != nil {
		t.Fatalf("query data source: %v", err)
	}

	assertJSONEqual(t, seenPageBody, `{
		"parent":{"type":"page_id","page_id":"parent-page-id"},
		"properties":{"Name":{"type":"title","title":[{"type":"text","text":{"content":"Roadmap"}}]}}
	}`)
	assertJSONEqual(t, seenQueryBody, `{
		"filter":{"property":"Name","title":{"contains":"roadmap"}},
		"page_size":10
	}`)
}

func TestDecodeListWithGenericResultType(t *testing.T) {
	t.Parallel()

	type typedPage struct {
		ID     string `json:"id"`
		Object string `json:"object"`
	}

	var list ListResponse
	if err := json.Unmarshal([]byte(`{
		"object":"list",
		"results":[{"object":"page","id":"page-id"}],
		"has_more":false,
		"next_cursor":null,
		"type":"page"
	}`), &list); err != nil {
		t.Fatalf("unmarshal list: %v", err)
	}

	typed, err := DecodeList[typedPage](&list)
	if err != nil {
		t.Fatalf("decode typed list: %v", err)
	}
	if len(typed.Results) != 1 || typed.Results[0].ID != "page-id" {
		t.Fatalf("typed results = %#v", typed.Results)
	}
}

func TestCollectPaginatedDecodeWithGenericResultType(t *testing.T) {
	t.Parallel()

	type typedPage struct {
		ID string `json:"id"`
	}

	next := "next"
	calls := 0
	results, err := CollectPaginatedDecode[typedPage](context.Background(), func(ctx context.Context, cursor string) (*ListResponse, error) {
		calls++
		switch cursor {
		case "":
			return &ListResponse{
				Results:    []Object{{"id": "a"}},
				HasMore:    true,
				NextCursor: &next,
			}, nil
		case "next":
			return &ListResponse{
				Results: []Object{{"id": "b"}},
			}, nil
		default:
			t.Fatalf("unexpected cursor %q", cursor)
			return nil, nil
		}
	})
	if err != nil {
		t.Fatalf("collect typed pages: %v", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
	if len(results) != 2 || results[0].ID != "a" || results[1].ID != "b" {
		t.Fatalf("results = %#v", results)
	}
}

func TestTypedRequestsKeepRawEscapeHatches(t *testing.T) {
	t.Parallel()

	body := CreatePageRequest{
		Parent: PageParent("parent-page-id"),
		Properties: PageProperties{
			"Custom": RawProperty(Object{"custom_type": Object{"value": "x"}}),
		},
		Children: []BlockRequest{
			RawBlock(Object{
				"object": "block",
				"type":   "unsupported_future_block",
				"unsupported_future_block": Object{
					"value": "x",
				},
			}),
		},
	}

	assertJSONEqual(t, mustJSON(t, body), `{
		"parent":{"type":"page_id","page_id":"parent-page-id"},
		"properties":{"Custom":{"custom_type":{"value":"x"}}},
		"children":[{"object":"block","type":"unsupported_future_block","unsupported_future_block":{"value":"x"}}]
	}`)
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	return data
}

func readRequestBody(t *testing.T, r *http.Request) []byte {
	t.Helper()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		t.Fatalf("read body: %v", err)
	}
	return buf.Bytes()
}

func assertJSONEqual(t *testing.T, got []byte, want string) {
	t.Helper()
	var gotValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("unmarshal got json %s: %v", string(got), err)
	}
	var wantValue any
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("unmarshal want json: %v", err)
	}
	if !jsonDeepEqual(gotValue, wantValue) {
		gotPretty, _ := json.MarshalIndent(gotValue, "", "  ")
		wantPretty, _ := json.MarshalIndent(wantValue, "", "  ")
		t.Fatalf("json mismatch\ngot:\n%s\nwant:\n%s", gotPretty, wantPretty)
	}
}

func jsonDeepEqual(a, b any) bool {
	return jsonEqualValue(a, b)
}

func jsonEqualValue(a, b any) bool {
	switch av := a.(type) {
	case map[string]any:
		bv, ok := b.(map[string]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		for key, avv := range av {
			if !jsonEqualValue(avv, bv[key]) {
				return false
			}
		}
		return true
	case []any:
		bv, ok := b.([]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !jsonEqualValue(av[i], bv[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
