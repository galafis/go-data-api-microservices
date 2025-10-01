package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParseUUID parses a string into a UUID
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// NewUUID generates a new UUID
func NewUUID() uuid.UUID {
	return uuid.New()
}

// GenerateRandomBytes generates random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomString generates a random string
func GenerateRandomString(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// ParseInt parses a string into an int with a default value
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return i
}

// ParseFloat parses a string into a float64 with a default value
func ParseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return f
}

// ParseBool parses a string into a bool with a default value
func ParseBool(s string, defaultValue bool) bool {
	if s == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return b
}

// ParseTime parses a string into a time.Time with a default value
func ParseTime(s, layout string, defaultValue time.Time) time.Time {
	if s == "" {
		return defaultValue
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return defaultValue
	}
	return t
}

// FormatTime formats a time.Time into a string
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// IsEmpty checks if a value is empty
func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return value.Len() == 0
	case reflect.Struct:
		return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
	default:
		return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
	}
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := false
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else if vIsLow && capNext {
			n.WriteByte(v - 32)
			capNext = false
		} else if i == 0 && vIsCap {
			n.WriteByte(v + 32)
		} else {
			n.WriteByte(v)
		}
	}
	return n.String()
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	s = ToCamelCase(s)
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Truncate truncates a string to a specified length
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Round rounds a float64 to a specified precision
func Round(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}

// Contains checks if a slice contains a value
func Contains(slice interface{}, val interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < sliceValue.Len(); i++ {
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), val) {
			return true
		}
	}
	return false
}

// Map applies a function to each element of a slice
func Map(slice interface{}, fn interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return nil
	}

	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return nil
	}

	resultType := fnValue.Type().Out(0)
	resultSlice := reflect.MakeSlice(reflect.SliceOf(resultType), sliceValue.Len(), sliceValue.Len())

	for i := 0; i < sliceValue.Len(); i++ {
		result := fnValue.Call([]reflect.Value{sliceValue.Index(i)})
		resultSlice.Index(i).Set(result[0])
	}

	return resultSlice.Interface()
}

// Filter filters a slice based on a predicate function
func Filter(slice interface{}, fn interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return nil
	}

	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return nil
	}

	resultType := sliceValue.Type().Elem()
	resultSlice := reflect.MakeSlice(reflect.SliceOf(resultType), 0, 0)

	for i := 0; i < sliceValue.Len(); i++ {
		item := sliceValue.Index(i)
		result := fnValue.Call([]reflect.Value{item})
		if result[0].Bool() {
			resultSlice = reflect.Append(resultSlice, item)
		}
	}

	return resultSlice.Interface()
}

// ToJSON converts a value to JSON
func ToJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON converts JSON to a value
func FromJSON(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

// PaginationParams extracts pagination parameters from a request
func PaginationParams(c *gin.Context) (page, pageSize int) {
	page = ParseInt(c.DefaultQuery("page", "1"), 1)
	if page < 1 {
		page = 1
	}

	pageSize = ParseInt(c.DefaultQuery("page_size", "10"), 10)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}

// RespondWithError responds with an error
func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// RespondWithJSON responds with JSON
func RespondWithJSON(c *gin.Context, code int, payload interface{}) {
	c.JSON(code, payload)
}

// GetRequestID gets the request ID from the context
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}
	return requestID.(string)
}

// GetClientIP gets the client IP address
func GetClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check for X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Use RemoteAddr as fallback
	return r.RemoteAddr
}

// FormatBytes formats bytes as a human-readable string
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration formats a duration as a human-readable string
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1f s", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1f m", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1f h", d.Hours())
	}
	return fmt.Sprintf("%.1f d", d.Hours()/24)
}

