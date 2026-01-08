package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(gateway *Gateway, environment string) *gin.Engine {
	// Set Gin mode based on environment
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if environment == "test" {
		gin.SetMode(gin.TestMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery()) // Panic recovery
	router.Use(LoggingMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(CORSMiddleware())
	router.Use(MetricsMiddleware())
	// router.Use(AuthMiddleware()) // Enable when implemented

	// Health check endpoint (non-S3)
	router.GET("/health", healthCheckHandler)
	router.GET("/metrics", metricsHandler) // Prometheus metrics

	// Admin API (non-S3)
	admin := router.Group("/admin")
	{
		admin.GET("/cluster/status", clusterStatusHandler)
		admin.GET("/nodes", nodesStatusHandler)
		admin.GET("/costs/by-bucket", costsByBucketHandler)
		admin.GET("/costs/top-objects", topObjectsHandler)
		admin.GET("/repair/status", repairStatusHandler)
	}

	// S3 API routes
	setupS3Routes(router, gateway)

	return router
}

func setupS3Routes(router *gin.Engine, gateway *Gateway) {
	// Root-level operations
	router.GET("/", gateway.ListBuckets)

	// Bucket-level operations
	bucket := router.Group("/:bucket")
	{
		// Bucket operations
		bucket.HEAD("", gateway.HeadBucket)
		bucket.PUT("", gateway.CreateBucket)
		bucket.DELETE("", gateway.DeleteBucket)
		bucket.GET("", handleBucketGet(gateway))

		// Object operations
		bucket.HEAD("/*key", gateway.HeadObject)
		bucket.GET("/*key", gateway.GetObject)
		bucket.PUT("/*key", handleObjectPut(gateway))
		bucket.DELETE("/*key", handleObjectDelete(gateway))
		bucket.POST("/*key", handleObjectPost(gateway))
	}
}

// Handler wrappers that check query parameters for operation type

func handleBucketGet(gateway *Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for multipart upload listing
		if c.Query("uploads") != "" {
			gateway.ListMultipartUploads(c)
			return
		}

		// Default: list objects
		gateway.ListObjects(c)
	}
}

func handleObjectPut(gateway *Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for multipart upload part
		if c.Query("uploadId") != "" {
			gateway.UploadPart(c)
			return
		}

		// Default: put object
		gateway.PutObject(c)
	}
}

func handleObjectDelete(gateway *Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for multipart upload abort
		if c.Query("uploadId") != "" {
			gateway.AbortMultipartUpload(c)
			return
		}

		// Default: delete object
		gateway.DeleteObject(c)
	}
}

func handleObjectPost(gateway *Gateway) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initiate multipart upload
		if c.Query("uploads") != "" {
			gateway.InitiateMultipartUpload(c)
			return
		}

		// Complete multipart upload
		if c.Query("uploadId") != "" {
			gateway.CompleteMultipartUpload(c)
			return
		}

		// Invalid POST request
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid POST operation",
		})
	}
}

// Non-S3 handlers

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "plinth-gateway",
		"version": "0.1.0-alpha",
	})
}

func metricsHandler(c *gin.Context) {
	// TODO: Return Prometheus metrics
	c.String(http.StatusOK, "# TODO: Prometheus metrics\n")
}

func clusterStatusHandler(c *gin.Context) {
	// TODO: Implement cluster status check
	c.JSON(http.StatusOK, gin.H{
		"cluster": "healthy",
		"nodes": []gin.H{
			{"id": "node1", "status": "healthy", "disk_usage": 0.45},
			{"id": "node2", "status": "healthy", "disk_usage": 0.32},
			{"id": "node3", "status": "healthy", "disk_usage": 0.28},
		},
		"replication": gin.H{
			"under_replicated": 0,
		},
	})
}

func nodesStatusHandler(c *gin.Context) {
	// TODO: Implement node status
	c.JSON(http.StatusOK, gin.H{
		"nodes": []gin.H{
			{
				"id":              "node1",
				"status":          "healthy",
				"total_disk_gb":   1000,
				"used_disk_gb":    450,
				"object_count":    15234,
				"last_heartbeat":  "2026-01-08T12:00:00Z",
			},
		},
	})
}

func costsByBucketHandler(c *gin.Context) {
	// TODO: Implement cost tracking query
	c.JSON(http.StatusOK, gin.H{
		"costs": []gin.H{
			{"bucket": "ml-datasets", "total_gb": 150.5, "monthly_cost": 3.46},
			{"bucket": "backups", "total_gb": 89.2, "monthly_cost": 2.05},
		},
	})
}

func topObjectsHandler(c *gin.Context) {
	// TODO: Implement top objects query
	limit := c.DefaultQuery("limit", "20")
	_ = limit

	c.JSON(http.StatusOK, gin.H{
		"objects": []gin.H{
			{"key": "ml-datasets/imagenet.tar", "size_gb": 45.2, "cost": 1.04},
			{"key": "ml-datasets/coco.tar", "size_gb": 28.3, "cost": 0.65},
		},
	})
}

func repairStatusHandler(c *gin.Context) {
	// TODO: Query repair worker status
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
		"last_run": "2026-01-08T11:55:00Z",
		"objects_repaired": 0,
		"under_replicated": 0,
	})
}

