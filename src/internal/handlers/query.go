package handlers

import (
	"net/http"
	"time"

	"github.com/galafis/go-data-api-microservices/internal/models"
	"github.com/galafis/go-data-api-microservices/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// QueryHandler handles data query operations
type QueryHandler struct {
	datasetRepository DatasetRepository
	queryService      QueryService
}

// QueryService defines the interface for query operations
type QueryService interface {
	ExecuteQuery(query *models.QueryRequest) ([]map[string]interface{}, int64, string, float64, error)
	ExecuteTransform(transform *models.TransformRequest) (*models.Dataset, error)
	ExecuteAggregate(aggregate *models.AggregateRequest) ([]map[string]interface{}, error)
	ExecuteJoin(join *models.JoinRequest) (*models.Dataset, error)
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(datasetRepository DatasetRepository, queryService QueryService) *QueryHandler {
	return &QueryHandler{
		datasetRepository: datasetRepository,
		queryService:      queryService,
	}
}

// QueryData handles querying data
// @Summary Query data
// @Description Query data from a dataset
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.QueryRequest true "Query request"
// @Success 200 {object} models.QueryResponse "Query executed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Dataset not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/query [post]
func (h *QueryHandler) QueryData(c *gin.Context) {
	// Parse request
	var req models.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if dataset exists
	dataset, err := h.datasetRepository.FindByID(req.DatasetID)
	if err != nil {
		logger.Errorf("Error finding dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	// Execute query
	data, total, rawSQL, executionTime, err := h.queryService.ExecuteQuery(&req)
	if err != nil {
		logger.Errorf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing query"})
		return
	}

	// Build response
	response := models.QueryResponse{
		Data:          data,
		Total:         total,
		Limit:         req.Limit,
		Offset:        req.Offset,
		ExecutionTime: executionTime,
	}

	// Include raw SQL if requested
	if req.IncludeRaw {
		response.RawSQL = rawSQL
	}

	c.JSON(http.StatusOK, response)
}

// TransformData handles transforming data
// @Summary Transform data
// @Description Transform data from a dataset
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.TransformRequest true "Transform request"
// @Success 200 {object} models.QueryResponse "Transform executed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Dataset not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/transform [post]
func (h *QueryHandler) TransformData(c *gin.Context) {
	// Parse request
	var req models.TransformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if dataset exists
	dataset, err := h.datasetRepository.FindByID(req.DatasetID)
	if err != nil {
		logger.Errorf("Error finding dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	// Execute transform
	start := time.Now()
	result, err := h.queryService.ExecuteTransform(&req)
	if err != nil {
		logger.Errorf("Error executing transform: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing transform"})
		return
	}
	executionTime := time.Since(start).Seconds()

	// Save result as new dataset if requested
	if req.SaveAs != "" {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Create new dataset
		now := time.Now()
		newDataset := &models.Dataset{
			ID:          uuid.New(),
			Name:        req.SaveAs,
			Description: "Transformed from " + dataset.Name,
			Schema:      result.Schema,
			Source:      "transform",
			Format:      dataset.Format,
			Size:        result.Size,
			RowCount:    int64(len(result.Data.([]map[string]interface{}))),
			Tags:        dataset.Tags,
			Metadata: map[string]interface{}{
				"source_dataset": dataset.ID.String(),
				"transform_steps": req.Steps,
			},
			CreatedBy: userID.(uuid.UUID),
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := h.datasetRepository.Create(newDataset); err != nil {
			logger.Errorf("Error saving transformed dataset: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving transformed dataset"})
			return
		}

		// Return dataset info
		c.JSON(http.StatusOK, gin.H{
			"message":       "Transform executed and saved successfully",
			"dataset_id":    newDataset.ID,
			"dataset_name":  newDataset.Name,
			"row_count":     newDataset.RowCount,
			"execution_time": executionTime,
		})
		return
	}

	// Convert result to response format
	data := make([]map[string]interface{}, len(result.Data.([]map[string]interface{})))
	for i, row := range result.Data.([]map[string]interface{}) {
		rowMap := make(map[string]interface{})
		for _, field := range result.Schema.Fields {
			rowMap[field.Name] = row[field.Name]
		}
		data[i] = rowMap
	}

	// Return data directly
	c.JSON(http.StatusOK, gin.H{
		"data":          data,
		"total":         len(data),
		"execution_time": executionTime,
	})
}

// AggregateData handles aggregating data
// @Summary Aggregate data
// @Description Aggregate data from a dataset
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AggregateRequest true "Aggregate request"
// @Success 200 {object} models.QueryResponse "Aggregate executed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Dataset not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/aggregate [post]
func (h *QueryHandler) AggregateData(c *gin.Context) {
	// Parse request
	var req models.AggregateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if dataset exists
	dataset, err := h.datasetRepository.FindByID(req.DatasetID)
	if err != nil {
		logger.Errorf("Error finding dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	// Execute aggregate
	start := time.Now()
	result, err := h.queryService.ExecuteAggregate(&req)
	if err != nil {
		logger.Errorf("Error executing aggregate: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing aggregate"})
		return
	}
	executionTime := time.Since(start).Seconds()

	// Save result as new dataset if requested
	if req.SaveAs != "" {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Create schema for aggregated data
		fields := make([]models.DataField, 0)
		
		// Add group by fields
		for _, field := range req.GroupBy {
			fields = append(fields, models.DataField{
				Name:     field,
				Type:     models.DataTypeString, // Assuming string for simplicity
				Required: true,
				Nullable: false,
			})
		}
		
		// Add aggregation fields
		for _, agg := range req.Aggregations {
			var dataType models.DataType
			switch agg.Type {
			case models.AggregationCount:
				dataType = models.DataTypeInteger
			case models.AggregationSum, models.AggregationAvg, models.AggregationMin, models.AggregationMax:
				dataType = models.DataTypeFloat
			default:
				dataType = models.DataTypeFloat
			}
			
			fields = append(fields, models.DataField{
				Name:     agg.OutputName,
				Type:     dataType,
				Required: true,
				Nullable: false,
			})
		}
		
		// Create new dataset
		now := time.Now()
		newDataset := &models.Dataset{
			ID:          uuid.New(),
			Name:        req.SaveAs,
			Description: "Aggregated from " + dataset.Name,
			Schema: models.DataSchema{
				Fields: fields,
			},
			Source:   "aggregate",
			Format:   dataset.Format,
			RowCount: int64(len(result)),
			Tags:     dataset.Tags,
			Metadata: map[string]interface{}{
				"source_dataset": dataset.ID.String(),
				"group_by":       req.GroupBy,
				"aggregations":   req.Aggregations,
			},
			CreatedBy: userID.(uuid.UUID),
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := h.datasetRepository.Create(newDataset); err != nil {
			logger.Errorf("Error saving aggregated dataset: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving aggregated dataset"})
			return
		}

		// Return dataset info
		c.JSON(http.StatusOK, gin.H{
			"message":       "Aggregate executed and saved successfully",
			"dataset_id":    newDataset.ID,
			"dataset_name":  newDataset.Name,
			"row_count":     newDataset.RowCount,
			"execution_time": executionTime,
		})
		return
	}

	// Return data directly
	c.JSON(http.StatusOK, gin.H{
		"data":          result,
		"total":         len(result),
		"execution_time": executionTime,
	})
}

// JoinData handles joining datasets
// @Summary Join datasets
// @Description Join two datasets
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.JoinRequest true "Join request"
// @Success 200 {object} models.QueryResponse "Join executed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Dataset not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/join [post]
func (h *QueryHandler) JoinData(c *gin.Context) {
	// Parse request
	var req models.JoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if left dataset exists
	leftDataset, err := h.datasetRepository.FindByID(req.LeftDatasetID)
	if err != nil {
		logger.Errorf("Error finding left dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if leftDataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Left dataset not found"})
		return
	}

	// Check if right dataset exists
	rightDataset, err := h.datasetRepository.FindByID(req.RightDatasetID)
	if err != nil {
		logger.Errorf("Error finding right dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if rightDataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Right dataset not found"})
		return
	}

	// Execute join
	start := time.Now()
	result, err := h.queryService.ExecuteJoin(&req)
	if err != nil {
		logger.Errorf("Error executing join: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error executing join"})
		return
	}
	executionTime := time.Since(start).Seconds()

	// Save result as new dataset if requested
	if req.SaveAs != "" {
		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Create new dataset
		now := time.Now()
		newDataset := &models.Dataset{
			ID:          uuid.New(),
			Name:        req.SaveAs,
			Description: "Joined from " + leftDataset.Name + " and " + rightDataset.Name,
			Schema:      result.Schema,
			Source:      "join",
			Format:      leftDataset.Format,
			Size:        result.Size,
			RowCount:    int64(len(result.Data.([]map[string]interface{}))),
			Tags:        append(leftDataset.Tags, rightDataset.Tags...),
			Metadata: map[string]interface{}{
				"left_dataset":  leftDataset.ID.String(),
				"right_dataset": rightDataset.ID.String(),
				"join_type":     req.JoinType,
				"conditions":    req.Conditions,
			},
			CreatedBy: userID.(uuid.UUID),
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := h.datasetRepository.Create(newDataset); err != nil {
			logger.Errorf("Error saving joined dataset: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving joined dataset"})
			return
		}

		// Return dataset info
		c.JSON(http.StatusOK, gin.H{
			"message":       "Join executed and saved successfully",
			"dataset_id":    newDataset.ID,
			"dataset_name":  newDataset.Name,
			"row_count":     newDataset.RowCount,
			"execution_time": executionTime,
		})
		return
	}

	// Convert result to response format
	data := make([]map[string]interface{}, len(result.Data.([]map[string]interface{})))
	for i, row := range result.Data.([]map[string]interface{}) {
		rowMap := make(map[string]interface{})
		for _, field := range result.Schema.Fields {
			rowMap[field.Name] = row[field.Name]
		}
		data[i] = rowMap
	}

	// Return data directly
	c.JSON(http.StatusOK, gin.H{
		"data":          data,
		"total":         len(data),
		"execution_time": executionTime,
	})
}

// QueryData is a placeholder handler for querying data
func QueryData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Query data endpoint"})
}

// TransformData is a placeholder handler for transforming data
func TransformData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Transform data endpoint"})
}

// AggregateData is a placeholder handler for aggregating data
func AggregateData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Aggregate data endpoint"})
}

// JoinData is a placeholder handler for joining data
func JoinData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Join data endpoint"})
}

