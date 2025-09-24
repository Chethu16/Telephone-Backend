package superadmin_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	superadmin_repo "github.com/Chethu16/Chethu/src/internals/repository/super_admin"
	"github.com/Chethu16/Chethu/src/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type SuperAdminMachineService struct {
	repo      superadmin_repo.SuperAdminMachineRepo
	validator *validator.Validate
}

func NewSuperAdminMAchineService(repo superadmin_repo.SuperAdminMachineRepo, v *validator.Validate) *SuperAdminMachineService {
	return &SuperAdminMachineService{
		repo:      repo,
		validator: v,
	}
}

func (sa *SuperAdminMachineService) CreateMachine(ctx context.Context, req super_admin.MachineCreateRequest) (string, error) {
	if err := sa.validator.Struct(req); err != nil {
		return "", fmt.Errorf("invalid machine details: %w", err)
	}
	if req.SuperAdminId == "" || req.CollegeId == "" {
		return "", errors.New("both superAdminId and CollegeId must be provided")
	}

	machineId := utils.GenerateUUID()
	machine := super_admin.Machine{
		MachineId:    machineId,
		MachineNo:    req.MachineNo,
		MachineName:  req.MachineName,
		CollegeId:    req.CollegeId,
		SuperAdminId: req.SuperAdminId,
		Balance:      "0",
	}

	if err := sa.repo.CreateMachine(ctx, machine); err != nil {
		return "", err // return raw repo error
	}

	return machineId, nil
}
