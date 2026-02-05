package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"context"
)

// Models
type Video struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	URL         string    `json:"url" gorm:"not null"`
	Thumbnail   string    `json:"thumbnail"`
	Duration    int       `json:"duration"`    // seconds
	Views       int       `json:"views"`
	ProducerID  string    `json:"producer_id"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Producer    *Producer `json:"producer,omitempty" gorm:"foreignKey:ProducerID"`
	Category    *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type Producer struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"not null"`
	Slug         string `json:"slug" gorm:"unique;not null"`
	Description  string `json:"description"`
	Avatar       string `json:"avatar"`
	Specialties  string `json:"specialties"` // JSON array as string
	Rating       float32 `json:"rating"`
	Followers    int    `json:"followers"`
	VideoCount   int    `json:"video_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"unique;not null"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	VideoCount  int       `json:"video_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Config
type Config struct {
	Port        string
	PostgresURL string
	RedisURL    string
	RabbitMQURL string
	APIKey      string
}

// Global instances
var (
	db    *gorm.DB
	rdb   *redis.Client
	cfg   Config
)

// Event structure
type Event struct {
	EventID   string      `json:"eventId"`
	EventType string      `json:"eventType"`
	Timestamp string      `json:"timestamp"`
	TraceID   string      `json:"traceId"`
	Payload   interface{} `json:"payload"`
}

