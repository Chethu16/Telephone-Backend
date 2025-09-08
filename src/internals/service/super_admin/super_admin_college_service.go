package superadmin_service

import (
	"context"
	"errors"

	"strings"
	"time"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	superadmin_repo "github.com/Chethu16/Chethu/src/internals/repository/super_admin"
	"github.com/Chethu16/Chethu/src/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type SuperAdminColllegeService struct {
	Repo     *superadmin_repo.SuperAdminCollegeRepository
	Validate *validator.Validate
}

func NewSuperAdminCollegeService(r *superadmin_repo.SuperAdminCollegeRepository, v *validator.Validate) *SuperAdminColllegeService {
	if v == nil {
		v = validator.New()
	}
	return &SuperAdminColllegeService{
		Repo:     r,
		Validate: v,
	}
}
func (sa *SuperAdminColllegeService) CreateCollege(ctx context.Context, req super_admin.CollegeCreateRequest, SuperadminId string) (*super_admin.CollegeResponse, error) {
	if sa.Validate == nil {
		return nil, errors.New("validator not initialized")
	}
	if err := sa.Validate.Struct(req); err != nil {
		return nil, err // return actual validation errors
	}
	superadmincollegeEmail := strings.ToLower(strings.TrimSpace(req.CollegeEmail))
	exists, err := sa.Repo.CheckCollegeEmailExists(ctx, superadmincollegeEmail)
	if err != nil {
		return nil, errors.New("unable to varify college email please try again later")
	}
	if exists {
		return nil, errors.New("this email registerd with another college")
	}
	hashPassword, err := utils.HashPassword(req.CollegePassword)
	if err != nil {
		return nil, errors.New("unable to process password, please try agin")
	}
	collegeId := utils.GenerateUUID()

	college := super_admin.SuperAdminCollege{
		SuperAdminId:    SuperadminId,
		CollegeId:       collegeId,
		CollegeName:     req.CollegeName,
		CollegeEmail:    superadmincollegeEmail,
		CollegePhone:    req.CollegePhone,
		CollegePassword: hashPassword,
		CollegeAddress:  req.CollegeAddress,
		Balance:         "0",
		CreatedAt:       time.Now().Format(time.RFC3339),
	}
	if err := sa.Repo.CreateCollege(ctx, college); err != nil {
		return nil, errors.New("failed to create college,try again")
	}
	return &super_admin.CollegeResponse{
		SuperAdminId: SuperadminId,
		CollegeId:    collegeId,
		CollegeName:  req.CollegeName,
		Balance:      "0",
	}, nil

}
