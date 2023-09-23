CREATE TABLE records
(
    id serial not null unique,
    username varchar(128) not null unique,
    movie_title varchar(256) not null,
    movie_id int
);