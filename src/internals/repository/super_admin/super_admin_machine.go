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
