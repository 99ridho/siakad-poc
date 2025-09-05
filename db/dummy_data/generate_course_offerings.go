package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Semester IDs from our generated data (30 semesters)
var semesterIDs = []string{
	"660e8400-e29b-41d4-a716-446655440001", "660e8400-e29b-41d4-a716-446655440002", "660e8400-e29b-41d4-a716-446655440003",
	"660e8400-e29b-41d4-a716-446655440004", "660e8400-e29b-41d4-a716-446655440005", "660e8400-e29b-41d4-a716-446655440006",
	"660e8400-e29b-41d4-a716-446655440007", "660e8400-e29b-41d4-a716-446655440008", "660e8400-e29b-41d4-a716-446655440009",
	"660e8400-e29b-41d4-a716-446655440010", "660e8400-e29b-41d4-a716-446655440011", "660e8400-e29b-41d4-a716-446655440012",
	"660e8400-e29b-41d4-a716-446655440013", "660e8400-e29b-41d4-a716-446655440014", "660e8400-e29b-41d4-a716-446655440015",
	"660e8400-e29b-41d4-a716-446655440016", "660e8400-e29b-41d4-a716-446655440017", "660e8400-e29b-41d4-a716-446655440018",
	"660e8400-e29b-41d4-a716-446655440019", "660e8400-e29b-41d4-a716-446655440020", "660e8400-e29b-41d4-a716-446655440021",
	"660e8400-e29b-41d4-a716-446655440022", "660e8400-e29b-41d4-a716-446655440023", "660e8400-e29b-41d4-a716-446655440024",
	"660e8400-e29b-41d4-a716-446655440025", "660e8400-e29b-41d4-a716-446655440026", "660e8400-e29b-41d4-a716-446655440027",
	"660e8400-e29b-41d4-a716-446655440028", "660e8400-e29b-41d4-a716-446655440029", "660e8400-e29b-41d4-a716-446655440030",
}

// Section codes for multiple sections of the same course
var sectionCodes = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

func generateCourseOfferingUUID(index int) string {
	return fmt.Sprintf("880e8400-e29b-41d4-a716-4466554%05d", index)
}

func generateCourseUUID(index int) string {
	return fmt.Sprintf("770e8400-e29b-41d4-a716-4466554%05d", index)
}

func randomChoice(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func generateCapacity() int {
	// Weighted capacity distribution
	weights := []struct {
		min, max int
		weight   int
	}{
		{20, 30, 15},   // Small classes
		{31, 50, 30},   // Medium classes
		{51, 80, 35},   // Large classes
		{81, 120, 15},  // Very large classes
		{121, 150, 5},  // Huge lectures
	}
	
	totalWeight := 0
	for _, w := range weights {
		totalWeight += w.weight
	}
	
	randVal := rand.Intn(totalWeight)
	cumulative := 0
	
	for _, w := range weights {
		cumulative += w.weight
		if randVal < cumulative {
			return w.min + rand.Intn(w.max-w.min+1)
		}
	}
	return 50 // fallback
}

func generateStartTime(semesterIndex int) string {
	// Map semester index to time periods
	year := 2015 + (semesterIndex / 3)
	semesterType := semesterIndex % 3
	
	var baseMonth, baseDay int
	switch semesterType {
	case 0: // Ganjil (August-December)
		baseMonth = 8
		baseDay = 1 + rand.Intn(30) // August 1-30
	case 1: // Genap (January-May)
		baseMonth = 1
		baseDay = 1 + rand.Intn(30) // January 1-30
	case 2: // Pendek (June-July)
		baseMonth = 6
		baseDay = 1 + rand.Intn(20) // June 1-20
	}
	
	// Add random hours for class start times (7 AM to 5 PM)
	hour := 7 + rand.Intn(11) // 7-17 (5 PM)
	minute := []int{0, 30}[rand.Intn(2)] // Either :00 or :30
	
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:00+00", year, baseMonth, baseDay, hour, minute)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	file, err := os.Create("04_course_offerings.sql")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()
	
	// Write header
	file.WriteString("-- Course Offerings Dummy Data (10,000 records)\n")
	file.WriteString("-- Course sections offered across different semesters with realistic scheduling\n\n")
	file.WriteString("INSERT INTO course_offerings (id, semester_id, course_id, section_code, capacity, start_time, created_at, updated_at, deleted_at) VALUES\n")
	
	// Track used combinations to avoid duplicates (semester_id, course_id, section_code)
	usedCombinations := make(map[string]bool)
	
	// Generate 10,000 course offerings
	recordCount := 0
	attempts := 0
	maxAttempts := 50000 // Prevent infinite loop
	
	for recordCount < 10000 && attempts < maxAttempts {
		attempts++
		
		// Select random semester and course
		semesterID := randomChoice(semesterIDs)
		courseIndex := 1 + rand.Intn(10000) // 1-10000 (matching our courses)
		courseID := generateCourseUUID(courseIndex)
		sectionCode := randomChoice(sectionCodes)
		
		// Create combination key
		combKey := fmt.Sprintf("%s-%s-%s", semesterID, courseID, sectionCode)
		
		// Skip if combination already exists
		if usedCombinations[combKey] {
			continue
		}
		usedCombinations[combKey] = true
		recordCount++
		
		uuid := generateCourseOfferingUUID(recordCount)
		capacity := generateCapacity()
		
		// Get semester index for start time generation
		semesterIndex := 0
		for i, sid := range semesterIDs {
			if sid == semesterID {
				semesterIndex = i
				break
			}
		}
		
		startTime := generateStartTime(semesterIndex)
		
		// Create timestamps (created before the semester starts)
		createdYear := 2015 + (semesterIndex / 3)
		createdMonth := []int{7, 12, 5}[semesterIndex%3] // Before each semester
		createdDay := 1 + (recordCount % 28)
		createdAt := fmt.Sprintf("%d-%02d-%02d 09:00:00+00", createdYear, createdMonth, createdDay)
		
		line := fmt.Sprintf("('%s', '%s', '%s', '%s', %d, '%s', '%s', '%s', NULL)",
			uuid, semesterID, courseID, sectionCode, capacity, startTime, createdAt, createdAt)
		
		if recordCount < 10000 {
			line += ","
		} else {
			line += ";"
		}
		
		file.WriteString(line + "\n")
		
		// Progress indicator
		if recordCount%1000 == 0 {
			fmt.Printf("Generated %d course offerings...\n", recordCount)
		}
	}
	
	if recordCount < 10000 {
		fmt.Printf("Warning: Only generated %d course offerings due to constraint limitations\n", recordCount)
	}
	
	// Write verification query
	file.WriteString("\n-- Verify insert count\n")
	file.WriteString("SELECT 'Course Offerings inserted: ' || COUNT(*) as summary FROM course_offerings;\n")
	
	fmt.Printf("Successfully generated 04_course_offerings.sql with %d course offering records!\n", recordCount)
}