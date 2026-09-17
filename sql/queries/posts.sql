-- name: CreatePost :exec
insert into posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
on conflict(url) do nothing;

-- name: GetPostsForUser :many
select posts.*
from posts
inner join feed_follows on posts.feed_id = feed_follows.feed_id
where feed_follows.user_id = $1
order by published_at desc nulls last
limit $2;
