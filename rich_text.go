package notion

// RichText is a Notion rich text request item.
type RichText struct {
	Type        string       `json:"type,omitempty"`
	Text        *TextContent `json:"text,omitempty"`
	Mention     Object       `json:"mention,omitempty"`
	Equation    *Equation    `json:"equation,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

// TextContent is the text payload inside a rich text item.
type TextContent struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

// Link is an inline rich text link.
type Link struct {
	URL string `json:"url"`
}

// Equation is an inline rich text equation.
type Equation struct {
	Expression string `json:"expression"`
}

// Annotations configures rich text styling.
type Annotations struct {
	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
	Underline     bool   `json:"underline,omitempty"`
	Code          bool   `json:"code,omitempty"`
	Color         string `json:"color,omitempty"`
}

// Text creates a text rich text item.
func Text(content string) RichText {
	return RichText{
		Type: "text",
		Text: &TextContent{Content: content},
	}
}

// TextLink creates a linked text rich text item.
func TextLink(content, url string) RichText {
	return RichText{
		Type: "text",
		Text: &TextContent{
			Content: content,
			Link:    &Link{URL: url},
		},
	}
}

// EquationText creates an equation rich text item.
func EquationText(expression string) RichText {
	return RichText{
		Type:     "equation",
		Equation: &Equation{Expression: expression},
	}
}

// WithAnnotations returns a copy of rt with annotations attached.
func (rt RichText) WithAnnotations(annotations Annotations) RichText {
	rt.Annotations = &annotations
	return rt
}
