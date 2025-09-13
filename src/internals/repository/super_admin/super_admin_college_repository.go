package superadmin_repo

import (
	"errors"
	"strconv"
	"strings"

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
func(repo *SuperAdminCollegeRepository)UpadateCollegeBalance(ctx context.Context,collegeId,newBalance string)error{
	filter := bson.M{"college_id":collegeId}
	update := bson.M{"$set":bson.M{"balance":newBalance}}
	res,err := repo.CollegeCollection.UpdateOne(ctx,filter,update)
	if err != nil{
		return errors.New("failed to update college balance ")
	}
	if res.MatchedCount == 0{
		return errors.New("college not found")
	}
	return nil
}
func(repo *SuperAdminCollegeRepository)GetCollegeById(ctx context.Context,collegeId string)(*super_admin.SuperAdminCollege,error){
	filter := bson.M{"college_id":collegeId}
	var college super_admin.SuperAdminCollege
	err := repo.CollegeCollection.FindOne(ctx,filter).Decode(&college)
	if err != nil{
		if errors.Is(err,mongo.ErrNoDocuments){
			return nil,errors.New("college not found")
		}
		return nil,errors.New("unable to retrive college details")
		
	}
	return &college,nil
}
func(repo *SuperAdminCollegeRepository) GetCollegesBySuperAdminID(ctx context.Context, adminID string) ([]super_admin.SuperAdminCollege, error) {
	filter := bson.M{"super_admin_id": adminID}
	cursor, err := repo.CollegeCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("failed to fetch colleges")
	}
	defer cursor.Close(ctx)

	var colleges []super_admin.SuperAdminCollege
	for cursor.Next(ctx) {
		var college super_admin.SuperAdminCollege
		if err := cursor.Decode(&college); err != nil {
			return nil, errors.New("error reading college data")
		}
		colleges = append(colleges, college)
	}

	if len(colleges) == 0 {
		return nil, errors.New("no colleges found for the given super admin")
	}

	return colleges, nil
}
func(repo *SuperAdminCollegeRepository)DeleteCollege(ctx context.Context,collegeId string)error{
	filter := bson.M{"college_id":collegeId}
	res,err:=repo.CollegeCollection.DeleteOne(ctx,filter)
	if err !=nil{
		return  errors.New("failed to delete college")
	}
	if res.DeletedCount == 0{
		return errors.New("college not found")
	}
	return nil
}
func(repo *SuperAdminCollegeRepository)RechargeCollege(ctx context.Context,recharge super_admin.CollegeRecharge)error{
	_,err:= repo.RechargeCollection.InsertOne(ctx,recharge)
	if err != nil{
		return errors.New("failed to record recharge transaction")
	}
	var college super_admin.SuperAdminCollege
	filter := bson.M{"college_id":recharge.CollegeId}
	err = repo.CollegeCollection.FindOne(ctx,filter).Decode(&college)
	if err !=nil{
		if errors.Is(err,mongo.ErrNoDocuments){
		return errors.New("college not found")
	}
	return errors.New("unable to retrive college details for recharge")
}
	oldAmount,err := strconv.Atoi(college.Balance)
	if err != nil{
		return errors.New("invalid balance formate in database") 
	}
	if strings.Contains(recharge.RechargeAmount,"."){
		return errors.New("recharge amount must be a whole number")
	}
	recharAmountInt,err := strconv.Atoi(recharge.RechargeAmount)
	if err != nil{
		return 	errors .New("invalid recharge amoount formate")
	}
	newBalanceInt := oldAmount+recharAmountInt
	if newBalanceInt < 10 {
		return errors.New("balance can not be less than 10")
	}
	newBalance := strconv.Itoa(newBalanceInt)
	update := bson.M{"$set":bson.M{"balance":newBalance}}
	_,err = repo.CollegeCollection.UpdateOne(ctx,filter,update)
	if err != nil{
		return errors.New("failed to update balance after recharge")
	}
	return nil
}
func(repo *SuperAdminCollegeRepository)GetRechargeHistoryByCollegeId(ctx context.Context,collegeId string)([]super_admin.CollegeRecharge,error){
	filter := bson.M{"college_id":collegeId}

	cursor,err := repo.RechargeCollection.Find(ctx,filter)
	if err != nil{
		return nil,errors.New("failed to retrive recharge cllection")
	}
	defer cursor.Close(ctx)

	var recherges []super_admin.CollegeRecharge
	for cursor.Next(ctx){
		var recharge super_admin.CollegeRecharge
		if err := cursor.Decode(&recharge); err!=nil{
			return nil, errors.New("error reading recharge history of data")
		}
		recherges = append(recherges,recharge)
	}
	if len(recherges) ==0{
		return nil,errors.New("no recharge history found for this college")
	}
	return recherges,nil
		
			
}
