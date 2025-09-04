-- name: GetUser :one
select * from users where id = $1;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: CreateUser :one
insert into users (id, email, password, role, created_at, updated_at)
values (gen_random_uuid(), $1, $2, $3, now(), now())
returning *;