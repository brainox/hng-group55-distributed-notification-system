package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"checks": gin.H{
				"postgres": "healthy (mock)",
				"redis":    "healthy (mock)",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Create Template
	r.POST("/api/v1/templates", func(c *gin.Context) {
		var req map[string]interface{}
		c.BindJSON(&req)

		templateID := uuid.New().String()
		versionID := uuid.New().String()

		c.JSON(201, gin.H{
			"success": true,
			"data": gin.H{
				"template": gin.H{
					"id":            templateID,
					"template_key":  req["template_key"],
					"name":          req["name"],
					"description":   req["description"],
					"template_type": req["template_type"],
					"is_active":     true,
					"created_at":    time.Now(),
					"updated_at":    time.Now(),
				},
				"version": gin.H{
					"id":             versionID,
					"template_id":    templateID,
					"version_number": 1,
					"language":       req["language"],
					"subject":        req["subject"],
					"body":           req["body"],
					"variables":      req["variables"],
					"is_published":   true,
					"created_by":     "system",
					"created_at":     time.Now(),
				},
			},
			"message": "Template created successfully",
		})
	})

	// Get Template by Key
	r.GET("/api/v1/templates/key/:key", func(c *gin.Context) {
		key := c.Param("key")
		templateID := uuid.New().String()
		versionID := uuid.New().String()

		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"template": gin.H{
					"id":            templateID,
					"template_key":  key,
					"name":          "Welcome Email",
					"description":   "Welcome email sent to new users",
					"template_type": "email",
					"is_active":     true,
					"created_at":    time.Now(),
					"updated_at":    time.Now(),
				},
				"version": gin.H{
					"id":             versionID,
					"template_id":    templateID,
					"version_number": 1,
					"language":       "en",
					"subject":        "Welcome to {{company_name}}, {{user_name}}!",
					"body":           "<html><body><h1>Welcome {{user_name}}!</h1><p>Thank you for joining {{company_name}}.</p></body></html>",
					"variables":      []string{"user_name", "company_name", "email"},
					"is_published":   true,
					"created_by":     "system",
					"created_at":     time.Now(),
				},
			},
			"message": "Template retrieved successfully",
		})
	})

	// Get Template by ID
	r.GET("/api/v1/templates/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"template": gin.H{
					"id":            id,
					"template_key":  "welcome_email",
					"name":          "Welcome Email",
					"description":   "Welcome email sent to new users",
					"template_type": "email",
					"is_active":     true,
					"created_at":    time.Now(),
					"updated_at":    time.Now(),
				},
				"version": gin.H{
					"id":             uuid.New().String(),
					"template_id":    id,
					"version_number": 1,
					"language":       "en",
					"subject":        "Welcome to {{company_name}}, {{user_name}}!",
					"body":           "<html><body><h1>Welcome {{user_name}}!</h1></body></html>",
					"variables":      []string{"user_name", "company_name", "email"},
					"is_published":   true,
					"created_by":     "system",
					"created_at":     time.Now(),
				},
			},
			"message": "Template retrieved successfully",
		})
	})

	// List Templates
	r.GET("/api/v1/templates", func(c *gin.Context) {
		templates := []gin.H{
			{
				"id":            uuid.New().String(),
				"template_key":  "welcome_email",
				"name":          "Welcome Email",
				"description":   "Welcome email sent to new users",
				"template_type": "email",
				"is_active":     true,
				"created_at":    time.Now(),
				"updated_at":    time.Now(),
			},
			{
				"id":            uuid.New().String(),
				"template_key":  "password_reset",
				"name":          "Password Reset Email",
				"description":   "Password reset instructions",
				"template_type": "email",
				"is_active":     true,
				"created_at":    time.Now(),
				"updated_at":    time.Now(),
			},
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    templates,
			"message": "Templates retrieved successfully",
			"meta": gin.H{
				"total":        2,
				"limit":        10,
				"page":         1,
				"total_pages":  1,
				"has_next":     false,
				"has_previous": false,
			},
		})
	})

	// Update Template
	r.PUT("/api/v1/templates/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"data":    nil,
			"message": "Template updated successfully",
		})
	})

	// Delete Template
	r.DELETE("/api/v1/templates/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"data":    nil,
			"message": "Template deleted successfully",
		})
	})

	// Create Version
	r.POST("/api/v1/templates/:id/versions", func(c *gin.Context) {
		templateID := c.Param("id")
		var req map[string]interface{}
		c.BindJSON(&req)

		c.JSON(201, gin.H{
			"success": true,
			"data": gin.H{
				"id":             uuid.New().String(),
				"template_id":    templateID,
				"version_number": 2,
				"language":       req["language"],
				"subject":        req["subject"],
				"body":           req["body"],
				"variables":      req["variables"],
				"is_published":   false,
				"created_by":     "system",
				"created_at":     time.Now(),
			},
			"message": "Version created successfully",
		})
	})

	// Get Version History
	r.GET("/api/v1/templates/:id/versions", func(c *gin.Context) {
		templateID := c.Param("id")

		versions := []gin.H{
			{
				"id":             uuid.New().String(),
				"template_id":    templateID,
				"version_number": 2,
				"language":       "en",
				"subject":        "Welcome aboard, {{user_name}}!",
				"body":           "<html><body><h1>Hi {{user_name}}!</h1></body></html>",
				"variables":      []string{"user_name", "company_name"},
				"is_published":   false,
				"created_by":     "system",
				"created_at":     time.Now(),
			},
			{
				"id":             uuid.New().String(),
				"template_id":    templateID,
				"version_number": 1,
				"language":       "en",
				"subject":        "Welcome to {{company_name}}, {{user_name}}!",
				"body":           "<html><body><h1>Welcome {{user_name}}!</h1></body></html>",
				"variables":      []string{"user_name", "company_name", "email"},
				"is_published":   true,
				"created_by":     "system",
				"created_at":     time.Now().Add(-24 * time.Hour),
			},
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    versions,
			"message": "Version history retrieved successfully",
		})
	})

	// Publish Version
	r.POST("/api/v1/templates/:id/versions/:version_id/publish", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"data":    nil,
			"message": "Version published successfully",
		})
	})

	// Validate Template
	r.POST("/api/v1/templates/validate", func(c *gin.Context) {
		var req map[string]interface{}
		c.BindJSON(&req)

		variables, _ := req["variables"].(map[string]interface{})
		required := []string{"user_name", "company_name", "email"}
		missing := []string{}

		for _, v := range required {
			if _, ok := variables[v]; !ok {
				missing = append(missing, v)
			}
		}

		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"valid":             len(missing) == 0,
				"missing_variables": missing,
			},
			"message": "Template validation completed",
		})
	})

	// Preview Template
	r.POST("/api/v1/templates/:id/preview", func(c *gin.Context) {
		var req map[string]interface{}
		c.BindJSON(&req)

		variables, _ := req["variables"].(map[string]interface{})

		// Simple variable replacement
		subject := "Welcome to Tech Corp, John Doe!"
		body := "<html><body><h1>Welcome John Doe!</h1><p>Email: john@example.com</p></body></html>"

		if name, ok := variables["user_name"].(string); ok {
			subject = "Welcome to Tech Corp, " + name + "!"
			body = "<html><body><h1>Welcome " + name + "!</h1></body></html>"
		}

		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"subject": subject,
				"body":    body,
			},
			"message": "Template preview generated successfully",
		})
	})

	println("╔════════════════════════════════════════════════════════════╗")
	println("║   🚀 Mock Template Service Running                        ║")
	println("╠════════════════════════════════════════════════════════════╣")
	println("║   URL: http://localhost:8081                               ║")
	println("║   Status: ✅ All 11 endpoints ready                        ║")
	println("║                                                            ║")
	println("║   Next Steps:                                              ║")
	println("║   1. Open: services/template_service/api-tests.http        ║")
	println("║   2. Click 'Send Request' above any HTTP request          ║")
	println("║   3. Start with Health Check (line ~10)                   ║")
	println("║                                                            ║")
	println("║   Press Ctrl+C to stop                                     ║")
	println("╚════════════════════════════════════════════════════════════╝")

	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
