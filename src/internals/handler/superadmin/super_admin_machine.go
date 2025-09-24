package superadmin_handler

import (
	"net/http"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	superadmin_service "github.com/Chethu16/Chethu/src/internals/service/super_admin"
	"github.com/labstack/echo/v4"
)

type SuperAdminMachineHandler struct {
	SuperAdminMachinService *superadmin_service.SuperAdminMachineService
}

func NewSuperAdminMachineHandler(SuperAdminMachineService *superadmin_service.SuperAdminMachineService) *SuperAdminMachineHandler {
	return &SuperAdminMachineHandler{
		SuperAdminMachinService: SuperAdminMachineService,
	}
}

func (h *SuperAdminMachineHandler) CreateMachine(c echo.Context) error {
	var req super_admin.MachineCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "error",
			Error:  "the machine creation data is invalid, please check and try again later",
		})
	}

	machineId, err := h.SuperAdminMachinService.CreateMachine(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "failed",
			Error:  err.Error(), // show real reason (duplicate, invalid, etc.)
		})
	}

	return c.JSON(http.StatusOK, super_admin.SuccesResponse{
		Status:  "success",
		Message: "Machine created successfully",
		Data: map[string]string{
			"machine_id": machineId,
		},
	})
}
