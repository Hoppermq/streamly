CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tenants
(
  id         UUID PRIMARY KEY             DEFAULT uuid_generate_v4(),
  identifier VARCHAR(255) UNIQUE NOT NULL,

  name       VARCHAR(255)        NOT NULL,

  created_at TIMESTAMP           NOT NULL DEFAULT now(),
  updated_at TIMESTAMP           NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
  id               UUID PRIMARY KEY             DEFAULT uuid_generate_v4(),
  identifier       VARCHAR(255) UNIQUE NOT NULL,

  primary_email_id VARCHAR(255),
  tenant_email     VARCHAR(255),

  first_name       VARCHAR(255),
  last_name        VARCHAR(255),
  username         VARCHAR(255),
  password         VARCHAR(255),

  tenants_id       UUID,

  created_at       TIMESTAMP           NOT NULL DEFAULT now(),
  updated_at       TIMESTAMP           NOT NULL
-- should have relation to zitadel auth table --
);

CREATE TABLE IF NOT EXISTS tenant_members
(
  id              UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
  identifier      VARCHAR(255),

  tenant_id       VARCHAR(255),
  zitadel_user_id VARCHAR(255), -- why ?? --

  joined_at       TIMESTAMP NOT NULL,

  created_at      TIMESTAMP NOT NULL DEFAULT now(),
  updated_at      TIMESTAMP NOT NULL
);
