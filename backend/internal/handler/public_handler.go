package handler

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// PublicHandler handles public API requests (no authentication required)
type PublicHandler struct {
	usageService  *service.UsageService
	apiKeyService *service.APIKeyService
}

// NewPublicHandler creates a new PublicHandler
func NewPublicHandler(usageService *service.UsageService, apiKeyService *service.APIKeyService) *PublicHandler {
	return &PublicHandler{
		usageService:  usageService,
		apiKeyService: apiKeyService,
	}
}

// PublicUsageStatsResponse represents the usage statistics response
type PublicUsageStatsResponse struct {
	TotalRequests            int64   `json:"total_requests"`
	TotalInputTokens         int64   `json:"total_input_tokens"`
	TotalOutputTokens        int64   `json:"total_output_tokens"`
	TotalCacheCreationTokens int64   `json:"total_cache_creation_tokens"`
	TotalCacheReadTokens     int64   `json:"total_cache_read_tokens"`
	TotalTokens              int64   `json:"total_tokens"`
	TotalCost                float64 `json:"total_cost"`
	TotalActualCost          float64 `json:"total_actual_cost"`
	AverageDurationMs        float64 `json:"average_duration_ms"`
}

// Usage handles getting usage statistics by API key
// GET /api/v1/public/usage?key=xxx&period=today
// GET /api/v1/public/usage?key=xxx&start_date=2024-01-01&end_date=2024-01-31
func (h *PublicHandler) Usage(c *gin.Context) {
	// Get API key from query parameter
	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "API key is required")
		return
	}

	// Validate API key exists and is active
	apiKey, err := h.apiKeyService.GetByKey(c.Request.Context(), key)
	if err != nil {
		response.NotFound(c, "API key not found")
		return
	}

	if !apiKey.IsActive() {
		response.BadRequest(c, "API key is not active")
		return
	}

	// Parse time range parameters
	userTZ := c.Query("timezone")
	now := timezone.NowInUserLocation(userTZ)
	var startTime, endTime time.Time

	// Check for custom date range first
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" && endDateStr != "" {
		// Use custom date range
		var parseErr error
		startTime, parseErr = timezone.ParseInUserLocation("2006-01-02", startDateStr, userTZ)
		if parseErr != nil {
			response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD")
			return
		}
		endTime, parseErr = timezone.ParseInUserLocation("2006-01-02", endDateStr, userTZ)
		if parseErr != nil {
			response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD")
			return
		}
		// Set end time to end of day
		endTime = endTime.Add(24*time.Hour - time.Nanosecond)
	} else {
		// Use period parameter
		period := c.DefaultQuery("period", "today")
		switch period {
		case "today":
			startTime = timezone.StartOfDayInUserLocation(now, userTZ)
		case "week":
			startTime = now.AddDate(0, 0, -7)
		case "month":
			startTime = now.AddDate(0, -1, 0)
		default:
			startTime = timezone.StartOfDayInUserLocation(now, userTZ)
		}
		endTime = now
	}

	// Get usage statistics for the API key
	stats, err := h.usageService.GetDetailedStatsByAPIKey(c.Request.Context(), apiKey.ID, startTime, endTime)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Build response
	resp := PublicUsageStatsResponse{
		TotalRequests:            stats.TotalRequests,
		TotalInputTokens:         stats.TotalInputTokens,
		TotalOutputTokens:        stats.TotalOutputTokens,
		TotalCacheCreationTokens: stats.TotalCacheCreationTokens,
		TotalCacheReadTokens:     stats.TotalCacheReadTokens,
		TotalTokens:              stats.TotalTokens,
		TotalCost:                stats.TotalCost,
		TotalActualCost:          stats.TotalActualCost,
		AverageDurationMs:        stats.AverageDurationMs,
	}

	response.Success(c, resp)
}
