package admin

type User struct{
	UserId string `json:"user_id" bson:"user_id"`
	CollegeId string `json:"college_id" bosn:"college_id"`
	MachineId string `json:"machine_id" bson:"machine_id"`
	UserName string `json:"user_name" bson:"user_name" validate:"required"`
	Email string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}

type UserAccessCreateRequest struct{
	CollegeId string `json:"college_id" validate:"required"`
	MachineId string `json:"machine_id" validate:"required"`
	UserId string `json:"user_id" bson:"user_id"`
	UserName string `json:"user_name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password"  validate:"required"`
}

type UserAccessLoginRequest struct{
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type UserAccesLoginResponse struct{
	UserId string `json:"user_id"`
	MachineId string `json:"machine_id"`
	CollegeId string `json:"college_id"`
	Balance string `json:"balance"`
}
type UserAccessDeleteRequest struct{
	UserId string `json:"user_id" validate:"required"`
}