ALTER TABLE IF EXISTS tenants
  ADD deleted BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE IF EXISTS tenants
  ADD deleted_at TIMESTAMP DEFAULT NULL;

CREATE INDEX idx_deleted_tenants ON tenants(deleted);
