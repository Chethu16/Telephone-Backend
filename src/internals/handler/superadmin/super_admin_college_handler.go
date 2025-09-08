package superadmin_handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	superadmin_service "github.com/Chethu16/Chethu/src/internals/service/super_admin"
	"github.com/labstack/echo/v4"
)

type SuperAdminCollegeHandler struct{
	service *superadmin_service.SuperAdminColllegeService
}

func NewSuperAdminCollegeHandler(service *superadmin_service.SuperAdminColllegeService)*SuperAdminCollegeHandler{
	return &SuperAdminCollegeHandler{service: service }
}
func(h *SuperAdminCollegeHandler)CreateCollege(c echo.Context)error{
	var req super_admin.CollegeCreateRequest
	if err := c.Bind(&req); err !=nil{
		log.Println("Create college bind error :",err)
		return c.JSON(http.StatusBadRequest,super_admin.ErrorResponse{
			Status: "error",
			Error: "We could not process your request. Please ensure all required fields are filled correctly.",
		})
	}
	adminId := strings.TrimSpace(c.Param("super_admin_id"))
	if adminId == ""{
		return c.JSON(http.StatusBadRequest,super_admin.ErrorResponse{
			Status: "error",
			Error: "super Admin ID is required to create a college.",
		})
	}
	res,err := h.service.CreateCollege(c.Request().Context(),req,adminId)
	if err != nil{
		log.Println("Create college service error:",err)
		return c.JSON(http.StatusBadRequest,super_admin.ErrorResponse{
			Status: "error",
			Error: "unable to create college :"+err.Error(),
		})
	}
	return c.JSON(http.StatusCreated,super_admin.SuccesResponse{
		Status: "succes",
		Message: "college created succesfully.",
		Data: res,
	})
}