package models

type MedicalRecord struct {
	UserID          string    `json:"user_id" bson:"user_id"`
	MedicalID       string    `json:"medical_id" bson:"medical_id"`
	SentAt          string    `json:"sent_at" bson:"sent_at"`
	HeartBeats      []float64 `json:"heart_beats" bson:"heart_beats"`
	Conclusion      string    `json:"overall_conclusion" bson:"overall_conclusion"`
	HospitalTreated *string   `json:"hospital_treated,omitempty" bson:"hospital_treated,omitempty"`
	TreatedHospital *Hospital `json:"treated_hospital,omitempty" bson:"-"`
}
