package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	superadmin "github.com/Chethu16/Chethu/src/internals/repository/super_admin"
	"github.com/Chethu16/Chethu/src/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type SuperAdminService struct {
	Repo     *superadmin.SuperAdminRepository
	Validate *validator.Validate
}

func NewSuperAdminService(r *superadmin.SuperAdminRepository, v *validator.Validate) *SuperAdminService {
	return &SuperAdminService{
		Repo:     r,
		Validate: v,
	}
}
func (sa *SuperAdminService) CreateSuperAdmin(ctx context.Context, req super_admin.SuperAminCreateRequest) (*super_admin.SuperAdminCreateResponse, error) {
	if err := sa.Validate.Struct(req); err != nil {
		return nil,errors.New("please fill the all the details")
	}
	superadminEmail := strings.ToLower(strings.TrimSpace(req.SuperAdminEmail))
	exists, err := sa.Repo.CheckSuperAdminEmailExists(ctx, superadminEmail)
	if err != nil {
		return nil, errors.New("")
	}
	if exists {
		return nil, errors.New("email id already exists")
	}

	hashPassword, err := utils.HashPassword(req.SuperAdminPassword)
	if err != nil {
		return nil, errors.New("")
	}
	superadminId := utils.GenerateUUID()
	superAdmin := super_admin.SuperAminCreateRequest{
		SuperAdminId:       superadminId,
		SuperAdminName:     req.SuperAdminName,
		SuperAdminEmail:    req.SuperAdminEmail,
		SuperAdminPassword: hashPassword,
	}
	if err := sa.Repo.CreateSuperAdmin(ctx, superAdmin); err != nil {
		return nil, errors.New("")
	}
	return &super_admin.SuperAdminCreateResponse{
		SuperAdminID: superadminId,
	},nil

}
