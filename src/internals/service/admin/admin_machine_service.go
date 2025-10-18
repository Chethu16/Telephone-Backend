package admin_service

import admin_repo "github.com/Chethu16/Chethu/src/internals/repository/admin"

type MachineRechargeService struct{
	Repo *admin_repo.MachineRechargeRepo
}
func NewMachineRechargeService(r *admin_repo.MachineRechargeRepo)*MachineRechargeService{
	return &MachineRechargeService{
		Repo: r,
	}
}
