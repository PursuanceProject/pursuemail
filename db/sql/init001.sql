create extension if not exists "uuid-ossp";
CREATE TABLE email_account (
  id          uuid      NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
  email       text      NOT NULL CHECK (email ~* '^[A-Za-z0-9_\.\-\+]+@[A-Za-z0-9\.\-]+\.[A-Za-z0-9]+$'), /* regex based on https://stackoverflow.com/a/10164872/197160 */
  created     timestamp WITH time zone DEFAULT now()
);
ALTER TABLE email_account OWNER TO pursuemail;
