-- +goose Up
CREATE TABLE records (
    id serial not null unique,
    username varchar(128) not null,
    movie_id int not null,
    is_watched boolean not null
);

CREATE TABLE movies (
    movie_id integer not null unique,
    title varchar(256) not null,
    year integer not null,
    description text not null,
    poster text not null,
    rating real not null,
    length integer not null
);

CREATE TABLE requests (
    id serial not null unique,
    request varchar(256)
);

CREATE INDEX records_username_mid_idx ON records(username, movie_id);

-- +goose Down
DROP INDEX records_username_mid_idx;
DROP TABLE records;
DROP TABLE movies;
DROP TABLE requests;
