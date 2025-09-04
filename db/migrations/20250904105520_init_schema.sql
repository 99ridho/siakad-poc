-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid not null,
    email varchar(255) not null,
    password varchar(255) not null,
    role numeric(2) not null, -- 1 admin, 2 koorprodi, 3 siswa
    created_at timestamptz not null default now(),
    updated_at timestamptz null,
    deleted_at timestamptz null,

    PRIMARY KEY (id)
);

CREATE TABLE academic_years (
    id uuid not null,
    code varchar(255) not null,
    start_time timestamptz not null,
    end_time timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz null,
    deleted_at timestamptz null,

    PRIMARY KEY (id)
);

CREATE TABLE semesters (
    id uuid not null,
    academic_year_id uuid not null,
    code varchar(255) not null,
    start_time timestamptz not null,
    end_time timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz null,
    deleted_at timestamptz null,

    PRIMARY KEY (id),
    FOREIGN KEY (academic_year_id) REFERENCES academic_years (id),
    UNIQUE (academic_year_id, code)
);

CREATE TABLE courses (
    id uuid not null,
    code varchar(255) not null,
    name varchar(255) not null,
    credit int not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz null,
    deleted_at timestamptz null,

    PRIMARY KEY (id)
);

CREATE TABLE course_offerings (
    id uuid not null,
    semester_id uuid not null,
    course_id uuid not null,
    section_code varchar(255) not null,
    capacity int not null,
    start_time timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz null,
    deleted_at timestamptz null,

    PRIMARY KEY (id),
    FOREIGN KEY (semester_id) REFERENCES semesters (id),
    FOREIGN KEY (course_id) REFERENCES courses (id),
    UNIQUE (semester_id, course_id, section_code)
);

CREATE TABLE course_registrations (
    id uuid not null,
    student_id uuid not null,
    course_offering_id uuid not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz null,
    deleted_at timestamptz null,

    PRIMARY KEY (id),
    FOREIGN KEY (student_id) REFERENCES users (id),
    FOREIGN KEY (course_offering_id) REFERENCES course_offerings (id),
    UNIQUE (student_id, course_offering_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE course_registrations;
DROP TABLE course_offerings;
DROP TABLE semesters;
DROP TABLE courses;
DROP TABLE academic_years;
DROP TABLE users;
-- +goose StatementEnd
