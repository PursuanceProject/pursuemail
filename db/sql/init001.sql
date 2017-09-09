create extension if not exists "uuid-ossp";
CREATE TABLE email_account (
  id          uuid      NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  email       text      NOT NULL,
  has_pubkey  bool      DEFAULT false,
  created     timestamp WITH time zone DEFAULT now()
);
ALTER TABLE email_account OWNER TO pursuemail;
