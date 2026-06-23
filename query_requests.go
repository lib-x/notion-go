package notion

import "encoding/json"

// Direction is a Notion sort direction.
type Direction string

const (
	Ascending  Direction = "ascending"
	Descending Direction = "descending"
)

const (
	ResultTypePage       = "page"
	ResultTypeDataSource = "data_source"
)

// DataSourceQueryRequest is a typed request body for querying a data source.
type DataSourceQueryRequest struct {
	Sorts       []Sort `json:"sorts,omitempty"`
	Filter      Filter `json:"filter,omitempty"`
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
	InTrash     *bool  `json:"in_trash,omitempty"`
	ResultType  string `json:"result_type,omitempty"`
}

// Sort is a typed Notion sort request.
type Sort struct {
	Property  string    `json:"property,omitempty"`
	Timestamp string    `json:"timestamp,omitempty"`
	Direction Direction `json:"direction"`
}

// PropertySort sorts by a property.
func PropertySort(property string, direction Direction) Sort {
	return Sort{Property: property, Direction: direction}
}

// TimestampSort sorts by a Notion timestamp such as created_time.
func TimestampSort(timestamp string, direction Direction) Sort {
	return Sort{Timestamp: timestamp, Direction: direction}
}

// Filter is implemented by typed Notion filter request values.
type Filter interface {
	filter()
}

// RawFilterValue preserves an arbitrary filter object inside typed query
// requests.
type RawFilterValue Object

func (RawFilterValue) filter() {}

// RawFilter wraps an arbitrary filter object.
func RawFilter(value Object) Filter {
	return RawFilterValue(value)
}

type compoundFilter struct {
	Operator string
	Filters  []Filter
}

func (compoundFilter) filter() {}

func (f compoundFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(Object{f.Operator: f.Filters})
}

type propertyFilter struct {
	Property  string
	FieldType string
	Condition any
}

func (propertyFilter) filter() {}

func (f propertyFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(Object{
		"property":  f.Property,
		f.FieldType: f.Condition,
	})
}

// And combines filters with a logical and.
func And(filters ...Filter) Filter {
	return compoundFilter{Operator: "and", Filters: filters}
}

// Or combines filters with a logical or.
func Or(filters ...Filter) Filter {
	return compoundFilter{Operator: "or", Filters: filters}
}

type TextFilterCondition struct {
	Equals         string `json:"equals,omitempty"`
	DoesNotEqual   string `json:"does_not_equal,omitempty"`
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	StartsWith     string `json:"starts_with,omitempty"`
	EndsWith       string `json:"ends_with,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type CheckboxFilterCondition struct {
	Equals       *bool `json:"equals,omitempty"`
	DoesNotEqual *bool `json:"does_not_equal,omitempty"`
}

type NumberFilterCondition struct {
	Equals               *float64 `json:"equals,omitempty"`
	DoesNotEqual         *float64 `json:"does_not_equal,omitempty"`
	GreaterThan          *float64 `json:"greater_than,omitempty"`
	LessThan             *float64 `json:"less_than,omitempty"`
	GreaterThanOrEqualTo *float64 `json:"greater_than_or_equal_to,omitempty"`
	LessThanOrEqualTo    *float64 `json:"less_than_or_equal_to,omitempty"`
	IsEmpty              bool     `json:"is_empty,omitempty"`
	IsNotEmpty           bool     `json:"is_not_empty,omitempty"`
}

type DateFilterCondition struct {
	Equals     string `json:"equals,omitempty"`
	Before     string `json:"before,omitempty"`
	After      string `json:"after,omitempty"`
	OnOrBefore string `json:"on_or_before,omitempty"`
	OnOrAfter  string `json:"on_or_after,omitempty"`
	IsEmpty    bool   `json:"is_empty,omitempty"`
	IsNotEmpty bool   `json:"is_not_empty,omitempty"`
}

// TitleContains filters a title property by substring.
func TitleContains(property, value string) Filter {
	return propertyFilter{Property: property, FieldType: "title", Condition: TextFilterCondition{Contains: value}}
}

// RichTextContains filters a rich_text property by substring.
func RichTextContains(property, value string) Filter {
	return propertyFilter{Property: property, FieldType: "rich_text", Condition: TextFilterCondition{Contains: value}}
}

// StatusEquals filters a status property by option name.
func StatusEquals(property, value string) Filter {
	return propertyFilter{Property: property, FieldType: "status", Condition: TextFilterCondition{Equals: value}}
}

// SelectEquals filters a select property by option name.
func SelectEquals(property, value string) Filter {
	return propertyFilter{Property: property, FieldType: "select", Condition: TextFilterCondition{Equals: value}}
}

// MultiSelectContains filters a multi_select property by option name.
func MultiSelectContains(property, value string) Filter {
	return propertyFilter{Property: property, FieldType: "multi_select", Condition: TextFilterCondition{Contains: value}}
}

// CheckboxEquals filters a checkbox property.
func CheckboxEquals(property string, value bool) Filter {
	return propertyFilter{Property: property, FieldType: "checkbox", Condition: CheckboxFilterCondition{Equals: &value}}
}

// NumberEquals filters a number property for equality.
func NumberEquals(property string, value float64) Filter {
	return propertyFilter{Property: property, FieldType: "number", Condition: NumberFilterCondition{Equals: &value}}
}

// DateOnOrAfter filters a date property.
func DateOnOrAfter(property, value string) Filter {
	return propertyFilter{Property: property, FieldType: "date", Condition: DateFilterCondition{OnOrAfter: value}}
}
