package models

import (
	"time"

	"github.com/google/uuid"
)

// DataType represents the type of a data field
type DataType string

const (
	DataTypeString   DataType = "string"
	DataTypeInteger  DataType = "integer"
	DataTypeFloat    DataType = "float"
	DataTypeBoolean  DataType = "boolean"
	DataTypeDateTime DataType = "datetime"
	DataTypeArray    DataType = "array"
	DataTypeObject   DataType = "object"
)

// DataField represents a field in a dataset schema
type DataField struct {
	Name        string   `json:"name" bson:"name"`
	Type        DataType `json:"type" bson:"type"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Required    bool     `json:"required" bson:"required"`
	Nullable    bool     `json:"nullable" bson:"nullable"`
	Unique      bool     `json:"unique,omitempty" bson:"unique,omitempty"`
	Default     any      `json:"default,omitempty" bson:"default,omitempty"`
	Metadata    any      `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// DataSchema represents the schema of a dataset
type DataSchema struct {
	Fields      []DataField          `json:"fields" bson:"fields"`
	PrimaryKey  string               `json:"primary_key,omitempty" bson:"primary_key,omitempty"`
	ForeignKeys map[string]string    `json:"foreign_keys,omitempty" bson:"foreign_keys,omitempty"`
	Indexes     []string             `json:"indexes,omitempty" bson:"indexes,omitempty"`
	Constraints map[string]string    `json:"constraints,omitempty" bson:"constraints,omitempty"`
	Metadata    map[string]any       `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// Dataset represents a dataset in the system
type Dataset struct {
	ID          uuid.UUID            `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description,omitempty" bson:"description,omitempty"`
	Schema      DataSchema           `json:"schema" bson:"schema"`
	Data        any                  `json:"data,omitempty" bson:"data,omitempty"` // Adicionado para armazenar o conteúdo do dataset
	Source      string               `json:"source,omitempty" bson:"source,omitempty"`
	Format      string               `json:"format,omitempty" bson:"format,omitempty"`
	Size        int64                `json:"size,omitempty" bson:"size,omitempty"`
	RowCount    int64                `json:"row_count,omitempty" bson:"row_count,omitempty"`
	Tags        []string             `json:"tags,omitempty" bson:"tags,omitempty"`
	Metadata    map[string]any       `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedBy   uuid.UUID            `json:"created_by" bson:"created_by"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}

// CreateDatasetRequest represents a request to create a new dataset
type CreateDatasetRequest struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description,omitempty"`
	Schema      DataSchema           `json:"schema" binding:"required"`
	Source      string               `json:"source,omitempty"`
	Format      string               `json:"format,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Metadata    map[string]any       `json:"metadata,omitempty"`
}

// UpdateDatasetRequest represents a request to update an existing dataset
type UpdateDatasetRequest struct {
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	Schema      *DataSchema          `json:"schema,omitempty"`
	Source      string               `json:"source,omitempty"`
	Format      string               `json:"format,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Metadata    map[string]any       `json:"metadata,omitempty"`
}

// DatasetResponse represents a dataset response
type DatasetResponse struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Schema      DataSchema           `json:"schema"`
	Source      string               `json:"source,omitempty"`
	Format      string               `json:"format,omitempty"`
	Size        int64                `json:"size,omitempty"`
	RowCount    int64                `json:"row_count,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Metadata    map[string]any       `json:"metadata,omitempty"`
	CreatedBy   uuid.UUID            `json:"created_by"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// DatasetListResponse represents a paginated list of datasets
type DatasetListResponse struct {
	Datasets []DatasetResponse `json:"datasets"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

