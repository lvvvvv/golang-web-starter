
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE SEQUENCE public.users_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;


CREATE TABLE public.users
(
  id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  name text,
  email text,
  mobile text,
  password text,
  role text,
  status integer,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  CONSTRAINT users_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

CREATE UNIQUE INDEX users_email_idx
  ON public.users
  USING btree
  (email COLLATE pg_catalog."default");

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP INDEX public.users_email_idx;

DROP TABLE public.users;

DROP SEQUENCE public.users_id_seq;
