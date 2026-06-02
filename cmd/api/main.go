package main

import (
	"log"

	"job-tracker-api/internal/config"
	"job-tracker-api/internal/db"
	"job-tracker-api/internal/handler"
	"job-tracker-api/internal/middleware"
	"job-tracker-api/internal/repository"
	"job-tracker-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	database, err := db.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer database.Close()

	// Repositories
	userRepo := repository.NewUserRepository(database)
	companyRepo := repository.NewCompanyRepository(database)
	jobRepo := repository.NewJobRepository(database)

	// Usecases
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	companyUC := usecase.NewCompanyUsecase(companyRepo)
	jobUC := usecase.NewJobUsecase(jobRepo)

	// Handlers
	authH := handler.NewAuthHandler(authUC)
	companyH := handler.NewCompanyHandler(companyUC)
	jobH := handler.NewJobHandler(jobUC)
	interviewH := handler.NewInterviewHandler()

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		{
			companies := protected.Group("/companies")
			{
				companies.GET("", companyH.GetAll)
				companies.POST("", companyH.Create)
				companies.GET("/:id", companyH.GetByID)
				companies.PUT("/:id", companyH.Update)
				companies.DELETE("/:id", companyH.Delete)
			}

			jobs := protected.Group("/jobs")
			{
				jobs.GET("", jobH.GetAll)
				jobs.POST("", jobH.Create)
				jobs.GET("/:id", jobH.GetByID)
				jobs.PUT("/:id/status", jobH.UpdateStatus)
			}

			interviews := protected.Group("/interviews")
			{
				interviews.GET("", interviewH.GetAll)
				interviews.POST("", interviewH.Create)
				interviews.GET("/:id", interviewH.GetByID)
				interviews.POST("/:id/notes", interviewH.CreateNote)
			}
		}
	}

	log.Println("server started on :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
