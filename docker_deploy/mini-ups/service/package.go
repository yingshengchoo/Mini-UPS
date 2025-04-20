package service

import (
	"mini-ups/dao"
	"mini-ups/model"
)

// Returns basic info of packages belonging to a user
func GetPackagesForUser(userID uint) ([]model.Package, error) {
	return dao.GetPackagesByUser(userID)
}

// Returns basic info of a single package
func GetPackageInfo(packageID uint) (*model.Package, error) {
	return dao.GetPackagesByPackageID(packageID)
}

// Attempts to update the delivery address; handles logic
func ChangePackageDestination(packageID uint, newCoord model.Coordinate) (string, error) {
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
func CreatePackage(pack *model.Package) error {
	return dao.AddPackage(pack)
}

// Assigns a truck to a package
func LinkTruckToPackage(packageID string, truckID uint) error {
	return dao.AssignTruckToPackage(packageID, truckID)
}

// Changes status of a package
func ChangePackageStatus(packageID string, newStatus model.PackageStatus) error {
	return dao.UpdatePackageStatus(packageID, newStatus)
}
