## 🧪 Quick Testing Guide (No Infrastructure Required)

### What You Can Test Right Now

**Template Service** - All 11 REST endpoints work with in-memory data (no database needed for testing)
**Email Service** - Health check only (actual email processing requires RabbitMQ)

---

## 🚀 Option 1: Test Template Service with Mock Data (Easiest)

I'll create a simple mock server that returns dummy data for all endpoints.

### Step 1: Create Mock Server

Create this file: `services/template_service/mock-server.go`

```go
package main

import (
	"encoding/json"
	"net/http"
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
				"postgres": "healthy",
				"redis":    "healthy",
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

		variables := req["variables"].(map[string]interface{})
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

		variables := req["variables"].(map[string]interface{})

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

	println("🚀 Mock Template Service running on http://localhost:8081")
	println("✅ All endpoints ready for testing with .http file!")
	println("\nPress Ctrl+C to stop")

	r.Run(":8081")
}
```

### Step 2: Start Mock Server

```bash
cd services/template_service
go run mock-server.go
```

You should see:
```
🚀 Mock Template Service running on http://localhost:8081
✅ All endpoints ready for testing with .http file!
```

### Step 3: Test with .http File

1. **Install REST Client Extension** in VS Code (if not installed):
   - Open Extensions (Cmd+Shift+X)
   - Search "REST Client"
   - Install by Huachao Mao

2. **Open the test file**:
   ```
   services/template_service/api-tests.http
   ```

3. **Run tests** by clicking "Send Request" above each HTTP request:

   ✅ **Test 1: Health Check** (line ~10)
   - Click "Send Request"
   - Expected: 200 OK with "healthy" status

   ✅ **Test 2: Create Welcome Email** (line ~23)
   - Click "Send Request"
   - Expected: 201 Created with template ID

   ✅ **Test 3: Get Template by Key** (line ~86)
   - Click "Send Request"
   - Expected: 200 OK with welcome_email data

   ✅ **Test 4: List Templates** (line ~67)
   - Click "Send Request"
   - Expected: 200 OK with array of templates

   ✅ **Test 5: Create Version** (line ~103)
   - Replace `{template_id}` with any UUID (e.g., `123e4567-e89b-12d3-a456-426614174000`)
   - Click "Send Request"
   - Expected: 201 Created with version 2

   ✅ **Test 6: Preview Template** (line ~165)
   - Replace `{template_id}` and `{version_id}` with any UUIDs
   - Click "Send Request"
   - Expected: Rendered subject and body

   ✅ **Test 7: Validate Template** (line ~136)
   - Click "Send Request"
   - Expected: valid=true, missing_variables=[]

---

## 📊 What Each Test Shows

| Endpoint | What It Does | Mock Response |
|----------|-------------|---------------|
| `GET /health` | Check service health | Returns healthy status |
| `POST /api/v1/templates` | Create new template | Returns template with generated UUID |
| `GET /api/v1/templates/key/:key` | Get template by key | Returns welcome_email template |
| `GET /api/v1/templates` | List all templates | Returns 2 templates with pagination |
| `PUT /api/v1/templates/:id` | Update template | Returns success message |
| `DELETE /api/v1/templates/:id` | Soft delete template | Returns success message |
| `POST /api/v1/templates/:id/versions` | Create new version | Returns version 2 |
| `GET /api/v1/templates/:id/versions` | Get version history | Returns 2 versions |
| `POST /api/v1/templates/:id/versions/:version_id/publish` | Publish version | Returns success |
| `POST /api/v1/templates/validate` | Validate variables | Checks for missing variables |
| `POST /api/v1/templates/:id/preview` | Preview with data | Returns rendered template |

---

## 🎯 Testing Checklist

- [ ] Health check returns healthy
- [ ] Create template returns 201 with UUID
- [ ] Get template by key returns template data
- [ ] List templates returns array
- [ ] Create version returns version 2
- [ ] Get version history returns multiple versions
- [ ] Publish version succeeds
- [ ] Validate template checks variables
- [ ] Preview renders template
- [ ] Update template succeeds
- [ ] Delete template succeeds

---

## ⚡ Quick Demo

After starting the mock server, try these in order:

1. Health check → See green "healthy" status
2. Create template → Get back a template with UUID
3. Get by key → Retrieve the welcome_email template
4. Preview → See variables replaced in output

**That's it!** No database, no Redis, no RabbitMQ needed. Just pure HTTP testing with dummy data.

---

## 📝 Notes

- The mock server returns **dummy data** - it doesn't persist anything
- UUIDs are randomly generated for each request
- All responses follow the same format as the real service
- Perfect for testing API contracts and response structures
- To test real functionality, you'd need PostgreSQL + Redis (see full setup)

---

## 🔄 Email Service Testing

For Email Service, since it's a RabbitMQ consumer (not REST API), you can only test:

```bash
cd services/email_service
go run cmd/server/main.go
```

Then test health check:
- Open `services/email_service/api-tests.http`
- Click "Send Request" on the health check endpoint
- It will show "unhealthy" for RabbitMQ/Redis (expected without infrastructure)

---

**That's the easiest way to test all endpoints with the .http file!** 🎉
