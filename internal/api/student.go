package api

import (
	"fmt"
	"time"
)

// Student represents the student data structure from the Node.js API
type Student struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Email              string  `json:"email"`
	SystemAccess       bool    `json:"systemAccess"`
	Phone              *string `json:"phone"`
	Gender             *string `json:"gender"`
	DOB                *string `json:"dob"`
	Class              *string `json:"class"`
	Section            *string `json:"section"`
	Roll               *int    `json:"roll"`
	FatherName         *string `json:"fatherName"`
	FatherPhone        *string `json:"fatherPhone"`
	MotherName         *string `json:"motherName"`
	MotherPhone        *string `json:"motherPhone"`
	GuardianName       *string `json:"guardianName"`
	GuardianPhone      *string `json:"guardianPhone"`
	RelationOfGuardian *string `json:"relationOfGuardian"`
	CurrentAddress     *string `json:"currentAddress"`
	PermanentAddress   *string `json:"permanentAddress"`
	AdmissionDate      *string `json:"admissionDate"`
	ReporterName       *string `json:"reporterName"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (s *Student) FormatDate(dateStr *string) string {
	if dateStr == nil || *dateStr == "" {
		return "N/A"
	}

	if t, err := time.Parse("2006-01-02", *dateStr); err == nil {
		return t.Format("January 2, 2006")
	}

	return *dateStr
}

// GetValueOrNA returns the string value or "N/A" if nil
func GetValueOrNA(str *string) string {
	if str == nil || *str == "" {
		return "N/A"
	}
	return *str
}

// GetIntValueOrNA returns the int value as string or "N/A" if nil
func GetIntValueOrNA(num *int) string {
	if num == nil {
		return "N/A"
	}
	return fmt.Sprintf("%d", *num)
}
