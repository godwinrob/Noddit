package main

import (
	"log"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/godwinrob/noddit/internal/config"
	"github.com/godwinrob/noddit/internal/database"
	"github.com/godwinrob/noddit/internal/handlers"
	"github.com/godwinrob/noddit/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from root .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load and validate configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize Clerk SDK (must be called before any token verification)
	if cfg.ClerkSecretKey != "" {
		clerk.SetKey(cfg.ClerkSecretKey)
		log.Println("Clerk SDK initialized with secret key")
	} else {
		log.Println("WARNING: CLERK_SECRET_KEY not set — Clerk token verification will fail")
	}

	// Initialize database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Gin router
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.FrontendURL}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(corsConfig))

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Add security headers
	router.Use(middleware.SecurityHeaders())

	// Initialize rate limiters
	// Username check: 10 requests per second, burst of 20 (more lenient for UX)
	usernameLimiter := middleware.NewIPRateLimiter(10, 20)
	usernameLimiter.Cleanup()

	// Public routes (no authentication)
	public := router.Group("/")
	{
		// Public API endpoints
		api := public.Group("/api/public")
		{
			// Subnoddits
			api.GET("/subnoddits", h.GetAllSubnoddits)
			api.GET("/subnoddits/active", h.GetActiveSubnoddits)
			api.GET("/subnoddits/:name", h.GetSubnodditByName)
			api.GET("/subnoddits/search/:term", h.SearchSubnoddits)
			api.GET("/moderators/:name", h.GetModerators)

			// Posts
			api.GET("/allposts", h.GetAllPosts)
			api.GET("/recentposts", h.GetRecentPosts)
			api.GET("/popularposts", h.GetPopularPosts)
			api.GET("/allposts/:subnodditName", h.GetPostsBySubnoddit)
			api.GET("/allpostspopular/:subnodditName", h.GetPopularPostsBySubnoddit)
			api.GET("/:subnodditName/:postId", h.GetPost)
			api.GET("/:subnodditName/:postId/replies", h.GetReplies)

			// Votes
			api.GET("/post/votes/:postId", h.GetVotes)

			// User posts
			api.GET("/post/user/:username", h.GetUserPosts)

			// Username availability check (rate limited)
			api.GET("/user/available/:username", middleware.RateLimitMiddleware(usernameLimiter), h.CheckUsernameAvailable)
		}

		// User profile (public)
		public.GET("/api/user/:username", h.GetUser)
	}

	// Create reply - protected route (moved from public group for clarity)
	router.POST("/:subnodditName/:postId/createreply", middleware.AuthMiddleware(), h.CreateReply)

	// Protected routes (require authentication)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Subnoddits
		protected.POST("/subnoddits/create", h.CreateSubnoddit)
		protected.PUT("/subnoddits/update", h.UpdateSubnoddit)
		protected.DELETE("/subnoddits/delete/:name", h.DeleteSubnoddit)

		// Posts
		protected.POST("/post/create", h.CreatePost)
		protected.PUT("/post/update/:postId", h.UpdatePost)
		protected.DELETE("/post/delete/:postId", h.DeletePost)

		// Votes
		protected.POST("/post/vote", h.VotePost)

		// User sync (post-auth hook — creates DB user for Clerk users)
		protected.POST("/user/sync", h.SyncUser)

		// User management
		protected.PUT("/user/update/email/:username", h.UpdateEmail)
		protected.PUT("/user/update/username/:username", h.UpdateUsername)
		protected.PUT("/user/update/name/:username", h.UpdateName)
		protected.PUT("/user/update/avatar/:username", h.UpdateAvatar)
		protected.DELETE("/user/delete/:username", middleware.AdminOnly(), h.DeleteUser)

		// Favorites
		protected.GET("/favorites/:username", h.GetFavorites)
		protected.POST("/favorites/create/post", h.FavoritePost)
		protected.POST("/favorites/create/subnoddit", h.FavoriteSubnoddit)
		protected.DELETE("/favorites/delete/post/:postId", h.UnfavoritePost)
		protected.DELETE("/favorites/subnoddit/:subnodditId", h.UnfavoriteSubnoddit)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
