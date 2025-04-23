package models

type User struct {
	UserID     string    `json:"user_id" bson:"user_id"`
	Username   string    `json:"username" bson:"username"`
	Password   string    `json:"password" bson:"password"`
	HospitalID string    `json:"hospital_id" bson:"hospital_id"` // Stores only the HospitalID
	MedicalID  string    `json:"medical_id" bson:"medical_id"`
	Devices    []*Device `json:"devices,omitempty" bson:"-"` // Assuming this is still needed
}
