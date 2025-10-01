package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/galafis/go-data-api-microservices/internal/models"
	"github.com/galafis/go-data-api-microservices/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DatasetHandler handles dataset operations
type DatasetHandler struct {
	datasetRepository DatasetRepository
}

// DatasetRepository defines the interface for dataset operations
type DatasetRepository interface {
	FindByID(id uuid.UUID) (*models.Dataset, error)
	FindAll(page, pageSize int, filters map[string]interface{}) ([]models.Dataset, int64, error)
	Create(dataset *models.Dataset) error
	Update(dataset *models.Dataset) error
	Delete(id uuid.UUID) error
}

// NewDatasetHandler creates a new dataset handler
func NewDatasetHandler(datasetRepository DatasetRepository) *DatasetHandler {
	return &DatasetHandler{
		datasetRepository: datasetRepository,
	}
}

// ListDatasets handles listing datasets
// @Summary List datasets
// @Description List all datasets with pagination
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param name query string false "Filter by name"
// @Param tag query string false "Filter by tag"
// @Success 200 {object} models.DatasetListResponse "Datasets retrieved successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/datasets [get]
func (h *DatasetHandler) ListDatasets(c *gin.Context) {
	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Parse filters
	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if tag := c.Query("tag"); tag != "" {
		filters["tag"] = tag
	}

	// Get datasets
	datasets, total, err := h.datasetRepository.FindAll(page, pageSize, filters)
	if err != nil {
		logger.Errorf("Error finding datasets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convert to response
	response := models.DatasetListResponse{
		Datasets: make([]models.DatasetResponse, len(datasets)),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	for i, dataset := range datasets {
		response.Datasets[i] = models.DatasetResponse{
			ID:          dataset.ID,
			Name:        dataset.Name,
			Description: dataset.Description,
			Schema:      dataset.Schema,
			Source:      dataset.Source,
			Format:      dataset.Format,
			Size:        dataset.Size,
			RowCount:    dataset.RowCount,
			Tags:        dataset.Tags,
			Metadata:    dataset.Metadata,
			CreatedBy:   dataset.CreatedBy,
			CreatedAt:   dataset.CreatedAt,
			UpdatedAt:   dataset.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetDataset handles getting a dataset by ID
// @Summary Get a dataset
// @Description Get a dataset by ID
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dataset ID"
// @Success 200 {object} models.DatasetResponse "Dataset retrieved successfully"
// @Failure 400 {object} ErrorResponse "Invalid dataset ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Dataset not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/datasets/{id} [get]
func (h *DatasetHandler) GetDataset(c *gin.Context) {
	// Parse dataset ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dataset ID"})
		return
	}

	// Get dataset
	dataset, err := h.datasetRepository.FindByID(id)
	if err != nil {
		logger.Errorf("Error finding dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	// Convert to response
	response := models.DatasetResponse{
		ID:          dataset.ID,
		Name:        dataset.Name,
		Description: dataset.Description,
		Schema:      dataset.Schema,
		Source:      dataset.Source,
		Format:      dataset.Format,
		Size:        dataset.Size,
		RowCount:    dataset.RowCount,
		Tags:        dataset.Tags,
		Metadata:    dataset.Metadata,
		CreatedBy:   dataset.CreatedBy,
		CreatedAt:   dataset.CreatedAt,
		UpdatedAt:   dataset.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// CreateDataset handles creating a new dataset
// @Summary Create a dataset
// @Description Create a new dataset
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateDatasetRequest true "Dataset creation request"
// @Success 201 {object} models.DatasetResponse "Dataset created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/datasets [post]
func (h *DatasetHandler) CreateDataset(c *gin.Context) {
	// Parse request
	var req models.CreateDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create dataset
	now := time.Now()
	dataset := &models.Dataset{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Schema:      req.Schema,
		Source:      req.Source,
		Format:      req.Format,
		Tags:        req.Tags,
		Metadata:    req.Metadata,
		CreatedBy:   userID.(uuid.UUID),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.datasetRepository.Create(dataset); err != nil {
		logger.Errorf("Error creating dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convert to response
	response := models.DatasetResponse{
		ID:          dataset.ID,
		Name:        dataset.Name,
		Description: dataset.Description,
		Schema:      dataset.Schema,
		Source:      dataset.Source,
		Format:      dataset.Format,
		Size:        dataset.Size,
		RowCount:    dataset.RowCount,
		Tags:        dataset.Tags,
		Metadata:    dataset.Metadata,
		CreatedBy:   dataset.CreatedBy,
		CreatedAt:   dataset.CreatedAt,
		UpdatedAt:   dataset.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateDataset handles updating a dataset
// @Summary Update a dataset
// @Description Update an existing dataset
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dataset ID"
// @Param request body models.UpdateDatasetRequest true "Dataset update request"
// @Success 200 {object} models.DatasetResponse "Dataset updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Dataset not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/datasets/{id} [put]
func (h *DatasetHandler) UpdateDataset(c *gin.Context) {
	// Parse dataset ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dataset ID"})
		return
	}

	// Parse request
	var req models.UpdateDatasetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get dataset
	dataset, err := h.datasetRepository.FindByID(id)
	if err != nil {
		logger.Errorf("Error finding dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	// Check if user is the owner
	if dataset.CreatedBy != userID.(uuid.UUID) {
		// Check if user is admin (would require role check)
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this dataset"})
		return
	}

	// Update dataset
	if req.Name != "" {
		dataset.Name = req.Name
	}
	if req.Description != "" {
		dataset.Description = req.Description
	}
	if req.Schema != nil {
		dataset.Schema = *req.Schema
	}
	if req.Source != "" {
		dataset.Source = req.Source
	}
	if req.Format != "" {
		dataset.Format = req.Format
	}
	if req.Tags != nil {
		dataset.Tags = req.Tags
	}
	if req.Metadata != nil {
		dataset.Metadata = req.Metadata
	}
	dataset.UpdatedAt = time.Now()

	if err := h.datasetRepository.Update(dataset); err != nil {
		logger.Errorf("Error updating dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Convert to response
	response := models.DatasetResponse{
		ID:          dataset.ID,
		Name:        dataset.Name,
		Description: dataset.Description,
		Schema:      dataset.Schema,
		Source:      dataset.Source,
		Format:      dataset.Format,
		Size:        dataset.Size,
		RowCount:    dataset.RowCount,
		Tags:        dataset.Tags,
		Metadata:    dataset.Metadata,
		CreatedBy:   dataset.CreatedBy,
		CreatedAt:   dataset.CreatedAt,
		UpdatedAt:   dataset.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteDataset handles deleting a dataset
// @Summary Delete a dataset
// @Description Delete a dataset by ID
// @Tags data
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Dataset ID"
// @Success 204 "Dataset deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid dataset ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Dataset not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /data/datasets/{id} [delete]
func (h *DatasetHandler) DeleteDataset(c *gin.Context) {
	// Parse dataset ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dataset ID"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get dataset
	dataset, err := h.datasetRepository.FindByID(id)
	if err != nil {
		logger.Errorf("Error finding dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	// Check if user is the owner
	if dataset.CreatedBy != userID.(uuid.UUID) {
		// Check if user is admin (would require role check)
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this dataset"})
		return
	}

	// Delete dataset
	if err := h.datasetRepository.Delete(id); err != nil {
		logger.Errorf("Error deleting dataset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListDatasets is a placeholder handler for listing datasets
func ListDatasets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List datasets endpoint"})
}

// GetDataset is a placeholder handler for getting a dataset
func GetDataset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get dataset endpoint"})
}

// CreateDataset is a placeholder handler for creating a dataset
func CreateDataset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create dataset endpoint"})
}

// UpdateDataset is a placeholder handler for updating a dataset
func UpdateDataset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update dataset endpoint"})
}

// DeleteDataset is a placeholder handler for deleting a dataset
func DeleteDataset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete dataset endpoint"})
}

