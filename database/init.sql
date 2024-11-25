DO
$$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_database WHERE datname = 'tradeapi_db'
   ) THEN
      CREATE DATABASE tradeapi_db;
   END IF;
END
$$;