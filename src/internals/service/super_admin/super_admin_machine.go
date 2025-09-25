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
func(sa *SuperAdminMachineService)GetAllMachinesByCollege(ctx context.Context,collegeId string)([]super_admin.Machine,string,error){
	if collegeId == ""{
			return nil,"",errors.New("collegeID must be provided")
	}
	balance ,err := sa.repo.GetCollageBalance(ctx,collegeId)
	if err != nil{
		return nil,"",fmt.Errorf("unable to retrive balance for collegeId %s: %w",collegeId,err)
	}
	machines,err := sa.repo.GetAllMachinesByCollege(ctx,collegeId)
	if err != nil{
		return nil,"",fmt.Errorf("unable to retrive machine list for collegeID %s:%w",collegeId,err)
	}
	return machines,balance,nil
}
func(sa *SuperAdminMachineService)DeleteMachine(ctx context.Context, machineId string)error{
	if machineId == ""{
		return  errors.New("machineID must be provided")
	}
	if err := sa.repo.DeleteMachine(ctx,machineId);err !=nil{
		return fmt.Errorf("failed to delete machine with ID %s:%w",machineId,err)
	}
	return nil
} 