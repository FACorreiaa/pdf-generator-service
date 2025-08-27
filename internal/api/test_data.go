package api

// GetMockStudent returns mock student data for testing
func GetMockStudent() *Student {
	name := "John Doe"
	email := "john.doe@school.com"
	phone := "+1234567890"
	gender := "Male"
	dob := "1995-05-15"
	class := "Grade 10"
	section := "A"
	roll := 15
	fatherName := "Robert Doe"
	fatherPhone := "+1234567891"
	motherName := "Jane Doe"
	motherPhone := "+1234567892"
	guardianName := "Robert Doe"
	guardianPhone := "+1234567891"
	relationOfGuardian := "Father"
	currentAddress := "123 Main St, City, State 12345"
	permanentAddress := "123 Main St, City, State 12345"
	admissionDate := "2020-09-01"
	reporterName := "Ms. Sarah Johnson"

	return &Student{
		ID:                 1,
		Name:               name,
		Email:              email,
		SystemAccess:       true,
		Phone:              &phone,
		Gender:             &gender,
		DOB:                &dob,
		Class:              &class,
		Section:            &section,
		Roll:               &roll,
		FatherName:         &fatherName,
		FatherPhone:        &fatherPhone,
		MotherName:         &motherName,
		MotherPhone:        &motherPhone,
		GuardianName:       &guardianName,
		GuardianPhone:      &guardianPhone,
		RelationOfGuardian: &relationOfGuardian,
		CurrentAddress:     &currentAddress,
		PermanentAddress:   &permanentAddress,
		AdmissionDate:      &admissionDate,
		ReporterName:       &reporterName,
	}
}
