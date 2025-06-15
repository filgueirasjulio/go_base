DO
$$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_database WHERE datname = 'base_db'
   ) THEN
      CREATE DATABASE base_db;
   END IF;
END
$$;