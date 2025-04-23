// models/device.go
package models

type Device struct {
	DeviceID   string `json:"device_id" bson:"device_id"`
	DeviceName string `json:"device_name" bson:"device_name"`
	UserID     string `json:"user_id" bson:"user_id"`
}
