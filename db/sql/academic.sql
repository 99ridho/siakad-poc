-- name: GetCourseOffering :one
select * from course_offerings where id = $1;

-- name: GetCourse :one
select * from courses where id = $1;

-- name: GetCourseOfferingWithCourse :one
select 
    co.id as course_offering_id,
    co.semester_id,
    co.course_id,
    co.section_code,
    co.capacity,
    co.start_time as course_offering_start_time,
    co.created_at as course_offering_created_at,
    co.updated_at as course_offering_updated_at,
    co.deleted_at as course_offering_deleted_at,
    c.id as course_id,
    c.code as course_code,
    c.name as course_name,
    c.credit,
    c.created_at as course_created_at,
    c.updated_at as course_updated_at,
    c.deleted_at as course_deleted_at
from course_offerings co
join courses c on co.course_id = c.id
where co.id = $1;

-- name: GetStudentEnrollmentsWithDetails :many
select 
    cr.id as registration_id,
    cr.student_id,
    cr.course_offering_id,
    cr.created_at as registration_created_at,
    co.start_time as course_offering_start_time,
    c.credit
from course_registrations cr
join course_offerings co on cr.course_offering_id = co.id
join courses c on co.course_id = c.id
where cr.student_id = $1;

-- name: CountCourseOfferingEnrollments :one
select count(*) from course_registrations where course_offering_id = $1;

-- name: CheckEnrollmentExists :one
select exists(
    select 1 from course_registrations 
    where student_id = $1 and course_offering_id = $2
);

-- name: CreateEnrollment :one
insert into course_registrations (id, student_id, course_offering_id, created_at, updated_at)
values (gen_random_uuid(), $1, $2, now(), now())
returning *;

-- name: GetCourseOfferingsWithPagination :many
select 
    co.id as course_offering_id,
    co.semester_id,
    co.course_id,
    co.section_code,
    co.capacity,
    co.start_time as course_offering_start_time,
    co.created_at as course_offering_created_at,
    co.updated_at as course_offering_updated_at,
    co.deleted_at as course_offering_deleted_at,
    c.id as course_id,
    c.code as course_code,
    c.name as course_name,
    c.credit,
    c.created_at as course_created_at,
    c.updated_at as course_updated_at,
    c.deleted_at as course_deleted_at
from course_offerings co
join courses c on co.course_id = c.id
where co.deleted_at IS NULL
order by co.created_at desc
limit $1 offset $2;

-- name: CountCourseOfferings :one
select count(*) 
from course_offerings co
where co.deleted_at IS NULL;

-- name: CreateCourseOffering :one
insert into course_offerings (id, semester_id, course_id, section_code, capacity, start_time, created_at, updated_at)
values (gen_random_uuid(), $1, $2, $3, $4, $5, now(), now())
returning *;

-- name: UpdateCourseOffering :one
update course_offerings 
set semester_id = $2, course_id = $3, section_code = $4, capacity = $5, start_time = $6, updated_at = now()
where id = $1 and deleted_at IS NULL
returning *;

-- name: DeleteCourseOffering :one
update course_offerings 
set deleted_at = now(), updated_at = now()
where id = $1 and deleted_at IS NULL
returning *;

-- name: GetCourseOfferingByIDWithDetails :one
select 
    co.id as course_offering_id,
    co.semester_id,
    co.course_id,
    co.section_code,
    co.capacity,
    co.start_time as course_offering_start_time,
    co.created_at as course_offering_created_at,
    co.updated_at as course_offering_updated_at,
    co.deleted_at as course_offering_deleted_at,
    c.id as course_id,
    c.code as course_code,
    c.name as course_name,
    c.credit,
    c.created_at as course_created_at,
    c.updated_at as course_updated_at,
    c.deleted_at as course_deleted_at
from course_offerings co
join courses c on co.course_id = c.id
where co.id = $1 and co.deleted_at IS NULL;