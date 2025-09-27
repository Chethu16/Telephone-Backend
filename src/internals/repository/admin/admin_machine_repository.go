package admin_repo

import (
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MachineRechargeRepo struct {
	SuperAdminMachineCollection          *mongo.Collection
	CollegeCollection                    *mongo.Collection
	RechargeHistoryCollection            *mongo.Collection
	UserMachineRechargeCollection        *mongo.Collection
	UserMachineRechargeHistoryCollection *mongo.Collection
	UsersColletion                       *mongo.Collection
}

func NewMachineRechargeRepo(db *mongo.Database) *MachineRechargeRepo {
	return &MachineRechargeRepo{
		SuperAdminMachineCollection:          db.Collection("super_admin_machine"),
		CollegeCollection:                    db.Collection("colleges"),
		RechargeHistoryCollection:            db.Collection("machine_recharge_history"),
		UserMachineRechargeCollection:        db.Collection("user_machine_recharge_request"),
		UserMachineRechargeHistoryCollection: db.Collection("user_machine_recharge_history"),
		UsersColletion:                       db.Collection("users"),
	}
}
func(repo *MachineRechargeRepo)GetCollegeBalance(ctx context.Context,collegeId string)(string,error){
	var college struct{
		Balance string `bson:"balance"`
	}
	err :=repo.CollegeCollection.FindOne(ctx,bson.M{"college_id":collegeId}).Decode(&college)
	if err != nil{
		if errors.Is(err,mongo.ErrNoDocuments){
			return "",errors.New("college not found")
		}
		return "",errors.New("unable to fetch college balance at the moment")
	}
	if _,err := strconv.ParseFloat(college.Balance,64); err != nil{
		return "",errors.New("stored college balance is an invalid formate")
	}
	return college.Balance,nil
}

