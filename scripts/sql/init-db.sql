-- Create project_admin role if it doesn't exist
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'project_admin') THEN
    CREATE ROLE project_admin WITH LOGIN SUPERUSER PASSWORD 'admin';
  END IF;
END
$$;

ALTER ROLE project_admin SUPERUSER;
