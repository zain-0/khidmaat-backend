package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zain-0/khidmaat-backend/config"
	"github.com/zain-0/khidmaat-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateHospital(c *gin.Context) {
	var hospital models.Hospital
	if err := c.ShouldBindJSON(&hospital); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hospital data"})
		return
	}

	hospital.HospitalID = uuid.New().String()

	fmt.Println("Hospital data:", hospital)

	_, err := config.DB.Collection("hospitals").InsertOne(context.TODO(), hospital)
	if err != nil {
		fmt.Println("Error creating hospital:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hospital"})
		return
	}

	c.JSON(http.StatusCreated, hospital)
}

func GetAllHospitals(c *gin.Context) {
	cursor, err := config.DB.Collection("hospitals").Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hospitals"})
		return
	}
	defer cursor.Close(context.TODO())

	var hospitals []models.Hospital
	if err := cursor.All(context.TODO(), &hospitals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse hospitals"})
		return
	}

	c.JSON(http.StatusOK, hospitals)
}

func GetHospitalByID(c *gin.Context) {
	id := c.Param("id")
	var hospital models.Hospital
	err := config.DB.Collection("hospitals").FindOne(context.TODO(), bson.M{"hospital_id": id}).Decode(&hospital)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hospital not found"})
		return
	}

	c.JSON(http.StatusOK, hospital)
}

func UpdateHospital(c *gin.Context) {
	id := c.Param("id")
	var update models.Hospital
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	filter := bson.M{"hospital_id": id}
	updateFields := bson.M{
		"$set": bson.M{
			"hospital_name": update.HospitalName,
			"location":      update.Location,
		},
	}

	_, err := config.DB.Collection("hospitals").UpdateOne(context.TODO(), filter, updateFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hospital"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hospital updated"})
}

func DeleteHospital(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Collection("hospitals").DeleteOne(context.TODO(), bson.M{"hospital_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hospital"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hospital deleted"})
}
