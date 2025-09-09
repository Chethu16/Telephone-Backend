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
func(repo *SuperAdminCollegeRepository)GetCollegeForLogin(ctx context.Context,email string)(collegeID,collegeName,hashedPassword,superAdminID  string,err error){
	filter := bson.M{"college_email":email}
	var college super_admin.SuperAdminCollege
	err = repo.CollegeCollection.FindOne(ctx,filter).Decode(&college)
	if err !=nil{
		if errors.Is(err,mongo.ErrNoDocuments){
			return "","","","",errors.New("college not found")
		}
		return "","","","",errors.New("unable to retrive college login details")
		
	}
	return college.CollegeId,college.CollegeName,college.CollegePassword,college.SuperAdminId,nil
}

func(repo *SuperAdminCollegeRepository)GetCollegeBalance(ctx context.Context,collegeId string)(string,error){
	filter := bson.M{"college_id":collegeId}
	var college super_admin.SuperAdminCollege
	err:= repo.CollegeCollection.FindOne(ctx,filter).Decode(&college)
	if err != nil{
		if errors.Is(err,mongo.ErrNoDocuments){
			return "0",errors.New("college not found")
		}
		return "0",errors.New("unable to fetch college balance")
	}
	return college.Balance,nil
}