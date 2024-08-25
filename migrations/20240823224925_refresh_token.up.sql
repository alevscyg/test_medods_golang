CREATE TABLE refresh_tokens (
  userid bigserial primary key,
  email varchar unique not null,
  refresh_token text not null
);
