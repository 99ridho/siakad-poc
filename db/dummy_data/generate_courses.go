package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Indonesian course names and subjects
var courseNames = []string{
	"Matematika", "Fisika", "Kimia", "Biologi", "Bahasa Indonesia", "Bahasa Inggris",
	"Sejarah", "Geografi", "Ekonomi", "Sosiologi", "Antropologi", "Psikologi",
	"Filsafat", "Agama", "Pancasila", "Kewarganegaraan", "Seni", "Olahraga",
	"Teknologi Informasi", "Sistem Informasi", "Rekayasa Perangkat Lunak",
	"Jaringan Komputer", "Database", "Algoritma", "Struktur Data", "Pemrograman",
	"Web Development", "Mobile Development", "Machine Learning", "Data Mining",
	"Kecerdasan Buatan", "Computer Vision", "Natural Language Processing",
	"Cyber Security", "Blockchain", "Cloud Computing", "Internet of Things",
	"Manajemen", "Akuntansi", "Keuangan", "Pemasaran", "Operasional",
	"Sumber Daya Manusia", "Entrepreneur", "Statistik", "Penelitian Operasi",
	"Logistik", "Supply Chain", "Quality Control", "Project Management",
	"Hukum", "Politik", "Hubungan Internasional", "Komunikasi", "Jurnalistik",
	"Public Relations", "Advertising", "Media Digital", "Broadcasting",
	"Teknik Sipil", "Teknik Mesin", "Teknik Elektro", "Teknik Industri",
	"Arsitektur", "Planologi", "Lingkungan", "Pertanian", "Perkebunan",
	"Peternakan", "Perikanan", "Kehutanan", "Geologi", "Pertambangan",
	"Kedokteran", "Keperawatan", "Farmasi", "Kesehatan Masyarakat",
	"Gizi", "Fisioterapi", "Radiologi", "Laboratorium", "Dental",
}

var levelPrefixes = []string{
	"Dasar", "Lanjut", "Menengah", "Tingkat Lanjut", "Spesialisasi",
	"Pengantar", "Fundamental", "Aplikasi", "Praktikum", "Seminar",
	"Workshop", "Proyek", "Tugas Akhir", "Skripsi", "Thesis",
}

var courseCodes = []string{
	"MAT", "FIS", "KIM", "BIO", "BIN", "ENG", "SEJ", "GEO", "EKO", "SOS",
	"ANT", "PSI", "FIL", "AGM", "PAN", "PKN", "SEN", "OLH", "TIF", "SIF",
	"RPL", "JKO", "DBS", "ALG", "STD", "PRG", "WEB", "MOB", "MLN", "DMN",
	"KCB", "CVS", "NLP", "CYB", "BCH", "CLD", "IOT", "MNG", "AKT", "KEU",
	"PMR", "OPS", "SDM", "ENT", "STA", "PEO", "LOG", "SCM", "QCT", "PMT",
	"HUK", "POL", "HIN", "KOM", "JUR", "PRE", "ADV", "MED", "BRC", "TSI",
	"TME", "TEL", "TIN", "ARS", "PLN", "LIN", "PTN", "PKB", "PTR", "PIK",
	"KEH", "GEL", "PTB", "DOK", "KEP", "FAR", "KMA", "GIZ", "FIS", "RAD",
	"LAB", "DNT",
}

func generateUUID(index int) string {
	return fmt.Sprintf("770e8400-e29b-41d4-a716-4466554%05d", index)
}

func randomChoice(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func generateCourseName(index int) (string, string, int) {
	baseName := randomChoice(courseNames)
	
	// Add level prefix for variety (30% chance)
	if rand.Float32() < 0.3 {
		prefix := randomChoice(levelPrefixes)
		baseName = fmt.Sprintf("%s %s", prefix, baseName)
	}
	
	// Add number suffix for sequence (50% chance)
	if rand.Float32() < 0.5 {
		level := rand.Intn(4) + 1 // 1-4
		baseName = fmt.Sprintf("%s %d", baseName, level)
	}
	
	// Generate course code
	codePrefix := randomChoice(courseCodes)
	codeNumber := 100 + (index % 400) // 100-499 range
	courseCode := fmt.Sprintf("%s%d", codePrefix, codeNumber)
	
	// Generate credits (1-4 SKS, weighted toward 2-3)
	weights := []int{10, 40, 35, 15} // 1:10%, 2:40%, 3:35%, 4:15%
	cumulative := 0
	randVal := rand.Intn(100)
	credits := 1
	for i, weight := range weights {
		cumulative += weight
		if randVal < cumulative {
			credits = i + 1
			break
		}
	}
	
	return courseCode, baseName, credits
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	file, err := os.Create("03_courses.sql")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()
	
	// Write header
	file.WriteString("-- Courses Dummy Data (10,000 records)\n")
	file.WriteString("-- Indonesian university courses with realistic names and credits\n\n")
	file.WriteString("INSERT INTO courses (id, code, name, credit, created_at, updated_at, deleted_at) VALUES\n")
	
	// Generate 10,000 courses
	for i := 1; i <= 10000; i++ {
		uuid := generateUUID(i)
		code, name, credit := generateCourseName(i)
		
		// Create timestamps (vary creation time over the years)
		year := 2015 + (i % 10) // 2015-2024
		month := 1 + (i % 12)   // 1-12
		day := 1 + (i % 28)     // 1-28
		
		createdAt := fmt.Sprintf("%d-%02d-%02d 10:00:00+00", year, month, day)
		
		// Escape single quotes in names
		name = strings.ReplaceAll(name, "'", "''")
		
		line := fmt.Sprintf("('%s', '%s', '%s', %d, '%s', '%s', NULL)",
			uuid, code, name, credit, createdAt, createdAt)
		
		if i < 10000 {
			line += ","
		} else {
			line += ";"
		}
		
		file.WriteString(line + "\n")
		
		// Progress indicator
		if i%1000 == 0 {
			fmt.Printf("Generated %d courses...\n", i)
		}
	}
	
	// Write verification query
	file.WriteString("\n-- Verify insert count\n")
	file.WriteString("SELECT 'Courses inserted: ' || COUNT(*) as summary FROM courses;\n")
	
	fmt.Println("Successfully generated 03_courses.sql with 10,000 course records!")
}