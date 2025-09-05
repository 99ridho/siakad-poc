package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func generateRegistrationUUID(index int) string {
	return fmt.Sprintf("990e8400-e29b-41d4-a716-4466554%05d", index)
}

func generateCourseOfferingUUID(index int) string {
	return fmt.Sprintf("880e8400-e29b-41d4-a716-4466554%05d", index)
}

// Generate student user UUIDs (assuming they exist in the users table)
// This would need to be adjusted based on actual student user IDs in the system
func generateStudentUUID(index int) string {
	return fmt.Sprintf("aa0e8400-e29b-41d4-a716-4466554%05d", index)
}

func randomChoice(slice []int) int {
	return slice[rand.Intn(len(slice))]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	file, err := os.Create("05_course_registrations.sql")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()
	
	// Write header and note about student IDs
	file.WriteString("-- Course Registrations Dummy Data (10,000 records)\n")
	file.WriteString("-- NOTE: This script assumes student user IDs exist in the users table with role=3\n")
	file.WriteString("-- You may need to adjust student_id values based on your actual user data\n")
	file.WriteString("-- The script generates fictional student UUIDs that should be replaced with real ones\n\n")
	
	file.WriteString("-- First, let's create some sample student users if they don't exist\n")
	file.WriteString("-- (Skip this if you already have student users)\n")
	file.WriteString("INSERT INTO users (id, email, password, role, created_at, updated_at, deleted_at)\n")
	file.WriteString("SELECT \n")
	file.WriteString("    'aa0e8400-e29b-41d4-a716-4466554' || LPAD(generate_series(1, 5000)::text, 5, '0'),\n")
	file.WriteString("    'student' || generate_series(1, 5000) || '@university.ac.id',\n")
	file.WriteString("    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: 'password'\n")
	file.WriteString("    3, -- student role\n")
	file.WriteString("    '2015-01-01 10:00:00+00'::timestamptz + (random() * interval '9 years'),\n")
	file.WriteString("    '2015-01-01 10:00:00+00'::timestamptz + (random() * interval '9 years'),\n")
	file.WriteString("    NULL\n")
	file.WriteString("WHERE NOT EXISTS (\n")
	file.WriteString("    SELECT 1 FROM users WHERE id = 'aa0e8400-e29b-41d4-a716-4466554' || LPAD(generate_series(1, 5000)::text, 5, '0')\n")
	file.WriteString(");\n\n")
	
	file.WriteString("-- Now insert course registrations\n")
	file.WriteString("INSERT INTO course_registrations (id, student_id, course_offering_id, created_at, updated_at, deleted_at) VALUES\n")
	
	// Track used combinations to avoid duplicates (student_id, course_offering_id)
	usedCombinations := make(map[string]bool)
	
	// Generate 10,000 course registrations
	recordCount := 0
	attempts := 0
	maxAttempts := 50000 // Prevent infinite loop
	
	for recordCount < 10000 && attempts < maxAttempts {
		attempts++
		
		// Select random student (1-5000) and course offering (1-10000)
		studentIndex := 1 + rand.Intn(5000)  // 1-5000 students
		courseOfferingIndex := 1 + rand.Intn(10000) // 1-10000 course offerings
		
		studentID := generateStudentUUID(studentIndex)
		courseOfferingID := generateCourseOfferingUUID(courseOfferingIndex)
		
		// Create combination key
		combKey := fmt.Sprintf("%s-%s", studentID, courseOfferingID)
		
		// Skip if combination already exists
		if usedCombinations[combKey] {
			continue
		}
		usedCombinations[combKey] = true
		recordCount++
		
		uuid := generateRegistrationUUID(recordCount)
		
		// Create realistic registration timestamp
		// Registration typically happens before semester starts
		year := 2015 + rand.Intn(10) // 2015-2024
		month := []int{7, 12, 5}[rand.Intn(3)] // July, December, May (before semesters)
		day := 15 + rand.Intn(15) // Mid to end of month
		hour := 8 + rand.Intn(12)  // 8 AM to 8 PM
		minute := rand.Intn(60)
		
		createdAt := fmt.Sprintf("%d-%02d-%02d %02d:%02d:00+00", year, month, day, hour, minute)
		
		line := fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', NULL)",
			uuid, studentID, courseOfferingID, createdAt, createdAt)
		
		if recordCount < 10000 {
			line += ","
		} else {
			line += ";"
		}
		
		file.WriteString(line + "\n")
		
		// Progress indicator
		if recordCount%1000 == 0 {
			fmt.Printf("Generated %d course registrations...\n", recordCount)
		}
	}
	
	if recordCount < 10000 {
		fmt.Printf("Warning: Only generated %d course registrations due to constraint limitations\n", recordCount)
	}
	
	// Write verification queries
	file.WriteString("\n-- Verify insert counts\n")
	file.WriteString("SELECT 'Students inserted: ' || COUNT(*) as student_summary FROM users WHERE role = 3;\n")
	file.WriteString("SELECT 'Course Registrations inserted: ' || COUNT(*) as registration_summary FROM course_registrations;\n")
	file.WriteString("\n-- Show some enrollment statistics\n")
	file.WriteString("SELECT \n")
	file.WriteString("    'Average enrollments per student: ' || ROUND(AVG(enrollment_count), 2) as avg_enrollments_per_student\n")
	file.WriteString("FROM (\n")
	file.WriteString("    SELECT student_id, COUNT(*) as enrollment_count\n")
	file.WriteString("    FROM course_registrations\n")
	file.WriteString("    GROUP BY student_id\n")
	file.WriteString(") student_enrollments;\n")
	
	fmt.Printf("Successfully generated 05_course_registrations.sql with %d course registration records!\n", recordCount)
}