package superadmin_handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	service "github.com/Chethu16/Chethu/src/internals/service/super_admin"
	"github.com/labstack/echo/v4"
)

type SuperAdminHandler struct {
	SuperAdminService *service.SuperAdminService
}

func NewSuperAdminHandler(SuperadminService *service.SuperAdminService) *SuperAdminHandler {
	return &SuperAdminHandler{
		SuperAdminService: SuperadminService,
	}
}
func (h *SuperAdminHandler) CreateSuperAdmin(c echo.Context) error {
	var req super_admin.SuperAminCreateRequest
	if err := c.Bind(&req); err != nil {
		log.Println("create superadmin bind error: ", err)
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "failed",
			Error:  "the request data is invalid please check and try again.",
		})
	}
	res, err := h.SuperAdminService.CreateSuperAdmin(c.Request().Context(), req)
	if err != nil {
		log.Println("create superadmin service error :", err)
		return c.JSON(http.StatusInternalServerError, super_admin.ErrorResponse{
			Status: "failed",
			Error:  "we encounterd an a issue while creating the super admin. please try again later.",
		})
	}
	return c.JSON(http.StatusCreated, super_admin.SuccesResponse{
		Status:  "succes",
		Message: "super admin account created succesfully",
		Data:    res,
	})

}
func (h *SuperAdminHandler) SuperAdminLogin(c echo.Context) error {
	var req super_admin.SuperAdminLoginRequest

	// Bind request
	if err := c.Bind(&req); err != nil {
		log.Println("superadmin login bind error:", err)
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "failed",
			Error:  "the request data is invalid, please check and try again.",
		})
	}

	// Call service
	res, err := h.SuperAdminService.SuperAdminLogin(c.Request().Context(), req)
	if err != nil {
		log.Println("superadmin login error:", err)

		// Return 401 if credentials are wrong
		if strings.Contains(err.Error(), "invalid email or password") {
			return c.JSON(http.StatusUnauthorized, super_admin.ErrorResponse{
				Status: "failed",
				Error:  "invalid email or password",
			})
		}

		// For validation issues
		if strings.Contains(err.Error(), "validator") {
			return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
				Status: "failed",
				Error:  "validation failed for request data",
			})
		}

		// Default: internal server error
		return c.JSON(http.StatusInternalServerError, super_admin.ErrorResponse{
			Status: "failed",
			Error:  "something went wrong on our side. please try again later.",
		})
	}

	// Success
	return c.JSON(http.StatusOK, super_admin.SuccesResponse{
		Status:  "success",
		Message: "login successful",
		Data:    res,
	})
}
