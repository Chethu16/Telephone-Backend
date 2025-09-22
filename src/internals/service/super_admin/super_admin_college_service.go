package superadmin_service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"strings"
	"time"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	superadmin_repo "github.com/Chethu16/Chethu/src/internals/repository/super_admin"
	"github.com/Chethu16/Chethu/src/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type SuperAdminCollegeService struct {
	Repo     *superadmin_repo.SuperAdminCollegeRepository
	Validate *validator.Validate
}

func NewSuperAdminCollegeService(r *superadmin_repo.SuperAdminCollegeRepository, v *validator.Validate) *SuperAdminCollegeService {
	if v == nil {
		v = validator.New()
	}
	return &SuperAdminCollegeService{
		Repo:     r,
		Validate: v,
	}
}
func (sa *SuperAdminCollegeService) CreateCollege(ctx context.Context, req super_admin.CollegeCreateRequest, SuperadminId string) (*super_admin.CollegeResponse, error) {
	if sa.Validate == nil {
		fmt.Println("")
		return nil, errors.New("validator not initialized")
	}
	if err := sa.Validate.Struct(req); err != nil {
		return nil, err
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
func (sa *SuperAdminCollegeService) CollegeLogin(ctx context.Context, req super_admin.CollegeLoginRequest) (*super_admin.CollegeTokenResponse, error) {
	if sa.Validate == nil {
		return nil, errors.New("validator not initialized")
	}
	if err := sa.Validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	collegeEmail := strings.ToLower(strings.TrimSpace(req.CollegeEmail))

	exists, err := sa.Repo.CheckCollegeEmailExists(ctx, collegeEmail)
	if err != nil {
		return nil, errors.New("failed to check email existence")
	}
	if !exists {
		return nil, errors.New("invalid email or password")
	}
	collegeID, name, hashPassword, superadminID, err := sa.Repo.GetCollegeForLogin(ctx, req.CollegeEmail)
	if err != nil {
		return nil, errors.New("college account not found")

	}
	if err := utils.CheckPasswordHash(req.CollegePassword, hashPassword); err != nil {
		return nil, errors.New("incorrect email or password")
	}
	balance, err := sa.Repo.GetCollegeBalance(ctx, collegeID)
	if err != nil {
		return nil, errors.New("unable to retrive account balance at the moment")
	}
	token, err := utils.GenerateCollegeToken(req.CollegeEmail, name, collegeID, superadminID)
	if err != nil {
		return nil, errors.New("failed to login, please try again")
	}
	return &super_admin.CollegeTokenResponse{
		Token:   token,
		Balance: balance,
	}, nil

}
func (sa *SuperAdminCollegeService) UpdateCollegeBalance(ctx context.Context, collegeId, amount string) error {
	amInt, err := strconv.Atoi(amount)
	if err != nil || amInt < 0 {
		return errors.New("recharge amount must be positive whole number")
	}
	curStr, err := sa.Repo.GetCollegeBalance(ctx, collegeId)
	if err != nil {
		return errors.New("unable to retrive current balance")
	}
	curInt, err := strconv.Atoi(curStr)
	if err != nil {
		return errors.New("account balance data invailde")
	}
	newBalance := curInt + amInt

	if err := sa.Repo.UpadateCollegeBalance(ctx, collegeId, strconv.Itoa(newBalance)); err != nil {
		return errors.New("unable to update alance at this time")
	}
	if _, err := sa.Repo.GetCollegeById(ctx, collegeId); err != nil {
		log.Printf("warning: could not fetch college for collegeid %s:%v", collegeId, err)
	}
	return nil
}
func (sa *SuperAdminCollegeService) GetCollegesBySuperadminId(ctx context.Context, adminId string) ([]super_admin.SuperAdminCollege, error) {
	if adminId == "" {
		return nil, errors.New("admin id required")
	}
	colleges, err := sa.Repo.GetCollegesBySuperAdminID(ctx, adminId)
	if err != nil {
		return nil, errors.New("unable to fetch colleges at this time")
	}
	return colleges, nil
}
func (sa *SuperAdminCollegeService) GetCollegeDetails(ctx context.Context, collegeID string) (*super_admin.SuperAdminCollege, error) {
	return sa.Repo.GetCollegeById(ctx, collegeID)
}
func (sa *SuperAdminCollegeService) DeleteCollege(ctx context.Context, collegeID string) error {

	err := sa.Repo.DeleteCollege(ctx, collegeID)
	if err != nil {
		return errors.New("failed to delete college")
	}
	return nil
}
func (sa *SuperAdminCollegeService) RechargeCollege(ctx context.Context, req super_admin.CollegeRechargeRequest) error {
	if err := sa.Validate.Struct(req); err != nil {
		return errors.New("please provide all recharge details correctly")
	}

	amtInt, err := strconv.Atoi(req.RechargeAmount)
	if err != nil || amtInt <= 0 {
		return errors.New("recharge amount must be greater than zero")
	}

	// 🔹 1. Check super_admin balance
	// superBalanceStr, err := sa.Repo.GetSuperAdminBalance(ctx, req.SuperAdminId)
	// if err != nil {
	// 	return errors.New("unable to retrieve super admin balance")
	// }
	// superBalance, err := strconv.Atoi(superBalanceStr)
	// if err != nil {
	// 	return errors.New("super admin balance data is invalid")
	// }
	// if superBalance < amtInt {
	// 	return errors.New("insufficient super admin balance for recharge")
	// }

	// 🔹 2. Check current college balance
	currStr, err := sa.Repo.GetCollegeBalance(ctx, req.CollegeId)
	if err != nil {
		return errors.New("unable to retrieve current college balance")
	}
	currInt, err := strconv.Atoi(currStr)
	if err != nil {
		return errors.New("college balance data is invalid")
	}
	newBalance := strconv.Itoa(currInt + amtInt)

	// 🔹 3. Save recharge history
	recharge := super_admin.CollegeRecharge{
		RechargeId:     utils.GenerateUUID(),
		CollegeId:      req.CollegeId,
		SuperAdminId:   req.SuperAdminId,
		RechargeAmount: req.RechargeAmount,
		RechargedAt:    time.Now().Format(time.RFC3339),
	}
	if err := sa.Repo.RechargeCollege(ctx, recharge, req.SuperAdminId); err != nil {
		return errors.New("failed to process recharge")
	}

	// 🔹 4. Update college balance
	if err := sa.Repo.UpadateCollegeBalance(ctx, req.CollegeId, newBalance); err != nil {
		return errors.New("failed to update college balance")
	}

	// // 🔹 5. Deduct from super admin balance
	// newSuperBalance := strconv.Itoa(superBalance - amtInt)
	// if err := sa.Repo.UpdateSuperAdminBalance(ctx, req.SuperAdminId, newSuperBalance); err != nil {
	// 	return errors.New("failed to update super admin balance")
	// }

	return nil
}
func (sa *SuperAdminCollegeService)GetRechargeHistory(ctx context.Context,collegeId string)([]super_admin.CollegeRecharge,error){
	if collegeId == ""{
		return nil,errors.New("college Id is required")
	}
	recharges,err := sa.Repo.GetRechargeHistoryByCollegeId(ctx,collegeId)
	if err !=nil{
		return nil,errors.New("unable to retrive recharge history")
	}	
	return recharges,nil

}