func loadConfig() Config {
	return Config{
		Port:        getEnv("PORT", "8081"),
		PostgresURL: getEnv("POSTGRES_URL", "postgresql://localhost:5432/nudex_catalog"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://localhost"),
		APIKey:      getEnv("API_KEY", "default_api_key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Initialize database
func initDB() {
	var err error
	
	db, err = gorm.Open(postgres.Open(cfg.PostgresURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&Video{}, &Producer{}, &Category{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed data
	seedData()
}

// Initialize Redis
func initRedis() {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}
	
	rdb = redis.NewClient(opt)
	
	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
}

// Seed initial data
func seedData() {
	// Check if data already exists
	var count int64
	db.Model(&Video{}).Count(&count)
	if count > 0 {
		return
	}

	// Create categories
	categories := []Category{
		{ID: uuid.New().String(), Name: "Action", Slug: "action", Description: "High-octane action videos", Icon: "💥"},
		{ID: uuid.New().String(), Name: "Comedy", Slug: "comedy", Description: "Funny and entertaining content", Icon: "😂"},
		{ID: uuid.New().String(), Name: "Drama", Slug: "drama", Description: "Dramatic storytelling", Icon: "🎭"},
		{ID: uuid.New().String(), Name: "Documentary", Slug: "documentary", Description: "Real-world stories", Icon: "📹"},
	}
	db.Create(&categories)

	// Create producers
	producers := []Producer{
		{
			ID: uuid.New().String(), 
			Name: "NUDEX Studios", 
			Slug: "nudex-studios",
			Description: "Premium content creators",
			Specialties: `["Action", "Drama"]`,
			Rating: 4.8,
			Followers: 125000,
		},
		{
			ID: uuid.New().String(), 
			Name: "RedCam Productions", 
			Slug: "redcam-productions",
			Description: "Independent filmmakers",
			Specialties: `["Documentary", "Comedy"]`,
			Rating: 4.6,
			Followers: 89000,
		},
	}
	db.Create(&producers)

	// Create sample videos
	videos := []Video{
		{
			ID: uuid.New().String(),
			Title: "Epic Action Sequence",
			Description: "Mind-blowing action with stunning visuals",
			URL: "https://example.com/video1.mp4",
			Thumbnail: "/placeholder-video.jpg",
			Duration: 180,
			Views: 1250,
			ProducerID: producers[0].ID,
			CategoryID: categories[0].ID,
		},
		{
			ID: uuid.New().String(),
			Title: "Comedy Gold",
			Description: "Hilarious comedy sketch",
			URL: "https://example.com/video2.mp4",
			Thumbnail: "/placeholder-video.jpg",
			Duration: 240,
			Views: 980,
			ProducerID: producers[1].ID,
			CategoryID: categories[1].ID,
		},
		// Add more videos...
	}

	// Generate more videos programmatically
	titles := []string{
		"Dramatic Masterpiece", "Action Packed", "Documentary Truth", 
		"Comedy Central", "Epic Adventure", "Thriller Night",
		"Romance Story", "Sci-Fi Future", "Horror Tales", "Musical Journey",
		"Sports Highlights", "Travel Diary", "Cooking Show", "Tech Review",
		"Gaming Session", "Art Tutorial", "Fashion Show", "News Report",
		"Interview Special", "Behind Scenes"
	}

	for i, title := range titles {
		if i >= 2 { // Skip first 2 as they're already created
			video := Video{
				ID: uuid.New().String(),
				Title: title,
				Description: fmt.Sprintf("Amazing %s content", title),
				URL: fmt.Sprintf("https://example.com/video%d.mp4", i+1),
				Thumbnail: "/placeholder-video.jpg",
				Duration: 120 + (i * 15),
				Views: 500 + (i * 100),
				ProducerID: producers[i%2].ID,
				CategoryID: categories[i%4].ID,
			}
			videos = append(videos, video)
		}
	}

	db.Create(&videos)
	
	log.Printf("Seeded %d videos, %d producers, %d categories", len(videos), len(producers), len(categories))
}

// Middleware for API Key authentication (internal endpoints)
func requireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != cfg.APIKey {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Health check
func healthCheck(c *gin.Context) {
	// Check database
	sqlDB, err := db.DB()
	dbStatus := "connected"
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "disconnected"
	}

	// Check Redis
	ctx := context.Background()
	redisStatus := "connected"
	if rdb.Ping(ctx).Err() != nil {
		redisStatus = "disconnected"
	}

	c.JSON(200, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"database":  dbStatus,
		"redis":     redisStatus,
	})
}

// Get video by ID
func getVideo(c *gin.Context) {
	id := c.Param("id")
	
	var video Video
	result := db.Preload("Producer").Preload("Category").First(&video, "id = ?", id)
	
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Video not found"})
		return
	}

	// Increment view count (async)
	go func() {
		db.Model(&video).Update("views", video.Views+1)
	}()

	c.JSON(200, video)
}

// Search videos
func searchVideos(c *gin.Context) {
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	var videos []Video
	
	dbQuery := db.Preload("Producer").Preload("Category")
	
	if query != "" {
		dbQuery = dbQuery.Where("title ILIKE ? OR description ILIKE ?", 
			"%"+query+"%", "%"+query+"%")
	}
	
	dbQuery.Limit(limit).Offset(offset).Find(&videos)

	c.JSON(200, gin.H{
		"videos": videos,
		"query":  query,
		"limit":  limit,
		"offset": offset,
	})
}

// Get videos by category
func getVideosByCategory(c *gin.Context) {
	slug := c.Param("slug")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	var videos []Video
	db.Joins("JOIN categories ON categories.id = videos.category_id").
		Where("categories.slug = ?", slug).
		Preload("Producer").Preload("Category").
		Limit(limit).Find(&videos)

	c.JSON(200, gin.H{"videos": videos, "category": slug})
}

// Get videos by producer
func getVideosByProducer(c *gin.Context) {
	slug := c.Param("slug")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	var videos []Video
	db.Joins("JOIN producers ON producers.id = videos.producer_id").
		Where("producers.slug = ?", slug).
		Preload("Producer").Preload("Category").
		Limit(limit).Find(&videos)

	c.JSON(200, gin.H{"videos": videos, "producer": slug})
}

// Get random videos
func getRandomVideos(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	var videos []Video
	db.Preload("Producer").Preload("Category").
		Order("RANDOM()").Limit(limit).Find(&videos)

	c.JSON(200, gin.H{"videos": videos})
}

// Get all producers
func getProducers(c *gin.Context) {
	var producers []Producer
	db.Find(&producers)
	c.JSON(200, gin.H{"producers": producers})
}

// Get all categories
func getCategories(c *gin.Context) {
	var categories []Category
	db.Find(&categories)
	c.JSON(200, gin.H{"categories": categories})
}

// Internal: Upsert video
func upsertVideo(c *gin.Context) {
	var video Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Generate ID if not provided
	if video.ID == "" {
		video.ID = uuid.New().String()
	}

	// Upsert
	result := db.Save(&video)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to save video"})
		return
	}

	c.JSON(200, gin.H{"video": video, "message": "Video upserted successfully"})
}

func main() {
	cfg = loadConfig()

	// Initialize services
	initDB()
	initRedis()

	// Setup Gin
	if os.Getenv("GO_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Public routes
	r.GET("/health", healthCheck)
	r.GET("/videos/:id", getVideo)
	r.GET("/videos/search", searchVideos)
	r.GET("/videos/category/:slug", getVideosByCategory)
	r.GET("/videos/producer/:slug", getVideosByProducer)
	r.GET("/videos", getRandomVideos)
	r.GET("/producers", getProducers)
	r.GET("/categories", getCategories)

	// Internal routes (require API key)
	internal := r.Group("/internal")
	internal.Use(requireAPIKey())
	{
		internal.POST("/videos/upsert", upsertVideo)
	}

	log.Printf("🎬 NUDEX Catalog Service starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}