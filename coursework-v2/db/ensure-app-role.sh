#!/bin/bash
set -euo pipefail

: "${POSTGRES_ADMIN_USER:?POSTGRES_ADMIN_USER is required}"
: "${POSTGRES_ADMIN_PASSWORD:?POSTGRES_ADMIN_PASSWORD is required}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"

PGPASSWORD="$POSTGRES_ADMIN_PASSWORD" psql \
  --host="${DB_HOST:-localhost}" \
  --port="${DB_PORT:-5434}" \
  --username="$POSTGRES_ADMIN_USER" \
  --dbname="${DB_NAME:-szi_registry}" \
  --set=app_password="$DB_PASSWORD" <<'SQL'
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'szi_app') THEN
        EXECUTE format('CREATE ROLE szi_app LOGIN PASSWORD %L', :'app_password');
    ELSE
        EXECUTE format('ALTER ROLE szi_app LOGIN PASSWORD %L', :'app_password');
    END IF;
END
$$;

GRANT CONNECT ON DATABASE szi_registry TO szi_app;
GRANT USAGE, CREATE ON SCHEMA public TO szi_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO szi_app;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO szi_app;
ALTER TABLE IF EXISTS public.szi_records ALTER COLUMN ac_class TYPE varchar(32);

DO $$
DECLARE
    object_name text;
BEGIN
    FOR object_name IN SELECT tablename FROM pg_tables WHERE schemaname = 'public' LOOP
        EXECUTE format('ALTER TABLE public.%I OWNER TO szi_app', object_name);
    END LOOP;
    FOR object_name IN SELECT sequence_name FROM information_schema.sequences WHERE sequence_schema = 'public' LOOP
        EXECUTE format('ALTER SEQUENCE public.%I OWNER TO szi_app', object_name);
    END LOOP;
END
$$;
SQL
