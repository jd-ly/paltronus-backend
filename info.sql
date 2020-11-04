
-- Database: paltronus_database

-- DROP DATABASE paltronus_database;

CREATE DATABASE paltronus_database
  WITH OWNER = postgres
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       CONNECTION LIMIT = -1;

-- Table: file

-- DROP TABLE file;

CREATE TABLE file
(
 id serial NOT NULL,
 title character varying NOT NULL,
 author character varying,
 description character varying,
 creation_date date,
 CONSTRAINT pk_file PRIMARY KEY (id )
)
WITH (
 OIDS=FALSE
);
ALTER TABLE file
 OWNER TO postgres;

-- Table: version

-- DROP TABLE version;

CREATE TABLE version
(
 id serial NOT NULL,
 raw_data character varying NOT NULL,
 author character varying,
 creation_date date,
 versionId integer,
 FOREIGN KEY (versionId) REFERENCES Version(versionId)
 CONSTRAINT pk_version PRIMARY KEY (id )
)
WITH (
 OIDS=FALSE
);
ALTER TABLE version
 OWNER TO postgres;