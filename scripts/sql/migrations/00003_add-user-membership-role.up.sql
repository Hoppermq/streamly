CREATE TABLE IF NOT EXISTS tenant_roles
(
  id          uuid PRIMARY KEY      DEFAULT uuid_generate_v4(),
  identifier  uuid UNIQUE  NOT NULL,

  role        VARCHAR(255) NOT NULL,

  permissions TEXT[],
  metadata    jsonb,

  created_at  TIMESTAMP    NOT NULL DEFAULT now(),
  updated_at  TIMESTAMP    NOT NULL DEFAULT now()
);

--bun:split

ALTER TABLE IF EXISTS tenant_members
  ADD COLUMN role_id VARCHAR,
  ADD CONSTRAINT fk_member_role
    FOREIGN KEY (role_id) REFERENCES tenant_roles ON DELETE CASCADE;

--bun:split

CREATE INDEX idx_role ON tenant_roles (identifier);
CREATE INDEX idx_user_role ON tenant_members (role_id);
