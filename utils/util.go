package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cant hash user password %w", err)
	}
	hash := string(hashedPassword)
	return hash, nil
}

func CheckPassword(password1, password2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
}

// ConvertStringToString - returns string to string conversion func
func ConvertStringToString() func(strData string) (interface{}, error) {
	return func(strData string) (interface{}, error) {
		return strData, nil
	}
}

// ConvertStringToInt - returns string to int conversion func
func ConvertStringToInt() func(strData string) (interface{}, error) {
	return func(strData string) (interface{}, error) {
		return strconv.Atoi(strData)
	}
}

// ConvertStringToFloat - returns string to float conversion func
func ConvertStringToFloat() func(strData string) (interface{}, error) {
	return func(strData string) (interface{}, error) {
		return strconv.ParseFloat(strData, 64)
	}
}

// ConvertStringToTime - returns string to time conversion func. Time format: YYYY-MM-DD
func ConvertStringToTime() func(strData string) (interface{}, error) {
	return func(strData string) (interface{}, error) {
		layout := "2006-01-02"
		return time.Parse(layout, strData)
	}
}

// ConvertStringToBool - returns string to bool conversion func.
func ConvertStringToBool() func(strData string) (interface{}, error) {
	return func(strData string) (interface{}, error) {
		return strconv.ParseBool(strData)
	}
}