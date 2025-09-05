# SIAKAD Dummy Data

This directory contains dummy data for the SIAKAD (Student Information Academic System) database with approximately 10,000 records per table (except users).

## Files Overview

| File | Description | Records |
|------|-------------|---------|
| `01_academic_years.sql` | Academic years (2015-2024) | 10 |
| `02_semesters.sql` | Semesters (Ganjil, Genap, Pendek) | 30 |
| `03_courses.sql` | Indonesian university courses | 10,000 |
| `04_course_offerings.sql` | Course sections per semester | 10,000 |
| `05_course_registrations.sql` | Student enrollments + student users | 10,000 + 5,000 |
| `load_all.sql` | Master script to load all data | - |

## Data Characteristics

### Academic Years (10 records)
- Covers 2015/2016 through 2024/2025
- Indonesian academic calendar format
- Each year runs August to July

### Semesters (30 records)
- 3 semesters per academic year:
  - **Ganjil** (Odd): August - December
  - **Genap** (Even): January - May  
  - **Pendek** (Short/Summer): June - July

### Courses (10,000 records)
- Realistic Indonesian course names (Matematika, Fisika, etc.)
- Course codes follow Indonesian convention (MAT101, FIS201, etc.)
- Credits distributed: 1 SKS (10%), 2 SKS (40%), 3 SKS (35%), 4 SKS (15%)
- Subject areas: Math, Science, Engineering, Business, Social Sciences, etc.

### Course Offerings (10,000 records)
- Multiple sections per popular course (A, B, C, etc.)
- Capacity ranges: 20-150 students per section
- Realistic scheduling across different semesters
- Start times spread throughout academic periods

### Course Registrations (10,000 records)
- **Auto-generates 5,000 student users** (role=3) if they don't exist
- Student email format: `student1@university.ac.id` to `student5000@university.ac.id`
- Default password: `password` (bcrypt hashed)
- Realistic enrollment patterns
- No duplicate enrollments (student + course offering unique)

## How to Load Data

### Prerequisites
1. Ensure your PostgreSQL database is running
2. Have appropriate database credentials
3. Make sure the SIAKAD schema is created (run migrations first)

### Loading All Data
```bash
# Navigate to the dummy data directory
cd ./db/dummy_data

# Load all data using the master script
psql -d siakad_database -f load_all.sql

# Or specify connection parameters
psql -h localhost -U your_user -d siakad_database -f load_all.sql
```

### Loading Individual Tables
You can also load tables individually in dependency order:
```bash
psql -d siakad_database -f 01_academic_years.sql
psql -d siakad_database -f 02_semesters.sql
psql -d siakad_database -f 03_courses.sql
psql -d siakad_database -f 04_course_offerings.sql
psql -d siakad_database -f 05_course_registrations.sql
```

## Data Validation

After loading, you can verify the data with these queries:

```sql
-- Check record counts
SELECT 'Academic Years' as table_name, COUNT(*) as records FROM academic_years
UNION ALL SELECT 'Semesters', COUNT(*) FROM semesters  
UNION ALL SELECT 'Courses', COUNT(*) FROM courses
UNION ALL SELECT 'Course Offerings', COUNT(*) FROM course_offerings
UNION ALL SELECT 'Students', COUNT(*) FROM users WHERE role = 3
UNION ALL SELECT 'Registrations', COUNT(*) FROM course_registrations;

-- Check data integrity
SELECT 'Orphaned Semesters' as check_type, COUNT(*) as issues
FROM semesters s LEFT JOIN academic_years ay ON s.academic_year_id = ay.id 
WHERE ay.id IS NULL;

-- View sample data
SELECT c.code, c.name, c.credit, co.section_code, co.capacity
FROM courses c 
JOIN course_offerings co ON c.id = co.course_id 
LIMIT 10;
```

## Using the Data

### For API Testing
- Use the generated student accounts to test enrollment endpoints
- Student emails: `student1@university.ac.id` to `student5000@university.ac.id`
- Password: `password`

### For Development
- Test course offering CRUD operations
- Test enrollment business logic (capacity limits, duplicates)
- Test pagination with large datasets
- Explore realistic academic scheduling scenarios

### Example Test Scenarios
```sql
-- Find courses with high enrollment
SELECT c.name, COUNT(*) as enrollments
FROM course_registrations cr
JOIN course_offerings co ON cr.course_offering_id = co.id
JOIN courses c ON co.course_id = c.id
GROUP BY c.id, c.name
ORDER BY enrollments DESC
LIMIT 10;

-- Find students with most enrollments
SELECT u.email, COUNT(*) as course_count
FROM course_registrations cr
JOIN users u ON cr.student_id = u.id
GROUP BY u.id, u.email
ORDER BY course_count DESC
LIMIT 10;

-- Check semester distribution
SELECT ay.code as academic_year, s.code as semester, COUNT(co.id) as offerings
FROM academic_years ay
JOIN semesters s ON ay.id = s.academic_year_id
JOIN course_offerings co ON s.id = co.semester_id
GROUP BY ay.code, s.code, ay.start_time, s.start_time
ORDER BY ay.start_time, s.start_time;
```

## Cleanup

To remove all dummy data while preserving schema:
```sql
-- Remove in reverse dependency order
DELETE FROM course_registrations;
DELETE FROM course_offerings;
DELETE FROM courses;
DELETE FROM semesters;
DELETE FROM academic_years;
DELETE FROM users WHERE role = 3; -- Only remove generated students
```

## Notes

- All generated UUIDs follow a predictable pattern for easy identification
- Timestamps are spread realistically across the time periods
- Data respects all foreign key constraints
- Student users are only created if they don't already exist
- The script handles duplicate prevention automatically

## Generated Files Info

- **Go Scripts**: Used to generate the large SQL files (`generate_*.go`)
- **SQL Files**: Ready-to-execute INSERT statements
- **Total Size**: Approximately 4MB of SQL data
- **Load Time**: ~30 seconds on typical hardware
- **Memory**: Minimal memory usage during generation