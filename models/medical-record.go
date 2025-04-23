package models

type MedicalRecord struct {
	UserID          string    `json:"user_id" bson:"user_id"`
	MedicalID       string    `json:"medical_id" bson:"medical_id"`
	SentAt          string    `json:"sent_at" bson:"sent_at"`
	HeartBeat1      float64   `json:"1st_heart_beat" bson:"1st_heart_beat"`
	HeartBeat2      float64   `json:"2nd_heart_beat" bson:"2nd_heart_beat"`
	HeartBeat3      float64   `json:"3rd_heart_beat" bson:"3rd_heart_beat"`
	HeartBeat4      float64   `json:"4th_heart_beat" bson:"4th_heart_beat"`
	HeartBeat5      float64   `json:"5th_heart_beat" bson:"5th_heart_beat"`
	Conclusion      string    `json:"overall_conclusion" bson:"overall_conclusion"`
	HospitalTreated *string   `json:"hospital_treated,omitempty" bson:"hospital_treated,omitempty"`
	TreatedHospital *Hospital `json:"treated_hospital,omitempty" bson:"-"`
}
