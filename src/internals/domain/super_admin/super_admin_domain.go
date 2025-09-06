package super_admin

type SuperAdmin struct{
	SuperAdminID string `json:"super_admin_id" bson:"super_admin_id"`
	SuperAdminName string `json:"super_admin_name" bson:"super_admin_name" validate:"required"`
	SuperAdminEmail string `json:"super_admin_email" bson:"super_admin_email" validate:"required,email"`
	SuperAdminPassword string `json:"super_admin_password" bson:"super_admin_password" validate:"required"`

}
type SuperAminCreateRequest struct{
	SuperAdminId string `json:"super_admin_id" bson:"super_admin_id"`
	SuperAdminName string `json:"super_admin_name" bson:"super_admin_name" validate:"required"`
	SuperAdminEmail string `json:"super_admin_email" bson:"super_admin_email" validate:"required,email"`
	SuperAdminPassword string `json:"super_admin_password" bson:"super_admin_password" validate:"required"`

}
type SuperAdminLoginRequest struct{
	SuperAdminEmail string `json:"super_admin_email" bson:"super_admin_email" validate:"required,email"`
	SuperAdminPassword string `json:"super_admin_password" bson:"super_admin_password" validate:"required"`
}
type SuperAdminCreateResponse struct{
	SuperAdminID string `json:"super_admin_id"`
}
type SuperAdminLoginResponse struct{
Token string `json:"token"`
}
