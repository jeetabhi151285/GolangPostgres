--docker exec -it b98e5a20d7cb bash
--docker exec -it b98e5a20d7cb psql -U postgres -c "CREATE USER maerskadmin WITH PASSWORD 'admin4test';"
--docker exec -it b98e5a20d7cb psql -U postgres -c "CREATE USER maerskapp WITH PASSWORD 'app4test';"
--docker exec -it b98e5a20d7cb psql -U postgres -c "CREATE DATABASE maerskdb OWNER maerskadmin ;"
--docker exec -it b98e5a20d7cb psql -U postgres -c "\connect maerskdb maerskadmin;"
--docker exec -it b98e5a20d7cb psql -U postgres -c "CREATE SCHEMA odsschema AUTHORIZATION maerskadmin;"

-- maerskadmin is db owner
CREATE USER maerskadmin WITH PASSWORD 'admin4test';
-- maerskapp is the application user
CREATE USER maerskapp WITH PASSWORD 'app4test';
-- maerskdb is the db name
CREATE DATABASE maerskdb OWNER maerskadmin ;
\connect maerskdb maerskadmin;
-- maersk schema will contain all the hr related tables
CREATE SCHEMA odsschema AUTHORIZATION maerskadmin;
GRANT ALL ON SCHEMA odsschema TO maerskapp;
   