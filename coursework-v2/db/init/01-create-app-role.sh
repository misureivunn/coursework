#!/bin/bash
set -euo pipefail

psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
  --set=app_user="$APP_DB_USER" \
  --set=app_password="$APP_DB_PASSWORD" \
  --set=db_name="$POSTGRES_DB" \
  --set=admin_user="$POSTGRES_USER" <<'SQL'
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = :'app_user') THEN
        EXECUTE format('CREATE ROLE %I LOGIN PASSWORD %L', :'app_user', :'app_password');
    ELSE
        EXECUTE format('ALTER ROLE %I PASSWORD %L', :'app_user', :'app_password');
    END IF;
END
$$;

GRANT CONNECT ON DATABASE :"db_name" TO :"app_user";
GRANT USAGE, CREATE ON SCHEMA public TO :"app_user";
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO :"app_user";
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO :"app_user";
ALTER DEFAULT PRIVILEGES FOR ROLE :"admin_user" IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO :"app_user";
ALTER DEFAULT PRIVILEGES FOR ROLE :"admin_user" IN SCHEMA public GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO :"app_user";
SQL