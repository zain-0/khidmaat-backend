// models/hospital.go
package models

type Hospital struct {
	HospitalID   string `json:"hospital_id" bson:"hospital_id"`
	HospitalName string `json:"hospital_name" bson:"hospital_name"`
	Location     string `json:"location" bson:"location"`
	Email        string `json:"email" bson:"email"`
}
