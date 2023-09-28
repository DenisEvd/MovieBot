CREATE TABLE records
(
    id serial not null unique,
    username varchar(128) not null,
    movie_title varchar(256) not null,
    movie_id int not null
);

CREATE TABLE requests
(
    id serial not null unique,
    request varchar(256)
)