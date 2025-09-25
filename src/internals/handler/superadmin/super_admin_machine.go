package superadmin_handler

import (
	"log"
	"net/http"
	"strings"

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
func (h *SuperAdminMachineHandler) GetAllMachine(c echo.Context) error {
	collegeId := c.Param("college_id")
	if collegeId == "" {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "failed",
			Error:  "the 'collegeId' parameter required to fetch machine",
		})
	}
	machines, balance, err := h.SuperAdminMachinService.GetAllMachinesByCollege(c.Request().Context(), collegeId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, super_admin.ErrorResponse{
			Status: "failed",
			Error:  "we could not fetch machines at this time . please try again later",
		})
	}
	return c.JSON(http.StatusOK, super_admin.SuccesResponse{
		Status:  "success",
		Message: "Machines retrived succesfully",
		Data: map[string]interface{}{
			"balance":  balance,
			"machines": machines,
		},
	})

}
func (h *SuperAdminMachineHandler) DeleteMachine(c echo.Context) error {
	machineId := strings.TrimSpace(c.Param("machine_id"))
	if machineId == "" {
		return c.JSON(http.StatusBadRequest, super_admin.ErrorResponse{
			Status: "failed",
			Error:  "the 'machine_id' parameter required for delete a machine",
		})
	}
	err := h.SuperAdminMachinService.DeleteMachine(c.Request().Context(), machineId)
	log.Println(err)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,super_admin.ErrorResponse{
			Status: "failed",
			Error: "we could not delete machine at this time . please try again later",
		})
	}
	return c.JSON(http.StatusOK,super_admin.SuccesResponse{
		Status: "success",
		Message: "Machine deleted succesfully.",
		
	})
}
