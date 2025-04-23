package service

import (
	"mini-ups/dao"
	"mini-ups/model"

	"gorm.io/datatypes"
)

// Returns basic info of packages belonging to a user
func GetPackagesForUser(username string) ([]model.Package, error) {
	return dao.GetPackagesByUser(username)
}

// Returns basic info of a single package
func GetPackageInfo(packageID string) (*model.Package, error) {
	return dao.GetPackagesByPackageID(packageID)
}

func GetPackageInfoByTruck(truckID int32) (*model.Package, error) {
	return dao.GetPackageInfoByTruckID(truckID)
}

// Attempts to update the delivery address; handles logic
func ChangePackageDestination(packageID string, newCoord model.Coordinate) (string, error) {
	rows, err := dao.UpdateDeliveryAddress(packageID, newCoord)
	if err != nil {
		return "", err
	}
	if rows == 0 {
		return "Package is already out for delivery and cannot be redirected.", nil
	}
	return "Package destination updated successfully.", nil
}

// Creates a new package
func CreatePackage(package_id string, username string, items datatypes.JSON, dest_x int, dest_y int, warehouse_id uint) (model.PackageID, error) {
	return dao.CreatePackage(model.NewPackage(
		model.PackageID(package_id),
		username,
		items,
		dest_x, dest_y,
		warehouse_id,
		model.StatusCreated,
	))
}

// Assigns a truck to a package
func LinkTruckToPackage(packageID string, truckID uint) error {
	return dao.AssignTruckToPackage(packageID, truckID)
}

// Changes status of a package
func ChangePackageStatus(packageID string, newStatus model.PackageStatus) error {
	return dao.UpdatePackageStatus(packageID, newStatus)
}

// retrieves warehouseID of the package with package_id
func GetWarehouseID(package_id string) (uint, error) {
	warehouse_id, err := dao.GetWareHouseIDByPackage(package_id)
	if err != nil {
		return 0, err //assuming 0 is not associated with any warehouse
	}
	return warehouse_id, nil
}
