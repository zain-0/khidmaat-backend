package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zain-0/khidmaat-backend/config"
	"github.com/zain-0/khidmaat-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// SignUp function with password hashing
func SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	if user.Username == "" || user.Password == "" || user.HospitalID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, password, and hospital_id are required"})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	user.UserID = uuid.New().String() // Generate unique ID

	// Validate if hospital exists
	var hospital models.Hospital
	if err := config.DB.Collection("hospitals").FindOne(c, bson.M{"hospital_id": user.HospitalID}).Decode(&hospital); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		return
	}

	// No need to attach the full hospital details to the user
	// We just store the HospitalID in the User

	// Create the user in the database
	collection := config.DB.Collection("users")
	_, err = collection.InsertOne(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login function with password comparison
func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	collection := config.DB.Collection("users")
	filter := bson.D{{Key: "username", Value: input.Username}}
	err := collection.FindOne(c, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// GetUserWithDetails function (unchanged)
func GetUserWithDetails(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	collection := config.DB.Collection("users")
	filter := bson.D{{Key: "_id", Value: id}}
	err := collection.FindOne(c, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		return
	}

	// Fetch related hospital
	var hospital models.Hospital
	hospitalCollection := config.DB.Collection("hospitals")
	hospitalFilter := bson.D{{Key: "hospital_id", Value: user.HospitalID}} // Use HospitalID
	err = hospitalCollection.FindOne(c, hospitalFilter).Decode(&hospital)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching hospital"})
		return
	}

	// Fetch related medical record
	var medicalRecord models.MedicalRecord
	medicalRecordCollection := config.DB.Collection("medical_records")
	medicalRecordFilter := bson.D{{Key: "medical_id", Value: user.MedicalID}} // Use MedicalID
	err = medicalRecordCollection.FindOne(c, medicalRecordFilter).Decode(&medicalRecord)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching medical record"})
		return
	}

	// Query the devices associated with the user
	deviceCollection := config.DB.Collection("devices")
	deviceFilter := bson.D{{Key: "user_id", Value: id}}
	cursor, err := deviceCollection.Find(c, deviceFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching devices"})
		return
	}
	defer cursor.Close(c)
	var devicePointers []*models.Device
	for cursor.Next(c) {
		var device models.Device
		if err := cursor.Decode(&device); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding device"})
			return
		}
		devicePointers = append(devicePointers, &device)
	}

	// Return user along with related details (hospital, medical record, and devices)
	c.JSON(http.StatusOK, gin.H{
		"user":           user,
		"hospital":       hospital,
		"medical_record": medicalRecord,
		"devices":        devicePointers,
	})
}

// GetUsersByQuery function (unchanged)
func GetUsersByQuery(c *gin.Context) {
	hospitalID := c.Query("hospital_id")
	deviceID := c.Query("device_id")

	var users []models.User
	var queryBuilder bson.A

	// Build the query filter
	if hospitalID != "" {
		queryBuilder = append(queryBuilder, bson.D{{Key: "hospital_id", Value: hospitalID}})
	}
	if deviceID != "" {

		// Find all devices by device ID
		deviceCollection := config.DB.Collection("devices")
		deviceFilter := bson.D{{Key: "device_id", Value: deviceID}}
		cursor, err := deviceCollection.Find(c, deviceFilter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching devices"})
			return
		}
		defer cursor.Close(c)
		var userIDs []string
		for cursor.Next(c) {
			var device models.Device
			if err := cursor.Decode(&device); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding device"})
				return
			}
			userIDs = append(userIDs, device.UserID)
		}
		if len(userIDs) > 0 {
			queryBuilder = append(queryBuilder, bson.D{{Key: "user_id", Value: bson.D{{Key: "$in", Value: userIDs}}}})
		}
	}

	if len(queryBuilder) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one filter (hospital_id or device_id) is required"})
		return
	}

	// Query users based on the constructed filter
	usersCollection := config.DB.Collection("users")
	cursor, err := usersCollection.Find(c, bson.D{{Key: "$and", Value: queryBuilder}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}
	defer cursor.Close(c)
	for cursor.Next(c) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
