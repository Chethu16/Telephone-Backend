package superadmin_repo

import (
	"errors"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type SuperAdminCollegeRepository struct{
	CollegeCollection *mongo.Collection
	RechargeCollection *mongo.Collection
}

func NewSuperAdminCollegeRepository(db *mongo.Database)*SuperAdminCollegeRepository{
	return &SuperAdminCollegeRepository{
		CollegeCollection: db.Collection("colleges"),
		RechargeCollection: db.Collection("recharges"),

	}
}
func(repo *SuperAdminCollegeRepository)CheckCollegeEmailExists(ctx context.Context,email string)(bool,error){
	filter := bson.M{"college_email":email}
	count,err := repo.CollegeCollection.CountDocuments(ctx,filter)
	if err != nil{
		return false,errors.New("failed to check college email")
	}
	return count>0,nil
}
func(repo *SuperAdminCollegeRepository)CreateCollege(ctx context.Context,college super_admin.SuperAdminCollege)error{
	_,err := repo.CollegeCollection.InsertOne(ctx,college)
	if err != nil{
		return errors.New("failed to create college")
	}
	return nil
}