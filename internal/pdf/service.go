package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"

	"pdf-generator/internal/api"
)

type PDFGenerator struct {
	pdf *gofpdf.Fpdf
}

func NewPDFGenerator() *PDFGenerator {
	pdf := gofpdf.New("P", "mm", "A4", "")
	return &PDFGenerator{pdf: pdf}
}

// GenerateStudentReport creates a PDF report for a student
func (pg *PDFGenerator) GenerateStudentReport(student *api.Student) ([]byte, error) {
	pg.pdf.AddPage()

	pg.pdf.SetFont("Arial", "B", 20)

	pg.pdf.CellFormat(190, 10, "STUDENT REPORT", "0", 1, "C", false, 0, "")
	pg.pdf.Ln(5)

	pg.pdf.SetFont("Arial", "", 12)
	pg.pdf.CellFormat(190, 8, "School Management System", "0", 1, "C", false, 0, "")
	pg.pdf.CellFormat(190, 8, fmt.Sprintf("Generated on: %s", time.Now().Format("January 2, 2006")), "0", 1, "C", false, 0, "")
	pg.pdf.Ln(10)

	pg.addSectionHeader("PERSONAL INFORMATION")
	pg.addInfoRow("Student ID:", fmt.Sprintf("%d", student.ID))
	pg.addInfoRow("Full Name:", student.Name)
	pg.addInfoRow("Email:", student.Email)
	pg.addInfoRow("Phone:", api.GetValueOrNA(student.Phone))
	pg.addInfoRow("Gender:", api.GetValueOrNA(student.Gender))
	pg.addInfoRow("Date of Birth:", student.FormatDate(student.DOB))
	pg.pdf.Ln(5)

	pg.addSectionHeader("ACADEMIC INFORMATION")
	pg.addInfoRow("Class:", api.GetValueOrNA(student.Class))
	pg.addInfoRow("Section:", api.GetValueOrNA(student.Section))
	pg.addInfoRow("Roll Number:", api.GetIntValueOrNA(student.Roll))
	pg.addInfoRow("Admission Date:", student.FormatDate(student.AdmissionDate))
	pg.addInfoRow("System Access:", fmt.Sprintf("%t", student.SystemAccess))
	pg.pdf.Ln(5)

	pg.addSectionHeader("FAMILY INFORMATION")
	pg.addInfoRow("Father's Name:", api.GetValueOrNA(student.FatherName))
	pg.addInfoRow("Father's Phone:", api.GetValueOrNA(student.FatherPhone))
	pg.addInfoRow("Mother's Name:", api.GetValueOrNA(student.MotherName))
	pg.addInfoRow("Mother's Phone:", api.GetValueOrNA(student.MotherPhone))
	pg.addInfoRow("Guardian's Name:", api.GetValueOrNA(student.GuardianName))
	pg.addInfoRow("Guardian's Phone:", api.GetValueOrNA(student.GuardianPhone))
	pg.addInfoRow("Relation to Guardian:", api.GetValueOrNA(student.RelationOfGuardian))
	pg.pdf.Ln(5)

	pg.addSectionHeader("ADDRESS INFORMATION")
	pg.addInfoRow("Current Address:", api.GetValueOrNA(student.CurrentAddress))
	pg.addInfoRow("Permanent Address:", api.GetValueOrNA(student.PermanentAddress))
	pg.pdf.Ln(5)

	pg.addSectionHeader("ADDITIONAL INFORMATION")
	pg.addInfoRow("Reporter/Class Teacher:", api.GetValueOrNA(student.ReporterName))

	pg.pdf.Ln(20)
	pg.pdf.SetFont("Arial", "I", 10)
	pg.pdf.CellFormat(190, 8, "This report was generated automatically by the School Management System", "0", 1, "C", false, 0, "")
	pg.pdf.CellFormat(190, 8, "For any queries, please contact the school administration", "0", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pg.pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// addSectionHeader adds a formatted section header
func (pg *PDFGenerator) addSectionHeader(title string) {
	pg.pdf.SetFont("Arial", "B", 14)
	pg.pdf.SetFillColor(240, 240, 240)
	pg.pdf.CellFormat(190, 8, title, "1", 1, "L", true, 0, "")
	pg.pdf.Ln(2)
	pg.pdf.SetFont("Arial", "", 11)
}

// addInfoRow adds a formatted information row
func (pg *PDFGenerator) addInfoRow(label, value string) {
	pg.pdf.SetFont("Arial", "B", 11)
	pg.pdf.CellFormat(60, 6, label, "0", 0, "L", false, 0, "")
	pg.pdf.SetFont("Arial", "", 11)
	pg.pdf.CellFormat(130, 6, value, "0", 1, "L", false, 0, "")
}
