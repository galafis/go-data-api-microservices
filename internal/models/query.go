package models

import "github.com/google/uuid"

// FilterOperator represents a filter operator
type FilterOperator string

const (
	FilterEQ     FilterOperator = "eq"     // Equal
	FilterNE     FilterOperator = "ne"     // Not equal
	FilterGT     FilterOperator = "gt"     // Greater than
	FilterGTE    FilterOperator = "gte"    // Greater than or equal
	FilterLT     FilterOperator = "lt"     // Less than
	FilterLTE    FilterOperator = "lte"    // Less than or equal
	FilterIN     FilterOperator = "in"     // In array
	FilterNIN    FilterOperator = "nin"    // Not in array
	FilterLIKE   FilterOperator = "like"   // Like (string pattern)
	FilterREGEX  FilterOperator = "regex"  // Regular expression
	FilterEXISTS FilterOperator = "exists" // Field exists
)

// SortDirection represents a sort direction
type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// FilterCondition represents a filter condition
type FilterCondition struct {
	Field    string         `json:"field" binding:"required"`
	Operator FilterOperator `json:"operator" binding:"required"`
	Value    any            `json:"value"`
}

// SortField represents a sort field
type SortField struct {
	Field     string        `json:"field" binding:"required"`
	Direction SortDirection `json:"direction" binding:"required"`
}

// QueryRequest represents a data query request
type QueryRequest struct {
	DatasetID  uuid.UUID         `json:"dataset_id" binding:"required"`
	Fields     []string          `json:"fields,omitempty"`
	Filters    []FilterCondition `json:"filters,omitempty"`
	Sort       []SortField       `json:"sort,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	Offset     int               `json:"offset,omitempty"`
	IncludeRaw bool              `json:"include_raw,omitempty"`
}

// TransformType represents a data transformation type
type TransformType string

const (
	TransformSelect    TransformType = "select"
	TransformRename    TransformType = "rename"
	TransformFilter    TransformType = "filter"
	TransformSort      TransformType = "sort"
	TransformAddColumn TransformType = "add_column"
	TransformCast      TransformType = "cast"
	TransformDrop      TransformType = "drop"
	TransformFill      TransformType = "fill"
	TransformReplace   TransformType = "replace"
	TransformNormalize TransformType = "normalize"
)

// TransformStep represents a data transformation step
type TransformStep struct {
	Type   TransformType `json:"type" binding:"required"`
	Params map[string]any `json:"params" binding:"required"`
}

// TransformRequest represents a data transformation request
type TransformRequest struct {
	DatasetID uuid.UUID      `json:"dataset_id" binding:"required"`
	Steps     []TransformStep `json:"steps" binding:"required"`
	SaveAs    string          `json:"save_as,omitempty"`
}

// AggregationType represents an aggregation type
type AggregationType string

const (
	AggregationCount AggregationType = "count"
	AggregationSum   AggregationType = "sum"
	AggregationAvg   AggregationType = "avg"
	AggregationMin   AggregationType = "min"
	AggregationMax   AggregationType = "max"
)

// AggregationField represents an aggregation field
type AggregationField struct {
	Type       AggregationType `json:"type" binding:"required"`
	Field      string          `json:"field"`
	OutputName string          `json:"output_name" binding:"required"`
}

// AggregateRequest represents a data aggregation request
type AggregateRequest struct {
	DatasetID    uuid.UUID         `json:"dataset_id" binding:"required"`
	GroupBy      []string          `json:"group_by,omitempty"`
	Aggregations []AggregationField `json:"aggregations" binding:"required"`
	Having       []FilterCondition `json:"having,omitempty"`
	Sort         []SortField       `json:"sort,omitempty"`
	Limit        int               `json:"limit,omitempty"`
	SaveAs       string            `json:"save_as,omitempty"`
}

// JoinType represents a join type
type JoinType string

const (
	JoinInner JoinType = "inner"
	JoinLeft  JoinType = "left"
	JoinRight JoinType = "right"
	JoinFull  JoinType = "full"
	JoinCross JoinType = "cross"
)

// JoinCondition represents a join condition
type JoinCondition struct {
	LeftField  string `json:"left_field" binding:"required"`
	RightField string `json:"right_field" binding:"required"`
}

// JoinRequest represents a data join request
type JoinRequest struct {
	LeftDatasetID  uuid.UUID       `json:"left_dataset_id" binding:"required"`
	RightDatasetID uuid.UUID       `json:"right_dataset_id" binding:"required"`
	JoinType       JoinType        `json:"join_type" binding:"required"`
	Conditions     []JoinCondition `json:"conditions" binding:"required"`
	Fields         []string        `json:"fields,omitempty"`
	SaveAs         string          `json:"save_as,omitempty"`
}

// QueryResponse represents a data query response
type QueryResponse struct {
	Data       []map[string]any `json:"data"`
	Total      int64            `json:"total"`
	Limit      int              `json:"limit,omitempty"`
	Offset     int              `json:"offset,omitempty"`
	RawSQL     string           `json:"raw_sql,omitempty"`
	ExecutionTime float64       `json:"execution_time"`
}

