package notion

// CreatePageRequest is a typed request body for creating a page.
type CreatePageRequest struct {
	Parent     *Parent        `json:"parent,omitempty"`
	Properties PageProperties `json:"properties,omitempty"`
	Icon       any            `json:"icon,omitempty"`
	Cover      any            `json:"cover,omitempty"`
	Content    []BlockRequest `json:"content,omitempty"`
	Children   []BlockRequest `json:"children,omitempty"`
	Markdown   string         `json:"markdown,omitempty"`
	Template   any            `json:"template,omitempty"`
	Position   Object         `json:"position,omitempty"`
}

// UpdatePageRequest is a typed request body for updating a page.
type UpdatePageRequest struct {
	Properties PageProperties `json:"properties,omitempty"`
	Icon       any            `json:"icon,omitempty"`
	Cover      any            `json:"cover,omitempty"`
	Archived   *bool          `json:"archived,omitempty"`
	InTrash    *bool          `json:"in_trash,omitempty"`
}

// ArchivePageRequest returns an update request that archives or restores a page.
func ArchivePageRequest(archived bool) UpdatePageRequest {
	return UpdatePageRequest{Archived: &archived}
}
