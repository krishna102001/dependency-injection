package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	UserID      bson.ObjectID `bson:"_id,omitempty" json:"-"`
	FName       string        `bson:"fname" json:"fname"`
	MName       string        `bson:"mname" json:"mname,omitempty"`
	LName       string        `bson:"lname" json:"lname"`
	DOB         string        `bson:"dob" json:"dob"`
	Phone       string        `bson:"phone" json:"phone"`
	CountryCode string        `bson:"country_code" json:"country_code"`
	Password    string        `bson:"password" json:"password"`
	Email       string        `bson:"email" json:"email"`
	IsVerified  bool          `bson:"is_verified" json:"is_verified"`
}
