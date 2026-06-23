package notion

// BlockRequest is implemented by typed Notion block request objects.
type BlockRequest interface {
	blockRequest()
}

// RawBlockRequest preserves an arbitrary block request object inside typed
// page and append-children requests.
type RawBlockRequest Object

func (RawBlockRequest) blockRequest() {}

// RawBlock wraps an arbitrary block request object.
func RawBlock(value Object) BlockRequest {
	return RawBlockRequest(value)
}

type RichTextBlockContent struct {
	RichText []RichText     `json:"rich_text"`
	Color    string         `json:"color,omitempty"`
	Children []BlockRequest `json:"children,omitempty"`
}

type ToDoBlockContent struct {
	RichText []RichText     `json:"rich_text"`
	Checked  bool           `json:"checked"`
	Color    string         `json:"color,omitempty"`
	Children []BlockRequest `json:"children,omitempty"`
}

type CodeBlockContent struct {
	RichText []RichText `json:"rich_text"`
	Language string     `json:"language"`
	Caption  []RichText `json:"caption,omitempty"`
}

type ChildPageBlockContent struct {
	Title string `json:"title"`
}

type EmptyObject struct{}

type blockRequest struct {
	Object           string                 `json:"object,omitempty"`
	Type             string                 `json:"type"`
	Paragraph        *RichTextBlockContent  `json:"paragraph,omitempty"`
	Heading1         *RichTextBlockContent  `json:"heading_1,omitempty"`
	Heading2         *RichTextBlockContent  `json:"heading_2,omitempty"`
	Heading3         *RichTextBlockContent  `json:"heading_3,omitempty"`
	BulletedListItem *RichTextBlockContent  `json:"bulleted_list_item,omitempty"`
	NumberedListItem *RichTextBlockContent  `json:"numbered_list_item,omitempty"`
	Quote            *RichTextBlockContent  `json:"quote,omitempty"`
	Toggle           *RichTextBlockContent  `json:"toggle,omitempty"`
	Callout          *RichTextBlockContent  `json:"callout,omitempty"`
	ToDo             *ToDoBlockContent      `json:"to_do,omitempty"`
	Code             *CodeBlockContent      `json:"code,omitempty"`
	ChildPage        *ChildPageBlockContent `json:"child_page,omitempty"`
	Divider          *EmptyObject           `json:"divider,omitempty"`
}

func (blockRequest) blockRequest() {}

// Paragraph creates a paragraph block.
func Paragraph(text ...RichText) BlockRequest {
	return blockWithRichText("paragraph", text)
}

// Heading1 creates a heading_1 block.
func Heading1(text ...RichText) BlockRequest {
	return blockWithRichText("heading_1", text)
}

// Heading2 creates a heading_2 block.
func Heading2(text ...RichText) BlockRequest {
	return blockWithRichText("heading_2", text)
}

// Heading3 creates a heading_3 block.
func Heading3(text ...RichText) BlockRequest {
	return blockWithRichText("heading_3", text)
}

// BulletedListItem creates a bulleted_list_item block.
func BulletedListItem(text ...RichText) BlockRequest {
	return blockWithRichText("bulleted_list_item", text)
}

// NumberedListItem creates a numbered_list_item block.
func NumberedListItem(text ...RichText) BlockRequest {
	return blockWithRichText("numbered_list_item", text)
}

// Quote creates a quote block.
func Quote(text ...RichText) BlockRequest {
	return blockWithRichText("quote", text)
}

// Toggle creates a toggle block.
func Toggle(text ...RichText) BlockRequest {
	return blockWithRichText("toggle", text)
}

// Callout creates a callout block.
func Callout(text ...RichText) BlockRequest {
	return blockWithRichText("callout", text)
}

// ToDo creates a to_do block.
func ToDo(checked bool, text ...RichText) BlockRequest {
	return blockRequest{
		Object: "block",
		Type:   "to_do",
		ToDo: &ToDoBlockContent{
			RichText: text,
			Checked:  checked,
		},
	}
}

// CodeBlock creates a code block.
func CodeBlock(language string, text ...RichText) BlockRequest {
	return blockRequest{
		Object: "block",
		Type:   "code",
		Code: &CodeBlockContent{
			RichText: text,
			Language: language,
		},
	}
}

// ChildPage creates a child_page block.
func ChildPage(title string) BlockRequest {
	return blockRequest{
		Object:    "block",
		Type:      "child_page",
		ChildPage: &ChildPageBlockContent{Title: title},
	}
}

// Divider creates a divider block.
func Divider() BlockRequest {
	return blockRequest{
		Object:  "block",
		Type:    "divider",
		Divider: &EmptyObject{},
	}
}

func blockWithRichText(blockType string, text []RichText) BlockRequest {
	content := &RichTextBlockContent{RichText: text}
	req := blockRequest{Object: "block", Type: blockType}
	switch blockType {
	case "paragraph":
		req.Paragraph = content
	case "heading_1":
		req.Heading1 = content
	case "heading_2":
		req.Heading2 = content
	case "heading_3":
		req.Heading3 = content
	case "bulleted_list_item":
		req.BulletedListItem = content
	case "numbered_list_item":
		req.NumberedListItem = content
	case "quote":
		req.Quote = content
	case "toggle":
		req.Toggle = content
	case "callout":
		req.Callout = content
	}
	return req
}
