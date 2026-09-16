-- name: CreateFeed :one
insert into feeds(id, created_at, updated_at, name, url, user_id)
values (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
returning *;

-- name: GetUsersFeeds :many
select *
from users
join feeds
    on users.id = feeds.user_id;
