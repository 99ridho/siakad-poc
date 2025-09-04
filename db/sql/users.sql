-- name: GetUser :one
select * from users where id = $1;

-- name: GetUserByEmail :one
select * from users where email = $1;