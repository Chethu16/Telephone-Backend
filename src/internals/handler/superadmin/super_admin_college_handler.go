package superadmin_handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	superadmin_service "github.com/Chethu16/Chethu/src/internals/service/super_admin"
	"github.com/labstack/echo/v4"
)

type SuperAdminCollegeHandler struct {
	service *superadmin_service.SuperAdminCollegeService
}

func NewSuperAdminCollegeHandler(service *superadmin_service.SuperAdminCollegeService) *SuperAdminCollegeHandler {
	return &SuperAdminCollegeHandler{service: service}
}
func (h *SuperAdminCollegeHandler) CreateCollege(c echo.Context) error {
	var req super_admin.CollegeCreateRequest
	if err := c.Bind(&req); err != nil {
		log.Println("Create college bind error :", err)
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "We could not process your request. Please ensure all required fields are filled correctly.",
		})
	}
	adminId := strings.TrimSpace(c.Param("super_admin_id"))
	if adminId == "" {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "super Admin ID is required to create a college.",
		})
	}
	res, err := h.service.CreateCollege(c.Request().Context(), req, adminId)
	if err != nil {
		log.Println("Create college service error:", err)
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "unable to create college :" + err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, super_admin.SuccesResponse{
		Status:  "succes",
		Message: "college created succesfully.",
		Data:    res,
	})
}
func (h *SuperAdminCollegeHandler) CollegeLogin(c echo.Context) error {
	var req super_admin.CollegeLoginRequest
	if err := c.Bind(&req); err != nil {
		log.Println("Collegelogin bind error:", err)
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "failed",
			Error:  "Invalid login request please check your credential formate.",
		})
	}
	res, err := h.service.CollegeLogin(c.Request().Context(), req)
	if err != nil {
		log.Println("college login error:", err)
		return c.JSON(http.StatusUnauthorized, super_admin.ErrorResponse{
			Status: "error",
			Error:  "login failed: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, super_admin.SuccesResponse{
		Status:  "succes",
		Message: "Login Succesful",
		Data:    res,
	})

}
func (h *SuperAdminCollegeHandler) GetCollegesBySuperAdminID(c echo.Context) error {
	adminId := strings.TrimSpace(c.Param("super_admin_id"))
	if adminId == "" {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "super admin required to fetch colleges",
		})
	}
	colleges, err := h.service.GetCollegesBySuperadminId(c.Request().Context(), adminId)
	if err != nil {
		log.Println("GetCollegesBySuperAdminID error:", err)
		return c.JSON(http.StatusInternalServerError, super_admin.ErrorResponse{
			Status: "error",
			Error:  "unable to fetch colleges: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, super_admin.SuccesResponse{
		Status:  "succes",
		Message: "colleges fetched succesfully.",
		Data:    colleges,
	})
}
func (h *SuperAdminCollegeHandler) GetCollegeDetails(c echo.Context) error {
	collegeID := strings.TrimSpace(c.Param("college_id"))
	if collegeID == "" {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "college id required for feth details",
		})
	}
	college, err := h.service.GetCollegeDetails(c.Request().Context(), collegeID)
	if err != nil {
		log.Println("Get college details error: ", err)
		return c.JSON(http.StatusNotFound, super_admin.ErrorResponse{
			Status: "error",
			Error:  "College not found :" + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, super_admin.SuccesResponse{
		Status:  "succes",
		Message: "college details recived succesfully",
		Data:    college,
	})
}
func (h *SuperAdminCollegeHandler) DeleteCollege(c echo.Context) error {
	collegeID := strings.TrimSpace(c.Param("college_id"))
	if collegeID == "" {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "college id required to delete",
		})
	}
	if err := h.service.DeleteCollege(c.Request().Context(), collegeID); err != nil {
		log.Println("Delete college error :", err)
		return c.JSON(http.StatusInternalServerError, super_admin.ErrorResponse{
			Status: "error",
			Error:  "unable to delete college :" + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, super_admin.SuccesResponse{
		Status:  "succes",
		Message: "college deteled succesfully",
	})
}
func (h *SuperAdminCollegeHandler) RechargeCollge(c echo.Context) error {
	type Requesst struct {
		RechargeAmount string `json:"recharge_amount"`
	}
	var req Requesst
	if err := c.Bind(&req); err != nil {
		log.Println("recharge college bind error : ", err)
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "invalid recharge request format.",
		})
	}

	collegeId := strings.TrimSpace(c.Param("college_id"))
	if collegeId == "" {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "college id required for recharge",
		})
	}
	superAdminId := strings.TrimSpace(c.Param("super_admin_id"))
	if superAdminId == "" {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "superAdminId required for recharge",
		})
	}
	if err := h.service.RechargeCollege(c.Request().Context(), super_admin.CollegeRechargeRequest{
		SuperAdminId:   superAdminId,
		CollegeId:      collegeId,
		RechargeAmount: req.RechargeAmount,
	}); err != nil {
		log.Println("Recharge coollege service error : ", err)
		return c.JSON(http.StatusInternalServerError, super_admin.ErrorResponse{
			Status: "error",
			Error:  "unable to process recharge :" + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, super_admin.SuccesResponse{
		Status:  "succes",
		Message: "college account recharged succesfuly",
	})

}
