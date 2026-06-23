package notion

// PropertyValue is implemented by typed page property request values.
type PropertyValue interface {
	propertyValue()
}

// PageProperties maps Notion property names or IDs to typed property values.
type PageProperties map[string]PropertyValue

// RawPropertyValue preserves an arbitrary property value object inside typed
// page requests.
type RawPropertyValue Object

func (RawPropertyValue) propertyValue() {}

// RawProperty wraps an arbitrary property value object.
func RawProperty(value Object) PropertyValue {
	return RawPropertyValue(value)
}

type TitlePropertyValue struct {
	Type  string     `json:"type,omitempty"`
	Title []RichText `json:"title"`
}

func (TitlePropertyValue) propertyValue() {}

type RichTextPropertyValue struct {
	Type     string     `json:"type,omitempty"`
	RichText []RichText `json:"rich_text"`
}

func (RichTextPropertyValue) propertyValue() {}

type NumberPropertyValue struct {
	Type   string   `json:"type,omitempty"`
	Number *float64 `json:"number"`
}

func (NumberPropertyValue) propertyValue() {}

type CheckboxPropertyValue struct {
	Type     string `json:"type,omitempty"`
	Checkbox bool   `json:"checkbox"`
}

func (CheckboxPropertyValue) propertyValue() {}

type SelectOption struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SelectPropertyValue struct {
	Type   string        `json:"type,omitempty"`
	Select *SelectOption `json:"select"`
}

func (SelectPropertyValue) propertyValue() {}

type MultiSelectPropertyValue struct {
	Type        string         `json:"type,omitempty"`
	MultiSelect []SelectOption `json:"multi_select"`
}

func (MultiSelectPropertyValue) propertyValue() {}

type StatusPropertyValue struct {
	Type   string        `json:"type,omitempty"`
	Status *SelectOption `json:"status"`
}

func (StatusPropertyValue) propertyValue() {}

type Date struct {
	Start    string `json:"start"`
	End      string `json:"end,omitempty"`
	TimeZone string `json:"time_zone,omitempty"`
}

type DatePropertyValue struct {
	Type string `json:"type,omitempty"`
	Date *Date  `json:"date"`
}

func (DatePropertyValue) propertyValue() {}

type URLPropertyValue struct {
	Type string  `json:"type,omitempty"`
	URL  *string `json:"url"`
}

func (URLPropertyValue) propertyValue() {}

type EmailPropertyValue struct {
	Type  string  `json:"type,omitempty"`
	Email *string `json:"email"`
}

func (EmailPropertyValue) propertyValue() {}

type PhoneNumberPropertyValue struct {
	Type        string  `json:"type,omitempty"`
	PhoneNumber *string `json:"phone_number"`
}

func (PhoneNumberPropertyValue) propertyValue() {}

// TitleProperty creates a title property value.
func TitleProperty(text ...RichText) PropertyValue {
	return TitlePropertyValue{Type: "title", Title: text}
}

// RichTextProperty creates a rich_text property value.
func RichTextProperty(text ...RichText) PropertyValue {
	return RichTextPropertyValue{Type: "rich_text", RichText: text}
}

// NumberProperty creates a number property value.
func NumberProperty(number float64) PropertyValue {
	return NumberPropertyValue{Type: "number", Number: &number}
}

// NullNumberProperty clears a number property value.
func NullNumberProperty() PropertyValue {
	return NumberPropertyValue{Type: "number"}
}

// CheckboxProperty creates a checkbox property value.
func CheckboxProperty(checked bool) PropertyValue {
	return CheckboxPropertyValue{Type: "checkbox", Checkbox: checked}
}

// SelectPropertyName creates a select property value by option name.
func SelectPropertyName(name string) PropertyValue {
	return SelectPropertyValue{Type: "select", Select: &SelectOption{Name: name}}
}

// SelectPropertyID creates a select property value by option ID.
func SelectPropertyID(id string) PropertyValue {
	return SelectPropertyValue{Type: "select", Select: &SelectOption{ID: id}}
}

// NullSelectProperty clears a select property value.
func NullSelectProperty() PropertyValue {
	return SelectPropertyValue{Type: "select"}
}

// MultiSelectProperty creates a multi_select property value by option names.
func MultiSelectProperty(names ...string) PropertyValue {
	options := make([]SelectOption, 0, len(names))
	for _, name := range names {
		options = append(options, SelectOption{Name: name})
	}
	return MultiSelectPropertyValue{Type: "multi_select", MultiSelect: options}
}

// StatusPropertyName creates a status property value by option name.
func StatusPropertyName(name string) PropertyValue {
	return StatusPropertyValue{Type: "status", Status: &SelectOption{Name: name}}
}

// StatusPropertyID creates a status property value by option ID.
func StatusPropertyID(id string) PropertyValue {
	return StatusPropertyValue{Type: "status", Status: &SelectOption{ID: id}}
}

// NullStatusProperty clears a status property value.
func NullStatusProperty() PropertyValue {
	return StatusPropertyValue{Type: "status"}
}

// DateProperty creates a date property value.
func DateProperty(date Date) PropertyValue {
	return DatePropertyValue{Type: "date", Date: &date}
}

// NullDateProperty clears a date property value.
func NullDateProperty() PropertyValue {
	return DatePropertyValue{Type: "date"}
}

// URLProperty creates a url property value.
func URLProperty(url string) PropertyValue {
	return URLPropertyValue{Type: "url", URL: &url}
}

// NullURLProperty clears a url property value.
func NullURLProperty() PropertyValue {
	return URLPropertyValue{Type: "url"}
}

// EmailProperty creates an email property value.
func EmailProperty(email string) PropertyValue {
	return EmailPropertyValue{Type: "email", Email: &email}
}

// NullEmailProperty clears an email property value.
func NullEmailProperty() PropertyValue {
	return EmailPropertyValue{Type: "email"}
}

// PhoneNumberProperty creates a phone_number property value.
func PhoneNumberProperty(phone string) PropertyValue {
	return PhoneNumberPropertyValue{Type: "phone_number", PhoneNumber: &phone}
}

// NullPhoneNumberProperty clears a phone_number property value.
func NullPhoneNumberProperty() PropertyValue {
	return PhoneNumberPropertyValue{Type: "phone_number"}
}
