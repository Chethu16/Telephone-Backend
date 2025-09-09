package super_admin

type SuperAdminCollege struct {
	SuperAdminId    string `json:"super_admin_id" bson:"super_admin_id"`
	CollegeId       string `json:"college_id" bson:"college_id"`
	CollegeName     string `json:"college_name" bson:"college_name" validate:"required"`
	CollegeEmail    string `json:"college_email" bson:"college_email" validate:"required,email"`
	CollegePhone    string `json:"college_phone" bosn:"college_phone" validate:"requrired"`
	CollegePassword string `json:"college_password" bson:"college_password" validate:"required"`
	CollegeAddress  string `json:"college_address" bson:"college_address" validate:"required"`
	Balance         string `json:"balance" bson:"balance"`
	CreatedAt       string `json:"created_at" bson:"created_at"`
}

type CollegeCreateRequest struct {
	CollegeName     string `json:"college_name" bson:"college_name" validate:"required"`
	CollegeEmail    string `json:"college_email" bson:"college_email" validate:"required,email"`
	CollegePhone    string `json:"college_phone" bosn:"college_phone" validate:"requrired"`
	CollegePassword string `json:"college_password" bson:"college_password" validate:"required"`
	CollegeAddress  string `json:"college_address" bson:"college_address" validate:"required"`
}
type CollegeLoginRequest struct {
	CollegeEmail    string `json:"college_email" bson:"college_email" validate:"required,email"`
	CollegePassword string `json:"college_password" bson:"college_password" validate:"required"`
}
type CollegeResponse struct {
	SuperAdminId string `json:"super_admin_id"`
	CollegeId    string `json:"college_id"`
	CollegeName  string `json:"college_name"`
	Balance      string `json:"balance"`
}
type CollegeTokenResponse struct {
	Token   string `json:"token"`
	Balance string `json:"balance"`
}
type CollegeRecharge struct {
	RechargeId     string `json:"recharge_id" bson:"recharge_id"`
	CollegeId      string `json:"college_id" bson:"college_id"`
	SuperAdminId   string `json:"super_admin_id" bson:"super_admin_id"`
	RechargeAmount string `json:"recharge_amount" bson:"recharge_amount"`
	RechargedAt    string `json:"recharged_at" bson:"recharged_at"`
}
type CollegeRechargeRequest struct {
	SuperAdminId   string `json:"super_admin_id" bson:"super_admin_id" validate:"required"`
	CollegeId      string `json:"college_id" bson:"college_id" validate:"required"`
	RechargeAmount string `json:"recharge_amount" bson:"recharge_amount"`
}
type CollegeRechargeHistoryResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []CollegeRecharge `json:"data"`
}
