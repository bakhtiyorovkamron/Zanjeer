package postgres

import "github.com/Projects/Zanjeer/models"

type PostgresI interface {
	SetLocation(req []models.Record) error
	// Login(req models.Login) (models.LoginResponse, error)
	// CreateAdmin(req models.Admin) (models.Admin, error)
	// GetAdmins(req models.GetAdmins) ([]models.Admin, error)
	// CreateOTP(req models.SmsOtp) (models.Sms, error)
	// ConfirmOTP(req models.ConfirmOTP) error

	// //Driver represents
	// CreateDriver(req models.Driver) (models.Driver, error)
	// UpdateDriverInfo(req models.Driver) (models.Driver, error)
	// GetDriverInfo(id string) (models.Driver, error)
	// DeleteDriver(id string) error
	// GetDriverList(limit, offset int64) (models.DriverList, error)
	// SearchDriver(req models.DriverSearchRequest) (models.DriverList, error)

	// //Devices' Type
	// CreateDeviceType(req models.DeviceType) (models.DeviceType, error)
	// GetDeviceTypeList(req models.GetDeviceTypeListRequest) ([]models.DeviceType, error)

	// CreateDevice(req models.CreateDeviceRequest) (models.CreateDeviceResponse, error)
}
