--the script to remove all tables in the database
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS pairs CASCADE;

create table users
(
    id uuid,
    email character varying(200) not null,
    password_digest character varying(1000) not null,
    name character varying(255) not null,
    create_time timestamp without time zone not null default current_timestamp,
    update_time timestamp without time zone not null default current_timestamp,

    CONSTRAINT "users_pk" PRIMARY KEY (id)
);
ALTER TABLE users ADD CONSTRAINT users_u1 UNIQUE (email);

create table pairs
(
    user_id_one uuid,
    user_id_two uuid,

    CONSTRAINT "pairs_pk" PRIMARY KEY (user_id_one, user_id_two)
);
ALTER TABLE pairs ADD CONSTRAINT pairs_u1 UNIQUE (user_id_one);
