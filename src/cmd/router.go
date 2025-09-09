package cmd

import (
	superadmin_handler "github.com/Chethu16/Chethu/src/internals/handler/superadmin"
	superadmin_repo "github.com/Chethu16/Chethu/src/internals/repository/super_admin"
	superadmin_service "github.com/Chethu16/Chethu/src/internals/service/super_admin"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(e *echo.Echo, db *mongo.Database, validate *validator.Validate) {
	SuperAdminRepo := superadmin_repo.NewSuperAdminRepository(db)
	SuperAdminService := superadmin_service.NewSuperAdminService(SuperAdminRepo, validate)
	SuperAdminHandler := superadmin_handler.NewSuperAdminHandler(SuperAdminService)
	SuperAdminRoute := e.Group("/superadmin")

	SuperAdminRoute.POST("/create", SuperAdminHandler.CreateSuperAdmin)
	SuperAdminRoute.POST("/login", SuperAdminHandler.SuperAdminLogin)

	SuperAdminCollegeRepo := superadmin_repo.NewSuperAdminCollegeRepository(db)
	SuperAdminCollegeService := superadmin_service.NewSuperAdminCollegeService(SuperAdminCollegeRepo, validate)
	SuperAdminCollegeHandler := superadmin_handler.NewSuperAdminCollegeHandler(SuperAdminCollegeService)
	SuperAdminCollegeRoute := e.Group("/superadmin/college")

	SuperAdminCollegeRoute.POST("/create/:super_admin_id", SuperAdminCollegeHandler.CreateCollege)
	SuperAdminCollegeRoute.POST("/login", SuperAdminCollegeHandler.CollegeLogin)
}
