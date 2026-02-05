-- NUDEX PostgreSQL Database Initialization

-- Create databases
CREATE DATABASE nudex_catalog OWNER nudex_user;
CREATE DATABASE nudex_users OWNER nudex_user;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE nudex_catalog TO nudex_user;
GRANT ALL PRIVILEGES ON DATABASE nudex_users TO nudex_user;

-- Connect to catalog database and create schemas
\c nudex_catalog;
CREATE SCHEMA IF NOT EXISTS public;
GRANT ALL ON SCHEMA public TO nudex_user;

-- Connect to users database and create schemas
\c nudex_users;
CREATE SCHEMA IF NOT EXISTS public;
GRANT ALL ON SCHEMA public TO nudex_user;

-- Log completion
SELECT 'NUDEX PostgreSQL databases initialized successfully' as status;