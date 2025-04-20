package dao

import (
	"mini-ups/db"
	"mini-ups/model"
)

func GetPackagesByUser(userID uint) ([]model.Package, error) {
	var packages []model.Package
	if err := db.DB.
		Select("ID", "Items", "Destination_X", "Destination_Y", "Status").
		Where("user_id = ?", userID).
		Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

func GetPackagesByPackageID(packageID uint) (*model.Package, error) {
	var pack model.Package
	if err := db.DB.
		Select("ID", "Items", "Destination_X", "Destination_Y", "Status").
		Where("ID = ?", packageID).
		First(&pack).Error; err != nil {
		return nil, err
	}
	return &pack, nil
}

func UpdateDeliveryAddress(packageID uint, newX, newY uint) (int64, error) {
	result := db.DB.Model(&model.Package{}).
		Where("id = ? AND truck_status != ?", packageID, "out_for_delivery").
		Updates(map[string]interface{}{
			"Destination_X": newX,
			"Destination_Y": newY,
		})
	return result.RowsAffected, result.Error
}

func AddPackage(pack *model.Package) error {
	return db.DB.Create(pack).Error
}

func AssignTruckToPackage(packageID string, truckID uint) error {
	return db.DB.Model(&model.Package{}).
		Where("id = ?", packageID).
		Update("truck_id", truckID).Error
}

func UpdatePackageStatus(packageID string, newStatus model.PackageStatus) error {
	return db.DB.Model(&model.Package{}).
		Where("id = ?", packageID).
		Update("status", newStatus).Error
}
