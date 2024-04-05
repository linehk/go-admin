CREATE TABLE app_user (
  id SERIAL PRIMARY KEY,
  username VARCHAR NOT NULL,
  password VARCHAR NOT NULL,
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  phone VARCHAR NOT NULL,
  remark VARCHAR NOT NULL,
  status VARCHAR NOT NULL,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL
);

CREATE TABLE user_role (
  id SERIAL PRIMARY KEY,
  user_id SERIAL NOT NULL,
  role_id SERIAL NOT NULL,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL
);

CREATE TABLE role (
  id SERIAL PRIMARY KEY,
  code VARCHAR NOT NULL,
  name VARCHAR NOT NULL,
  description VARCHAR NOT NULL,
  sequence SMALLINT NOT NULL,
  status VARCHAR NOT NULL,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL
);

CREATE TABLE role_menu (
  id SERIAL PRIMARY KEY,
  role_id SERIAL NOT NULL,
  menu_id SERIAL NOT NULL,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL
);

CREATE TABLE menu (
  id SERIAL PRIMARY KEY,
  code VARCHAR NOT NULL,
  name VARCHAR NOT NULL,
  description VARCHAR NOT NULL,
  sequence SMALLINT NOT NULL,
  type VARCHAR NOT NULL,
  path VARCHAR NOT NULL,
  property VARCHAR NOT NULL,
  parent_id SERIAL NOT NULL,
  parent_path VARCHAR NOT NULL,
  status VARCHAR NOT NULL,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL
);

CREATE TABLE resource (
  id SERIAL PRIMARY KEY,
  menu_id SERIAL NOT NULL,
  method VARCHAR NOT NULL,
  path VARCHAR NOT NULL,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL
);