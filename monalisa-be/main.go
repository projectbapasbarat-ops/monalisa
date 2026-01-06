package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"monalisa-be/internal/handler"
	"monalisa-be/internal/middleware"
	"monalisa-be/internal/repository"
	"monalisa-be/internal/service"
)

func main() {
	// =========================
	// LOAD ENV
	// =========================
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	dsn := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// =========================
	// DATABASE
	// =========================
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// =========================
	// GIN SETUP
	// =========================
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// =========================
	// CORS (FRONTEND READY)
	// =========================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// =========================
	// REPOSITORIES
	// =========================
	userRepo := &repository.UserRepository{DB: db}
	roleRepo := &repository.RoleRepository{DB: db}
	auditRepo := &repository.AuditRepository{DB: db}

	// =========================
	// SERVICES
	// =========================
	adminUserService := service.NewAdminUserService(
		userRepo,
		roleRepo,
		auditRepo,
	)

	adminRoleService := service.NewAdminRoleService(roleRepo)

	// =========================
	// HANDLERS
	// =========================
	authHandler := handler.NewAuthHandler(userRepo)
	adminUserHandler := handler.NewAdminUserHandler(adminUserService)
	adminRoleHandler := handler.NewAdminRoleHandler(adminRoleService)

	// =========================
	// PUBLIC ROUTES
	// =========================
	r.POST("/api/v1/auth/login", authHandler.Login)

	// =========================
	// PROTECTED ADMIN ROUTES
	// =========================
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.RequirePermission("user.manage"))
	{
		admin.GET("/users", adminUserHandler.ListUsers)
		admin.POST("/users/:id/roles", adminUserHandler.AssignRole)
		admin.DELETE("/users/:id/roles/:role", adminUserHandler.RemoveRole)
		admin.GET("/roles", adminRoleHandler.ListRoles)
	}

	// =========================
	// START SERVER
	// =========================
	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
