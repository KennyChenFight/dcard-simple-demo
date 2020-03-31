/*
	dcard_admin is the database owner.
	It owns all objects(table, trigger, etc) in the database.
	Usually, it should be used ONLY when DBA wants to change the database schema.
	dcard_user is the account used by the golang executable.
	Thus, it is NOT allowed to create / change / any object in database.
	For normal table, only CURD privilege is granted, truncate table should NOT be granted.
	dcard_readonly is used during debugging.
	Trusted software developer will use this account to view the data in production database directly.
	Thus it should have select privilege.
*/
CREATE ROLE dcard_admin LOGIN PASSWORD 'admin_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE dcard_user LOGIN PASSWORD 'user_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE dcard_readonly LOGIN PASSWORD 'readonly_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

/*
    setup db
*/
CREATE DATABASE dcard_db with ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' CONNECTION LIMIT = -1 template=template0;
ALTER DATABASE dcard_db OWNER TO dcard_admin;
ALTER DATABASE dcard_db SET timezone TO 'UTC';
REVOKE USAGE ON SCHEMA public FROM PUBLIC;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public to dcard_admin;
GRANT CREATE ON SCHEMA public to dcard_admin;
/* grant the schema access privilege to normal users. Without schema right, user will unable to see the tables. */
GRANT USAGE ON SCHEMA public to dcard_user;
GRANT USAGE ON SCHEMA public to dcard_readonly;

/*
    create table
*/
--the script to remove all tables in the database
\connect dcard_db;
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

\connect dcard_db;
ALTER TABLE pairs ADD CONSTRAINT pairs_fk1 FOREIGN KEY (user_id_one) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE pairs ADD CONSTRAINT pairs_fk2 FOREIGN KEY (user_id_two) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

\connect dcard_db;

/*
    grant_table_privilege
*/
/*for normal tables */
ALTER TABLE users OWNER TO dcard_admin;
ALTER TABLE pairs OWNER TO dcard_admin;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE users to dcard_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE pairs to dcard_user;
GRANT SELECT ON TABLE users to dcard_readonly;
GRANT SELECT ON TABLE pairs to dcard_readonly;

/*
    insert testing data
*/
\connect dcard_db;
-- each test user's plain password is 0000 --
insert into users(id, email, password_digest, name)
values('97327413-6b65-486f-b299-91be0871f898', 'kenny@example.com', '$2a$10$gVtjNk4YL.O4I//ZBtvfN.YEebwR1Ci3.5OBHan4PWFzniSFqpzce', 'kenny');

insert into users(id, email, password_digest, name)
values('eb3c75df-b0df-4e06-a02f-e2ba77eba68a', 'nicole@example.com', '$2a$10$6tsb.2dRzV5gSTEJmtwkgeKpPIMO0VbMv2E6hP9xuAytwFlf0trVm', 'nicole');

insert into users(id, email, password_digest, name)
values('80695811-0bf2-44fd-980d-1635de7734a8', 'jack@example.com', '$2a$10$WkWwIpCbMyB1A2OuMC9LI.4LtQZtxNb1djcYqzeP0IayazJQgVkHG', 'jack');

insert into pairs(user_id_one, user_id_two)
values('97327413-6b65-486f-b299-91be0871f898', 'eb3c75df-b0df-4e06-a02f-e2ba77eba68a')