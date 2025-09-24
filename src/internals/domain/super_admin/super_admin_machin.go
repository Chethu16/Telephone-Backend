package super_admin

type Machine struct{
	MachineId string `json:"machine_id" bson:"machine_id"`
	MachineNo string `json:"machine_no" bson:"machine_no" validate:"required"`
	MachineName string `json:"machine_name" bson:"machine_name" validate:"required"`
	CollegeId string `json:"college_id" bson:"college_id"`
	SuperAdminId string `json:"super_admin_id" bson:"super_admin_id"`
	Balance string `json:"balance" bson:"balance"`
}
type MachineCreateRequest struct{
	MachineNo string `json:"machine_no" validate:"required"`
	MachineName string `json:"machine_name" validate:"required"`
	CollegeId string `json:"college_id" validate:"required"`
	SuperAdminId string `json:"super_admin_id" validate:"required"`
}
type MachineDeleteRequest struct {
	MachineId string `json:"machine_id" validate:"required"`
}
type GetMachinesByCollegerequest struct{
	CollegeId string `json:"college_id" vaidate:"required"`
}
type GetAllRechargeHistory struct{
	CollegeId string `json:"college_id" validate:"required"`
}
