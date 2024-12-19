-- name: GetUserById :one
select * from "user"
where id = $1;

-- name: GetUserByGithubRemoteId :one
select * from "user"
where github_remote_id = $1;

-- name: ListUsers :many
select * from "user"
where id > sqlc.arg('offset') and deleted_at is null
order by id
limit sqlc.arg('limit');

-- name: CreateUser :one
insert into "user" (id, name, email, phone, github_remote_id, github_handle, x_remote_id, x_handle)
values ($1, $2, $3, $4, $5, $6, $7, $8) returning *;

-- name: UpdateUser :execresult
update "user"
set name = $2
where id = $1;

-- name: UpsertUserByGithubRemoteId :one
insert into "user" (id, name, email, phone, github_remote_id, github_handle, x_remote_id, x_handle)
values ($1, $2, $3, $4, $5, $6, $7, $8)
on conflict (github_remote_id)
where deleted_at is null
do update set
  name = EXCLUDED.name,
  email = EXCLUDED.email,
  phone = EXCLUDED.phone,
  github_handle = EXCLUDED.github_handle,
  x_remote_id = EXCLUDED.x_remote_id,
  x_handle = EXCLUDED.x_handle,
  updated_at = now()
returning *;

-- name: DeleteUser :execresult
delete from "user"
where id = $1;
