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

type SuperAdminCollegeRepository struct {
	CollegeCollection    *mongo.Collection
	RechargeCollection   *mongo.Collection
	SuperAdminCollection *mongo.Collection
}

func NewSuperAdminCollegeRepository(db *mongo.Database) *SuperAdminCollegeRepository {
	return &SuperAdminCollegeRepository{
		CollegeCollection:    db.Collection("colleges"),
		RechargeCollection:   db.Collection("recharges"),
		SuperAdminCollection: db.Collection("super_admin"),
	}
}
func (repo *SuperAdminCollegeRepository) CheckCollegeEmailExists(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"college_email": email}
	count, err := repo.CollegeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.New("failed to check college email")
	}
	return count > 0, nil
}
func (repo *SuperAdminCollegeRepository) CreateCollege(ctx context.Context, college super_admin.SuperAdminCollege) error {
	_, err := repo.CollegeCollection.InsertOne(ctx, college)
	if err != nil {
		return errors.New("failed to create college")
	}
	return nil
}
func (repo *SuperAdminCollegeRepository) GetCollegeForLogin(ctx context.Context, email string) (collegeID, collegeName, hashedPassword, superAdminID string, err error) {
	filter := bson.M{"college_email": email}
	var college super_admin.SuperAdminCollege
	err = repo.CollegeCollection.FindOne(ctx, filter).Decode(&college)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", "", "", "", errors.New("college not found")
		}
		return "", "", "", "", errors.New("unable to retrive college login details")

	}
	return college.CollegeId, college.CollegeName, college.CollegePassword, college.SuperAdminId, nil
}

func (repo *SuperAdminCollegeRepository) GetCollegeBalance(ctx context.Context, collegeId string) (string, error) {
	filter := bson.M{"college_id": collegeId}
	var college super_admin.SuperAdminCollege
	err := repo.CollegeCollection.FindOne(ctx, filter).Decode(&college)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "0", errors.New("college not found")
		}
		return "0", errors.New("unable to fetch college balance")
	}
	return college.Balance, nil
}
func (repo *SuperAdminCollegeRepository) UpadateCollegeBalance(ctx context.Context, collegeId, newBalance string) error {
	filter := bson.M{"college_id": collegeId}
	update := bson.M{"$set": bson.M{"balance": newBalance}}
	res, err := repo.CollegeCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to update college balance ")
	}
	if res.MatchedCount == 0 {
		return errors.New("college not found")
	}
	return nil
}
func (repo *SuperAdminCollegeRepository) GetCollegeById(ctx context.Context, collegeId string) (*super_admin.SuperAdminCollege, error) {
	filter := bson.M{"college_id": collegeId}
	var college super_admin.SuperAdminCollege
	err := repo.CollegeCollection.FindOne(ctx, filter).Decode(&college)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("college not found")
		}
		return nil, errors.New("unable to retrive college details")

	}
	return &college, nil
}
func (repo *SuperAdminCollegeRepository) GetCollegesBySuperAdminID(ctx context.Context, adminID string) ([]super_admin.SuperAdminCollege, error) {
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
func (repo *SuperAdminCollegeRepository) DeleteCollege(ctx context.Context, collegeId string) error {
	filter := bson.M{"college_id": collegeId}
	res, err := repo.CollegeCollection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("failed to delete college")
	}
	if res.DeletedCount == 0 {
		return errors.New("college not found")
	}
	return nil
}
func (repo *SuperAdminCollegeRepository) RechargeCollege(ctx context.Context, recharge super_admin.CollegeRecharge, superAdminId string) error {
	superAdminFilter := bson.M{"super_admin_id": superAdminId}
	var superAdmin super_admin.SuperAdmin
	err := repo.CollegeCollection.Database().Collection("super_admin").FindOne(ctx, superAdminFilter).Decode(&superAdmin)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("superadmin not found")
		}
		return errors.New("unable to fetch superadmin balance")
	}

	// 2. Validate recharge amount (must be integer, not negative or float)
	if strings.Contains(recharge.RechargeAmount, ".") {
		return errors.New("recharge amount must be a whole number")
	}

	rechargeAmountInt, err := strconv.Atoi(recharge.RechargeAmount)
	if err != nil || rechargeAmountInt <= 0 {
		return errors.New("invalid recharge amount, must be positive integer")
	}

	// // 3. Check superadmin balance
	// superBalanceInt, err := strconv.Atoi(superAdmin.Balance)
	// if err != nil {
	// 	return errors.New("invalid superadmin balance format in database")
	// }
	// if rechargeAmountInt > superBalanceInt {
	// 	return errors.New("insufficient superadmin balance for recharge")
	// }

	// 4. Fetch college
	var college super_admin.SuperAdminCollege
	collegeFilter := bson.M{"college_id": recharge.CollegeId}
	err = repo.CollegeCollection.FindOne(ctx, collegeFilter).Decode(&college)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("college not found")
		}
		return errors.New("unable to retrieve college details")
	}

	// 5. Calculate new balances
	oldCollegeBalanceInt, err := strconv.Atoi(college.Balance)
	if err != nil {
		return errors.New("invalid college balance format in database")
	}

	newCollegeBalance := strconv.Itoa(oldCollegeBalanceInt + rechargeAmountInt)
	// newSuperBalance := strconv.Itoa(superBalanceInt - rechargeAmountInt)

	// 6. Update both balances in DB
	_, err = repo.CollegeCollection.UpdateOne(ctx, collegeFilter, bson.M{"$set": bson.M{"balance": newCollegeBalance}})
	if err != nil {
		return errors.New("failed to update college balance")
	}

	// _, err = repo.CollegeCollection.Database().Collection("super_admin").UpdateOne(ctx, superAdminFilter, bson.M{"$set": bson.M{"balance": newSuperBalance}})
	// if err != nil {
	// 	return errors.New("failed to update superadmin balance")
	// }

	// 7. Record recharge transaction
	_, err = repo.RechargeCollection.InsertOne(ctx, recharge)
	if err != nil {
		return errors.New("failed to record recharge transaction")
	}

	return nil
}

