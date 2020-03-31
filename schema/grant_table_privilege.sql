/*for normal tables */
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE users to dcard_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE pairs to dcard_user;

GRANT SELECT ON TABLE users to dcard_readonly;
GRANT SELECT ON TABLE pairs to dcard_readonly;
