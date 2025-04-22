package dao

import (
	"fmt"
	"mini-ups/db"
	"mini-ups/model"
)

// retrievs all the packages belonging to userID
func GetPackagesByUser(username string) ([]model.Package, error) {
	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	var packages []model.Package
	if err := db.DB.
		Where("username = ?", username).
		Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

// retrives the pacakge of the given packageID
func GetPackagesByPackageID(packageID string) (*model.Package, error) {
	var pack model.Package
	if err := db.DB.
		Where("id = ?", packageID).
		First(&pack).Error; err != nil {
		return nil, err
	}
	return &pack, nil
}

// Updates the delivery address to newCoord of the package with packageID
func UpdateDeliveryAddress(packageID string, newCoord model.Coordinate) (int64, error) {
	result := db.DB.Model(&model.Package{}).
		Where("id = ? AND status != ?", packageID, "out_for_delivery").
		Updates(model.Coordinate{
			X: newCoord.X,
			Y: newCoord.Y,
		})
	return result.RowsAffected, result.Error
}

// Creates a new Package
func CreatePackage(pack *model.Package) (model.PackageID, error) {
	result := db.DB.Create(pack)
	return pack.ID, result.Error
}

// Assigns Truck with TruckID to package with PackageID
func AssignTruckToPackage(packageID string, truckID uint) error {
	return db.DB.Model(&model.Package{}).
		Where("id = ?", packageID).
		Update("truck_id", truckID).Error
}

// Updates package with PackagedID to the new PacakgeStatus
func UpdatePackageStatus(packageID string, newStatus model.PackageStatus) error {
	return db.DB.Model(&model.Package{}).
		Where("id = ?", packageID).
		Update("status", newStatus).Error
}

// Gets the WarehouseID of the Pacakge
func GetWareHouseIDByPackage(packageID string) (uint, error) {
	var warehouseID uint
	err := db.DB.Model(&model.Package{}).
		Select("warehouse_id").
		Where("id = ?", packageID).
		Scan(&warehouseID).Error
	if err != nil {
		return 0, err
	}
	return warehouseID, nil
}
