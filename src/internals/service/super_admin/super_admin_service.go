package superadmin_service

import (
	"context"
	"errors"
	"fmt"
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
func (sa *SuperAdminService) SuperAdminLogin(ctx context.Context, req super_admin.SuperAdminLoginRequest) (*super_admin.SuperAdminLoginResponse, error) {
	if sa.Validate == nil {
		return nil, errors.New("validator not initialized")
	}

	// Validate input struct
	if err := sa.Validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	superadminemail := strings.ToLower(strings.TrimSpace(req.SuperAdminEmail))

	// Check if email exists
	exists, err := sa.Repo.CheckSuperAdminEmailExists(ctx, superadminemail)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if !exists {
		return nil, errors.New("invalid email or password") // ✅ hide whether it's email or password
	}

	// Fetch details for login
	superadminId, superadminName, hashedPassword, err := sa.Repo.CheckSuperAdminEmailForLogin(ctx, superadminemail)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch login details: %w", err)
	}

	// Compare password
	if err := utils.CheckPasswordHash(req.SuperAdminPassword, hashedPassword); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateSuperAdminToken(superadminemail, superadminName, superadminId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &super_admin.SuperAdminLoginResponse{
		Token: token,
	}, nil
}
