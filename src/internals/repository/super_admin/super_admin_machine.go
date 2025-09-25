package superadmin_repo

import (
	"errors"

	"github.com/Chethu16/Chethu/src/internals/domain/super_admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type SuperAdminMachineRepo struct {
	SuperAdminMachineCollection *mongo.Collection
	CollegeCollection           *mongo.Collection
}

func NewSuperAdminMachineRepository(db *mongo.Database) *SuperAdminMachineRepo {
	return &SuperAdminMachineRepo{
		SuperAdminMachineCollection: db.Collection("super_admin_machine"),
		CollegeCollection:           db.Collection("colleges"),
	}
}

func (repo *SuperAdminMachineRepo) CreateMachine(ctx context.Context, machine super_admin.Machine) error {
	filter := bson.M{
		"machine_no": machine.MachineNo,
		"college_id": machine.CollegeId,
	}

	var existing super_admin.Machine
	err := repo.SuperAdminMachineCollection.FindOne(ctx, filter).Decode(&existing)

	if err == nil {
		// Found an existing machine
		return errors.New("a machine with this number already exists for the selected college")
	}

	if err != mongo.ErrNoDocuments {
		// Some other DB error
		return errors.New("unable to check for existing machine. please try again later")
	}

	// Safe to insert
	machine.Balance = "0"
	_, err = repo.SuperAdminMachineCollection.InsertOne(ctx, machine)
	if err != nil {
		return errors.New("failed to create machine. please try again later")
	}
	return nil
}
func(repo *SuperAdminMachineRepo)GetAllMachinesByCollege(ctx context.Context,collegeId string)([]super_admin.Machine,error){
	filter := bson.M{"collegge_id":collegeId}

	cursor,err := repo.SuperAdminMachineCollection.Find(ctx,filter)
	if err !=nil{
		return nil,errors.New("unable to fetch machine list. please try again later")

	}
	defer cursor.Close(ctx)
	var machines []super_admin.Machine
	for cursor.Next(ctx){
		var machine super_admin.Machine
		if err := cursor.Decode(&machine);err != nil{
			return nil, errors.New("error reading machine data. please contact support")
		}
		machines = append(machines,machine)
	}
	if err := cursor.Err();err !=nil{
		return 	nil, errors.New("error while proccesing machine list . please try again later")
	}

	return machines,nil
}
func(repo *SuperAdminMachineRepo)GetCollageBalance(ctx context.Context,collgeId string)(string,error){
	filter := bson.M{"college_id":collgeId}
	var college super_admin.SuperAdminCollege
		err:=repo.CollegeCollection.FindOne(ctx,filter).Decode(&college)
		if err != nil{
		if err == mongo.ErrNoDocuments{
			return "", errors.New("no college found provided college id")
		}
		return "",errors.New("unable to retrive the collge balance . please try again later")
		}
		return college.Balance,nil
}
func(repo *SuperAdminMachineRepo)DeleteMachine(ctx context.Context,machineId string)error{
	filter := bson.M{"machine_id":machineId}

	res , err := repo.SuperAdminMachineCollection.DeleteOne(ctx,filter)
	if err !=nil{
		return errors.New("unable to delete the machine . please try again later")
	}
	if res.DeletedCount == 0{
		return errors.New("no machine found with provided machineId")
	}
	return nil
}
