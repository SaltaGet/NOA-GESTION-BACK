--
-- PostgreSQL database dump
--

\restrict aXhPVpZh0ortDLcOBostXiIOpKvYqUKdJBaT6l19efXzJBS6dU2o1KsHdHBdI2p

-- Dumped from database version 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)

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

ALTER TABLE IF EXISTS ONLY public.user_tenants DROP CONSTRAINT IF EXISTS fk_user_tenants_user;
ALTER TABLE IF EXISTS ONLY public.user_tenants DROP CONSTRAINT IF EXISTS fk_tenants_user_tenants;
ALTER TABLE IF EXISTS ONLY public.setting_tenants DROP CONSTRAINT IF EXISTS fk_tenants_setting;
ALTER TABLE IF EXISTS ONLY public.pay_tenants DROP CONSTRAINT IF EXISTS fk_tenants_pay_tenant;
ALTER TABLE IF EXISTS ONLY public.tenant_modules DROP CONSTRAINT IF EXISTS fk_tenants_modules;
ALTER TABLE IF EXISTS ONLY public.credentials DROP CONSTRAINT IF EXISTS fk_tenants_credentials;
ALTER TABLE IF EXISTS ONLY public.tenants DROP CONSTRAINT IF EXISTS fk_plans_tenants;
ALTER TABLE IF EXISTS ONLY public.pay_details DROP CONSTRAINT IF EXISTS fk_pay_tenants_pay_detail;
ALTER TABLE IF EXISTS ONLY public.pay_tenants DROP CONSTRAINT IF EXISTS fk_pay_tenants_admin;
ALTER TABLE IF EXISTS ONLY public.tenant_modules DROP CONSTRAINT IF EXISTS fk_modules_tenants;
ALTER TABLE IF EXISTS ONLY public.audit_log_admins DROP CONSTRAINT IF EXISTS fk_audit_log_admins_admin;
DROP TRIGGER IF EXISTS tr_audit_users ON public.users;
DROP TRIGGER IF EXISTS tr_audit_user_tenants ON public.user_tenants;
DROP TRIGGER IF EXISTS tr_audit_tenants ON public.tenants;
DROP TRIGGER IF EXISTS tr_audit_tenant_modules ON public.tenant_modules;
DROP TRIGGER IF EXISTS tr_audit_setting_tenants ON public.setting_tenants;
DROP TRIGGER IF EXISTS tr_audit_schema_migrations ON public.schema_migrations;
DROP TRIGGER IF EXISTS tr_audit_plans ON public.plans;
DROP TRIGGER IF EXISTS tr_audit_pay_tenants ON public.pay_tenants;
DROP TRIGGER IF EXISTS tr_audit_pay_details ON public.pay_details;
DROP TRIGGER IF EXISTS tr_audit_news ON public.news;
DROP TRIGGER IF EXISTS tr_audit_modules ON public.modules;
DROP TRIGGER IF EXISTS tr_audit_feedbacks ON public.feedbacks;
DROP TRIGGER IF EXISTS tr_audit_credentials ON public.credentials;
DROP TRIGGER IF EXISTS tr_audit_admins ON public.admins;
DROP INDEX IF EXISTS public.idx_tenants_email;
DROP INDEX IF EXISTS public.idx_tenants_deleted_at;
DROP INDEX IF EXISTS public.idx_tenants_cuit_pdv;
DROP INDEX IF EXISTS public.idx_tenant_modules_deleted_at;
DROP INDEX IF EXISTS public.idx_tenant_module;
DROP INDEX IF EXISTS public.idx_setting_tenants_tenant_id;
DROP INDEX IF EXISTS public.idx_plans_name;
DROP INDEX IF EXISTS public.idx_modules_name;
DROP INDEX IF EXISTS public.idx_modules_deleted_at;
DROP INDEX IF EXISTS public.idx_credentials_tenant_id;
DROP INDEX IF EXISTS public.idx_audit_log_admins_transaction_id;
DROP INDEX IF EXISTS public.idx_admins_username;
DROP INDEX IF EXISTS public.idx_admins_deleted_at;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.user_tenants DROP CONSTRAINT IF EXISTS user_tenants_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS uni_users_username;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS uni_users_email;
ALTER TABLE IF EXISTS ONLY public.tenants DROP CONSTRAINT IF EXISTS uni_tenants_identifier;
ALTER TABLE IF EXISTS ONLY public.credentials DROP CONSTRAINT IF EXISTS uni_credentials_cuit;
ALTER TABLE IF EXISTS ONLY public.admins DROP CONSTRAINT IF EXISTS uni_admins_email;
ALTER TABLE IF EXISTS ONLY public.tenants DROP CONSTRAINT IF EXISTS tenants_pkey;
ALTER TABLE IF EXISTS ONLY public.tenant_modules DROP CONSTRAINT IF EXISTS tenant_modules_pkey;
ALTER TABLE IF EXISTS ONLY public.setting_tenants DROP CONSTRAINT IF EXISTS setting_tenants_pkey;
ALTER TABLE IF EXISTS ONLY public.schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.plans DROP CONSTRAINT IF EXISTS plans_pkey;
ALTER TABLE IF EXISTS ONLY public.pay_tenants DROP CONSTRAINT IF EXISTS pay_tenants_pkey;
ALTER TABLE IF EXISTS ONLY public.pay_details DROP CONSTRAINT IF EXISTS pay_details_pkey;
ALTER TABLE IF EXISTS ONLY public.news DROP CONSTRAINT IF EXISTS news_pkey;
ALTER TABLE IF EXISTS ONLY public.modules DROP CONSTRAINT IF EXISTS modules_pkey;
ALTER TABLE IF EXISTS ONLY public.feedbacks DROP CONSTRAINT IF EXISTS feedbacks_pkey;
ALTER TABLE IF EXISTS ONLY public.credentials DROP CONSTRAINT IF EXISTS credentials_pkey;
ALTER TABLE IF EXISTS ONLY public.audit_log_admins DROP CONSTRAINT IF EXISTS audit_log_admins_pkey;
ALTER TABLE IF EXISTS ONLY public.admins DROP CONSTRAINT IF EXISTS admins_pkey;
ALTER TABLE IF EXISTS public.users ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.tenants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.tenant_modules ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.setting_tenants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.plans ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.pay_tenants ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.pay_details ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.news ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.modules ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.feedbacks ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.credentials ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.audit_log_admins ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.admins ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.users_id_seq;
DROP TABLE IF EXISTS public.users;
DROP TABLE IF EXISTS public.user_tenants;
DROP SEQUENCE IF EXISTS public.tenants_id_seq;
DROP TABLE IF EXISTS public.tenants;
DROP SEQUENCE IF EXISTS public.tenant_modules_id_seq;
DROP TABLE IF EXISTS public.tenant_modules;
DROP SEQUENCE IF EXISTS public.setting_tenants_id_seq;
DROP TABLE IF EXISTS public.setting_tenants;
DROP TABLE IF EXISTS public.schema_migrations;
DROP SEQUENCE IF EXISTS public.plans_id_seq;
DROP TABLE IF EXISTS public.plans;
DROP SEQUENCE IF EXISTS public.pay_tenants_id_seq;
DROP TABLE IF EXISTS public.pay_tenants;
DROP SEQUENCE IF EXISTS public.pay_details_id_seq;
DROP TABLE IF EXISTS public.pay_details;
DROP SEQUENCE IF EXISTS public.news_id_seq;
DROP TABLE IF EXISTS public.news;
DROP SEQUENCE IF EXISTS public.modules_id_seq;
DROP TABLE IF EXISTS public.modules;
DROP SEQUENCE IF EXISTS public.feedbacks_id_seq;
DROP TABLE IF EXISTS public.feedbacks;
DROP SEQUENCE IF EXISTS public.credentials_id_seq;
DROP TABLE IF EXISTS public.credentials;
DROP SEQUENCE IF EXISTS public.audit_log_admins_id_seq;
DROP TABLE IF EXISTS public.audit_log_admins;
DROP SEQUENCE IF EXISTS public.admins_id_seq;
DROP TABLE IF EXISTS public.admins;
DROP FUNCTION IF EXISTS public.audit_trigger_function_admin();
--
-- Name: audit_trigger_function_admin(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.audit_trigger_function_admin() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
	DECLARE
			current_member TEXT;
			current_tx_id BIGINT;
	BEGIN
			current_member := current_setting('app.current_member_id', true);
			current_tx_id := txid_current();

			IF current_member IS NULL OR current_member = '' OR current_member = '0' THEN
					RETURN NEW;
			END IF;

			INSERT INTO audit_log_admins (
					transaction_id,
					admin_id,
					method,
					path,
					old_value,
					new_value,
					created_at
			)
			VALUES (
					current_tx_id,
					current_member::BIGINT,
					LOWER(TG_OP),
					TG_TABLE_NAME,
					CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE to_jsonb(OLD) END,
					CASE WHEN TG_OP = 'DELETE' THEN NULL ELSE to_jsonb(NEW) END,
					NOW()
			);

			RETURN NULL;
	END;
	$$;


ALTER FUNCTION public.audit_trigger_function_admin() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: admins; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.admins (
    id bigint NOT NULL,
    first_name character varying(30) NOT NULL,
    last_name character varying(30) NOT NULL,
    username character varying(30) NOT NULL,
    email character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    is_super_admin boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.admins OWNER TO postgres;

--
-- Name: admins_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.admins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.admins_id_seq OWNER TO postgres;

--
-- Name: admins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.admins_id_seq OWNED BY public.admins.id;


--
-- Name: audit_log_admins; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_log_admins (
    id bigint NOT NULL,
    transaction_id bigint,
    admin_id bigint NOT NULL,
    method character varying(10) NOT NULL,
    path character varying(255) NOT NULL,
    old_value jsonb,
    new_value jsonb,
    created_at timestamp with time zone
);


ALTER TABLE public.audit_log_admins OWNER TO postgres;

--
-- Name: audit_log_admins_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.audit_log_admins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.audit_log_admins_id_seq OWNER TO postgres;

--
-- Name: audit_log_admins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.audit_log_admins_id_seq OWNED BY public.audit_log_admins.id;


--
-- Name: credentials; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.credentials (
    id bigint NOT NULL,
    tenant_id bigint NOT NULL,
    access_token_mp character varying(255),
    access_token_test_mp character varying(255),
    social_reason character varying(255),
    business_name character varying(255),
    address character varying(255),
    responsibility_front_iva character varying(255),
    gross_income character varying(255),
    start_activities text,
    cuit character varying(255),
    concept character varying(255),
    arca_certificate text,
    arca_key text,
    token_arca character varying(255),
    sign_arca character varying(255),
    expire_token_arca timestamp with time zone,
    token_email text
);


ALTER TABLE public.credentials OWNER TO postgres;

--
-- Name: credentials_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.credentials_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.credentials_id_seq OWNER TO postgres;

--
-- Name: credentials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.credentials_id_seq OWNED BY public.credentials.id;


--
-- Name: feedbacks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.feedbacks (
    id bigint NOT NULL,
    title character varying(255),
    content text,
    is_read boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.feedbacks OWNER TO postgres;

--
-- Name: feedbacks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.feedbacks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.feedbacks_id_seq OWNER TO postgres;

--
-- Name: feedbacks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.feedbacks_id_seq OWNED BY public.feedbacks.id;


--
-- Name: modules; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.modules (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    price_monthly numeric(10,2),
    price_yearly numeric(10,2),
    description text,
    features text,
    amount_images_per_product integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.modules OWNER TO postgres;

--
-- Name: modules_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.modules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.modules_id_seq OWNER TO postgres;

--
-- Name: modules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.modules_id_seq OWNED BY public.modules.id;


--
-- Name: news; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.news (
    id bigint NOT NULL,
    title character varying(255),
    content text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.news OWNER TO postgres;

--
-- Name: news_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.news_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.news_id_seq OWNER TO postgres;

--
-- Name: news_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.news_id_seq OWNED BY public.news.id;


--
-- Name: pay_details; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pay_details (
    id bigint NOT NULL,
    pay_tenant_id bigint NOT NULL,
    pay_id character varying(255),
    amount numeric NOT NULL,
    method_pay character varying(30) DEFAULT 'cash'::character varying NOT NULL,
    state_pay character varying(30) DEFAULT 'pending'::character varying NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.pay_details OWNER TO postgres;

--
-- Name: pay_details_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pay_details_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pay_details_id_seq OWNER TO postgres;

--
-- Name: pay_details_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pay_details_id_seq OWNED BY public.pay_details.id;


--
-- Name: pay_tenants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pay_tenants (
    id bigint NOT NULL,
    tenant_id bigint NOT NULL,
    admin_id bigint NOT NULL,
    amount_month bigint NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.pay_tenants OWNER TO postgres;

--
-- Name: pay_tenants_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pay_tenants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pay_tenants_id_seq OWNER TO postgres;

--
-- Name: pay_tenants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pay_tenants_id_seq OWNED BY public.pay_tenants.id;


--
-- Name: plans; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plans (
    id bigint NOT NULL,
    name character varying(255),
    price_mounthly numeric(10,2),
    price_yearly numeric(10,2),
    description text,
    features text,
    amount_point_sale bigint NOT NULL,
    amount_member bigint NOT NULL,
    amount_product bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.plans OWNER TO postgres;

--
-- Name: plans_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.plans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.plans_id_seq OWNER TO postgres;

--
-- Name: plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.plans_id_seq OWNED BY public.plans.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: setting_tenants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.setting_tenants (
    id bigint NOT NULL,
    tenant_id bigint NOT NULL,
    logo character varying(255),
    front_page character varying(255),
    title character varying(255),
    slogan text,
    primary_color character varying(255),
    secondary_color character varying(255),
    phone character varying(255),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.setting_tenants OWNER TO postgres;

--
-- Name: setting_tenants_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.setting_tenants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.setting_tenants_id_seq OWNER TO postgres;

--
-- Name: setting_tenants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.setting_tenants_id_seq OWNED BY public.setting_tenants.id;


--
-- Name: tenant_modules; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tenant_modules (
    id bigint NOT NULL,
    module_id bigint NOT NULL,
    tenant_id bigint NOT NULL,
    expiration timestamp with time zone,
    accepted_terms boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.tenant_modules OWNER TO postgres;

--
-- Name: tenant_modules_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tenant_modules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tenant_modules_id_seq OWNER TO postgres;

--
-- Name: tenant_modules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tenant_modules_id_seq OWNED BY public.tenant_modules.id;


--
-- Name: tenants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tenants (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    identifier character varying(50) NOT NULL,
    address character varying(255) NOT NULL,
    phone character varying(20) NOT NULL,
    email character varying(100) NOT NULL,
    cuit_pdv character varying(50) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    plan_id bigint NOT NULL,
    connection character varying(255) NOT NULL,
    expiration timestamp with time zone,
    accepted_terms boolean DEFAULT false NOT NULL,
    ip character varying(255) DEFAULT NULL::character varying,
    date_accepted timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.tenants OWNER TO postgres;

--
-- Name: tenants_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tenants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tenants_id_seq OWNER TO postgres;

--
-- Name: tenants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tenants_id_seq OWNED BY public.tenants.id;


--
-- Name: user_tenants; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_tenants (
    user_id bigint NOT NULL,
    tenant_id bigint NOT NULL
);


ALTER TABLE public.user_tenants OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    first_name character varying(30) NOT NULL,
    last_name character varying(30) NOT NULL,
    email character varying(100) NOT NULL,
    username character varying(50) NOT NULL,
    address character varying(255) DEFAULT NULL::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
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
-- Name: admins id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admins ALTER COLUMN id SET DEFAULT nextval('public.admins_id_seq'::regclass);


--
-- Name: audit_log_admins id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_log_admins ALTER COLUMN id SET DEFAULT nextval('public.audit_log_admins_id_seq'::regclass);


--
-- Name: credentials id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.credentials ALTER COLUMN id SET DEFAULT nextval('public.credentials_id_seq'::regclass);


--
-- Name: feedbacks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feedbacks ALTER COLUMN id SET DEFAULT nextval('public.feedbacks_id_seq'::regclass);


--
-- Name: modules id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.modules ALTER COLUMN id SET DEFAULT nextval('public.modules_id_seq'::regclass);


--
-- Name: news id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news ALTER COLUMN id SET DEFAULT nextval('public.news_id_seq'::regclass);


--
-- Name: pay_details id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_details ALTER COLUMN id SET DEFAULT nextval('public.pay_details_id_seq'::regclass);


--
-- Name: pay_tenants id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_tenants ALTER COLUMN id SET DEFAULT nextval('public.pay_tenants_id_seq'::regclass);


--
-- Name: plans id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plans ALTER COLUMN id SET DEFAULT nextval('public.plans_id_seq'::regclass);


--
-- Name: setting_tenants id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.setting_tenants ALTER COLUMN id SET DEFAULT nextval('public.setting_tenants_id_seq'::regclass);


--
-- Name: tenant_modules id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tenant_modules ALTER COLUMN id SET DEFAULT nextval('public.tenant_modules_id_seq'::regclass);


--
-- Name: tenants id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tenants ALTER COLUMN id SET DEFAULT nextval('public.tenants_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: admins admins_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT admins_pkey PRIMARY KEY (id);


--
-- Name: audit_log_admins audit_log_admins_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_log_admins
    ADD CONSTRAINT audit_log_admins_pkey PRIMARY KEY (id);


--
-- Name: credentials credentials_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.credentials
    ADD CONSTRAINT credentials_pkey PRIMARY KEY (id);


--
-- Name: feedbacks feedbacks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.feedbacks
    ADD CONSTRAINT feedbacks_pkey PRIMARY KEY (id);


--
-- Name: modules modules_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.modules
    ADD CONSTRAINT modules_pkey PRIMARY KEY (id);


--
-- Name: news news_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.news
    ADD CONSTRAINT news_pkey PRIMARY KEY (id);


--
-- Name: pay_details pay_details_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_details
    ADD CONSTRAINT pay_details_pkey PRIMARY KEY (id);


--
-- Name: pay_tenants pay_tenants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_tenants
    ADD CONSTRAINT pay_tenants_pkey PRIMARY KEY (id);


--
-- Name: plans plans_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plans
    ADD CONSTRAINT plans_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: setting_tenants setting_tenants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.setting_tenants
    ADD CONSTRAINT setting_tenants_pkey PRIMARY KEY (id);


--
-- Name: tenant_modules tenant_modules_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tenant_modules
    ADD CONSTRAINT tenant_modules_pkey PRIMARY KEY (id);


--
-- Name: tenants tenants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tenants
    ADD CONSTRAINT tenants_pkey PRIMARY KEY (id);


--
-- Name: admins uni_admins_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admins
    ADD CONSTRAINT uni_admins_email UNIQUE (email);


--
-- Name: credentials uni_credentials_cuit; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.credentials
    ADD CONSTRAINT uni_credentials_cuit UNIQUE (cuit);


--
-- Name: tenants uni_tenants_identifier; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tenants
    ADD CONSTRAINT uni_tenants_identifier UNIQUE (identifier);


--
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- Name: users uni_users_username; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_username UNIQUE (username);


--
-- Name: user_tenants user_tenants_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_tenants
    ADD CONSTRAINT user_tenants_pkey PRIMARY KEY (user_id, tenant_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_admins_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_admins_deleted_at ON public.admins USING btree (deleted_at);


--
-- Name: idx_admins_username; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_admins_username ON public.admins USING btree (username);


--
-- Name: idx_audit_log_admins_transaction_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_audit_log_admins_transaction_id ON public.audit_log_admins USING btree (transaction_id);


--
-- Name: idx_credentials_tenant_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_credentials_tenant_id ON public.credentials USING btree (tenant_id);


--
-- Name: idx_modules_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_modules_deleted_at ON public.modules USING btree (deleted_at);


--
-- Name: idx_modules_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_modules_name ON public.modules USING btree (name);


--
-- Name: idx_plans_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_plans_name ON public.plans USING btree (name);


--
-- Name: idx_setting_tenants_tenant_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_setting_tenants_tenant_id ON public.setting_tenants USING btree (tenant_id);


--
-- Name: idx_tenant_module; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_tenant_module ON public.tenant_modules USING btree (module_id, tenant_id);


--
-- Name: idx_tenant_modules_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_tenant_modules_deleted_at ON public.tenant_modules USING btree (deleted_at);


--
-- Name: idx_tenants_cuit_pdv; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_tenants_cuit_pdv ON public.tenants USING btree (cuit_pdv);


--
-- Name: idx_tenants_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_tenants_deleted_at ON public.tenants USING btree (deleted_at);


--
-- Name: idx_tenants_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_tenants_email ON public.tenants USING btree (email);


--
-- Name: admins tr_audit_admins; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_admins AFTER INSERT OR DELETE OR UPDATE ON public.admins FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: credentials tr_audit_credentials; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_credentials AFTER INSERT OR DELETE OR UPDATE ON public.credentials FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: feedbacks tr_audit_feedbacks; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_feedbacks AFTER INSERT OR DELETE OR UPDATE ON public.feedbacks FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: modules tr_audit_modules; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_modules AFTER INSERT OR DELETE OR UPDATE ON public.modules FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: news tr_audit_news; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_news AFTER INSERT OR DELETE OR UPDATE ON public.news FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: pay_details tr_audit_pay_details; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_pay_details AFTER INSERT OR DELETE OR UPDATE ON public.pay_details FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: pay_tenants tr_audit_pay_tenants; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_pay_tenants AFTER INSERT OR DELETE OR UPDATE ON public.pay_tenants FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: plans tr_audit_plans; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_plans AFTER INSERT OR DELETE OR UPDATE ON public.plans FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: schema_migrations tr_audit_schema_migrations; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_schema_migrations AFTER INSERT OR DELETE OR UPDATE ON public.schema_migrations FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: setting_tenants tr_audit_setting_tenants; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_setting_tenants AFTER INSERT OR DELETE OR UPDATE ON public.setting_tenants FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: tenant_modules tr_audit_tenant_modules; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_tenant_modules AFTER INSERT OR DELETE OR UPDATE ON public.tenant_modules FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: tenants tr_audit_tenants; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_tenants AFTER INSERT OR DELETE OR UPDATE ON public.tenants FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: user_tenants tr_audit_user_tenants; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_user_tenants AFTER INSERT OR DELETE OR UPDATE ON public.user_tenants FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: users tr_audit_users; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_users AFTER INSERT OR DELETE OR UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function_admin();


--
-- Name: audit_log_admins fk_audit_log_admins_admin; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_log_admins
    ADD CONSTRAINT fk_audit_log_admins_admin FOREIGN KEY (admin_id) REFERENCES public.admins(id);


--
-- Name: tenant_modules fk_modules_tenants; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tenant_modules
    ADD CONSTRAINT fk_modules_tenants FOREIGN KEY (module_id) REFERENCES public.modules(id);


--
-- Name: pay_tenants fk_pay_tenants_admin; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_tenants
    ADD CONSTRAINT fk_pay_tenants_admin FOREIGN KEY (admin_id) REFERENCES public.admins(id);


--
-- Name: pay_details fk_pay_tenants_pay_detail; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_details
    ADD CONSTRAINT fk_pay_tenants_pay_detail FOREIGN KEY (pay_tenant_id) REFERENCES public.pay_tenants(id);


--
-- Name: tenants fk_plans_tenants; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tenants
    ADD CONSTRAINT fk_plans_tenants FOREIGN KEY (plan_id) REFERENCES public.plans(id);


--
-- Name: credentials fk_tenants_credentials; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.credentials
    ADD CONSTRAINT fk_tenants_credentials FOREIGN KEY (tenant_id) REFERENCES public.tenants(id);


--
-- Name: tenant_modules fk_tenants_modules; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tenant_modules
    ADD CONSTRAINT fk_tenants_modules FOREIGN KEY (tenant_id) REFERENCES public.tenants(id);


--
-- Name: pay_tenants fk_tenants_pay_tenant; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_tenants
    ADD CONSTRAINT fk_tenants_pay_tenant FOREIGN KEY (tenant_id) REFERENCES public.tenants(id);


--
-- Name: setting_tenants fk_tenants_setting; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.setting_tenants
    ADD CONSTRAINT fk_tenants_setting FOREIGN KEY (tenant_id) REFERENCES public.tenants(id);


--
-- Name: user_tenants fk_tenants_user_tenants; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_tenants
    ADD CONSTRAINT fk_tenants_user_tenants FOREIGN KEY (tenant_id) REFERENCES public.tenants(id);


--
-- Name: user_tenants fk_user_tenants_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_tenants
    ADD CONSTRAINT fk_user_tenants_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict aXhPVpZh0ortDLcOBostXiIOpKvYqUKdJBaT6l19efXzJBS6dU2o1KsHdHBdI2p