// func (repo *SuperAdminCollegeRepository) GetSuperAdminBalance(ctx context.Context, superAdminID string) (string, error) {
//     filter := bson.M{"super_admin_id": superAdminID}
//     var superAdmin super_admin.SuperAdmin

//     err := repo.SuperAdminCollection.FindOne(ctx, filter).Decode(&superAdmin)
//     if err != nil {
//         if errors.Is(err, mongo.ErrNoDocuments) {
//             return "0", errors.New("superadmin not found")
//         }
//         return "0", errors.New("unable to fetch superadmin balance")
//     }

//	    return superAdmin.Balance, nil
//	}
//
//	func (repo *SuperAdminCollegeRepository) UpdateSuperAdminBalance(ctx context.Context, superAdminID, newBalance string) error {
//		filter := bson.M{"super_admin_id": superAdminID}
//		update := bson.M{"$set": bson.M{"balance": newBalance}}
//		res, err := repo.SuperAdminCollection.UpdateOne(ctx, filter, update)
//		if err != nil {
//			return errors.New("failed to update superadmin balance")
//		}
//		if res.MatchedCount == 0 {
//			return errors.New("superadmin not found")
//		}
//		return nil
//	}
func (repo *SuperAdminCollegeRepository) GetRechargeHistoryByCollegeId(ctx context.Context, collegeId string) ([]super_admin.CollegeRecharge, error) {
	filter := bson.M{"college_id": collegeId}

	cursor, err := repo.RechargeCollection.Find(ctx, filter)
	if err != nil {
		return nil, errors.New("failed to retrive recharge cllection")
	}
	defer cursor.Close(ctx)

	var recherges []super_admin.CollegeRecharge
	for cursor.Next(ctx) {
		var recharge super_admin.CollegeRecharge
		if err := cursor.Decode(&recharge); err != nil {
			return nil, errors.New("error reading recharge history of data")
		}
		recherges = append(recherges, recharge)
	}
	if len(recherges) == 0 {
		return nil, errors.New("no recharge history found for this college")
	}
	return recherges, nil

}
