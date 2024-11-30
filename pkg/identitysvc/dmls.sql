-- name: GetUserById :one
select * from "user"
where id = $1 and deleted_at is null;

-- name: GetUserByGithubRemoteId :one
select * from "user"
where github_remote_id = $1 and deleted_at is null;

-- name: ListUsers :many
select * from "user"
where deleted_at is null;

-- name: CreateUser :one
insert into "user" (github_remote_id, name, email)
values ($1, $2, $3) returning *;

-- name: UpdateUser :execresult
update "user"
set name = $2, email = $3
where id = $1 and deleted_at is null;

-- name: DeleteUser :execresult
update "user"
set deleted_at = CURRENT_TIMESTAMP
where id = $1 and deleted_at is null;
