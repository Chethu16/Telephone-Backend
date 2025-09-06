package superadmin_service

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
	// Safety: if caller forgot to pass validator, initialize it
	if v == nil {
		v = validator.New()
	}
	return &SuperAdminService{
		Repo:     r,
		Validate: v,
	}
}

func (sa *SuperAdminService) CreateSuperAdmin(ctx context.Context, req super_admin.SuperAminCreateRequest) (*super_admin.SuperAdminCreateResponse, error) {
	// Validate request
	if sa.Validate == nil {
		return nil, errors.New("validator not initialized")
	}
	if err := sa.Validate.Struct(req); err != nil {
		return nil, err // return actual validation errors
	}

	// Normalize email
	superadminEmail := strings.ToLower(strings.TrimSpace(req.SuperAdminEmail))

	// Check if email exists
	exists, err := sa.Repo.CheckSuperAdminEmailExists(ctx, superadminEmail)
	if err != nil {
		return nil, errors.New("failed to check if email exists")
	}
	if exists {
		return nil, errors.New("email id already exists")
	}

	// Hash password
	hashPassword, err := utils.HashPassword(req.SuperAdminPassword)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate UUID
	superadminId := utils.GenerateUUID()

	// Prepare object for DB
	superAdmin := super_admin.SuperAminCreateRequest{
		SuperAdminId:       superadminId,
		SuperAdminName:     req.SuperAdminName,
		SuperAdminEmail:    superadminEmail,
		SuperAdminPassword: hashPassword,
	}

	// Save to repository
	if err := sa.Repo.CreateSuperAdmin(ctx, superAdmin); err != nil {
		return nil, errors.New("failed to create super admin")
	}

	// Return response
	return &super_admin.SuperAdminCreateResponse{
		SuperAdminID: superadminId,
	}, nil
}
