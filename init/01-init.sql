CREATE DATABASE gator;

\connect gator

CREATE TABLE users (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text unique not null
);

CREATE TABLE feeds (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text unique not null,
    url text unique not null,
    user_id uuid not null references users(id) on delete cascade,
    last_fetched_at timestamp
);

CREATE TABLE feed_follows (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null references users(id) on delete cascade,
    feed_id uuid not null references feeds(id) on delete cascade,
    constraint feed_follows_user_id_feed_id unique (user_id, feed_id)
);

CREATE TABLE posts (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title text not null,
    url text not null unique,
    description text,
    published_at timestamp,
    feed_id uuid not null references feeds(id) on delete cascade
);
