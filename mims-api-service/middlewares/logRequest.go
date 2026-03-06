package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Read request body
		reqBody, _ := ioutil.ReadAll(c.Request.Body)

		// Reset the request body to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		// Process request
		c.Next()

		// Calculate request processing duration
		duration := time.Since(startTime)

		queryParams := c.Request.URL.Query()
		
		// Unmarshal and marshal the JSON request body
		var requestBodyMap map[string]interface{}
		_ = json.Unmarshal(reqBody, &requestBodyMap)
		formattedReqBody, _ := json.Marshal(requestBodyMap)
		// Remove outer curly braces
		formattedReqBodyStr := strings.Trim(string(formattedReqBody), "{}")
		removeStr2Backslash := strings.ReplaceAll(formattedReqBodyStr, "\\", "")
		removeStr2Backslash = strings.ReplaceAll(removeStr2Backslash, "  ", " ")
		removeStr3Backslash := strings.ReplaceAll(removeStr2Backslash, "\"", "")
		removeStr3Backslash = strings.ReplaceAll(removeStr3Backslash, "  ", " ")
		


		formattedQueryParams := formatQueryParams(queryParams)
		// Log request information
		logger.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query_params", formattedQueryParams),
			zap.String("ip", c.ClientIP()),
			zap.Duration("duration", duration),
			zap.String("request_body", removeStr3Backslash),
		)
	}
}

func formatQueryParams(params url.Values) string {
	var formattedParams []string
	for key, values := range params {
		for _, value := range values {
			formattedParams = append(formattedParams, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(formattedParams, " ")
}





