package routes

import (
	"roottrack-backend/controllers"
	"roottrack-backend/middleware"
	"roottrack-backend/repositories"
	"roottrack-backend/services"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// External middleware (Optional)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // allow all for now (MVP)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize Repositories and Services
	routineRepo := repositories.RoutineRepository{}
	userRepo := repositories.UserRepository{}
	notificationService := &services.NotificationService{RoutineRepo: routineRepo}

	// Initialize Controllers
	authCtrl := &controllers.AuthController{
		UserRepo:            userRepo,
		RoutineRepo:         routineRepo,
		NotificationService: notificationService,
	}
	userCtrl := &controllers.UserController{UserRepo: userRepo}
	routineCtrl := &controllers.RoutineController{Repo: routineRepo}
	sheddingCtrl := &controllers.SheddingController{}
	productCtrl := &controllers.ProductController{}
	photoCtrl := &controllers.PhotoController{}
	analyticsCtrl := &controllers.AnalyticsController{
		RoutineRepo:  routineCtrl.Repo,
		SheddingRepo: sheddingCtrl.Repo,
		ProductRepo:  productCtrl.Repo,
		PhotoRepo:    photoCtrl.Repo,
	}

	// Start notification worker
	notificationWorker := &services.NotificationWorker{
		UserRepo:        userRepo,
		NotificationSvc: notificationService,
		WorkerLimit:     10,
	}
	go notificationWorker.Start()

	api := r.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
		}

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User Routes
			protected.PUT("/user/update", userCtrl.Update)
			protected.POST("/push-tokens", userCtrl.UpdatePushToken)

			// Routine Routes
			protected.POST("/routines", routineCtrl.Create)
			protected.GET("/routines", routineCtrl.GetAll)
			protected.GET("/routines/today", routineCtrl.GetToday)
			protected.GET("/routines/:id", routineCtrl.GetByID)
			protected.PUT("/routines/:id", routineCtrl.Update)
			protected.DELETE("/routines/:id", routineCtrl.Delete)

			// Shedding Routes
			protected.POST("/shedding", sheddingCtrl.Create)
			protected.GET("/shedding", sheddingCtrl.GetAll)
			protected.GET("/shedding/summary", sheddingCtrl.GetSummary)

			// Product Routes
			protected.POST("/products", productCtrl.Create)
			protected.GET("/products", productCtrl.GetAll)
			protected.PUT("/products/:id", productCtrl.Update)
			protected.DELETE("/products/:id", productCtrl.Delete)

			// Photo Routes
			protected.POST("/photos", photoCtrl.Create)
			protected.GET("/photos", photoCtrl.GetAll)
			protected.DELETE("/photos/:id", photoCtrl.Delete)

			// Analytics Routes
			protected.GET("/analytics/dashboard", analyticsCtrl.GetDashboard)
		}
	}

	return r
}
