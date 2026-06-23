package notion

// Parent is a typed Notion parent request object.
type Parent struct {
	Type         string `json:"type,omitempty"`
	PageID       string `json:"page_id,omitempty"`
	DatabaseID   string `json:"database_id,omitempty"`
	DataSourceID string `json:"data_source_id,omitempty"`
	Workspace    *bool  `json:"workspace,omitempty"`
}

// PageParent creates a page parent reference.
func PageParent(pageID string) *Parent {
	return &Parent{Type: "page_id", PageID: pageID}
}

// DatabaseParent creates a database parent reference.
func DatabaseParent(databaseID string) *Parent {
	return &Parent{Type: "database_id", DatabaseID: databaseID}
}

// DataSourceParent creates a data source parent reference.
func DataSourceParent(dataSourceID string) *Parent {
	return &Parent{Type: "data_source_id", DataSourceID: dataSourceID}
}

// WorkspaceParent creates a workspace parent reference.
func WorkspaceParent() *Parent {
	workspace := true
	return &Parent{Type: "workspace", Workspace: &workspace}
}
