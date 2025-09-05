-- Master Script to Load All Dummy Data
-- Execute this script to load approximately 10,000 records per table
-- Total records: ~40,000+ (excluding users which are handled separately)

-- Disable foreign key checks temporarily for faster loading (PostgreSQL)
SET session_replication_role = replica;

BEGIN;

\echo 'Loading dummy data for SIAKAD system...'
\echo ''

-- 1. Load Academic Years (10 records)
\echo 'Loading Academic Years...'
\i 01_academic_years.sql
\echo ''

-- 2. Load Semesters (30 records) 
\echo 'Loading Semesters...'
\i 02_semesters.sql
\echo ''

-- 3. Load Courses (10,000 records)
\echo 'Loading Courses...'
\i 03_courses.sql
\echo ''

-- 4. Load Course Offerings (10,000 records)
\echo 'Loading Course Offerings...'
\i 04_course_offerings.sql
\echo ''

-- 5. Load Course Registrations (10,000 records + 5,000 student users)
\echo 'Loading Course Registrations and Student Users...'
\i 05_course_registrations.sql
\echo ''

COMMIT;

-- Re-enable foreign key checks
SET session_replication_role = DEFAULT;

-- Final summary
\echo '=== LOADING COMPLETE ==='
\echo ''
\echo 'Summary of loaded data:'

SELECT 'Academic Years: ' || COUNT(*) as summary FROM academic_years
UNION ALL
SELECT 'Semesters: ' || COUNT(*) FROM semesters  
UNION ALL
SELECT 'Courses: ' || COUNT(*) FROM courses
UNION ALL  
SELECT 'Course Offerings: ' || COUNT(*) FROM course_offerings
UNION ALL
SELECT 'Student Users: ' || COUNT(*) FROM users WHERE role = 3
UNION ALL
SELECT 'Course Registrations: ' || COUNT(*) FROM course_registrations
UNION ALL
SELECT 'Total Records: ' || (
    (SELECT COUNT(*) FROM academic_years) +
    (SELECT COUNT(*) FROM semesters) +
    (SELECT COUNT(*) FROM courses) +
    (SELECT COUNT(*) FROM course_offerings) +
    (SELECT COUNT(*) FROM users WHERE role = 3) +
    (SELECT COUNT(*) FROM course_registrations)
);

\echo ''
\echo 'Data consistency checks:'

-- Check for any orphaned records
SELECT 
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ No orphaned semesters'
        ELSE '✗ Found orphaned semesters: ' || COUNT(*)
    END as semester_check
FROM semesters s 
LEFT JOIN academic_years ay ON s.academic_year_id = ay.id 
WHERE ay.id IS NULL;

SELECT 
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ No orphaned course offerings'
        ELSE '✗ Found orphaned course offerings: ' || COUNT(*)
    END as course_offering_check
FROM course_offerings co 
LEFT JOIN semesters s ON co.semester_id = s.id 
LEFT JOIN courses c ON co.course_id = c.id
WHERE s.id IS NULL OR c.id IS NULL;

SELECT 
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ No orphaned course registrations'
        ELSE '✗ Found orphaned course registrations: ' || COUNT(*)
    END as registration_check
FROM course_registrations cr 
LEFT JOIN users u ON cr.student_id = u.id 
LEFT JOIN course_offerings co ON cr.course_offering_id = co.id
WHERE u.id IS NULL OR co.id IS NULL;

\echo ''
\echo '=== DUMMY DATA LOADING COMPLETED SUCCESSFULLY ==='
\echo ''
\echo 'You can now:'
\echo '1. Test the course enrollment API endpoints'
\echo '2. Test course offering CRUD operations' 
\echo '3. Explore the data with SQL queries'
\echo '4. Use the data for development and testing'
\echo ''
\echo 'Example queries to explore the data:'
\echo '  -- View course offerings with course details:'
\echo '  SELECT co.section_code, c.name, c.credit, co.capacity FROM course_offerings co JOIN courses c ON co.course_id = c.id LIMIT 10;'
\echo ''
\echo '  -- View student enrollment counts:'  
\echo '  SELECT student_id, COUNT(*) as enrollments FROM course_registrations GROUP BY student_id ORDER BY enrollments DESC LIMIT 10;'
\echo ''
\echo '  -- View popular courses by enrollment:'
\echo '  SELECT c.name, COUNT(*) as total_enrollments FROM course_registrations cr JOIN course_offerings co ON cr.course_offering_id = co.id JOIN courses c ON co.course_id = c.id GROUP BY c.id, c.name ORDER BY total_enrollments DESC LIMIT 10;'