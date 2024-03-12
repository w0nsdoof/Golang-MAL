--
-- PostgreSQL database dump
--

-- Dumped from database version 16rc1
-- Dumped by pg_dump version 16rc1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: temp; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA temp;


ALTER SCHEMA temp OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: animes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.animes (
    id integer NOT NULL,
    rating double precision,
    title text NOT NULL,
    genres text
);


ALTER TABLE public.animes OWNER TO postgres;

--
-- Name: animes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.animes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.animes_id_seq OWNER TO postgres;

--
-- Name: animes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.animes_id_seq OWNED BY public.animes.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: user_and_anime; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_and_anime (
    id integer NOT NULL,
    user_id integer NOT NULL,
    anime_id integer NOT NULL,
    status text,
    user_rating double precision
);


ALTER TABLE public.user_and_anime OWNER TO postgres;

--
-- Name: user_and_anime_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_and_anime_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_and_anime_id_seq OWNER TO postgres;

--
-- Name: user_and_anime_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_and_anime_id_seq OWNED BY public.user_and_anime.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_and_animes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_and_animes (
    user_id integer NOT NULL,
    anime_id integer NOT NULL,
    status text,
    rating double precision
);


ALTER TABLE public.users_and_animes OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: animes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.animes ALTER COLUMN id SET DEFAULT nextval('public.animes_id_seq'::regclass);


--
-- Name: user_and_anime id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_and_anime ALTER COLUMN id SET DEFAULT nextval('public.user_and_anime_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: animes animes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.animes
    ADD CONSTRAINT animes_pkey PRIMARY KEY (id);


--
-- Name: schema_migration schema_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migration
    ADD CONSTRAINT schema_migration_pkey PRIMARY KEY (version);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: user_and_anime user_and_anime_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_and_anime
    ADD CONSTRAINT user_and_anime_pkey PRIMARY KEY (id);


--
-- Name: user_and_anime user_anime_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_and_anime
    ADD CONSTRAINT user_anime_unique UNIQUE (user_id, anime_id);


--
-- Name: users_and_animes users_and_animes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_and_animes
    ADD CONSTRAINT users_and_animes_pkey PRIMARY KEY (user_id, anime_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: user_and_anime user_and_anime_anime_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_and_anime
    ADD CONSTRAINT user_and_anime_anime_id_fkey FOREIGN KEY (anime_id) REFERENCES public.animes(id);


--
-- Name: user_and_anime user_and_anime_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_and_anime
    ADD CONSTRAINT user_and_anime_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

