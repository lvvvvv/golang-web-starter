
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE SEQUENCE public.user_sessions_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;


CREATE TABLE public.user_sessions
(
  id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  referer text,
  user_agent text,
  client_ip text,
  user_id bigint,
  status integer,
  created_at timestamp with time zone,
  updated_at timestamp with time zone,
  CONSTRAINT user_sessions_pkey PRIMARY KEY (id)
)
WITH (
  OIDS=FALSE
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE public.user_sessions;

DROP SEQUENCE public.user_sessions_id_seq;
