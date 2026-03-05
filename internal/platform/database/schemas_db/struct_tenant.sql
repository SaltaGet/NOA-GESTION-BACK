--
-- PostgreSQL database dump
--

\restrict UOEnI47Ah6HCC8yWO2iriunLk7bDorTRt1aqDVDZ7AGMGZ8ka9GexCzga4Tcgme

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

ALTER TABLE IF EXISTS ONLY public.stock_point_sales DROP CONSTRAINT IF EXISTS fk_stock_point_sales_point_sale;
ALTER TABLE IF EXISTS ONLY public.role_permissions DROP CONSTRAINT IF EXISTS fk_role_permissions_role;
ALTER TABLE IF EXISTS ONLY public.role_permissions DROP CONSTRAINT IF EXISTS fk_role_permissions_permission;
ALTER TABLE IF EXISTS ONLY public.stock_point_sales DROP CONSTRAINT IF EXISTS fk_products_stock_point_sales;
ALTER TABLE IF EXISTS ONLY public.deposits DROP CONSTRAINT IF EXISTS fk_products_stock_deposit;
ALTER TABLE IF EXISTS ONLY public.products DROP CONSTRAINT IF EXISTS fk_products_category;
ALTER TABLE IF EXISTS ONLY public.pay_incomes DROP CONSTRAINT IF EXISTS fk_pay_incomes_cash_register;
ALTER TABLE IF EXISTS ONLY public.pay_expense_others DROP CONSTRAINT IF EXISTS fk_pay_expense_others_expense_other;
ALTER TABLE IF EXISTS ONLY public.pay_expense_others DROP CONSTRAINT IF EXISTS fk_pay_expense_others_cash_register;
ALTER TABLE IF EXISTS ONLY public.pay_expense_buys DROP CONSTRAINT IF EXISTS fk_pay_expense_buys_cash_register;
ALTER TABLE IF EXISTS ONLY public.movement_stocks DROP CONSTRAINT IF EXISTS fk_movement_stocks_product;
ALTER TABLE IF EXISTS ONLY public.movement_stocks DROP CONSTRAINT IF EXISTS fk_movement_stocks_member;
ALTER TABLE IF EXISTS ONLY public.members DROP CONSTRAINT IF EXISTS fk_members_role;
ALTER TABLE IF EXISTS ONLY public.member_point_sales DROP CONSTRAINT IF EXISTS fk_member_point_sales_point_sale;
ALTER TABLE IF EXISTS ONLY public.member_point_sales DROP CONSTRAINT IF EXISTS fk_member_point_sales_member;
ALTER TABLE IF EXISTS ONLY public.income_sales DROP CONSTRAINT IF EXISTS fk_income_sales_point_sale;
ALTER TABLE IF EXISTS ONLY public.pay_incomes DROP CONSTRAINT IF EXISTS fk_income_sales_pay;
ALTER TABLE IF EXISTS ONLY public.income_sales DROP CONSTRAINT IF EXISTS fk_income_sales_member;
ALTER TABLE IF EXISTS ONLY public.income_sale_items DROP CONSTRAINT IF EXISTS fk_income_sales_items;
ALTER TABLE IF EXISTS ONLY public.income_sales DROP CONSTRAINT IF EXISTS fk_income_sales_invoice;
ALTER TABLE IF EXISTS ONLY public.income_sales DROP CONSTRAINT IF EXISTS fk_income_sales_client;
ALTER TABLE IF EXISTS ONLY public.income_sales DROP CONSTRAINT IF EXISTS fk_income_sales_cash_register;
ALTER TABLE IF EXISTS ONLY public.income_sale_items DROP CONSTRAINT IF EXISTS fk_income_sale_items_product;
ALTER TABLE IF EXISTS ONLY public.income_others DROP CONSTRAINT IF EXISTS fk_income_others_type_income;
ALTER TABLE IF EXISTS ONLY public.income_others DROP CONSTRAINT IF EXISTS fk_income_others_point_sale;
ALTER TABLE IF EXISTS ONLY public.income_others DROP CONSTRAINT IF EXISTS fk_income_others_member;
ALTER TABLE IF EXISTS ONLY public.income_others DROP CONSTRAINT IF EXISTS fk_income_others_cash_register;
ALTER TABLE IF EXISTS ONLY public.income_ecommerce_items DROP CONSTRAINT IF EXISTS fk_income_ecommerces_items;
ALTER TABLE IF EXISTS ONLY public.income_ecommerce_items DROP CONSTRAINT IF EXISTS fk_income_ecommerce_items_product;
ALTER TABLE IF EXISTS ONLY public.expense_others DROP CONSTRAINT IF EXISTS fk_expense_others_type_expense;
ALTER TABLE IF EXISTS ONLY public.expense_others DROP CONSTRAINT IF EXISTS fk_expense_others_point_sale;
ALTER TABLE IF EXISTS ONLY public.expense_others DROP CONSTRAINT IF EXISTS fk_expense_others_member;
ALTER TABLE IF EXISTS ONLY public.expense_others DROP CONSTRAINT IF EXISTS fk_expense_others_cash_register;
ALTER TABLE IF EXISTS ONLY public.expense_buys DROP CONSTRAINT IF EXISTS fk_expense_buys_supplier;
ALTER TABLE IF EXISTS ONLY public.pay_expense_buys DROP CONSTRAINT IF EXISTS fk_expense_buys_pay_expense_buy;
ALTER TABLE IF EXISTS ONLY public.expense_buys DROP CONSTRAINT IF EXISTS fk_expense_buys_member;
ALTER TABLE IF EXISTS ONLY public.expense_buy_items DROP CONSTRAINT IF EXISTS fk_expense_buys_expense_buy_item;
ALTER TABLE IF EXISTS ONLY public.expense_buy_items DROP CONSTRAINT IF EXISTS fk_expense_buy_items_product;
ALTER TABLE IF EXISTS ONLY public.pay_incomes DROP CONSTRAINT IF EXISTS fk_clients_pay;
ALTER TABLE IF EXISTS ONLY public.clients DROP CONSTRAINT IF EXISTS fk_clients_member_create;
ALTER TABLE IF EXISTS ONLY public.cash_registers DROP CONSTRAINT IF EXISTS fk_cash_registers_point_sale;
ALTER TABLE IF EXISTS ONLY public.cash_registers DROP CONSTRAINT IF EXISTS fk_cash_registers_member_open;
ALTER TABLE IF EXISTS ONLY public.cash_registers DROP CONSTRAINT IF EXISTS fk_cash_registers_member_close;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS fk_audit_logs_member;
DROP TRIGGER IF EXISTS tr_audit_type_incomes ON public.type_incomes;
DROP TRIGGER IF EXISTS tr_audit_type_expenses ON public.type_expenses;
DROP TRIGGER IF EXISTS tr_audit_suppliers ON public.suppliers;
DROP TRIGGER IF EXISTS tr_audit_stock_point_sales ON public.stock_point_sales;
DROP TRIGGER IF EXISTS tr_audit_roles ON public.roles;
DROP TRIGGER IF EXISTS tr_audit_role_permissions ON public.role_permissions;
DROP TRIGGER IF EXISTS tr_audit_products ON public.products;
DROP TRIGGER IF EXISTS tr_audit_point_sales ON public.point_sales;
DROP TRIGGER IF EXISTS tr_audit_permissions ON public.permissions;
DROP TRIGGER IF EXISTS tr_audit_pay_incomes ON public.pay_incomes;
DROP TRIGGER IF EXISTS tr_audit_pay_expense_others ON public.pay_expense_others;
DROP TRIGGER IF EXISTS tr_audit_pay_expense_buys ON public.pay_expense_buys;
DROP TRIGGER IF EXISTS tr_audit_movement_stocks ON public.movement_stocks;
DROP TRIGGER IF EXISTS tr_audit_members ON public.members;
DROP TRIGGER IF EXISTS tr_audit_member_point_sales ON public.member_point_sales;
DROP TRIGGER IF EXISTS tr_audit_invoices ON public.invoices;
DROP TRIGGER IF EXISTS tr_audit_income_sales ON public.income_sales;
DROP TRIGGER IF EXISTS tr_audit_income_sale_items ON public.income_sale_items;
DROP TRIGGER IF EXISTS tr_audit_income_others ON public.income_others;
DROP TRIGGER IF EXISTS tr_audit_income_ecommerces ON public.income_ecommerces;
DROP TRIGGER IF EXISTS tr_audit_income_ecommerce_items ON public.income_ecommerce_items;
DROP TRIGGER IF EXISTS tr_audit_expense_others ON public.expense_others;
DROP TRIGGER IF EXISTS tr_audit_expense_buys ON public.expense_buys;
DROP TRIGGER IF EXISTS tr_audit_expense_buy_items ON public.expense_buy_items;
DROP TRIGGER IF EXISTS tr_audit_deposits ON public.deposits;
DROP TRIGGER IF EXISTS tr_audit_clients ON public.clients;
DROP TRIGGER IF EXISTS tr_audit_categories ON public.categories;
DROP TRIGGER IF EXISTS tr_audit_cash_registers ON public.cash_registers;
DROP INDEX IF EXISTS public.idx_suppliers_deleted_at;
DROP INDEX IF EXISTS public.idx_products_code;
DROP INDEX IF EXISTS public.idx_point_sales_name;
DROP INDEX IF EXISTS public.idx_point_sales_delete_at;
DROP INDEX IF EXISTS public.idx_pay_incomes_income_sale_id;
DROP INDEX IF EXISTS public.idx_pay_incomes_client_id;
DROP INDEX IF EXISTS public.idx_pay_incomes_cash_register_id;
DROP INDEX IF EXISTS public.idx_pay_expense_others_expense_other_id;
DROP INDEX IF EXISTS public.idx_pay_expense_others_cash_register_id;
DROP INDEX IF EXISTS public.idx_pay_expense_buys_expense_buy_id;
DROP INDEX IF EXISTS public.idx_pay_expense_buys_cash_register_id;
DROP INDEX IF EXISTS public.idx_members_deleted_at;
DROP INDEX IF EXISTS public.idx_income_sale_items_product_id;
DROP INDEX IF EXISTS public.idx_income_sale_items_income_sale_id;
DROP INDEX IF EXISTS public.idx_income_ecommerces_payment_id;
DROP INDEX IF EXISTS public.idx_income_ecommerces_external_reference;
DROP INDEX IF EXISTS public.idx_income_ecommerce_items_product_id;
DROP INDEX IF EXISTS public.idx_income_ecommerce_items_income_ecommerce_id;
DROP INDEX IF EXISTS public.idx_clients_deleted_at;
DROP INDEX IF EXISTS public.idx_categories_name;
DROP INDEX IF EXISTS public.idx_categories_delete_at;
DROP INDEX IF EXISTS public.idx_audit_logs_transaction_id;
ALTER TABLE IF EXISTS ONLY public.type_incomes DROP CONSTRAINT IF EXISTS uni_type_incomes_name;
ALTER TABLE IF EXISTS ONLY public.type_expenses DROP CONSTRAINT IF EXISTS uni_type_expenses_name;
ALTER TABLE IF EXISTS ONLY public.suppliers DROP CONSTRAINT IF EXISTS uni_suppliers_identifier;
ALTER TABLE IF EXISTS ONLY public.suppliers DROP CONSTRAINT IF EXISTS uni_suppliers_email;
ALTER TABLE IF EXISTS ONLY public.roles DROP CONSTRAINT IF EXISTS uni_roles_name;
ALTER TABLE IF EXISTS ONLY public.members DROP CONSTRAINT IF EXISTS uni_members_username;
ALTER TABLE IF EXISTS ONLY public.members DROP CONSTRAINT IF EXISTS uni_members_email;
ALTER TABLE IF EXISTS ONLY public.clients DROP CONSTRAINT IF EXISTS uni_clients_identifier;
ALTER TABLE IF EXISTS ONLY public.clients DROP CONSTRAINT IF EXISTS uni_clients_email;
ALTER TABLE IF EXISTS ONLY public.type_incomes DROP CONSTRAINT IF EXISTS type_incomes_pkey;
ALTER TABLE IF EXISTS ONLY public.type_expenses DROP CONSTRAINT IF EXISTS type_expenses_pkey;
ALTER TABLE IF EXISTS ONLY public.suppliers DROP CONSTRAINT IF EXISTS suppliers_pkey;
ALTER TABLE IF EXISTS ONLY public.stock_point_sales DROP CONSTRAINT IF EXISTS stock_point_sales_pkey;
ALTER TABLE IF EXISTS ONLY public.schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.roles DROP CONSTRAINT IF EXISTS roles_pkey;
ALTER TABLE IF EXISTS ONLY public.role_permissions DROP CONSTRAINT IF EXISTS role_permissions_pkey;
ALTER TABLE IF EXISTS ONLY public.products DROP CONSTRAINT IF EXISTS products_pkey;
ALTER TABLE IF EXISTS ONLY public.point_sales DROP CONSTRAINT IF EXISTS point_sales_pkey;
ALTER TABLE IF EXISTS ONLY public.permissions DROP CONSTRAINT IF EXISTS permissions_pkey;
ALTER TABLE IF EXISTS ONLY public.pay_incomes DROP CONSTRAINT IF EXISTS pay_incomes_pkey;
ALTER TABLE IF EXISTS ONLY public.pay_expense_others DROP CONSTRAINT IF EXISTS pay_expense_others_pkey;
ALTER TABLE IF EXISTS ONLY public.pay_expense_buys DROP CONSTRAINT IF EXISTS pay_expense_buys_pkey;
ALTER TABLE IF EXISTS ONLY public.movement_stocks DROP CONSTRAINT IF EXISTS movement_stocks_pkey;
ALTER TABLE IF EXISTS ONLY public.members DROP CONSTRAINT IF EXISTS members_pkey;
ALTER TABLE IF EXISTS ONLY public.member_point_sales DROP CONSTRAINT IF EXISTS member_point_sales_pkey;
ALTER TABLE IF EXISTS ONLY public.invoices DROP CONSTRAINT IF EXISTS invoices_pkey;
ALTER TABLE IF EXISTS ONLY public.income_sales DROP CONSTRAINT IF EXISTS income_sales_pkey;
ALTER TABLE IF EXISTS ONLY public.income_sale_items DROP CONSTRAINT IF EXISTS income_sale_items_pkey;
ALTER TABLE IF EXISTS ONLY public.income_others DROP CONSTRAINT IF EXISTS income_others_pkey;
ALTER TABLE IF EXISTS ONLY public.income_ecommerces DROP CONSTRAINT IF EXISTS income_ecommerces_pkey;
ALTER TABLE IF EXISTS ONLY public.income_ecommerce_items DROP CONSTRAINT IF EXISTS income_ecommerce_items_pkey;
ALTER TABLE IF EXISTS ONLY public.expense_others DROP CONSTRAINT IF EXISTS expense_others_pkey;
ALTER TABLE IF EXISTS ONLY public.expense_buys DROP CONSTRAINT IF EXISTS expense_buys_pkey;
ALTER TABLE IF EXISTS ONLY public.expense_buy_items DROP CONSTRAINT IF EXISTS expense_buy_items_pkey;
ALTER TABLE IF EXISTS ONLY public.deposits DROP CONSTRAINT IF EXISTS deposits_pkey;
ALTER TABLE IF EXISTS ONLY public.clients DROP CONSTRAINT IF EXISTS clients_pkey;
ALTER TABLE IF EXISTS ONLY public.categories DROP CONSTRAINT IF EXISTS categories_pkey;
ALTER TABLE IF EXISTS ONLY public.cash_registers DROP CONSTRAINT IF EXISTS cash_registers_pkey;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS audit_logs_pkey;
ALTER TABLE IF EXISTS public.type_incomes ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.type_expenses ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.suppliers ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.stock_point_sales ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.roles ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.products ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.point_sales ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.permissions ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.pay_incomes ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.pay_expense_others ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.pay_expense_buys ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.movement_stocks ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.members ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.invoices ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.income_sales ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.income_sale_items ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.income_others ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.income_ecommerces ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.income_ecommerce_items ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.expense_others ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.expense_buys ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.expense_buy_items ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.deposits ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.clients ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.categories ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.cash_registers ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.audit_logs ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.type_incomes_id_seq;
DROP TABLE IF EXISTS public.type_incomes;
DROP SEQUENCE IF EXISTS public.type_expenses_id_seq;
DROP TABLE IF EXISTS public.type_expenses;
DROP SEQUENCE IF EXISTS public.suppliers_id_seq;
DROP TABLE IF EXISTS public.suppliers;
DROP SEQUENCE IF EXISTS public.stock_point_sales_id_seq;
DROP TABLE IF EXISTS public.stock_point_sales;
DROP TABLE IF EXISTS public.schema_migrations;
DROP SEQUENCE IF EXISTS public.roles_id_seq;
DROP TABLE IF EXISTS public.roles;
DROP TABLE IF EXISTS public.role_permissions;
DROP SEQUENCE IF EXISTS public.products_id_seq;
DROP TABLE IF EXISTS public.products;
DROP SEQUENCE IF EXISTS public.point_sales_id_seq;
DROP TABLE IF EXISTS public.point_sales;
DROP SEQUENCE IF EXISTS public.permissions_id_seq;
DROP TABLE IF EXISTS public.permissions;
DROP SEQUENCE IF EXISTS public.pay_incomes_id_seq;
DROP TABLE IF EXISTS public.pay_incomes;
DROP SEQUENCE IF EXISTS public.pay_expense_others_id_seq;
DROP TABLE IF EXISTS public.pay_expense_others;
DROP SEQUENCE IF EXISTS public.pay_expense_buys_id_seq;
DROP TABLE IF EXISTS public.pay_expense_buys;
DROP SEQUENCE IF EXISTS public.movement_stocks_id_seq;
DROP TABLE IF EXISTS public.movement_stocks;
DROP SEQUENCE IF EXISTS public.members_id_seq;
DROP TABLE IF EXISTS public.members;
DROP TABLE IF EXISTS public.member_point_sales;
DROP SEQUENCE IF EXISTS public.invoices_id_seq;
DROP TABLE IF EXISTS public.invoices;
DROP SEQUENCE IF EXISTS public.income_sales_id_seq;
DROP TABLE IF EXISTS public.income_sales;
DROP SEQUENCE IF EXISTS public.income_sale_items_id_seq;
DROP TABLE IF EXISTS public.income_sale_items;
DROP SEQUENCE IF EXISTS public.income_others_id_seq;
DROP TABLE IF EXISTS public.income_others;
DROP SEQUENCE IF EXISTS public.income_ecommerces_id_seq;
DROP TABLE IF EXISTS public.income_ecommerces;
DROP SEQUENCE IF EXISTS public.income_ecommerce_items_id_seq;
DROP TABLE IF EXISTS public.income_ecommerce_items;
DROP SEQUENCE IF EXISTS public.expense_others_id_seq;
DROP TABLE IF EXISTS public.expense_others;
DROP SEQUENCE IF EXISTS public.expense_buys_id_seq;
DROP TABLE IF EXISTS public.expense_buys;
DROP SEQUENCE IF EXISTS public.expense_buy_items_id_seq;
DROP TABLE IF EXISTS public.expense_buy_items;
DROP SEQUENCE IF EXISTS public.deposits_id_seq;
DROP TABLE IF EXISTS public.deposits;
DROP SEQUENCE IF EXISTS public.clients_id_seq;
DROP TABLE IF EXISTS public.clients;
DROP SEQUENCE IF EXISTS public.categories_id_seq;
DROP TABLE IF EXISTS public.categories;
DROP SEQUENCE IF EXISTS public.cash_registers_id_seq;
DROP TABLE IF EXISTS public.cash_registers;
DROP SEQUENCE IF EXISTS public.audit_logs_id_seq;
DROP TABLE IF EXISTS public.audit_logs;
DROP FUNCTION IF EXISTS public.audit_trigger_function();
--
-- Name: audit_trigger_function(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.audit_trigger_function() RETURNS trigger
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

    INSERT INTO audit_logs (
				transaction_id,
        member_id,
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


ALTER FUNCTION public.audit_trigger_function() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_logs (
    id bigint NOT NULL,
    transaction_id bigint,
    member_id bigint NOT NULL,
    method character varying(10) NOT NULL,
    path character varying(255) NOT NULL,
    old_value jsonb,
    new_value jsonb,
    created_at timestamp with time zone
);


ALTER TABLE public.audit_logs OWNER TO postgres;

--
-- Name: audit_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.audit_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.audit_logs_id_seq OWNER TO postgres;

--
-- Name: audit_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.audit_logs_id_seq OWNED BY public.audit_logs.id;


--
-- Name: cash_registers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cash_registers (
    id bigint NOT NULL,
    point_sale_id bigint NOT NULL,
    member_open_id bigint NOT NULL,
    open_amount numeric,
    hour_open timestamp with time zone,
    member_close_id bigint,
    close_amount numeric,
    hour_close timestamp with time zone,
    is_close boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.cash_registers OWNER TO postgres;

--
-- Name: cash_registers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cash_registers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cash_registers_id_seq OWNER TO postgres;

--
-- Name: cash_registers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cash_registers_id_seq OWNED BY public.cash_registers.id;


--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    delete_at timestamp with time zone
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: clients; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.clients (
    id bigint NOT NULL,
    first_name character varying(30) NOT NULL,
    last_name character varying(30) NOT NULL,
    company_name character varying(255),
    identifier character varying(20),
    email character varying(100),
    phone character varying(20),
    address character varying(255),
    responsability_front_iva character varying(255),
    member_create_id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.clients OWNER TO postgres;

--
-- Name: clients_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.clients_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.clients_id_seq OWNER TO postgres;

--
-- Name: clients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.clients_id_seq OWNED BY public.clients.id;


--
-- Name: deposits; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.deposits (
    id bigint NOT NULL,
    product_id bigint NOT NULL,
    stock numeric DEFAULT 0 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.deposits OWNER TO postgres;

--
-- Name: deposits_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.deposits_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.deposits_id_seq OWNER TO postgres;

--
-- Name: deposits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.deposits_id_seq OWNED BY public.deposits.id;


--
-- Name: expense_buy_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.expense_buy_items (
    id bigint NOT NULL,
    expense_buy_id bigint NOT NULL,
    product_id bigint NOT NULL,
    amount numeric NOT NULL,
    price numeric NOT NULL,
    discount numeric DEFAULT 0 NOT NULL,
    type_discount character varying(20) DEFAULT 'percent'::character varying NOT NULL,
    subtotal numeric NOT NULL,
    total numeric NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.expense_buy_items OWNER TO postgres;

--
-- Name: expense_buy_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.expense_buy_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.expense_buy_items_id_seq OWNER TO postgres;

--
-- Name: expense_buy_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.expense_buy_items_id_seq OWNED BY public.expense_buy_items.id;


--
-- Name: expense_buys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.expense_buys (
    id bigint NOT NULL,
    member_id bigint NOT NULL,
    supplier_id bigint NOT NULL,
    details text,
    subtotal numeric,
    discount numeric DEFAULT 0 NOT NULL,
    type_discount text DEFAULT 'percent'::text NOT NULL,
    total numeric,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.expense_buys OWNER TO postgres;

--
-- Name: expense_buys_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.expense_buys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.expense_buys_id_seq OWNER TO postgres;

--
-- Name: expense_buys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.expense_buys_id_seq OWNED BY public.expense_buys.id;


--
-- Name: expense_others; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.expense_others (
    id bigint NOT NULL,
    point_sale_id bigint,
    member_id bigint NOT NULL,
    cash_register_id bigint,
    details character varying(255),
    type_expense_id bigint NOT NULL,
    total numeric NOT NULL,
    pay_method character varying(30) DEFAULT 'efectivo'::character varying,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.expense_others OWNER TO postgres;

--
-- Name: expense_others_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.expense_others_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.expense_others_id_seq OWNER TO postgres;

--
-- Name: expense_others_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.expense_others_id_seq OWNED BY public.expense_others.id;


--
-- Name: income_ecommerce_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_ecommerce_items (
    id bigint NOT NULL,
    income_ecommerce_id bigint,
    product_id bigint,
    amount numeric NOT NULL,
    price_cost numeric NOT NULL,
    price numeric NOT NULL,
    discount numeric DEFAULT 0 NOT NULL,
    type_discount character varying(20) DEFAULT 'percent'::character varying NOT NULL,
    subtotal numeric NOT NULL,
    total numeric NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.income_ecommerce_items OWNER TO postgres;

--
-- Name: income_ecommerce_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.income_ecommerce_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.income_ecommerce_items_id_seq OWNER TO postgres;

--
-- Name: income_ecommerce_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.income_ecommerce_items_id_seq OWNED BY public.income_ecommerce_items.id;


--
-- Name: income_ecommerces; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_ecommerces (
    id bigint NOT NULL,
    payment_id character varying(255) NOT NULL,
    external_reference character varying(255) NOT NULL,
    status character varying(255) NOT NULL,
    total numeric NOT NULL,
    delivery_status character varying(255) NOT NULL,
    delivery_id character varying(255) NOT NULL,
    date_created character varying(255) NOT NULL,
    date_approved character varying(255) NOT NULL,
    transaction_amount numeric,
    net_received_amount numeric NOT NULL,
    payer_first_name character varying(255) NOT NULL,
    payer_last_name character varying(255) NOT NULL,
    payer_email character varying(255) NOT NULL,
    pay_method character varying(255) NOT NULL,
    operation_type character varying(255) NOT NULL,
    message character varying(255),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.income_ecommerces OWNER TO postgres;

--
-- Name: income_ecommerces_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.income_ecommerces_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.income_ecommerces_id_seq OWNER TO postgres;

--
-- Name: income_ecommerces_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.income_ecommerces_id_seq OWNED BY public.income_ecommerces.id;


--
-- Name: income_others; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_others (
    id bigint NOT NULL,
    point_sale_id bigint,
    member_id bigint,
    cash_register_id bigint,
    total numeric NOT NULL,
    type_income_id bigint NOT NULL,
    details character varying(255),
    method_income character varying(30) DEFAULT 'cash'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.income_others OWNER TO postgres;

--
-- Name: income_others_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.income_others_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.income_others_id_seq OWNER TO postgres;

--
-- Name: income_others_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.income_others_id_seq OWNED BY public.income_others.id;


--
-- Name: income_sale_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_sale_items (
    id bigint NOT NULL,
    income_sale_id bigint,
    product_id bigint,
    amount numeric NOT NULL,
    price_cost numeric NOT NULL,
    price numeric NOT NULL,
    discount numeric DEFAULT 0 NOT NULL,
    type_discount character varying(20) DEFAULT 'percent'::character varying NOT NULL,
    subtotal numeric NOT NULL,
    total numeric NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.income_sale_items OWNER TO postgres;

--
-- Name: income_sale_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.income_sale_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.income_sale_items_id_seq OWNER TO postgres;

--
-- Name: income_sale_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.income_sale_items_id_seq OWNED BY public.income_sale_items.id;


--
-- Name: income_sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.income_sales (
    id bigint NOT NULL,
    point_sale_id bigint NOT NULL,
    member_id bigint NOT NULL,
    client_id bigint NOT NULL,
    cash_register_id bigint NOT NULL,
    subtotal numeric NOT NULL,
    discount numeric DEFAULT 0 NOT NULL,
    type character varying(20) DEFAULT 'percent'::character varying NOT NULL,
    total numeric NOT NULL,
    is_budget boolean DEFAULT false NOT NULL,
    invoice_id bigint,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.income_sales OWNER TO postgres;

--
-- Name: income_sales_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.income_sales_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.income_sales_id_seq OWNER TO postgres;

--
-- Name: income_sales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.income_sales_id_seq OWNED BY public.income_sales.id;


--
-- Name: invoices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.invoices (
    id bigint NOT NULL,
    invoice_data jsonb,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.invoices OWNER TO postgres;

--
-- Name: invoices_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.invoices_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.invoices_id_seq OWNER TO postgres;

--
-- Name: invoices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.invoices_id_seq OWNED BY public.invoices.id;


--
-- Name: member_point_sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.member_point_sales (
    point_sale_id bigint NOT NULL,
    member_id bigint NOT NULL
);


ALTER TABLE public.member_point_sales OWNER TO postgres;

--
-- Name: members; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.members (
    id bigint NOT NULL,
    first_name character varying(30) NOT NULL,
    last_name character varying(30) NOT NULL,
    username character varying(30) NOT NULL,
    email character varying(100) NOT NULL,
    password character varying(255) NOT NULL,
    address character varying(255) DEFAULT NULL::character varying,
    phone character varying(20) DEFAULT NULL::character varying,
    is_admin boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    role_id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.members OWNER TO postgres;

--
-- Name: members_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.members_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.members_id_seq OWNER TO postgres;

--
-- Name: members_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.members_id_seq OWNED BY public.members.id;


--
-- Name: movement_stocks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.movement_stocks (
    id bigint NOT NULL,
    member_id bigint NOT NULL,
    product_id bigint NOT NULL,
    amount numeric NOT NULL,
    from_id bigint NOT NULL,
    from_type character varying(20) NOT NULL,
    to_id bigint NOT NULL,
    to_type character varying(20) NOT NULL,
    ignore_stock boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.movement_stocks OWNER TO postgres;

--
-- Name: movement_stocks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.movement_stocks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.movement_stocks_id_seq OWNER TO postgres;

--
-- Name: movement_stocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.movement_stocks_id_seq OWNED BY public.movement_stocks.id;


--
-- Name: pay_expense_buys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pay_expense_buys (
    id bigint NOT NULL,
    expense_buy_id bigint,
    cash_register_id bigint,
    total numeric NOT NULL,
    method_pay character varying(30) DEFAULT 'cash'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.pay_expense_buys OWNER TO postgres;

--
-- Name: pay_expense_buys_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pay_expense_buys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pay_expense_buys_id_seq OWNER TO postgres;

--
-- Name: pay_expense_buys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pay_expense_buys_id_seq OWNED BY public.pay_expense_buys.id;


--
-- Name: pay_expense_others; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pay_expense_others (
    id bigint NOT NULL,
    expense_other_id bigint,
    cash_register_id bigint,
    total numeric NOT NULL,
    method_pay character varying(30) DEFAULT 'cash'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.pay_expense_others OWNER TO postgres;

--
-- Name: pay_expense_others_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pay_expense_others_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pay_expense_others_id_seq OWNER TO postgres;

--
-- Name: pay_expense_others_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pay_expense_others_id_seq OWNED BY public.pay_expense_others.id;


--
-- Name: pay_incomes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pay_incomes (
    id bigint NOT NULL,
    income_sale_id bigint,
    cash_register_id bigint,
    client_id bigint,
    total numeric NOT NULL,
    method_pay character varying(30) DEFAULT 'cash'::character varying NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.pay_incomes OWNER TO postgres;

--
-- Name: pay_incomes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.pay_incomes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.pay_incomes_id_seq OWNER TO postgres;

--
-- Name: pay_incomes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.pay_incomes_id_seq OWNED BY public.pay_incomes.id;


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.permissions (
    id bigint NOT NULL,
    code character varying(50),
    name character varying(100) NOT NULL,
    details text NOT NULL,
    "group" character varying(50) NOT NULL,
    environment character varying(20) NOT NULL
);


ALTER TABLE public.permissions OWNER TO postgres;

--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.permissions_id_seq OWNER TO postgres;

--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: point_sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.point_sales (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    description character varying(200),
    number bigint NOT NULL,
    is_deposit boolean DEFAULT false NOT NULL,
    is_main boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    delete_at timestamp with time zone
);


ALTER TABLE public.point_sales OWNER TO postgres;

--
-- Name: point_sales_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.point_sales_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.point_sales_id_seq OWNER TO postgres;

--
-- Name: point_sales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.point_sales_id_seq OWNED BY public.point_sales.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    price numeric,
    category_id bigint NOT NULL,
    primary_image character varying(255) DEFAULT NULL::character varying,
    secondary_images text,
    is_visible boolean DEFAULT false NOT NULL,
    notifier boolean DEFAULT false NOT NULL,
    min_amount numeric DEFAULT 0 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.role_permissions (
    role_id bigint NOT NULL,
    permission_id bigint NOT NULL
);


ALTER TABLE public.role_permissions OWNER TO postgres;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id bigint NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_seq OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: stock_point_sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock_point_sales (
    id bigint NOT NULL,
    product_id bigint NOT NULL,
    point_sale_id bigint NOT NULL,
    stock numeric DEFAULT 0 NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.stock_point_sales OWNER TO postgres;

--
-- Name: stock_point_sales_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stock_point_sales_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.stock_point_sales_id_seq OWNER TO postgres;

--
-- Name: stock_point_sales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stock_point_sales_id_seq OWNED BY public.stock_point_sales.id;


--
-- Name: suppliers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.suppliers (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    company_name character varying(100) NOT NULL,
    identifier character varying(20),
    address character varying(255),
    debt_limit numeric,
    email character varying(100),
    phone character varying(20),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.suppliers OWNER TO postgres;

--
-- Name: suppliers_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.suppliers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.suppliers_id_seq OWNER TO postgres;

--
-- Name: suppliers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.suppliers_id_seq OWNED BY public.suppliers.id;


--
-- Name: type_expenses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.type_expenses (
    id bigint NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.type_expenses OWNER TO postgres;

--
-- Name: type_expenses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.type_expenses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.type_expenses_id_seq OWNER TO postgres;

--
-- Name: type_expenses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.type_expenses_id_seq OWNED BY public.type_expenses.id;


--
-- Name: type_incomes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.type_incomes (
    id bigint NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.type_incomes OWNER TO postgres;

--
-- Name: type_incomes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.type_incomes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.type_incomes_id_seq OWNER TO postgres;

--
-- Name: type_incomes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.type_incomes_id_seq OWNED BY public.type_incomes.id;


--
-- Name: audit_logs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs ALTER COLUMN id SET DEFAULT nextval('public.audit_logs_id_seq'::regclass);


--
-- Name: cash_registers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cash_registers ALTER COLUMN id SET DEFAULT nextval('public.cash_registers_id_seq'::regclass);


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: clients id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clients ALTER COLUMN id SET DEFAULT nextval('public.clients_id_seq'::regclass);


--
-- Name: deposits id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deposits ALTER COLUMN id SET DEFAULT nextval('public.deposits_id_seq'::regclass);


--
-- Name: expense_buy_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_buy_items ALTER COLUMN id SET DEFAULT nextval('public.expense_buy_items_id_seq'::regclass);


--
-- Name: expense_buys id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_buys ALTER COLUMN id SET DEFAULT nextval('public.expense_buys_id_seq'::regclass);


--
-- Name: expense_others id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_others ALTER COLUMN id SET DEFAULT nextval('public.expense_others_id_seq'::regclass);


--
-- Name: income_ecommerce_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_ecommerce_items ALTER COLUMN id SET DEFAULT nextval('public.income_ecommerce_items_id_seq'::regclass);


--
-- Name: income_ecommerces id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_ecommerces ALTER COLUMN id SET DEFAULT nextval('public.income_ecommerces_id_seq'::regclass);


--
-- Name: income_others id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_others ALTER COLUMN id SET DEFAULT nextval('public.income_others_id_seq'::regclass);


--
-- Name: income_sale_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sale_items ALTER COLUMN id SET DEFAULT nextval('public.income_sale_items_id_seq'::regclass);


--
-- Name: income_sales id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sales ALTER COLUMN id SET DEFAULT nextval('public.income_sales_id_seq'::regclass);


--
-- Name: invoices id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices ALTER COLUMN id SET DEFAULT nextval('public.invoices_id_seq'::regclass);


--
-- Name: members id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.members ALTER COLUMN id SET DEFAULT nextval('public.members_id_seq'::regclass);


--
-- Name: movement_stocks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.movement_stocks ALTER COLUMN id SET DEFAULT nextval('public.movement_stocks_id_seq'::regclass);


--
-- Name: pay_expense_buys id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_expense_buys ALTER COLUMN id SET DEFAULT nextval('public.pay_expense_buys_id_seq'::regclass);


--
-- Name: pay_expense_others id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_expense_others ALTER COLUMN id SET DEFAULT nextval('public.pay_expense_others_id_seq'::regclass);


--
-- Name: pay_incomes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_incomes ALTER COLUMN id SET DEFAULT nextval('public.pay_incomes_id_seq'::regclass);


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: point_sales id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.point_sales ALTER COLUMN id SET DEFAULT nextval('public.point_sales_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: stock_point_sales id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_point_sales ALTER COLUMN id SET DEFAULT nextval('public.stock_point_sales_id_seq'::regclass);


--
-- Name: suppliers id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suppliers ALTER COLUMN id SET DEFAULT nextval('public.suppliers_id_seq'::regclass);


--
-- Name: type_expenses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.type_expenses ALTER COLUMN id SET DEFAULT nextval('public.type_expenses_id_seq'::regclass);


--
-- Name: type_incomes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.type_incomes ALTER COLUMN id SET DEFAULT nextval('public.type_incomes_id_seq'::regclass);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: cash_registers cash_registers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cash_registers
    ADD CONSTRAINT cash_registers_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: clients clients_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (id);


--
-- Name: deposits deposits_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deposits
    ADD CONSTRAINT deposits_pkey PRIMARY KEY (id);


--
-- Name: expense_buy_items expense_buy_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_buy_items
    ADD CONSTRAINT expense_buy_items_pkey PRIMARY KEY (id);


--
-- Name: expense_buys expense_buys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_buys
    ADD CONSTRAINT expense_buys_pkey PRIMARY KEY (id);


--
-- Name: expense_others expense_others_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_others
    ADD CONSTRAINT expense_others_pkey PRIMARY KEY (id);


--
-- Name: income_ecommerce_items income_ecommerce_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_ecommerce_items
    ADD CONSTRAINT income_ecommerce_items_pkey PRIMARY KEY (id);


--
-- Name: income_ecommerces income_ecommerces_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_ecommerces
    ADD CONSTRAINT income_ecommerces_pkey PRIMARY KEY (id);


--
-- Name: income_others income_others_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_others
    ADD CONSTRAINT income_others_pkey PRIMARY KEY (id);


--
-- Name: income_sale_items income_sale_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sale_items
    ADD CONSTRAINT income_sale_items_pkey PRIMARY KEY (id);


--
-- Name: income_sales income_sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sales
    ADD CONSTRAINT income_sales_pkey PRIMARY KEY (id);


--
-- Name: invoices invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (id);


--
-- Name: member_point_sales member_point_sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.member_point_sales
    ADD CONSTRAINT member_point_sales_pkey PRIMARY KEY (point_sale_id, member_id);


--
-- Name: members members_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_pkey PRIMARY KEY (id);


--
-- Name: movement_stocks movement_stocks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.movement_stocks
    ADD CONSTRAINT movement_stocks_pkey PRIMARY KEY (id);


--
-- Name: pay_expense_buys pay_expense_buys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_expense_buys
    ADD CONSTRAINT pay_expense_buys_pkey PRIMARY KEY (id);


--
-- Name: pay_expense_others pay_expense_others_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_expense_others
    ADD CONSTRAINT pay_expense_others_pkey PRIMARY KEY (id);


--
-- Name: pay_incomes pay_incomes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_incomes
    ADD CONSTRAINT pay_incomes_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: point_sales point_sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.point_sales
    ADD CONSTRAINT point_sales_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (role_id, permission_id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: stock_point_sales stock_point_sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_point_sales
    ADD CONSTRAINT stock_point_sales_pkey PRIMARY KEY (id);


--
-- Name: suppliers suppliers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT suppliers_pkey PRIMARY KEY (id);


--
-- Name: type_expenses type_expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.type_expenses
    ADD CONSTRAINT type_expenses_pkey PRIMARY KEY (id);


--
-- Name: type_incomes type_incomes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.type_incomes
    ADD CONSTRAINT type_incomes_pkey PRIMARY KEY (id);


--
-- Name: clients uni_clients_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT uni_clients_email UNIQUE (email);


--
-- Name: clients uni_clients_identifier; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT uni_clients_identifier UNIQUE (identifier);


--
-- Name: members uni_members_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT uni_members_email UNIQUE (email);


--
-- Name: members uni_members_username; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT uni_members_username UNIQUE (username);


--
-- Name: roles uni_roles_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT uni_roles_name UNIQUE (name);


--
-- Name: suppliers uni_suppliers_email; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT uni_suppliers_email UNIQUE (email);


--
-- Name: suppliers uni_suppliers_identifier; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.suppliers
    ADD CONSTRAINT uni_suppliers_identifier UNIQUE (identifier);


--
-- Name: type_expenses uni_type_expenses_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.type_expenses
    ADD CONSTRAINT uni_type_expenses_name UNIQUE (name);


--
-- Name: type_incomes uni_type_incomes_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.type_incomes
    ADD CONSTRAINT uni_type_incomes_name UNIQUE (name);


--
-- Name: idx_audit_logs_transaction_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_audit_logs_transaction_id ON public.audit_logs USING btree (transaction_id);


--
-- Name: idx_categories_delete_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_categories_delete_at ON public.categories USING btree (delete_at);


--
-- Name: idx_categories_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_categories_name ON public.categories USING btree (name);


--
-- Name: idx_clients_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_clients_deleted_at ON public.clients USING btree (deleted_at);


--
-- Name: idx_income_ecommerce_items_income_ecommerce_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_income_ecommerce_items_income_ecommerce_id ON public.income_ecommerce_items USING btree (income_ecommerce_id);


--
-- Name: idx_income_ecommerce_items_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_income_ecommerce_items_product_id ON public.income_ecommerce_items USING btree (product_id);


--
-- Name: idx_income_ecommerces_external_reference; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_income_ecommerces_external_reference ON public.income_ecommerces USING btree (external_reference);


--
-- Name: idx_income_ecommerces_payment_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_income_ecommerces_payment_id ON public.income_ecommerces USING btree (payment_id);


--
-- Name: idx_income_sale_items_income_sale_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_income_sale_items_income_sale_id ON public.income_sale_items USING btree (income_sale_id);


--
-- Name: idx_income_sale_items_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_income_sale_items_product_id ON public.income_sale_items USING btree (product_id);


--
-- Name: idx_members_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_members_deleted_at ON public.members USING btree (deleted_at);


--
-- Name: idx_pay_expense_buys_cash_register_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pay_expense_buys_cash_register_id ON public.pay_expense_buys USING btree (cash_register_id);


--
-- Name: idx_pay_expense_buys_expense_buy_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pay_expense_buys_expense_buy_id ON public.pay_expense_buys USING btree (expense_buy_id);


--
-- Name: idx_pay_expense_others_cash_register_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pay_expense_others_cash_register_id ON public.pay_expense_others USING btree (cash_register_id);


--
-- Name: idx_pay_expense_others_expense_other_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pay_expense_others_expense_other_id ON public.pay_expense_others USING btree (expense_other_id);


--
-- Name: idx_pay_incomes_cash_register_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pay_incomes_cash_register_id ON public.pay_incomes USING btree (cash_register_id);


--
-- Name: idx_pay_incomes_client_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pay_incomes_client_id ON public.pay_incomes USING btree (client_id);


--
-- Name: idx_pay_incomes_income_sale_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pay_incomes_income_sale_id ON public.pay_incomes USING btree (income_sale_id);


--
-- Name: idx_point_sales_delete_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_point_sales_delete_at ON public.point_sales USING btree (delete_at);


--
-- Name: idx_point_sales_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_point_sales_name ON public.point_sales USING btree (name);


--
-- Name: idx_products_code; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX idx_products_code ON public.products USING btree (code);


--
-- Name: idx_suppliers_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_suppliers_deleted_at ON public.suppliers USING btree (deleted_at);


--
-- Name: cash_registers tr_audit_cash_registers; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_cash_registers AFTER INSERT OR DELETE OR UPDATE ON public.cash_registers FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: categories tr_audit_categories; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_categories AFTER INSERT OR DELETE OR UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: clients tr_audit_clients; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_clients AFTER INSERT OR DELETE OR UPDATE ON public.clients FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: deposits tr_audit_deposits; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_deposits AFTER INSERT OR DELETE OR UPDATE ON public.deposits FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: expense_buy_items tr_audit_expense_buy_items; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_expense_buy_items AFTER INSERT OR DELETE OR UPDATE ON public.expense_buy_items FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: expense_buys tr_audit_expense_buys; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_expense_buys AFTER INSERT OR DELETE OR UPDATE ON public.expense_buys FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: expense_others tr_audit_expense_others; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_expense_others AFTER INSERT OR DELETE OR UPDATE ON public.expense_others FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: income_ecommerce_items tr_audit_income_ecommerce_items; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_income_ecommerce_items AFTER INSERT OR DELETE OR UPDATE ON public.income_ecommerce_items FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: income_ecommerces tr_audit_income_ecommerces; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_income_ecommerces AFTER INSERT OR DELETE OR UPDATE ON public.income_ecommerces FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: income_others tr_audit_income_others; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_income_others AFTER INSERT OR DELETE OR UPDATE ON public.income_others FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: income_sale_items tr_audit_income_sale_items; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_income_sale_items AFTER INSERT OR DELETE OR UPDATE ON public.income_sale_items FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: income_sales tr_audit_income_sales; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_income_sales AFTER INSERT OR DELETE OR UPDATE ON public.income_sales FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: invoices tr_audit_invoices; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_invoices AFTER INSERT OR DELETE OR UPDATE ON public.invoices FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: member_point_sales tr_audit_member_point_sales; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_member_point_sales AFTER INSERT OR DELETE OR UPDATE ON public.member_point_sales FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: members tr_audit_members; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_members AFTER INSERT OR DELETE OR UPDATE ON public.members FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: movement_stocks tr_audit_movement_stocks; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_movement_stocks AFTER INSERT OR DELETE OR UPDATE ON public.movement_stocks FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: pay_expense_buys tr_audit_pay_expense_buys; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_pay_expense_buys AFTER INSERT OR DELETE OR UPDATE ON public.pay_expense_buys FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: pay_expense_others tr_audit_pay_expense_others; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_pay_expense_others AFTER INSERT OR DELETE OR UPDATE ON public.pay_expense_others FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: pay_incomes tr_audit_pay_incomes; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_pay_incomes AFTER INSERT OR DELETE OR UPDATE ON public.pay_incomes FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: permissions tr_audit_permissions; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_permissions AFTER INSERT OR DELETE OR UPDATE ON public.permissions FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: point_sales tr_audit_point_sales; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_point_sales AFTER INSERT OR DELETE OR UPDATE ON public.point_sales FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: products tr_audit_products; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_products AFTER INSERT OR DELETE OR UPDATE ON public.products FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: role_permissions tr_audit_role_permissions; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_role_permissions AFTER INSERT OR DELETE OR UPDATE ON public.role_permissions FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: roles tr_audit_roles; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_roles AFTER INSERT OR DELETE OR UPDATE ON public.roles FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: stock_point_sales tr_audit_stock_point_sales; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_stock_point_sales AFTER INSERT OR DELETE OR UPDATE ON public.stock_point_sales FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: suppliers tr_audit_suppliers; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_suppliers AFTER INSERT OR DELETE OR UPDATE ON public.suppliers FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: type_expenses tr_audit_type_expenses; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_type_expenses AFTER INSERT OR DELETE OR UPDATE ON public.type_expenses FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: type_incomes tr_audit_type_incomes; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER tr_audit_type_incomes AFTER INSERT OR DELETE OR UPDATE ON public.type_incomes FOR EACH ROW EXECUTE FUNCTION public.audit_trigger_function();


--
-- Name: audit_logs fk_audit_logs_member; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT fk_audit_logs_member FOREIGN KEY (member_id) REFERENCES public.members(id);


--
-- Name: cash_registers fk_cash_registers_member_close; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cash_registers
    ADD CONSTRAINT fk_cash_registers_member_close FOREIGN KEY (member_close_id) REFERENCES public.members(id);


--
-- Name: cash_registers fk_cash_registers_member_open; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cash_registers
    ADD CONSTRAINT fk_cash_registers_member_open FOREIGN KEY (member_open_id) REFERENCES public.members(id);


--
-- Name: cash_registers fk_cash_registers_point_sale; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cash_registers
    ADD CONSTRAINT fk_cash_registers_point_sale FOREIGN KEY (point_sale_id) REFERENCES public.point_sales(id);


--
-- Name: clients fk_clients_member_create; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT fk_clients_member_create FOREIGN KEY (member_create_id) REFERENCES public.members(id);


--
-- Name: pay_incomes fk_clients_pay; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_incomes
    ADD CONSTRAINT fk_clients_pay FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: expense_buy_items fk_expense_buy_items_product; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_buy_items
    ADD CONSTRAINT fk_expense_buy_items_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: expense_buy_items fk_expense_buys_expense_buy_item; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_buy_items
    ADD CONSTRAINT fk_expense_buys_expense_buy_item FOREIGN KEY (expense_buy_id) REFERENCES public.expense_buys(id);


--
-- Name: expense_buys fk_expense_buys_member; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_buys
    ADD CONSTRAINT fk_expense_buys_member FOREIGN KEY (member_id) REFERENCES public.members(id);


--
-- Name: pay_expense_buys fk_expense_buys_pay_expense_buy; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_expense_buys
    ADD CONSTRAINT fk_expense_buys_pay_expense_buy FOREIGN KEY (expense_buy_id) REFERENCES public.expense_buys(id);


--
-- Name: expense_buys fk_expense_buys_supplier; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_buys
    ADD CONSTRAINT fk_expense_buys_supplier FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id);


--
-- Name: expense_others fk_expense_others_cash_register; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_others
    ADD CONSTRAINT fk_expense_others_cash_register FOREIGN KEY (cash_register_id) REFERENCES public.cash_registers(id);


--
-- Name: expense_others fk_expense_others_member; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_others
    ADD CONSTRAINT fk_expense_others_member FOREIGN KEY (member_id) REFERENCES public.members(id);


--
-- Name: expense_others fk_expense_others_point_sale; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_others
    ADD CONSTRAINT fk_expense_others_point_sale FOREIGN KEY (point_sale_id) REFERENCES public.point_sales(id);


--
-- Name: expense_others fk_expense_others_type_expense; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.expense_others
    ADD CONSTRAINT fk_expense_others_type_expense FOREIGN KEY (type_expense_id) REFERENCES public.type_expenses(id);


--
-- Name: income_ecommerce_items fk_income_ecommerce_items_product; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_ecommerce_items
    ADD CONSTRAINT fk_income_ecommerce_items_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: income_ecommerce_items fk_income_ecommerces_items; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_ecommerce_items
    ADD CONSTRAINT fk_income_ecommerces_items FOREIGN KEY (income_ecommerce_id) REFERENCES public.income_ecommerces(id);


--
-- Name: income_others fk_income_others_cash_register; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_others
    ADD CONSTRAINT fk_income_others_cash_register FOREIGN KEY (cash_register_id) REFERENCES public.cash_registers(id);


--
-- Name: income_others fk_income_others_member; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_others
    ADD CONSTRAINT fk_income_others_member FOREIGN KEY (member_id) REFERENCES public.members(id);


--
-- Name: income_others fk_income_others_point_sale; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_others
    ADD CONSTRAINT fk_income_others_point_sale FOREIGN KEY (point_sale_id) REFERENCES public.point_sales(id);


--
-- Name: income_others fk_income_others_type_income; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_others
    ADD CONSTRAINT fk_income_others_type_income FOREIGN KEY (type_income_id) REFERENCES public.type_incomes(id);


--
-- Name: income_sale_items fk_income_sale_items_product; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sale_items
    ADD CONSTRAINT fk_income_sale_items_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: income_sales fk_income_sales_cash_register; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sales
    ADD CONSTRAINT fk_income_sales_cash_register FOREIGN KEY (cash_register_id) REFERENCES public.cash_registers(id);


--
-- Name: income_sales fk_income_sales_client; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sales
    ADD CONSTRAINT fk_income_sales_client FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- Name: income_sales fk_income_sales_invoice; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sales
    ADD CONSTRAINT fk_income_sales_invoice FOREIGN KEY (invoice_id) REFERENCES public.invoices(id);


--
-- Name: income_sale_items fk_income_sales_items; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sale_items
    ADD CONSTRAINT fk_income_sales_items FOREIGN KEY (income_sale_id) REFERENCES public.income_sales(id);


--
-- Name: income_sales fk_income_sales_member; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sales
    ADD CONSTRAINT fk_income_sales_member FOREIGN KEY (member_id) REFERENCES public.members(id);


--
-- Name: pay_incomes fk_income_sales_pay; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_incomes
    ADD CONSTRAINT fk_income_sales_pay FOREIGN KEY (income_sale_id) REFERENCES public.income_sales(id);


--
-- Name: income_sales fk_income_sales_point_sale; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.income_sales
    ADD CONSTRAINT fk_income_sales_point_sale FOREIGN KEY (point_sale_id) REFERENCES public.point_sales(id);


--
-- Name: member_point_sales fk_member_point_sales_member; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.member_point_sales
    ADD CONSTRAINT fk_member_point_sales_member FOREIGN KEY (member_id) REFERENCES public.members(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: member_point_sales fk_member_point_sales_point_sale; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.member_point_sales
    ADD CONSTRAINT fk_member_point_sales_point_sale FOREIGN KEY (point_sale_id) REFERENCES public.point_sales(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: members fk_members_role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT fk_members_role FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: movement_stocks fk_movement_stocks_member; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.movement_stocks
    ADD CONSTRAINT fk_movement_stocks_member FOREIGN KEY (member_id) REFERENCES public.members(id);


--
-- Name: movement_stocks fk_movement_stocks_product; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.movement_stocks
    ADD CONSTRAINT fk_movement_stocks_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: pay_expense_buys fk_pay_expense_buys_cash_register; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_expense_buys
    ADD CONSTRAINT fk_pay_expense_buys_cash_register FOREIGN KEY (cash_register_id) REFERENCES public.cash_registers(id);


--
-- Name: pay_expense_others fk_pay_expense_others_cash_register; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_expense_others
    ADD CONSTRAINT fk_pay_expense_others_cash_register FOREIGN KEY (cash_register_id) REFERENCES public.cash_registers(id);


--
-- Name: pay_expense_others fk_pay_expense_others_expense_other; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_expense_others
    ADD CONSTRAINT fk_pay_expense_others_expense_other FOREIGN KEY (expense_other_id) REFERENCES public.expense_others(id);


--
-- Name: pay_incomes fk_pay_incomes_cash_register; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pay_incomes
    ADD CONSTRAINT fk_pay_incomes_cash_register FOREIGN KEY (cash_register_id) REFERENCES public.cash_registers(id);


--
-- Name: products fk_products_category; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: deposits fk_products_stock_deposit; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.deposits
    ADD CONSTRAINT fk_products_stock_deposit FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: stock_point_sales fk_products_stock_point_sales; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_point_sales
    ADD CONSTRAINT fk_products_stock_point_sales FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: role_permissions fk_role_permissions_permission; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: role_permissions fk_role_permissions_role; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.role_permissions
    ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES public.roles(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: stock_point_sales fk_stock_point_sales_point_sale; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock_point_sales
    ADD CONSTRAINT fk_stock_point_sales_point_sale FOREIGN KEY (point_sale_id) REFERENCES public.point_sales(id);


--
-- PostgreSQL database dump complete
--

\unrestrict UOEnI47Ah6HCC8yWO2iriunLk7bDorTRt1aqDVDZ7AGMGZ8ka9GexCzga4Tcgme

