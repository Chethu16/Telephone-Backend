package superadmin_repo

import (
	"context"
	"errors"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SuperAdminRepository struct {
	SuperAdminCollection *mongo.Collection
}

func NewSuperAdminRepository(db *mongo.Database) *SuperAdminRepository {
	return &SuperAdminRepository{
		SuperAdminCollection: db.Collection("super_admin"),
	}
}
func (repo *SuperAdminRepository) CheckSuperAdminExists(ctx context.Context, superadminid string) (bool, error) {
	filter := bson.M{"super_admin_id": superadminid}
	count, err := repo.SuperAdminCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.New("unable to find superadmin")
	}
	return count > 0, nil
}
func (repo *SuperAdminRepository) CheckSuperAdminEmailExists(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"super_admin_email": email}
	count, err := repo.SuperAdminCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.New("unable to find email")
	}
	return count > 0, nil
}
func (repo *SuperAdminRepository) CreateSuperAdmin(ctx context.Context, req super_admin.SuperAminCreateRequest) error {
	_, err := repo.SuperAdminCollection.InsertOne(ctx, bson.M{
		"super_admin_id":       req.SuperAdminId,
		"super_admin_name":     req.SuperAdminName,
		"super_admin_email":    req.SuperAdminEmail,
		"super_admin_password": req.SuperAdminPassword,
	})
	if err != nil {
		return errors.New("failed to connect super admin")
	}
	return nil
}
func (repo *SuperAdminRepository) CheckSuperAdminEmailForLogin(ctx context.Context, superadminEmail string) (string, string, string, error) {
	filter := bson.M{"super_admin_email": superadminEmail}
	var superadmindetails super_admin.SuperAdmin
	err := repo.SuperAdminCollection.FindOne(ctx, filter).Decode(&superadmindetails)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", "", errors.New("")
		}
		return "", "", "", errors.New("")
	}
	return superadmindetails.SuperAdminID, superadmindetails.SuperAdminName, superadmindetails.SuperAdminPassword, nil
}
