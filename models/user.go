package models

type User struct {
	UserID        string         `json:"user_id" bson:"user_id"`
	Username      string         `json:"username" bson:"username"`
	Password      string         `json:"password" bson:"password"`
	HospitalID    string         `json:"hospital_id" bson:"hospital_id"`
	MedicalID     string         `json:"medical_id" bson:"medical_id"`
	Hospital      *Hospital      `json:"hospital,omitempty" bson:"-"`
	MedicalRecord *MedicalRecord `json:"medical_record,omitempty" bson:"-"`
	Devices       []*Device      `json:"devices,omitempty" bson:"-"`
}
