-- name: GetChannelById :one
select * from "channel"
where id = $1;

-- name: GetChannelByName :one
select * from "channel"
where name = $1;

-- name: ListChannels :many
select * from "channel";

-- name: CreateChannel :one
insert into "channel" (id, name)
values ($1, $2) returning *;

-- name: UpdateChannel :execrows
update "channel"
set name = $2
where id = $1;

-- name: DeleteChannel :execrows
delete from "channel"
where id = $1;
