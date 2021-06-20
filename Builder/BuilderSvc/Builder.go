package BuilderSvc

import (
	"github.com/ARMmaster17/Captain/Builder/drivers"
	"github.com/ARMmaster17/Captain/Shared"
	"gorm.io/gorm"
)

func BuildPlane(db *gorm.DB, planeID int64) error {
	plane, formation, err := getPlaneData(db, planeID)
	if err != nil {
		return err
	}
	netID, err := drivers.BuildPlaneOnAnyProvider(plane, formation)
	if err != nil {
		return err
	}
	plane.NetID = netID
	result := db.Save(&plane)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func getPlaneData(db *gorm.DB, planeID int64) (Shared.Plane, Shared.Formation, error) {
	var planeObject Shared.Plane
	result := db.First(&planeObject, planeID)
	if result.Error != nil {
		return Shared.Plane{}, Shared.Formation{}, result.Error
	}
	var formationObject Shared.Formation
	result = db.First(&formationObject, planeObject.FormationID)
	if result.Error != nil {
		return Shared.Plane{}, Shared.Formation{}, result.Error
	}
	return planeObject, formationObject, nil
}
