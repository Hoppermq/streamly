CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--bun:split

CREATE TYPE roles AS ENUM ('owner', 'admin', 'user');

--bun:split

CREATE TABLE IF NOT EXISTS tenants
(
  id         UUID PRIMARY KEY             DEFAULT uuid_generate_v4(),
  identifier VARCHAR(255) UNIQUE NOT NULL,

  name       VARCHAR(255)        NOT NULL,

  created_at TIMESTAMP           NOT NULL DEFAULT now(),
  updated_at TIMESTAMP           NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users
(
  id               UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
  identifier       UUID UNIQUE NOT NULL,
  zitadel_user_id  VARCHAR(25)        NOT NULL,


  first_name       VARCHAR(255),
  last_name        VARCHAR(255),
  username         VARCHAR(255),
  primary_email VARCHAR(255),

  role             roles                default 'user',

  created_at       TIMESTAMP   NOT NULL DEFAULT now(),
  updated_at       TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE TABLE emails
(
  id         uuid PRIMARY KEY             DEFAULT uuid_generate_v4(),
  identifier VARCHAR(255) UNIQUE NOT NULL,

  user_id    UUID NOT NULL,
  email      VARCHAR(255) UNIQUE NOT NULL,
  is_primary BOOLEAN             NOT NULL DEFAULT false,
  verifier   BOOLEAN             NOT NULL DEFAULT false,

  tenant_id  UUID NOT NULL,
  domain     VARCHAR(255) GENERATED ALWAYS AS (
    SUBSTRING(email FROM '@(.*)$')
    ) STORED,

  created_at TIMESTAMP           NOT NULL DEFAULT now(),
  updated_at TIMESTAMP           NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tenant_members
(
  id         UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
  identifier VARCHAR(255),

  tenant_id  UUID NOT NULL,
  user_id    UUID NOT NULL,

  joined_at  TIMESTAMP NOT NULL DEFAULT now(),

  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now()
);

--bun:split

ALTER TABLE IF EXISTS tenant_members
  ADD CONSTRAINT fk_membership_user
    FOREIGN KEY (user_id) REFERENCES users (identifier) ON DELETE CASCADE,
  ADD CONSTRAINT fk_membership_tenant
    FOREIGN KEY (tenant_id) REFERENCES tenants (identifier) ON DELETE CASCADE;

ALTER TABLE IF EXISTS emails
  ADD CONSTRAINT fk_user_emails_user
    FOREIGN KEY (user_id) REFERENCES users (identifier) ON DELETE CASCADE;

--bun:split

CREATE INDEX idx_users_zitadel ON users(zitadel_user_id);
CREATE INDEX idx_user_email ON emails(user_id);
CREATE INDEX idx_user_domain ON emails(domain);
CREATE INDEX idx_memberships_user ON tenant_members(user_id);
CREATE INDEX idx_memberships_tenant ON tenant_members(tenant_id);
