package models

import "github.com/google/uuid"

// StatisticsType represents a type of statistical analysis
type StatisticsType string

const (
	StatsMean       StatisticsType = "mean"
	StatsMedian     StatisticsType = "median"
	StatsMode       StatisticsType = "mode"
	StatsStdDev     StatisticsType = "std_dev"
	StatsVariance   StatisticsType = "variance"
	StatsMin        StatisticsType = "min"
	StatsMax        StatisticsType = "max"
	StatsRange      StatisticsType = "range"
	StatsQuantile   StatisticsType = "quantile"
	StatsHistogram  StatisticsType = "histogram"
	StatsBoxPlot    StatisticsType = "box_plot"
	StatsFrequency  StatisticsType = "frequency"
	StatsDistribution StatisticsType = "distribution"
)

// StatisticsRequest represents a request for statistical analysis
type StatisticsRequest struct {
	DatasetID uuid.UUID      `json:"dataset_id" binding:"required"`
	Type      StatisticsType `json:"type" binding:"required"`
	Fields    []string       `json:"fields" binding:"required"`
	Filters   []FilterCondition `json:"filters,omitempty"`
	Params    map[string]any `json:"params,omitempty"`
}

// CorrelationMethod represents a correlation method
type CorrelationMethod string

const (
	CorrelationPearson  CorrelationMethod = "pearson"
	CorrelationSpearman CorrelationMethod = "spearman"
	CorrelationKendall  CorrelationMethod = "kendall"
)

// CorrelationRequest represents a request for correlation analysis
type CorrelationRequest struct {
	DatasetID uuid.UUID        `json:"dataset_id" binding:"required"`
	Fields    []string         `json:"fields" binding:"required,min=2"`
	Method    CorrelationMethod `json:"method" binding:"required"`
	Filters   []FilterCondition `json:"filters,omitempty"`
}

// TimeSeriesAggregation represents a time series aggregation
type TimeSeriesAggregation string

const (
	TimeSeriesSum   TimeSeriesAggregation = "sum"
	TimeSeriesAvg   TimeSeriesAggregation = "avg"
	TimeSeriesMin   TimeSeriesAggregation = "min"
	TimeSeriesMax   TimeSeriesAggregation = "max"
	TimeSeriesCount TimeSeriesAggregation = "count"
)

// TimeSeriesInterval represents a time series interval
type TimeSeriesInterval string

const (
	TimeSeriesMinute TimeSeriesInterval = "minute"
	TimeSeriesHour   TimeSeriesInterval = "hour"
	TimeSeriesDay    TimeSeriesInterval = "day"
	TimeSeriesWeek   TimeSeriesInterval = "week"
	TimeSeriesMonth  TimeSeriesInterval = "month"
	TimeSeriesQuarter TimeSeriesInterval = "quarter"
	TimeSeriesYear   TimeSeriesInterval = "year"
)

// TimeSeriesRequest represents a request for time series analysis
type TimeSeriesRequest struct {
	DatasetID   uuid.UUID           `json:"dataset_id" binding:"required"`
	TimeField   string              `json:"time_field" binding:"required"`
	ValueField  string              `json:"value_field" binding:"required"`
	Aggregation TimeSeriesAggregation `json:"aggregation" binding:"required"`
	Interval    TimeSeriesInterval  `json:"interval" binding:"required"`
	StartTime   string              `json:"start_time,omitempty"`
	EndTime     string              `json:"end_time,omitempty"`
	GroupBy     []string            `json:"group_by,omitempty"`
	Filters     []FilterCondition   `json:"filters,omitempty"`
}

// ForecastMethod represents a forecasting method
type ForecastMethod string

const (
	ForecastARIMA     ForecastMethod = "arima"
	ForecastExponential ForecastMethod = "exponential"
	ForecastLinear    ForecastMethod = "linear"
	ForecastProphet   ForecastMethod = "prophet"
)

// ForecastRequest represents a request for time series forecasting
type ForecastRequest struct {
	DatasetID   uuid.UUID           `json:"dataset_id" binding:"required"`
	TimeField   string              `json:"time_field" binding:"required"`
	ValueField  string              `json:"value_field" binding:"required"`
	Method      ForecastMethod      `json:"method" binding:"required"`
	Horizon     int                 `json:"horizon" binding:"required,min=1"`
	Interval    TimeSeriesInterval  `json:"interval" binding:"required"`
	StartTime   string              `json:"start_time,omitempty"`
	EndTime     string              `json:"end_time,omitempty"`
	Filters     []FilterCondition   `json:"filters,omitempty"`
	Params      map[string]any      `json:"params,omitempty"`
}

// DataSummary represents a summary of a dataset
type DataSummary struct {
	DatasetID   uuid.UUID         `json:"dataset_id"`
	Name        string            `json:"name"`
	RowCount    int64             `json:"row_count"`
	ColumnCount int               `json:"column_count"`
	NumericColumns []string       `json:"numeric_columns"`
	CategoricalColumns []string   `json:"categorical_columns"`
	DateColumns []string          `json:"date_columns"`
	MissingValues map[string]int64 `json:"missing_values"`
	NumericStats map[string]map[string]float64 `json:"numeric_stats"`
	CategoricalStats map[string]map[string]int64 `json:"categorical_stats"`
}

// StatisticsResult represents the result of a statistical analysis
type StatisticsResult struct {
	Type   StatisticsType     `json:"type"`
	Fields []string           `json:"fields"`
	Results map[string]any    `json:"results"`
}

// CorrelationResult represents the result of a correlation analysis
type CorrelationResult struct {
	Method     CorrelationMethod `json:"method"`
	Fields     []string          `json:"fields"`
	Correlation [][]float64      `json:"correlation"`
	PValues    [][]float64       `json:"p_values,omitempty"`
}

// TimeSeriesPoint represents a point in a time series
type TimeSeriesPoint struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
	Group     string  `json:"group,omitempty"`
}

// TimeSeriesResult represents the result of a time series analysis
type TimeSeriesResult struct {
	TimeField   string           `json:"time_field"`
	ValueField  string           `json:"value_field"`
	Aggregation TimeSeriesAggregation `json:"aggregation"`
	Interval    TimeSeriesInterval `json:"interval"`
	Points      []TimeSeriesPoint `json:"points"`
}

// ForecastPoint represents a point in a forecast
type ForecastPoint struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
	LowerBound float64 `json:"lower_bound,omitempty"`
	UpperBound float64 `json:"upper_bound,omitempty"`
}

// ForecastResult represents the result of a forecast
type ForecastResult struct {
	Method      ForecastMethod   `json:"method"`
	Horizon     int              `json:"horizon"`
	Interval    TimeSeriesInterval `json:"interval"`
	Historical  []TimeSeriesPoint `json:"historical"`
	Forecast    []ForecastPoint  `json:"forecast"`
	Metrics     map[string]float64 `json:"metrics,omitempty"`
}

