--
-- PostgreSQL database dump
--

-- Dumped from database version 16.13 (Ubuntu 16.13-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 17.0

-- Started on 2026-03-15 17:12:25

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
-- TOC entry 7 (class 2615 OID 34065)
-- Name: public; Type: SCHEMA; Schema: -; Owner: super_shop_dev_user
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO super_shop_dev_user;

--
-- TOC entry 3871 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: super_shop_dev_user
--

COMMENT ON SCHEMA public IS '';


--
-- TOC entry 3 (class 3079 OID 34077)
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- TOC entry 3873 (class 0 OID 0)
-- Dependencies: 3
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- TOC entry 2 (class 3079 OID 34066)
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 3874 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- TOC entry 941 (class 1247 OID 34355)
-- Name: order_status; Type: TYPE; Schema: public; Owner: super_shop_dev_user
--

CREATE TYPE public.order_status AS ENUM (
    'pending',
    'confirmed',
    'processing',
    'shipped',
    'delivered',
    'cancelled',
    'returned'
);


ALTER TYPE public.order_status OWNER TO super_shop_dev_user;

--
-- TOC entry 944 (class 1247 OID 34370)
-- Name: payment_status; Type: TYPE; Schema: public; Owner: super_shop_dev_user
--

CREATE TYPE public.payment_status AS ENUM (
    'pending',
    'paid',
    'failed',
    'refunded',
    'partially_refunded'
);


ALTER TYPE public.payment_status OWNER TO super_shop_dev_user;

--
-- TOC entry 309 (class 1255 OID 34442)
-- Name: generate_order_number(); Type: FUNCTION; Schema: public; Owner: super_shop_dev_user
--

CREATE FUNCTION public.generate_order_number() RETURNS text
    LANGUAGE plpgsql
    AS $$
DECLARE
    prefix TEXT := 'ORD';
    date_part TEXT;
    seq_num INT;
    order_num TEXT;
BEGIN
    date_part := TO_CHAR(NOW(), 'YYYYMMDD');
    
    -- Get count of orders today + 1
    SELECT COUNT(*) + 1 INTO seq_num
    FROM orders
    WHERE DATE(created_at) = CURRENT_DATE;
    
    order_num := prefix || '-' || date_part || '-' || LPAD(seq_num::TEXT, 4, '0');
    
    RETURN order_num;
END;
$$;


ALTER FUNCTION public.generate_order_number() OWNER TO super_shop_dev_user;

--
-- TOC entry 308 (class 1255 OID 34440)
-- Name: update_orders_updated_at(); Type: FUNCTION; Schema: public; Owner: super_shop_dev_user
--

CREATE FUNCTION public.update_orders_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_orders_updated_at() OWNER TO super_shop_dev_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 242 (class 1259 OID 34837)
-- Name: attributes; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.attributes (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.attributes OWNER TO super_shop_dev_user;

--
-- TOC entry 241 (class 1259 OID 34836)
-- Name: attributes_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.attributes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.attributes_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3875 (class 0 OID 0)
-- Dependencies: 241
-- Name: attributes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.attributes_id_seq OWNED BY public.attributes.id;


--
-- TOC entry 238 (class 1259 OID 34470)
-- Name: branches; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.branches (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    country character varying(100) DEFAULT ''::character varying NOT NULL,
    city character varying(100) NOT NULL,
    address text NOT NULL,
    mobile character varying(50) DEFAULT ''::character varying NOT NULL,
    telephone character varying(50) DEFAULT ''::character varying NOT NULL,
    email character varying(255) DEFAULT ''::character varying NOT NULL,
    latitude double precision NOT NULL,
    longitude double precision NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.branches OWNER TO super_shop_dev_user;

--
-- TOC entry 237 (class 1259 OID 34469)
-- Name: branches_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.branches_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.branches_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3876 (class 0 OID 0)
-- Dependencies: 237
-- Name: branches_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.branches_id_seq OWNED BY public.branches.id;


--
-- TOC entry 250 (class 1259 OID 34971)
-- Name: brands; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.brands (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    thumbnail character varying(255),
    priority smallint DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.brands OWNER TO super_shop_dev_user;

--
-- TOC entry 249 (class 1259 OID 34970)
-- Name: brands_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.brands_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.brands_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3877 (class 0 OID 0)
-- Dependencies: 249
-- Name: brands_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.brands_id_seq OWNED BY public.brands.id;


--
-- TOC entry 222 (class 1259 OID 34160)
-- Name: categories; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.categories (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    thumbnail character varying(255),
    priority smallint DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.categories OWNER TO super_shop_dev_user;

--
-- TOC entry 221 (class 1259 OID 34159)
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3878 (class 0 OID 0)
-- Dependencies: 221
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- TOC entry 230 (class 1259 OID 34331)
-- Name: customers; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.customers (
    id bigint NOT NULL,
    name character varying(255),
    f_name character varying(255),
    l_name character varying(255),
    phone character varying(25) NOT NULL,
    image text DEFAULT 'def.png'::character varying NOT NULL,
    email character varying(255),
    email_verified_at timestamp with time zone,
    password text NOT NULL,
    remember_token character varying(100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    street_address character varying(250),
    country character varying(50),
    city character varying(50),
    zip character varying(20),
    house_no character varying(50),
    apartment_no character varying(50),
    cm_firebase_token character varying(255),
    is_active boolean DEFAULT true NOT NULL,
    payment_card_last_four character varying(255),
    payment_card_brand character varying(255),
    payment_card_fawry_token text,
    login_medium character varying(255),
    social_id character varying(255),
    is_phone_verified boolean DEFAULT false NOT NULL,
    temporary_token character varying(255),
    is_email_verified boolean DEFAULT false NOT NULL,
    wallet_balance numeric(8,2),
    loyalty_point numeric(18,4) DEFAULT 0.0000,
    login_hit_count smallint DEFAULT 0 NOT NULL,
    is_temp_blocked boolean DEFAULT false NOT NULL,
    temp_block_time timestamp with time zone,
    referral_code character varying(255),
    referred_by integer,
    app_language character varying(255) DEFAULT 'en'::character varying NOT NULL,
    is_retailer boolean DEFAULT false NOT NULL
);


ALTER TABLE public.customers OWNER TO super_shop_dev_user;

--
-- TOC entry 229 (class 1259 OID 34330)
-- Name: customers_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.customers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.customers_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3879 (class 0 OID 0)
-- Dependencies: 229
-- Name: customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.customers_id_seq OWNED BY public.customers.id;


--
-- TOC entry 218 (class 1259 OID 34115)
-- Name: employees; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.employees (
    id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    mobile character varying(100) DEFAULT ''::character varying NOT NULL,
    role character varying(100) DEFAULT ''::character varying NOT NULL,
    branch_id integer DEFAULT 1 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.employees OWNER TO super_shop_dev_user;

--
-- TOC entry 217 (class 1259 OID 34114)
-- Name: employees_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.employees_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.employees_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3880 (class 0 OID 0)
-- Dependencies: 217
-- Name: employees_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.employees_id_seq OWNED BY public.employees.id;


--
-- TOC entry 236 (class 1259 OID 34448)
-- Name: hero_sections; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.hero_sections (
    id integer NOT NULL,
    main_banner character varying(255) DEFAULT ''::character varying,
    main_title character varying(100) DEFAULT ''::character varying,
    main_subtitle character varying(150) DEFAULT ''::character varying,
    side_top_banner character varying(255) DEFAULT ''::character varying,
    side_top_title character varying(100) DEFAULT ''::character varying,
    side_top_tag character varying(50) DEFAULT ''::character varying,
    mini_banner_1 character varying(255) DEFAULT ''::character varying,
    mini_banner_1_title character varying(100) DEFAULT ''::character varying,
    mini_banner_2 character varying(255) DEFAULT ''::character varying,
    mini_banner_2_title character varying(100) DEFAULT ''::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.hero_sections OWNER TO super_shop_dev_user;

--
-- TOC entry 235 (class 1259 OID 34447)
-- Name: hero_sections_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.hero_sections_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.hero_sections_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3881 (class 0 OID 0)
-- Dependencies: 235
-- Name: hero_sections_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.hero_sections_id_seq OWNED BY public.hero_sections.id;


--
-- TOC entry 228 (class 1259 OID 34267)
-- Name: legacy_products; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.legacy_products (
    id bigint NOT NULL,
    added_by character varying(191) DEFAULT NULL::character varying,
    user_id bigint,
    name character varying(80) DEFAULT NULL::character varying,
    slug character varying(120) DEFAULT NULL::character varying,
    product_type character varying(20) DEFAULT 'physical'::character varying NOT NULL,
    category_ids character varying(80) DEFAULT NULL::character varying,
    category_id character varying(191) DEFAULT NULL::character varying,
    sub_category_id character varying(191) DEFAULT NULL::character varying,
    sub_sub_category_id character varying(191) DEFAULT NULL::character varying,
    brand_id bigint,
    unit character varying(191) DEFAULT NULL::character varying,
    min_qty integer DEFAULT 1 NOT NULL,
    refundable boolean DEFAULT true NOT NULL,
    digital_product_type character varying(30) DEFAULT NULL::character varying,
    digital_file_ready character varying(191) DEFAULT NULL::character varying,
    digital_file_ready_storage_type character varying(10) DEFAULT 'public'::character varying,
    images text,
    color_image text DEFAULT ''::text NOT NULL,
    thumbnail character varying(255) DEFAULT NULL::character varying,
    thumbnail_storage_type character varying(10) DEFAULT 'public'::character varying,
    preview_file character varying(255) DEFAULT NULL::character varying,
    preview_file_storage_type character varying(255) DEFAULT 'public'::character varying,
    featured character varying(255) DEFAULT NULL::character varying,
    flash_deal character varying(255) DEFAULT NULL::character varying,
    video_provider character varying(30) DEFAULT NULL::character varying,
    video_url character varying(150) DEFAULT NULL::character varying,
    colors character varying(150) DEFAULT NULL::character varying,
    variant_product boolean DEFAULT false NOT NULL,
    attributes character varying(255) DEFAULT NULL::character varying,
    choice_options text,
    variation text,
    digital_product_file_types text,
    digital_product_extensions text,
    published boolean DEFAULT false NOT NULL,
    unit_price double precision DEFAULT 0 NOT NULL,
    purchase_price double precision DEFAULT 0 NOT NULL,
    tax character varying(191) DEFAULT '0.00'::character varying NOT NULL,
    tax_type character varying(80) DEFAULT NULL::character varying,
    tax_model character varying(20) DEFAULT 'exclude'::character varying NOT NULL,
    discount character varying(191) DEFAULT '0.00'::character varying NOT NULL,
    discount_type character varying(80) DEFAULT NULL::character varying,
    current_stock integer,
    minimum_order_qty integer DEFAULT 1 NOT NULL,
    details text,
    free_shipping boolean DEFAULT false NOT NULL,
    attachment character varying(191) DEFAULT NULL::character varying,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    status boolean DEFAULT true NOT NULL,
    featured_status boolean DEFAULT true NOT NULL,
    meta_title character varying(191) DEFAULT NULL::character varying,
    meta_description character varying(191) DEFAULT NULL::character varying,
    meta_image character varying(191) DEFAULT NULL::character varying,
    request_status boolean DEFAULT false NOT NULL,
    denied_note text,
    shipping_cost numeric(8,2) DEFAULT NULL::numeric,
    multiply_qty boolean,
    temp_shipping_cost numeric(8,2) DEFAULT NULL::numeric,
    is_shipping_cost_updated boolean,
    code character varying(191) DEFAULT NULL::character varying
);


ALTER TABLE public.legacy_products OWNER TO super_shop_dev_user;

--
-- TOC entry 234 (class 1259 OID 34408)
-- Name: order_items; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.order_items (
    id bigint NOT NULL,
    order_id bigint NOT NULL,
    product_id bigint,
    product_name character varying(255) NOT NULL,
    quantity integer DEFAULT 1 NOT NULL,
    unit_price numeric(10,2) DEFAULT 0 NOT NULL,
    total_price numeric(12,2) DEFAULT 0 NOT NULL,
    discount numeric(10,2) DEFAULT 0 NOT NULL,
    tax numeric(10,2) DEFAULT 0 NOT NULL,
    variation_info jsonb,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.order_items OWNER TO super_shop_dev_user;

--
-- TOC entry 3882 (class 0 OID 0)
-- Dependencies: 234
-- Name: TABLE order_items; Type: COMMENT; Schema: public; Owner: super_shop_dev_user
--

COMMENT ON TABLE public.order_items IS 'Stores individual items within an order';


--
-- TOC entry 233 (class 1259 OID 34407)
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.order_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_items_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3883 (class 0 OID 0)
-- Dependencies: 233
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- TOC entry 232 (class 1259 OID 34382)
-- Name: orders; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    order_number character varying(50) NOT NULL,
    customer_id bigint,
    customer_name character varying(100) NOT NULL,
    customer_mobile character varying(20) NOT NULL,
    customer_email character varying(100),
    customer_area character varying(255),
    customer_city character varying(100),
    payment_method character varying(50) DEFAULT 'COD'::character varying NOT NULL,
    payment_status public.payment_status DEFAULT 'pending'::public.payment_status NOT NULL,
    order_status public.order_status DEFAULT 'pending'::public.order_status NOT NULL,
    subtotal numeric(12,2) DEFAULT 0 NOT NULL,
    shipping_cost numeric(10,2) DEFAULT 0 NOT NULL,
    discount numeric(10,2) DEFAULT 0 NOT NULL,
    tax numeric(10,2) DEFAULT 0 NOT NULL,
    total numeric(12,2) DEFAULT 0 NOT NULL,
    order_note text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    delivered_at timestamp with time zone,
    cancelled_at timestamp with time zone,
    cancelled_reason text
);


ALTER TABLE public.orders OWNER TO super_shop_dev_user;

--
-- TOC entry 3884 (class 0 OID 0)
-- Dependencies: 232
-- Name: TABLE orders; Type: COMMENT; Schema: public; Owner: super_shop_dev_user
--

COMMENT ON TABLE public.orders IS 'Stores customer orders';


--
-- TOC entry 3885 (class 0 OID 0)
-- Dependencies: 232
-- Name: COLUMN orders.order_number; Type: COMMENT; Schema: public; Owner: super_shop_dev_user
--

COMMENT ON COLUMN public.orders.order_number IS 'Unique order identifier for display (e.g., ORD-20260205-0001)';


--
-- TOC entry 3886 (class 0 OID 0)
-- Dependencies: 232
-- Name: COLUMN orders.payment_method; Type: COMMENT; Schema: public; Owner: super_shop_dev_user
--

COMMENT ON COLUMN public.orders.payment_method IS 'Payment method: COD, ONLINE, CARD, etc.';


--
-- TOC entry 231 (class 1259 OID 34381)
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3887 (class 0 OID 0)
-- Dependencies: 231
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- TOC entry 248 (class 1259 OID 34936)
-- Name: product_reviews; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.product_reviews (
    id bigint NOT NULL,
    product_id bigint NOT NULL,
    customer_id bigint NOT NULL,
    rating smallint NOT NULL,
    title character varying(255),
    comment text,
    review_images text[],
    is_verified_purchase boolean DEFAULT false,
    status smallint DEFAULT 1,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT product_reviews_rating_check CHECK (((rating >= 1) AND (rating <= 5)))
);


ALTER TABLE public.product_reviews OWNER TO super_shop_dev_user;

--
-- TOC entry 247 (class 1259 OID 34935)
-- Name: product_reviews_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.product_reviews_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_reviews_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3888 (class 0 OID 0)
-- Dependencies: 247
-- Name: product_reviews_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.product_reviews_id_seq OWNED BY public.product_reviews.id;


--
-- TOC entry 246 (class 1259 OID 34913)
-- Name: product_variations; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.product_variations (
    id bigint NOT NULL,
    product_id bigint,
    variation_attributes jsonb DEFAULT '{}'::jsonb NOT NULL,
    sku text NOT NULL,
    price numeric(12,2) NOT NULL,
    stock_qty numeric(12,2) DEFAULT 0 NOT NULL,
    thumbnail text DEFAULT ''::text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.product_variations OWNER TO super_shop_dev_user;

--
-- TOC entry 245 (class 1259 OID 34912)
-- Name: product_variations_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.product_variations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.product_variations_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3889 (class 0 OID 0)
-- Dependencies: 245
-- Name: product_variations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.product_variations_id_seq OWNED BY public.product_variations.id;


--
-- TOC entry 244 (class 1259 OID 34848)
-- Name: products; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    name text NOT NULL,
    description text DEFAULT ''::text,
    category_id bigint,
    sub_category_id bigint,
    sub_sub_category_id bigint,
    brand_id bigint,
    sku text,
    status smallint DEFAULT 1 NOT NULL,
    unit_id integer,
    tags text DEFAULT ''::text,
    thumbnail text DEFAULT ''::text,
    gallery_images text[],
    retail_price numeric(12,2) DEFAULT 0 NOT NULL,
    purchase_price numeric(12,2) DEFAULT 0,
    min_retail_order_qty numeric(12,2) DEFAULT 1 NOT NULL,
    current_stock_qty numeric(12,2) DEFAULT 0 NOT NULL,
    stock_alert_qty numeric(12,2) DEFAULT 0 NOT NULL,
    total_sold numeric(12,2) DEFAULT 0 NOT NULL,
    discount_type character varying(20) DEFAULT 'percentage'::character varying,
    discount_amount numeric(12,2) DEFAULT 0,
    tax_amount numeric(12,2) DEFAULT 0,
    tax_type character varying(20) DEFAULT 'exclusive'::character varying,
    shipping_cost numeric(12,2) DEFAULT 0,
    shipping_type character varying(20) DEFAULT 'static'::character varying,
    has_variation boolean DEFAULT false NOT NULL,
    variation_attributes jsonb DEFAULT '{}'::jsonb NOT NULL,
    total_reviews bigint DEFAULT 0,
    avg_rating numeric(3,2) DEFAULT 0,
    five_star_count bigint DEFAULT 0 NOT NULL,
    four_star_count bigint DEFAULT 0 NOT NULL,
    three_star_count bigint DEFAULT 0 NOT NULL,
    two_star_count bigint DEFAULT 0 NOT NULL,
    one_star_count bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    wholesale_price numeric(12,2) DEFAULT 0 NOT NULL,
    min_wholesale_order_qty numeric(12,2) DEFAULT 1 NOT NULL
);


ALTER TABLE public.products OWNER TO super_shop_dev_user;

--
-- TOC entry 227 (class 1259 OID 34266)
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3890 (class 0 OID 0)
-- Dependencies: 227
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.legacy_products.id;


--
-- TOC entry 243 (class 1259 OID 34847)
-- Name: products_id_seq1; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.products_id_seq1
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq1 OWNER TO super_shop_dev_user;

--
-- TOC entry 3891 (class 0 OID 0)
-- Dependencies: 243
-- Name: products_id_seq1; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.products_id_seq1 OWNED BY public.products.id;


--
-- TOC entry 220 (class 1259 OID 34137)
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.refresh_tokens (
    id bigint NOT NULL,
    employee_id bigint NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    revoked_at timestamp with time zone
);


ALTER TABLE public.refresh_tokens OWNER TO super_shop_dev_user;

--
-- TOC entry 219 (class 1259 OID 34136)
-- Name: refresh_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.refresh_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.refresh_tokens_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3892 (class 0 OID 0)
-- Dependencies: 219
-- Name: refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;


--
-- TOC entry 224 (class 1259 OID 34175)
-- Name: sub_categories; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.sub_categories (
    id bigint NOT NULL,
    category_id bigint,
    name character varying(255) NOT NULL,
    priority smallint DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.sub_categories OWNER TO super_shop_dev_user;

--
-- TOC entry 223 (class 1259 OID 34174)
-- Name: sub_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.sub_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sub_categories_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3893 (class 0 OID 0)
-- Dependencies: 223
-- Name: sub_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.sub_categories_id_seq OWNED BY public.sub_categories.id;


--
-- TOC entry 226 (class 1259 OID 34194)
-- Name: sub_sub_categories; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.sub_sub_categories (
    id bigint NOT NULL,
    sub_category_id bigint,
    name character varying(255) NOT NULL,
    priority smallint DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.sub_sub_categories OWNER TO super_shop_dev_user;

--
-- TOC entry 225 (class 1259 OID 34193)
-- Name: sub_sub_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.sub_sub_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sub_sub_categories_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3894 (class 0 OID 0)
-- Dependencies: 225
-- Name: sub_sub_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.sub_sub_categories_id_seq OWNED BY public.sub_sub_categories.id;


--
-- TOC entry 240 (class 1259 OID 34628)
-- Name: units; Type: TABLE; Schema: public; Owner: super_shop_dev_user
--

CREATE TABLE public.units (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    symbol character varying(20) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.units OWNER TO super_shop_dev_user;

--
-- TOC entry 239 (class 1259 OID 34627)
-- Name: units_id_seq; Type: SEQUENCE; Schema: public; Owner: super_shop_dev_user
--

CREATE SEQUENCE public.units_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.units_id_seq OWNER TO super_shop_dev_user;

--
-- TOC entry 3895 (class 0 OID 0)
-- Dependencies: 239
-- Name: units_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: super_shop_dev_user
--

ALTER SEQUENCE public.units_id_seq OWNED BY public.units.id;


--
-- TOC entry 3515 (class 2604 OID 34840)
-- Name: attributes id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.attributes ALTER COLUMN id SET DEFAULT nextval('public.attributes_id_seq'::regclass);


--
-- TOC entry 3505 (class 2604 OID 34473)
-- Name: branches id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.branches ALTER COLUMN id SET DEFAULT nextval('public.branches_id_seq'::regclass);


--
-- TOC entry 3559 (class 2604 OID 34974)
-- Name: brands id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.brands ALTER COLUMN id SET DEFAULT nextval('public.brands_id_seq'::regclass);


--
-- TOC entry 3398 (class 2604 OID 34163)
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- TOC entry 3462 (class 2604 OID 34334)
-- Name: customers id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.customers ALTER COLUMN id SET DEFAULT nextval('public.customers_id_seq'::regclass);


--
-- TOC entry 3385 (class 2604 OID 34118)
-- Name: employees id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.employees ALTER COLUMN id SET DEFAULT nextval('public.employees_id_seq'::regclass);


--
-- TOC entry 3492 (class 2604 OID 34451)
-- Name: hero_sections id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.hero_sections ALTER COLUMN id SET DEFAULT nextval('public.hero_sections_id_seq'::regclass);


--
-- TOC entry 3413 (class 2604 OID 34270)
-- Name: legacy_products id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.legacy_products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- TOC entry 3485 (class 2604 OID 34411)
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- TOC entry 3474 (class 2604 OID 34385)
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- TOC entry 3554 (class 2604 OID 34939)
-- Name: product_reviews id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_reviews ALTER COLUMN id SET DEFAULT nextval('public.product_reviews_id_seq'::regclass);


--
-- TOC entry 3548 (class 2604 OID 34916)
-- Name: product_variations id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_variations ALTER COLUMN id SET DEFAULT nextval('public.product_variations_id_seq'::regclass);


--
-- TOC entry 3518 (class 2604 OID 34851)
-- Name: products id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq1'::regclass);


--
-- TOC entry 3395 (class 2604 OID 34140)
-- Name: refresh_tokens id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);


--
-- TOC entry 3403 (class 2604 OID 34178)
-- Name: sub_categories id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.sub_categories ALTER COLUMN id SET DEFAULT nextval('public.sub_categories_id_seq'::regclass);


--
-- TOC entry 3408 (class 2604 OID 34197)
-- Name: sub_sub_categories id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.sub_sub_categories ALTER COLUMN id SET DEFAULT nextval('public.sub_sub_categories_id_seq'::regclass);


--
-- TOC entry 3512 (class 2604 OID 34631)
-- Name: units id; Type: DEFAULT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.units ALTER COLUMN id SET DEFAULT nextval('public.units_id_seq'::regclass);


--
-- TOC entry 3857 (class 0 OID 34837)
-- Dependencies: 242
-- Data for Name: attributes; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--



--
-- TOC entry 3853 (class 0 OID 34470)
-- Dependencies: 238
-- Data for Name: branches; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.branches VALUES (1, 'Head Quarter', 'Bangladesh', 'Nangalkot, Comilla', 'Nangalkot Bot-tola, Beside TNT, Srifolia Road', '+8801641-967080', '', 'noorsupermart.info@gmail.com', 23.1696608, 91.2011241, '2026-02-14 11:45:08.029834+00', '2026-03-12 03:34:54.061888+00');


--
-- TOC entry 3865 (class 0 OID 34971)
-- Dependencies: 250
-- Data for Name: brands; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.brands VALUES (1, 'Noor Super Mart', 'public/images/brands/noor-super-mart.jpg', 1, true, '2026-03-06 15:01:31.919461+00', '2026-03-06 15:01:31.919461+00');


--
-- TOC entry 3837 (class 0 OID 34160)
-- Dependencies: 222
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.categories VALUES (4, 'Fruit Item', 'public/images/categories/fruit-item.jpeg', 2, true, '2026-02-08 09:05:35.159546+00', '2026-02-08 09:05:35.159546+00');
INSERT INTO public.categories VALUES (5, 'মাছ,মাংস এবং সবজি', 'public/images/categories/--.jpeg', 3, true, '2026-02-09 10:23:09.488759+00', '2026-02-09 10:23:36.389717+00');
INSERT INTO public.categories VALUES (1, 'Ramadaan Special', 'public/images/categories/ramadaan-special.jpg', 5, true, '2026-02-05 05:38:02.767297+00', '2026-03-06 14:45:03.368728+00');
INSERT INTO public.categories VALUES (3, 'Electronics', 'public/images/categories/electronics.jpg', 4, true, '2026-02-05 08:13:26.098487+00', '2026-03-06 14:47:50.030899+00');
INSERT INTO public.categories VALUES (2, 'Grocery', 'public/images/categories/grocery.jpg', 1, true, '2026-02-05 06:16:55.147556+00', '2026-03-06 15:00:19.769742+00');


--
-- TOC entry 3845 (class 0 OID 34331)
-- Dependencies: 230
-- Data for Name: customers; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.customers VALUES (1, 'Manik Mia', 'Manik Mia', NULL, '01684665955', '', 'imranmanik9@gmail.com', NULL, '$2a$10$sma1HAgwaHsEZs2G18NsWupltScxa.7NbVgvqRVicbAM9NmUQ3PHe', NULL, '2026-02-05 09:13:54.117862+00', '2026-02-06 21:36:03.103859+00', NULL, NULL, 'Nangolkot', NULL, NULL, NULL, NULL, true, NULL, NULL, NULL, NULL, NULL, false, NULL, false, NULL, 0.0000, 0, false, NULL, NULL, NULL, '', false);
INSERT INTO public.customers VALUES (3, 'Riyad', 'Riyad', 'Hosen', '01741594938', 'def.png', 'mdriyadhosen3570@gmail.com', NULL, '$2a$10$QEMAYSjdrGRSYjKUenJrweoDxh8/Jpo3PIirM8eMvVNVQf16nYWbS', NULL, '2026-02-08 11:51:18.922733+00', '2026-02-08 11:51:18.922733+00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, NULL, NULL, NULL, NULL, NULL, false, NULL, false, NULL, 0.0000, 0, false, NULL, NULL, NULL, 'en', false);
INSERT INTO public.customers VALUES (9, 'মনি', 'মনি', NULL, '01765053886', '', 'moninoyon285@gmail.com', NULL, '$2a$10$nvNvzDbW4.hpzgplJJ0/6ObmWqHL1WFO/GbxpZX9fgJ1FOdqJDqTy', NULL, '2026-03-02 07:15:00.219982+00', '2026-03-02 07:15:00.219982+00', NULL, NULL, 'কুমিল্লা', NULL, NULL, NULL, NULL, true, NULL, NULL, NULL, NULL, NULL, false, NULL, false, NULL, 0.0000, 0, false, NULL, NULL, NULL, '', false);


--
-- TOC entry 3833 (class 0 OID 34115)
-- Dependencies: 218
-- Data for Name: employees; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.employees VALUES (1, 'b5046bdb-7b0a-41da-9c5b-9b300b09991a', 'noorsupermart@gmail.com', '$2a$10$Dcwf7EbwRiUDfKzuc1i8Lu7POM0BkPEhJiApCD7ZQLLhYHlt8sM3W', '', '', '', 1, true, true, '2026-02-01 23:43:38.458+00', '2026-02-01 23:43:38.458+00');


--
-- TOC entry 3851 (class 0 OID 34448)
-- Dependencies: 236
-- Data for Name: hero_sections; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.hero_sections VALUES (1, 'public/images/hero/main.webp', 'Online Super Market', 'Quality Products for Everyday Life', 'public/images/hero/side-top.webp', 'Best Deals', 'Popular Picks Across All Categories', 'public/images/hero/mini-1.webp', 'Fresh Fruits', 'public/images/hero/mini-2.webp', 'Premium Fishes', '2026-02-14 14:46:15.251664', '2026-02-15 19:31:08.614643');


--
-- TOC entry 3843 (class 0 OID 34267)
-- Dependencies: 228
-- Data for Name: legacy_products; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.legacy_products VALUES (1, NULL, NULL, 'Ramadaan Special', 'ramadaan-special', 'physical', NULL, '1', NULL, NULL, NULL, '20 bags', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/ramadaan-special.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 2850, 2500, '0.00', NULL, 'exclude', '350', 'flat', 10, 5, NULL, true, NULL, '2026-02-05 05:46:26.405123+00', '2026-02-05 05:46:26.405123+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112233');
INSERT INTO public.legacy_products VALUES (3, NULL, NULL, 'রসুন-Garlic', '-garlic', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-garlic.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 180, 165, '0.00', NULL, 'exclude', '15', 'flat', 50, 10, 'চায়না রসুন', true, NULL, '2026-02-05 06:20:16.760816+00', '2026-02-05 06:20:16.760816+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '১১২২৩৫');
INSERT INTO public.legacy_products VALUES (4, NULL, NULL, '2LTR. Soyabean Oil', '2ltr-soyabean-oil', 'physical', NULL, '2', NULL, NULL, NULL, '2Ltr', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/2ltr.-soyabean-oil.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 400, 400, '0.00', NULL, 'exclude', '0.00', NULL, 50, 10, NULL, true, NULL, '2026-02-05 06:24:38.097684+00', '2026-02-05 10:30:50.153964+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112236');
INSERT INTO public.legacy_products VALUES (30, NULL, NULL, 'জিরা/Cumin', 'cumin', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/cumin.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 980, 950, '0.00', NULL, 'exclude', '30', 'flat', 100, 10, NULL, true, NULL, '2026-02-07 04:42:47.94981+00', '2026-02-07 04:44:24.192875+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112264');
INSERT INTO public.legacy_products VALUES (5, NULL, NULL, 'salt', 'salt', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/salt.webp', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, NULL, NULL, NULL, NULL, NULL, true, 42, 40.74, '0.00', NULL, 'exclude', '3', 'percent', 100, 20, NULL, true, NULL, '2026-02-05 06:36:44.480948+00', '2026-02-05 06:36:44.480948+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112237');
INSERT INTO public.legacy_products VALUES (25, NULL, NULL, 'chikpeas/ছোলা', 'chikpeas', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/chikpeas.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 110, 105, '0.00', NULL, 'exclude', '05', 'flat', 50, 10, NULL, true, NULL, '2026-02-07 03:51:11.506687+00', '2026-02-07 03:51:11.506687+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112257');
INSERT INTO public.legacy_products VALUES (6, NULL, NULL, 'চিনিগুড়া চাল', '-', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, NULL, NULL, NULL, NULL, NULL, true, 180, 150, '0.00', NULL, 'exclude', '30', 'flat', 50, 10, NULL, true, NULL, '2026-02-05 06:43:31.170873+00', '2026-02-05 06:43:31.170873+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112238');
INSERT INTO public.legacy_products VALUES (2, NULL, NULL, 'মুড়ি', '', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/.webp', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 80, 70, '0.00', NULL, 'exclude', '10', 'flat', 100, 100, 'বাজারের ভালো মানের এবং বাছাইকৃত', true, NULL, '2026-02-05 06:02:41.317414+00', '2026-02-05 06:45:53.658593+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112234');
INSERT INTO public.legacy_products VALUES (7, NULL, NULL, 'Fundools ইনস্ট্যান্ট নুডলস', 'fundools-', 'physical', NULL, '2', NULL, NULL, NULL, '496gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/fundools--.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 170, 155, '0.00', NULL, 'exclude', '15', 'flat', 100, 10, NULL, true, NULL, '2026-02-05 06:55:41.634358+00', '2026-02-05 06:57:30.19847+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112239');
INSERT INTO public.legacy_products VALUES (9, NULL, NULL, 'ছোলা/chikpeas', 'chikpeas', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/chikpeas.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, NULL, NULL, NULL, NULL, NULL, true, 120, 110, '0.00', NULL, 'exclude', '10', 'flat', 100, 50, NULL, true, NULL, '2026-02-05 07:01:37.823959+00', '2026-02-05 07:01:37.823959+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112240');
INSERT INTO public.legacy_products VALUES (10, NULL, NULL, 'বেসন/gram flour', 'gram-flour', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/gram-flour.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 120, 110, '0.00', NULL, 'exclude', '10', NULL, 100, 10, NULL, true, NULL, '2026-02-05 07:09:00.130338+00', '2026-02-05 07:09:00.130338+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112241');
INSERT INTO public.legacy_products VALUES (11, NULL, NULL, 'পেঁয়াজ/Onion', 'onion', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/onion.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 60, 55, '0.00', NULL, 'exclude', '5', 'flat', 200, 20, 'দেশী ভালো মানের ', true, NULL, '2026-02-05 07:11:56.204528+00', '2026-02-05 07:11:56.204528+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112242');
INSERT INTO public.legacy_products VALUES (12, NULL, NULL, 'আদা/Ginger', 'ginger', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/ginger.webp', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 200, 165, '0.00', NULL, 'exclude', '35', 'flat', 50, 10, 'চায়না অরিজিনাল', true, NULL, '2026-02-05 07:15:46.060248+00', '2026-02-05 07:15:46.060248+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112243');
INSERT INTO public.legacy_products VALUES (13, NULL, NULL, 'আলু/Potato', 'potato', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/potato.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 25, 25, '0.00', NULL, 'exclude', '0.00', NULL, 1000, 50, 'নতুন দেশি আলু', true, NULL, '2026-02-05 07:20:08.384836+00', '2026-02-05 07:20:08.384836+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112244');
INSERT INTO public.legacy_products VALUES (14, NULL, NULL, 'বাহরাইন অরেঞ্জ/আমের ট্যাং', '-', 'physical', NULL, '2', NULL, NULL, NULL, '2kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/--.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 1950, 1800, '0.00', NULL, 'exclude', '150', 'flat', 20, 10, NULL, true, NULL, '2026-02-05 07:24:25.146056+00', '2026-02-05 07:24:25.146056+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112245');
INSERT INTO public.legacy_products VALUES (15, NULL, NULL, 'হলুদ গুরা/turmeric powder', '-turmeric-powder', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-turmeric-powder.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 550, 520, '0.00', NULL, 'exclude', '30', 'flat', 25, 10, NULL, true, NULL, '2026-02-05 07:28:08.418878+00', '2026-02-05 07:28:08.418878+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112246');
INSERT INTO public.legacy_products VALUES (16, NULL, NULL, 'মরিচ গুরা/Chilli powder', '-chilli-powder', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-chilli-powder.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 590, 550, '0.00', NULL, 'exclude', '40', 'flat', 20, 10, NULL, true, NULL, '2026-02-05 07:33:08.703543+00', '2026-02-05 07:33:08.703543+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112247');
INSERT INTO public.legacy_products VALUES (17, NULL, NULL, 'মরিয়ম খেজুর/Dates', '-dates', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-dates.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 1800, 1600, '0.00', NULL, 'exclude', '200', 'flat', 10, 10, NULL, true, NULL, '2026-02-05 07:37:03.565212+00', '2026-02-05 07:37:03.565212+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112248');
INSERT INTO public.legacy_products VALUES (19, NULL, NULL, 'চিড়া/Chira', 'chira', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/chira.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 100, 90, '0.00', NULL, 'exclude', '10', 'flat', 100, 10, NULL, true, NULL, '2026-02-05 07:42:55.43568+00', '2026-02-05 07:42:55.43568+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112250');
INSERT INTO public.legacy_products VALUES (21, NULL, NULL, 'সরিষার তেল/mustard oil', '-mustard-oil', 'physical', NULL, '2', NULL, NULL, NULL, '1Ltr', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-mustard-oil.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 350, 315, '0.00', NULL, 'exclude', '10', 'percent', 120, 10, NULL, true, NULL, '2026-02-05 07:51:40.466434+00', '2026-02-05 07:51:40.466434+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112252');
INSERT INTO public.legacy_products VALUES (22, NULL, NULL, 'বনফুল লাচ্চা সেমাই', '-', 'physical', NULL, '2', NULL, NULL, NULL, '200gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/--.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 50, 50, '0.00', NULL, 'exclude', '0.00', NULL, 50, 10, NULL, true, NULL, '2026-02-05 07:54:07.837328+00', '2026-02-05 07:54:45.274023+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112253');
INSERT INTO public.legacy_products VALUES (23, NULL, NULL, 'ডানো ফুল ক্রিম মিল্ক পাওডার', '-', 'physical', NULL, '2', NULL, NULL, NULL, '500gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/----.webp', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 480, 460, '0.00', NULL, 'exclude', '20', 'flat', 20, 10, NULL, true, NULL, '2026-02-05 08:04:46.392214+00', '2026-02-05 08:04:46.392214+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112254');
INSERT INTO public.legacy_products VALUES (26, NULL, NULL, 'Gas stove glass single', 'gas-stove-glass-single', 'physical', NULL, '3', NULL, NULL, NULL, '1pcs', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/gas-stove-glass-single.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 2250, 2150, '0.00', NULL, 'exclude', '100', 'flat', 20, 10, NULL, true, NULL, '2026-02-07 04:02:46.331448+00', '2026-02-07 04:02:46.331448+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112260');
INSERT INTO public.legacy_products VALUES (27, NULL, NULL, 'Tea Bag', 'tea-bag', 'physical', NULL, '2', NULL, NULL, NULL, '1box', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/tea-bag.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 90, 90, '0.00', NULL, 'exclude', '0.00', NULL, 50, 10, NULL, true, NULL, '2026-02-07 04:11:33.258691+00', '2026-02-07 04:11:33.258691+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112261');
INSERT INTO public.legacy_products VALUES (34, NULL, NULL, 'এলাচ/Cardamom', 'cardamom', 'physical', NULL, '2', NULL, NULL, NULL, '25gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/cardamom.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 150, 150, '0.00', NULL, 'exclude', '0.00', NULL, 100, 10, NULL, true, NULL, '2026-02-07 07:44:30.659243+00', '2026-02-07 07:44:30.659243+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112268');
INSERT INTO public.legacy_products VALUES (24, NULL, NULL, 'দাব্বাস খেজুর', '-', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-.png', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 350, 320, '0.00', NULL, 'exclude', '30', 'flat', 20, 10, NULL, true, NULL, '2026-02-05 10:27:13.583991+00', '2026-02-05 10:28:13.732859+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112252');
INSERT INTO public.legacy_products VALUES (20, NULL, NULL, 'খেসারি ডাল', '-', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 120, 110, '0.00', NULL, 'exclude', '10', 'flat', 100, 10, NULL, true, NULL, '2026-02-05 07:46:30.962635+00', '2026-02-05 10:30:22.044669+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112251');
INSERT INTO public.legacy_products VALUES (29, NULL, NULL, 'ধনিয়া /Coriander', '-coriander', 'physical', NULL, '2', NULL, NULL, NULL, '250gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-coriander.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 480, 450, '0.00', NULL, 'exclude', '30', 'flat', 100, 10, NULL, true, NULL, '2026-02-07 04:38:00.008498+00', '2026-02-07 04:38:00.008498+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112263');
INSERT INTO public.legacy_products VALUES (32, NULL, NULL, 'দারুচিনি/Cinnamon', 'cinnamon', 'physical', NULL, '2', NULL, NULL, NULL, '100', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/cinnamon.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 100, 80, '0.00', NULL, 'exclude', '20', 'flat', 100, 10, NULL, true, NULL, '2026-02-07 05:16:41.765053+00', '2026-02-07 05:16:41.765053+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112266');
INSERT INTO public.legacy_products VALUES (31, NULL, NULL, 'গরম মসলা/Garam Masala', '-garam-masala', 'physical', NULL, '2', NULL, NULL, NULL, '25gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-garam-masala.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 80, 50, '0.00', NULL, 'exclude', '30', NULL, 100, 10, NULL, true, NULL, '2026-02-07 05:11:22.808539+00', '2026-02-07 05:19:12.801763+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112265');
INSERT INTO public.legacy_products VALUES (28, NULL, NULL, 'গোল মরিচ', '-', 'physical', NULL, '2', NULL, NULL, NULL, '250gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 350, 300, '0.00', NULL, 'exclude', '50', 'flat', 100, 10, NULL, true, NULL, '2026-02-07 04:31:28.932876+00', '2026-02-07 05:31:07.157413+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112262');
INSERT INTO public.legacy_products VALUES (35, NULL, NULL, 'তেজপাতা/Bay Leaf', 'bay-leaf', 'physical', NULL, '2', NULL, NULL, NULL, '25gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/bay-leaf.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 30, 25, '0.00', NULL, 'exclude', '5', 'flat', 20, 10, NULL, true, NULL, '2026-02-07 08:34:25.509604+00', '2026-02-07 08:34:25.509604+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '1122690');
INSERT INTO public.legacy_products VALUES (36, NULL, NULL, 'কফি/coffee', 'coffee', 'physical', NULL, '1', NULL, NULL, NULL, '200gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/coffee.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 1000, 950, '0.00', NULL, 'exclude', '50', 'flat', 100, 10, NULL, true, NULL, '2026-02-07 09:00:55.924135+00', '2026-02-07 09:00:55.924135+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112270');
INSERT INTO public.legacy_products VALUES (37, NULL, NULL, 'Nutty Biscuit', 'nutty-biscuit', 'physical', NULL, '2', NULL, NULL, NULL, '175gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/nutty-biscuit.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 50, 50, '0.00', NULL, 'exclude', '0.00', NULL, 200, 10, NULL, true, NULL, '2026-02-07 09:54:38.26073+00', '2026-02-07 09:54:38.26073+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112272');
INSERT INTO public.legacy_products VALUES (38, NULL, NULL, 'বারবিকিউ চানাচুর', '-', 'physical', NULL, '2', NULL, NULL, NULL, '150gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-.webp', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 50, 0, '0.00', NULL, 'exclude', '0.00', NULL, 100, 10, NULL, true, NULL, '2026-02-07 09:58:50.987286+00', '2026-02-07 09:58:50.987286+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112273');
INSERT INTO public.legacy_products VALUES (53, NULL, NULL, 'fresh apple', 'fresh-apple', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/fresh-apple.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 350, 300, '0.00', NULL, 'exclude', '50', 'flat', 100, 10, NULL, true, NULL, '2026-02-09 09:05:39.704174+00', '2026-02-09 10:10:06.614203+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112296');
INSERT INTO public.legacy_products VALUES (39, NULL, NULL, 'চিপস', '', 'physical', NULL, '2', NULL, NULL, NULL, '10gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 10, 10, '0.00', NULL, 'exclude', '0.00', NULL, 500, 10, NULL, true, NULL, '2026-02-07 10:04:16.182258+00', '2026-02-07 10:08:50.470664+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112274');
INSERT INTO public.legacy_products VALUES (40, NULL, NULL, 'Bag', 'bag', 'physical', NULL, '2', NULL, NULL, NULL, '1 Pcs', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/bag.webp', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 2889, 2449, '0.00', NULL, 'exclude', '440', 'flat', 100, 10, NULL, true, NULL, '2026-02-07 10:31:16.910894+00', '2026-02-07 10:31:16.910894+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, NULL);
INSERT INTO public.legacy_products VALUES (41, NULL, NULL, 'নেহা মেহিদি', '-', 'physical', NULL, '2', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 30, 30, '0.00', NULL, 'exclude', '0.00', NULL, 100, 10, NULL, true, NULL, '2026-02-08 04:05:05.530129+00', '2026-02-08 04:05:05.530129+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112278');
INSERT INTO public.legacy_products VALUES (43, NULL, NULL, 'কাবেরি মেহেদি/Kaberi mehedi', '-kaberi-mehedi', 'physical', NULL, '2', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-kaberi-mehedi.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 70, 50, '0.00', NULL, 'exclude', '20', 'flat', 100, 10, NULL, true, NULL, '2026-02-08 04:15:26.74473+00', '2026-02-08 04:15:26.74473+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112280');
INSERT INTO public.legacy_products VALUES (56, NULL, NULL, 'আনারস/ Pineapple', '-pineapple', 'physical', NULL, '2', NULL, NULL, NULL, '1pes', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-pineapple.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 90, 70, '0.00', NULL, 'exclude', '20', 'flat', 100, 10, NULL, true, NULL, '2026-02-09 09:29:29.706798+00', '2026-02-10 04:23:28.138781+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112297');
INSERT INTO public.legacy_products VALUES (33, NULL, NULL, 'লবঙ্গ', '', 'physical', NULL, '2', NULL, NULL, NULL, '100', 1, true, NULL, NULL, NULL, NULL, '', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 110, 100, '0.00', NULL, 'exclude', '10', 'flat', 100, 10, NULL, true, NULL, '2026-02-07 05:28:51.606231+00', '2026-02-08 04:19:28.181358+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112267');
INSERT INTO public.legacy_products VALUES (42, NULL, NULL, 'স্মার্ট কোন মেহেদি', '-', 'physical', NULL, '2', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 70, 50, '0.00', NULL, 'exclude', '20', 'flat', 100, 10, NULL, true, NULL, '2026-02-08 04:11:41.307109+00', '2026-02-08 04:19:39.732709+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112279');
INSERT INTO public.legacy_products VALUES (44, NULL, NULL, 'ডিম/Egg', 'egg', 'physical', NULL, '2', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/egg.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 10, 10, '0.00', NULL, 'exclude', '0.00', NULL, 1000, 100, NULL, true, NULL, '2026-02-08 04:23:35.824201+00', '2026-02-08 04:23:35.824201+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112281');
INSERT INTO public.legacy_products VALUES (45, NULL, NULL, 'চিনি/suger', 'suger', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/suger.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 110, 105, '0.00', NULL, 'exclude', '5', 'flat', 100, 10, NULL, true, NULL, '2026-02-08 05:24:09.541363+00', '2026-02-08 05:24:09.541363+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112282');
INSERT INTO public.legacy_products VALUES (46, NULL, NULL, 'IGT Gas sensor Regulator', 'igt-gas-sensor-regulator', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/igt-gas-sensor-regulator.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 1350, 1150, '0.00', NULL, 'exclude', '200', 'flat', 50, 10, NULL, true, NULL, '2026-02-08 05:25:55.493441+00', '2026-02-08 05:25:55.493441+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112283');
INSERT INTO public.legacy_products VALUES (47, NULL, NULL, 'Fruit Combo', 'fruit-combo', 'physical', NULL, '4', NULL, NULL, NULL, '1 Dala', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/fruit-combo.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 1100, 1000, '0.00', NULL, 'exclude', '100', 'flat', 20, 5, NULL, true, NULL, '2026-02-08 09:07:22.44218+00', '2026-02-08 09:07:22.44218+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112284');
INSERT INTO public.legacy_products VALUES (48, NULL, NULL, 'আপেল/Apple', 'apple', 'physical', NULL, '4', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/apple.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 300, 280, '0.00', NULL, 'exclude', '20', 'flat', 20, 10, NULL, true, NULL, '2026-02-08 09:10:51.429173+00', '2026-02-08 09:10:51.429173+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112285');
INSERT INTO public.legacy_products VALUES (49, NULL, NULL, 'মালতা/Malta', 'malta', 'physical', NULL, '4', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/malta.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 280, 270, '0.00', NULL, 'exclude', '0.00', 'flat', 100, 10, NULL, true, NULL, '2026-02-08 09:13:32.632812+00', '2026-02-08 09:13:32.632812+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112286');
INSERT INTO public.legacy_products VALUES (50, NULL, NULL, 'আনার', '', 'physical', NULL, '4', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 450, 430, '0.00', NULL, 'exclude', '0.00', 'flat', 200, 100, NULL, true, NULL, '2026-02-08 09:15:58.931718+00', '2026-02-19 03:48:16.300375+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112287');
INSERT INTO public.legacy_products VALUES (51, NULL, NULL, 'ড্রাগন ফল', '-', 'physical', NULL, '4', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 450, 420, '0.00', NULL, 'exclude', '30', 'flat', 100, 10, NULL, true, NULL, '2026-02-08 09:21:21.76798+00', '2026-02-09 06:35:01.706698+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112288');
INSERT INTO public.legacy_products VALUES (52, NULL, NULL, 'পেয়ারা/Guava', 'guava', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/guava.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 40, 30, '0.00', NULL, 'exclude', '10', 'flat', 100, 10, NULL, true, NULL, '2026-02-09 09:00:58.058993+00', '2026-02-09 09:00:58.058993+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112295');
INSERT INTO public.legacy_products VALUES (65, NULL, NULL, 'মসুর ডাল', '-', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 180, 180, '0.00', NULL, 'exclude', '0.00', NULL, 50, 10, NULL, true, NULL, '2026-02-10 05:29:19.651525+00', '2026-02-19 04:37:19.757963+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112297');
INSERT INTO public.legacy_products VALUES (83, NULL, NULL, 'Lux Soap', 'lux-soap', 'physical', NULL, '2', NULL, NULL, NULL, '75gm,', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/lux-soap.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 50, 50, '0.00', NULL, 'exclude', '0.00', NULL, 100, 10, NULL, true, NULL, '2026-02-21 08:40:55.833407+00', '2026-02-21 08:40:55.833407+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '1122366');
INSERT INTO public.legacy_products VALUES (62, NULL, NULL, 'তেলাপিয়া মাছ', '-', 'physical', NULL, '5', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 200, 190, '0.00', NULL, 'exclude', '10', 'flat', 20, 10, NULL, true, NULL, '2026-02-09 10:24:47.504191+00', '2026-02-09 10:26:04.722625+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112297');
INSERT INTO public.legacy_products VALUES (63, NULL, NULL, 'গরুর মাংস', '-', 'physical', NULL, '5', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 800, 750, '0.00', NULL, 'exclude', '50', 'flat', 100, 50, NULL, true, NULL, '2026-02-09 10:29:17.217523+00', '2026-02-09 10:29:44.794071+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112298');
INSERT INTO public.legacy_products VALUES (64, NULL, NULL, 'মুরগীর মাংস', '-', 'physical', NULL, '5', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 180, 170, '0.00', NULL, 'exclude', '0.00', 'flat', 100, 100, NULL, true, NULL, '2026-02-09 10:31:22.968172+00', '2026-02-09 10:31:35.76625+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112299');
INSERT INTO public.legacy_products VALUES (67, NULL, NULL, 'মাস কলাইর ডাল', '-', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/--.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 225, 225, '0.00', NULL, 'exclude', '0.00', NULL, 20, 10, NULL, true, NULL, '2026-02-10 05:48:17.772579+00', '2026-02-10 05:48:17.772579+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112301');
INSERT INTO public.legacy_products VALUES (69, NULL, NULL, 'জয়ত্রী', '', 'physical', NULL, '2', NULL, NULL, NULL, '50gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 210, 200, '0.00', NULL, 'exclude', '10', 'flat', 100, 10, NULL, true, NULL, '2026-02-10 06:29:21.702793+00', '2026-02-10 06:29:21.702793+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112303');
INSERT INTO public.legacy_products VALUES (70, NULL, NULL, 'মেথি/Methi', 'methi', 'physical', NULL, '2', NULL, NULL, NULL, '250', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/methi.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 150, 150, '0.00', NULL, 'exclude', '0.00', NULL, 100, 10, NULL, true, NULL, '2026-02-10 06:34:34.344978+00', '2026-02-10 06:34:34.344978+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112304');
INSERT INTO public.legacy_products VALUES (71, NULL, NULL, '22mm Regulator', '22mm-regulator', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/22mm-regulator.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 350, 300, '0.00', NULL, 'exclude', '50', 'flat', 1000, 100, NULL, true, NULL, '2026-02-14 23:29:54.622226+00', '2026-02-14 23:29:54.622226+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112305');
INSERT INTO public.legacy_products VALUES (66, NULL, NULL, 'বুটের ডাল', '-', 'physical', NULL, '2', NULL, NULL, NULL, '1kg', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/-.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 70, 70, '0.00', NULL, 'exclude', '0.00', NULL, 50, 10, NULL, true, NULL, '2026-02-10 05:36:38.48077+00', '2026-02-14 23:30:35.107813+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112298');
INSERT INTO public.legacy_products VALUES (72, NULL, NULL, 'LG SS Double Burner Gas Stove', 'lg-ss-double-burner-gas-stove', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/lg-ss-double-burner-gas-stove.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 2500, 2250, '0.00', NULL, 'exclude', '250', 'flat', 50, 10, NULL, true, NULL, '2026-02-19 04:09:34.906361+00', '2026-02-19 04:09:34.906361+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112235');
INSERT INTO public.legacy_products VALUES (73, NULL, NULL, 'SS Pan Support', 'ss-pan-support', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/ss-pan-support.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 150, 120, '0.00', NULL, 'exclude', '30', 'flat', 50, 10, NULL, true, NULL, '2026-02-19 04:12:27.615654+00', '2026-02-19 04:12:27.615654+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '1122336');
INSERT INTO public.legacy_products VALUES (74, NULL, NULL, 'Honeycomb Burner', 'honeycomb-burner', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/honeycomb-burner.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 250, 220, '0.00', NULL, 'exclude', '30', 'flat', 100, 100, NULL, true, NULL, '2026-02-19 04:13:54.119099+00', '2026-02-19 04:13:54.119099+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '1122356');
INSERT INTO public.legacy_products VALUES (75, NULL, NULL, 'Pinium Burner', 'pinium-burner', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/pinium-burner.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 280, 260, '0.00', NULL, 'exclude', '20', 'flat', 100, 10, NULL, true, NULL, '2026-02-19 04:15:17.452394+00', '2026-02-19 04:15:17.452394+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '11223569');
INSERT INTO public.legacy_products VALUES (76, NULL, NULL, '100 MM Pinium', '100-mm-pinium', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/100-mm-pinium.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 150, 140, '0.00', NULL, 'exclude', '10', 'flat', 100, 50, NULL, true, NULL, '2026-02-19 04:16:49.251511+00', '2026-02-19 04:16:49.251511+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '1122365');
INSERT INTO public.legacy_products VALUES (77, NULL, NULL, 'Gas Saver', 'gas-saver', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/gas-saver.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 220, 200, '0.00', NULL, 'exclude', '20', 'flat', 1000, 100, NULL, true, NULL, '2026-02-19 04:19:25.571029+00', '2026-02-19 04:19:25.571029+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '11223654');
INSERT INTO public.legacy_products VALUES (78, NULL, NULL, 'LG Glass Single Gas Stove', 'lg-glass-single-gas-stove', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/lg-glass-single-gas-stove.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 2300, 2070, '0.00', NULL, 'exclude', '10', 'percent', 50, 10, NULL, true, NULL, '2026-02-19 04:21:45.415254+00', '2026-02-19 04:21:45.415254+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '11223654');
INSERT INTO public.legacy_products VALUES (79, NULL, NULL, 'Pan Support', 'pan-support', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/pan-support.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 150, 150, '0.00', NULL, 'exclude', '0.00', NULL, 100, 100, NULL, true, NULL, '2026-02-19 04:24:35.085857+00', '2026-02-19 04:24:35.085857+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '225863');
INSERT INTO public.legacy_products VALUES (80, NULL, NULL, 'Serbo Meter Regulator', 'serbo-meter-regulator', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/serbo-meter-regulator.jpeg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 450, 380, '0.00', NULL, 'exclude', '70', 'flat', 100, 100, NULL, true, NULL, '2026-02-19 04:26:55.663955+00', '2026-02-19 04:26:55.663955+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '45000');
INSERT INTO public.legacy_products VALUES (81, NULL, NULL, 'LG Double Glass Gas Stove', 'lg-double-glass-gas-stove', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/lg-double-glass-gas-stove.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 4500, 4200, '0.00', NULL, 'exclude', '300', 'flat', 50, 10, NULL, true, NULL, '2026-02-19 04:30:16.174721+00', '2026-02-19 04:30:16.174721+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112236654');
INSERT INTO public.legacy_products VALUES (82, NULL, NULL, 'Digital Scale 40kg', 'digital-scale-40kg', 'physical', NULL, '3', NULL, NULL, NULL, '1psc', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/digital-scale-40kg.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 2700, 2550, '0.00', NULL, 'exclude', '150', 'flat', 50, 10, NULL, true, NULL, '2026-02-19 04:33:55.243546+00', '2026-02-19 04:33:55.243546+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '1122365');
INSERT INTO public.legacy_products VALUES (68, NULL, NULL, 'জাইফল/Nutmeg', 'nutmeg', 'physical', NULL, '2', NULL, NULL, NULL, '50gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/nutmeg.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 100, 100, '0.00', NULL, 'exclude', '0.00', NULL, 100, 10, NULL, true, NULL, '2026-02-10 06:20:18.492204+00', '2026-02-19 04:36:52.834537+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '112302');
INSERT INTO public.legacy_products VALUES (85, NULL, NULL, 'Sunsilk Lusciously Thick & Long Shampoo', 'sunsilk-lusciously-thick-long-shampoo', 'physical', NULL, '2', NULL, NULL, NULL, '360ml', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/sunsilk-lusciously-thick--long-shampoo.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 650, 550, '0.00', NULL, 'exclude', '100', 'flat', 100, 10, NULL, true, NULL, '2026-02-21 08:52:13.471505+00', '2026-02-21 08:52:13.471505+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '1122368');
INSERT INTO public.legacy_products VALUES (84, NULL, NULL, 'Lux Nourished Glow Soap', 'lux-nourished-glow-soap', 'physical', NULL, '2', NULL, NULL, NULL, '98gm', 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/lux-nourished-glow-soap.webp', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, true, 110, 110, '0.00', NULL, 'exclude', '0.00', NULL, 100, 10, NULL, true, NULL, '2026-02-21 08:47:11.398267+00', '2026-02-21 08:53:26.504274+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, '1122367');
INSERT INTO public.legacy_products VALUES (86, NULL, NULL, 'Bahrain tang 500gm pack', 'bahrain-tang-500gm-pack', 'physical', NULL, '2', NULL, NULL, NULL, NULL, 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/bahrain-tang-500gm-pack.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, false, 400, 350, '0.00', 'exclusive', 'exclude', '0.00', 'flat', NULL, 1, NULL, false, NULL, '2026-02-22 03:52:08.21865+00', '2026-02-22 04:23:57.810964+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, NULL);
INSERT INTO public.legacy_products VALUES (87, NULL, NULL, 'Indian Tang 500mg', 'indian-tang-500mg', 'physical', NULL, '2', NULL, NULL, NULL, NULL, 1, true, NULL, NULL, NULL, NULL, '', 'public/images/products/indian-tang-500mg.jpg', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, NULL, NULL, NULL, NULL, NULL, false, 350, 320, '0.00', 'exclusive', 'exclude', '0.00', 'flat', NULL, 1, NULL, false, NULL, '2026-02-22 16:17:57.427164+00', '2026-02-22 16:17:57.427164+00', true, true, NULL, NULL, NULL, false, NULL, 0.00, NULL, NULL, NULL, NULL);


--
-- TOC entry 3849 (class 0 OID 34408)
-- Dependencies: 234
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.order_items VALUES (1, 1, NULL, 'মুড়ি', 1, 70.00, 70.00, 0.00, 0.00, NULL, '2026-02-05 06:07:51.184161+00');
INSERT INTO public.order_items VALUES (2, 2, NULL, 'salt', 2, 40.74, 81.48, 0.00, 0.00, NULL, '2026-02-05 07:04:34.996676+00');
INSERT INTO public.order_items VALUES (3, 3, NULL, 'আলু/Potato', 1, 25.00, 25.00, 0.00, 0.00, NULL, '2026-02-05 09:13:54.120992+00');
INSERT INTO public.order_items VALUES (4, 4, NULL, 'আলু/Potato', 1, 25.00, 25.00, 0.00, 0.00, NULL, '2026-02-05 09:30:32.211812+00');
INSERT INTO public.order_items VALUES (5, 4, NULL, 'মুড়ি', 1, 70.00, 70.00, 0.00, 0.00, NULL, '2026-02-05 09:30:32.211812+00');
INSERT INTO public.order_items VALUES (6, 5, NULL, 'বারবিকিউ চানাচুর', 1, 50.00, 50.00, 0.00, 0.00, NULL, '2026-02-09 03:23:36.579247+00');
INSERT INTO public.order_items VALUES (7, 8, NULL, 'বারবিকিউ চানাচুর', 1, 50.00, 50.00, 0.00, 0.00, NULL, '2026-02-09 03:23:36.615407+00');
INSERT INTO public.order_items VALUES (8, 9, NULL, 'বারবিকিউ চানাচুর', 1, 50.00, 50.00, 0.00, 0.00, NULL, '2026-02-09 03:23:36.868684+00');
INSERT INTO public.order_items VALUES (9, 10, NULL, 'Olympic foodie instant noodles masala 12pack 744gm', 1, 210.00, 210.00, 0.00, 0.00, NULL, '2026-02-28 08:14:46.279816+00');
INSERT INTO public.order_items VALUES (10, 11, NULL, 'fair & lovely winter fairness cream 80gm', 1, 180.00, 180.00, 0.00, 0.00, NULL, '2026-03-01 05:21:21.538725+00');
INSERT INTO public.order_items VALUES (11, 12, NULL, 'নাগাল খেজুর Nagal Dates | 1kg, 500gm | সৌদি আরব', 1, 350.00, 350.00, 0.00, 0.00, NULL, '2026-03-02 07:15:00.226146+00');
INSERT INTO public.order_items VALUES (12, 12, NULL, 'Lux Nourished Glow Soap', 1, 110.00, 110.00, 0.00, 0.00, NULL, '2026-03-02 07:15:00.226146+00');
INSERT INTO public.order_items VALUES (13, 13, NULL, 'Kaveri No.1 Indian Cone Mehedi (12PCS PACK)', 3, 360.00, 1080.00, 0.00, 0.00, NULL, '2026-03-10 05:53:26.516278+00');
INSERT INTO public.order_items VALUES (14, 14, NULL, 'Kaveri No.1 Indian Cone Mehedi (12PCS PACK)', 1, 360.00, 360.00, 0.00, 0.00, NULL, '2026-03-10 06:11:46.776594+00');
INSERT INTO public.order_items VALUES (15, 15, NULL, 'আদা/Ginger', 1, 200.00, 200.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (16, 15, NULL, 'পেঁয়াজ/Onion', 2, 60.00, 120.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (17, 15, NULL, 'রসুন-Garlic', 1, 180.00, 180.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (18, 15, NULL, 'মুড়ি', 1, 80.00, 80.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (19, 15, NULL, 'চিনি/suger', 10, 110.00, 1100.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (20, 15, NULL, 'বেসন/gram flour', 1, 120.00, 120.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (21, 15, NULL, 'আলু/Potato', 5, 25.00, 125.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (22, 15, NULL, 'মুরগীর মাংস(Chicken)', 2, 250.00, 500.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (23, 15, NULL, 'মালতা/Malta', 1, 350.00, 350.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (24, 15, NULL, 'surf exel', 1, 120.00, 120.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (25, 15, NULL, 'radhuni murgir masala', 1, 95.00, 95.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (26, 15, NULL, 'atta', 2, 65.00, 130.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (27, 15, NULL, 'parasuit oil 200ml', 1, 200.00, 200.00, 0.00, 0.00, NULL, '2026-03-12 04:18:29.723314+00');
INSERT INTO public.order_items VALUES (28, 16, NULL, 'মুরগীর মাংস(Chicken)', 2, 250.00, 500.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (29, 16, NULL, 'রসুন-Garlic', 1, 180.00, 180.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (30, 16, NULL, 'বেসন/gram flour', 1, 120.00, 120.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (31, 16, NULL, 'আদা/Ginger', 1, 200.00, 200.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (32, 16, NULL, 'atta', 2, 65.00, 130.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (33, 16, NULL, 'parasuit oil 200ml', 1, 200.00, 200.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (34, 16, NULL, 'আলু/Potato', 5, 25.00, 125.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (35, 16, NULL, 'মালতা/Malta', 1, 350.00, 350.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (36, 16, NULL, 'চিনি/suger', 1, 110.00, 110.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (37, 16, NULL, 'পেঁয়াজ/Onion', 2, 60.00, 120.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (38, 16, NULL, 'radhuni murgir masala', 1, 95.00, 95.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (39, 16, NULL, 'surf exel', 1, 120.00, 120.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (40, 16, NULL, 'মুড়ি', 1, 80.00, 80.00, 0.00, 0.00, NULL, '2026-03-12 04:28:07.294342+00');
INSERT INTO public.order_items VALUES (41, 17, NULL, 'Omera 12kg Cylinder', 1, 1800.00, 1800.00, 0.00, 0.00, NULL, '2026-03-14 15:32:08.235135+00');
INSERT INTO public.order_items VALUES (42, 18, NULL, 'Sunsilk Lusciously Thick & Long Shampoo', 1, 550.00, 550.00, 0.00, 0.00, NULL, '2026-03-15 04:15:38.168496+00');
INSERT INTO public.order_items VALUES (43, 19, NULL, '20 g Cricket Plastic Ball', 6, 60.00, 360.00, 0.00, 0.00, NULL, '2026-03-15 05:45:16.982148+00');
INSERT INTO public.order_items VALUES (44, 20, NULL, 'WATER TRAY CHULA', 8, 45.00, 360.00, 0.00, 0.00, NULL, '2026-03-15 07:53:22.29571+00');
INSERT INTO public.order_items VALUES (45, 20, NULL, 'SINGLE AUTO IGNITION GD', 3, 140.00, 420.00, 0.00, 0.00, NULL, '2026-03-15 07:53:22.29571+00');
INSERT INTO public.order_items VALUES (46, 20, NULL, 'DOUBLE GD AUTO', 3, 170.00, 510.00, 0.00, 0.00, NULL, '2026-03-15 07:53:22.29571+00');
INSERT INTO public.order_items VALUES (47, 20, NULL, 'DOUBLE IGNITION SS CHULA', 3, 120.00, 360.00, 0.00, 0.00, NULL, '2026-03-15 07:53:22.29571+00');
INSERT INTO public.order_items VALUES (48, 20, NULL, 'SINGLE IGNOITION SS CHULA', 3, 105.00, 315.00, 0.00, 0.00, NULL, '2026-03-15 07:53:22.29571+00');
INSERT INTO public.order_items VALUES (49, 21, NULL, 'SS Pan Support', 1, 110.00, 110.00, 0.00, 0.00, NULL, '2026-03-15 07:57:12.770355+00');
INSERT INTO public.order_items VALUES (50, 22, NULL, 'NOGEL', 10, 20.00, 200.00, 0.00, 0.00, NULL, '2026-03-15 08:21:27.029136+00');
INSERT INTO public.order_items VALUES (51, 22, NULL, 'ELBO', 10, 30.00, 300.00, 0.00, 0.00, NULL, '2026-03-15 08:21:27.029136+00');
INSERT INTO public.order_items VALUES (52, 22, NULL, '22mm Regulator', 2, 210.00, 420.00, 0.00, 0.00, NULL, '2026-03-15 08:21:27.029136+00');
INSERT INTO public.order_items VALUES (53, 22, NULL, '22 MM MINI', 1, 200.00, 200.00, 0.00, 0.00, NULL, '2026-03-15 08:21:27.029136+00');
INSERT INTO public.order_items VALUES (54, 23, NULL, 'NOGEL', 10, 20.00, 200.00, 0.00, 0.00, NULL, '2026-03-15 08:21:27.11294+00');
INSERT INTO public.order_items VALUES (55, 23, NULL, 'ELBO', 10, 30.00, 300.00, 0.00, 0.00, NULL, '2026-03-15 08:21:27.11294+00');
INSERT INTO public.order_items VALUES (56, 23, NULL, '22mm Regulator', 2, 210.00, 420.00, 0.00, 0.00, NULL, '2026-03-15 08:21:27.11294+00');
INSERT INTO public.order_items VALUES (57, 23, NULL, '22 MM MINI', 1, 200.00, 200.00, 0.00, 0.00, NULL, '2026-03-15 08:21:27.11294+00');
INSERT INTO public.order_items VALUES (58, 24, NULL, 'SINGLE IGNOITION SS CHULA', 6, 105.00, 630.00, 0.00, 0.00, NULL, '2026-03-15 08:23:35.466905+00');
INSERT INTO public.order_items VALUES (59, 24, NULL, 'DOUBLE IGNITION SS CHULA', 6, 120.00, 720.00, 0.00, 0.00, NULL, '2026-03-15 08:23:35.466905+00');
INSERT INTO public.order_items VALUES (60, 24, NULL, 'Honeycomb Burner', 10, 150.00, 1500.00, 0.00, 0.00, NULL, '2026-03-15 08:23:35.466905+00');
INSERT INTO public.order_items VALUES (61, 24, NULL, 'Pinium Burner', 6, 180.00, 1080.00, 0.00, 0.00, NULL, '2026-03-15 08:23:35.466905+00');
INSERT INTO public.order_items VALUES (62, 25, NULL, 'SS Pan Support', 12, 110.00, 1320.00, 0.00, 0.00, NULL, '2026-03-15 08:26:02.716424+00');
INSERT INTO public.order_items VALUES (63, 25, NULL, 'ELBO', 20, 30.00, 600.00, 0.00, 0.00, NULL, '2026-03-15 08:26:02.716424+00');
INSERT INTO public.order_items VALUES (64, 25, NULL, 'Honeycomb Burner', 10, 150.00, 1500.00, 0.00, 0.00, NULL, '2026-03-15 08:26:02.716424+00');
INSERT INTO public.order_items VALUES (65, 25, NULL, 'CALM', 1, 380.00, 380.00, 0.00, 0.00, NULL, '2026-03-15 08:26:02.716424+00');
INSERT INTO public.order_items VALUES (66, 25, NULL, 'SINGLE IGNOITION SS CHULA', 10, 105.00, 1050.00, 0.00, 0.00, NULL, '2026-03-15 08:26:02.716424+00');
INSERT INTO public.order_items VALUES (67, 25, NULL, 'CHANNEL SINGLE', 6, 35.00, 210.00, 0.00, 0.00, NULL, '2026-03-15 08:26:02.716424+00');
INSERT INTO public.order_items VALUES (68, 26, NULL, 'মিনিকেট চাল স্ট্যান্ডার্ড (৫ কেজি) |', 1, 350.00, 350.00, 0.00, 0.00, NULL, '2026-03-15 08:29:38.55933+00');
INSERT INTO public.order_items VALUES (69, 27, NULL, '20 g Cricket Plastic Ball', 12, 20.00, 240.00, 0.00, 0.00, NULL, '2026-03-15 08:42:37.71163+00');
INSERT INTO public.order_items VALUES (70, 27, NULL, 'Orignal OSAKA Pvc Tape -Pack of 12 (White, Black, Red)', 12, 13.33, 159.96, 0.00, 0.00, NULL, '2026-03-15 08:42:37.71163+00');
INSERT INTO public.order_items VALUES (71, 28, NULL, 'Mondete Rainbow Pencils', 2, 70.00, 140.00, 0.00, 0.00, NULL, '2026-03-15 08:47:03.806906+00');
INSERT INTO public.order_items VALUES (72, 28, NULL, '20 g Cricket Plastic Ball', 12, 20.00, 240.00, 0.00, 0.00, NULL, '2026-03-15 08:47:03.806906+00');
INSERT INTO public.order_items VALUES (73, 28, NULL, 'Orignal OSAKA Pvc Tape -Pack of 12 (White, Black, Red)', 12, 13.33, 159.96, 0.00, 0.00, NULL, '2026-03-15 08:47:03.806906+00');
INSERT INTO public.order_items VALUES (74, 28, NULL, 'Magnet Premium Quality Strong Super Glue 1pcs', 12, 11.00, 132.00, 0.00, 0.00, NULL, '2026-03-15 08:47:03.806906+00');
INSERT INTO public.order_items VALUES (75, 29, NULL, '20 g Cricket Plastic Ball', 12, 20.00, 240.00, 0.00, 0.00, NULL, '2026-03-15 08:48:25.241627+00');
INSERT INTO public.order_items VALUES (76, 30, NULL, 'রাইজ ব্লেড[rise bleet]', 1, 220.00, 220.00, 0.00, 0.00, NULL, '2026-03-15 08:50:45.172191+00');
INSERT INTO public.order_items VALUES (77, 31, NULL, 'রাইজ ব্লেড[rise bleet]', 1, 220.00, 220.00, 0.00, 0.00, NULL, '2026-03-15 08:51:41.080088+00');
INSERT INTO public.order_items VALUES (78, 31, NULL, 'Magnet Premium Quality Strong Super Glue 1pcs', 12, 11.00, 132.00, 0.00, 0.00, NULL, '2026-03-15 08:51:41.080088+00');
INSERT INTO public.order_items VALUES (79, 32, NULL, 'SS Pan Support', 6, 110.00, 660.00, 0.00, 0.00, NULL, '2026-03-15 09:28:37.228239+00');


--
-- TOC entry 3847 (class 0 OID 34382)
-- Dependencies: 232
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.orders VALUES (1, 'ORD-20260205-0001', NULL, 'রিয়াদ', '01741594938', NULL, 'bottoli', 'nanbgoalcoat', 'COD', 'pending', 'cancelled', 70.00, 0.00, 0.00, 0.00, 70.00, NULL, '2026-02-05 06:07:51.184161+00', '2026-02-05 06:15:33.784069+00', NULL, '2026-02-05 06:15:33.784069+00', 'never');
INSERT INTO public.orders VALUES (2, 'ORD-20260205-0002', NULL, 'Riyad', '01641967080', NULL, 'nangoalcoat bazar', 'Nangoalcoat', 'COD', 'pending', 'cancelled', 81.48, 0.00, 0.00, 0.00, 81.48, NULL, '2026-02-05 07:04:34.996676+00', '2026-02-05 07:05:16.170978+00', NULL, '2026-02-05 07:05:16.170978+00', NULL);
INSERT INTO public.orders VALUES (4, 'ORD-20260205-0004', NULL, 'Manik', '01684665955', 'imranmanik9@gmail.com', 'Sahodeya', 'Nangolkot', 'COD', 'pending', 'cancelled', 95.00, 0.00, 0.00, 0.00, 95.00, NULL, '2026-02-05 09:30:32.211812+00', '2026-02-05 10:15:09.88822+00', NULL, '2026-02-05 10:15:09.88822+00', NULL);
INSERT INTO public.orders VALUES (3, 'ORD-20260205-0003', 1, 'Manik', '01684665955', 'imranmanik9@gmail.com', 'Sahodeya', 'Nangolkot', 'COD', 'pending', 'cancelled', 25.00, 0.00, 0.00, 0.00, 25.00, NULL, '2026-02-05 09:13:54.120992+00', '2026-02-05 10:15:36.406655+00', NULL, '2026-02-05 10:15:36.406655+00', NULL);
INSERT INTO public.orders VALUES (9, 'ORD-20260209-0003', NULL, 'Ummy hany', '01629897043', 'Ummy11@', 'Bishnopur', 'nangalkot', 'COD', 'pending', 'cancelled', 50.00, 0.00, 0.00, 0.00, 50.00, NULL, '2026-02-09 03:23:36.868684+00', '2026-02-09 06:35:50.642315+00', NULL, '2026-02-09 06:35:50.642315+00', NULL);
INSERT INTO public.orders VALUES (8, 'ORD-20260209-0002', NULL, 'Ummy hany', '01629897043', 'Ummy11@', 'Bishnopur', 'nangalkot', 'COD', 'pending', 'cancelled', 50.00, 0.00, 0.00, 0.00, 50.00, NULL, '2026-02-09 03:23:36.615407+00', '2026-02-09 06:36:03.184013+00', NULL, '2026-02-09 06:36:03.184013+00', NULL);
INSERT INTO public.orders VALUES (5, 'ORD-20260209-0001', NULL, 'Ummy hany', '01629897043', 'Ummy11@', 'Bishnopur', 'nangalkot', 'COD', 'pending', 'cancelled', 50.00, 0.00, 0.00, 0.00, 50.00, NULL, '2026-02-09 03:23:36.579247+00', '2026-02-09 06:36:16.595407+00', NULL, '2026-02-09 06:36:16.595407+00', NULL);
INSERT INTO public.orders VALUES (10, 'ORD-20260228-0001', NULL, 'Habib', '01798004936', NULL, 'jodda bazar', 'Nangoalcoat', 'COD', 'pending', 'cancelled', 210.00, 0.00, 0.00, 0.00, 210.00, NULL, '2026-02-28 08:14:46.279816+00', '2026-03-01 05:16:58.710907+00', '2026-03-01 05:16:27.516778+00', '2026-03-01 05:16:58.710907+00', NULL);
INSERT INTO public.orders VALUES (18, 'ORD-20260315-0001', NULL, 'manik', '016846659555', NULL, 'Nangolkut', 'Nangoalcoat', 'COD', 'pending', 'cancelled', 550.00, 0.00, 0.00, 0.00, 550.00, NULL, '2026-03-15 04:15:38.168496+00', '2026-03-15 04:18:41.175285+00', NULL, '2026-03-15 04:18:41.175285+00', NULL);
INSERT INTO public.orders VALUES (11, 'ORD-20260301-0001', NULL, 'Riyad', '01641967080', NULL, 'nangoalcoat bazar', 'Nangoalcoat', 'COD', 'pending', 'cancelled', 180.00, 0.00, 0.00, 0.00, 180.00, NULL, '2026-03-01 05:21:21.538725+00', '2026-03-03 07:09:01.673777+00', '2026-03-01 05:22:33.408294+00', '2026-03-03 07:09:01.673777+00', NULL);
INSERT INTO public.orders VALUES (20, 'ORD-20260315-0003', NULL, 'RAKIB HARDWARE', '01875350537', NULL, 'SHIHASSO BAZAR.JODDA ROAD', 'nangoalcoat', 'COD', 'pending', 'pending', 1965.00, 0.00, 0.00, 0.00, 1965.00, NULL, '2026-03-15 07:53:22.29571+00', '2026-03-15 07:53:22.29571+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (13, 'ORD-20260310-0001', NULL, 'Ummay Hany', '01629897043', NULL, 'Bishnopur', 'nangalkot', 'COD', 'pending', 'cancelled', 1080.00, 0.00, 0.00, 0.00, 1080.00, NULL, '2026-03-10 05:53:26.516278+00', '2026-03-10 05:55:24.012787+00', NULL, '2026-03-10 05:55:24.012787+00', NULL);
INSERT INTO public.orders VALUES (22, 'ORD-20260315-0005', NULL, 'MIRAZ MOHIMA ENTERPRISE', '01815400040', NULL, 'MIYAR BAZAR', 'nangoalcoat', 'COD', 'pending', 'pending', 1120.00, 0.00, 0.00, 0.00, 1120.00, NULL, '2026-03-15 08:21:27.029136+00', '2026-03-15 08:21:27.029136+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (14, 'ORD-20260310-0002', NULL, 'Ummy Hany Akter', '01785988248', NULL, 'Nangolkot', 'Cumilla', 'COD', 'pending', 'delivered', 360.00, 0.00, 0.00, 0.00, 360.00, NULL, '2026-03-10 06:11:46.776594+00', '2026-03-12 03:37:33.684697+00', '2026-03-12 03:37:33.684697+00', NULL, NULL);
INSERT INTO public.orders VALUES (12, 'ORD-20260302-0001', 9, 'মনি', '01765053886', 'moninoyon285@gmail.com', 'থানা নাঙ্গলকোট জেলা কুমিল্লা গ্রাম মক্রমপুর', 'কুমিল্লা', 'COD', 'pending', 'delivered', 460.00, 0.00, 0.00, 0.00, 460.00, NULL, '2026-03-02 07:15:00.226146+00', '2026-03-12 03:37:40.115602+00', '2026-03-12 03:37:40.115602+00', '2026-03-03 02:54:15.020278+00', NULL);
INSERT INTO public.orders VALUES (15, 'ORD-20260312-0001', NULL, 'jannatul RaFITA', '01891543329', NULL, 'NANGOLCOT BAZAR', 'nangoalcoat', 'COD', 'pending', 'cancelled', 3320.00, 0.00, 0.00, 0.00, 3320.00, NULL, '2026-03-12 04:18:29.723314+00', '2026-03-12 04:28:46.793478+00', NULL, '2026-03-12 04:28:46.793478+00', NULL);
INSERT INTO public.orders VALUES (16, 'ORD-20260312-0002', NULL, 'jannatul RaFITA', '01891543329', NULL, 'NANGOLCOT BAZAR', 'nangoalcoat', 'COD', 'pending', 'delivered', 2330.00, 0.00, 0.00, 0.00, 2330.00, NULL, '2026-03-12 04:28:07.294342+00', '2026-03-13 13:04:32.558341+00', '2026-03-13 13:04:32.558341+00', NULL, NULL);
INSERT INTO public.orders VALUES (17, 'ORD-20260314-0001', NULL, 'nangoalcoat kamil madrasa', '01825101003', NULL, 'NANGOLCOT BAZAR', 'nangoalcoat', 'COD', 'pending', 'delivered', 1800.00, 0.00, 0.00, 0.00, 1800.00, NULL, '2026-03-14 15:32:08.235135+00', '2026-03-14 15:35:52.517102+00', '2026-03-14 15:35:52.517102+00', NULL, NULL);
INSERT INTO public.orders VALUES (24, 'ORD-20260315-0007', NULL, 'ISMAIL TRADERS', '01834017816', NULL, 'MUNSIRHAAT BAZAR', 'LAKSAM', 'COD', 'pending', 'pending', 3930.00, 0.00, 0.00, 0.00, 3930.00, NULL, '2026-03-15 08:23:35.466905+00', '2026-03-15 08:23:35.466905+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (25, 'ORD-20260315-0008', NULL, 'SOBUJ HARDWARE', '01880344708', NULL, 'HAZIPURA', 'MONOHORGANJ', 'COD', 'pending', 'pending', 5060.00, 0.00, 0.00, 0.00, 5060.00, NULL, '2026-03-15 08:26:02.716424+00', '2026-03-15 08:26:02.716424+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (23, 'ORD-20260315-0006', NULL, 'MIRAZ MOHIMA ENTERPRISE', '01815400040', NULL, 'MIYAR BAZAR', 'nangoalcoat', 'COD', 'pending', 'cancelled', 1120.00, 0.00, 0.00, 0.00, 1120.00, NULL, '2026-03-15 08:21:27.11294+00', '2026-03-15 08:27:50.455179+00', NULL, '2026-03-15 08:27:50.455179+00', '2 LINE DEEWA HOICE');
INSERT INTO public.orders VALUES (26, 'ORD-20260315-0009', NULL, 'JANNATUL NAIM', '01616933182', NULL, 'NANGOLCOT', 'nangoalcoat', 'COD', 'pending', 'pending', 350.00, 0.00, 0.00, 0.00, 350.00, NULL, '2026-03-15 08:29:38.55933+00', '2026-03-15 08:29:38.55933+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (19, 'ORD-20260315-0002', NULL, 'Riyad Store', '01879709354', NULL, 'বাগমারা বাজার', 'নাঙ্গলকোট, জোড্ডা রোড', 'COD', 'pending', 'cancelled', 360.00, 0.00, 0.00, 0.00, 360.00, NULL, '2026-03-15 05:45:16.982148+00', '2026-03-15 08:35:44.84784+00', NULL, '2026-03-15 08:35:44.84784+00', NULL);
INSERT INTO public.orders VALUES (27, 'ORD-20260315-0010', NULL, 'RIYAD STORE', '01879709354', NULL, 'BAGMARA ,JODDA ROAD', 'nanbgoalcoat', 'COD', 'pending', 'pending', 399.96, 0.00, 0.00, 0.00, 399.96, NULL, '2026-03-15 08:42:37.71163+00', '2026-03-15 08:42:37.71163+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (28, 'ORD-20260315-0011', NULL, 'AKRAMUL STORE', '01822892734', NULL, 'BAGMARA ,JODDA ROAD', 'nangoalcoat', 'COD', 'pending', 'pending', 671.96, 0.00, 0.00, 0.00, 671.96, NULL, '2026-03-15 08:47:03.806906+00', '2026-03-15 08:47:03.806906+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (29, 'ORD-20260315-0012', NULL, 'MOHIN STORE', '01895255931', NULL, 'SRIHASSO BAZAR', 'nanbgoalcoat', 'COD', 'pending', 'pending', 240.00, 0.00, 0.00, 0.00, 240.00, NULL, '2026-03-15 08:48:25.241627+00', '2026-03-15 08:48:25.241627+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (30, 'ORD-20260315-0013', NULL, 'AHSAN STORE', '01760923353', NULL, 'TALTOTA BAZAR', 'NANGOLKOT BAZAR', 'COD', 'pending', 'pending', 220.00, 0.00, 0.00, 0.00, 220.00, NULL, '2026-03-15 08:50:45.172191+00', '2026-03-15 08:50:45.172191+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (31, 'ORD-20260315-0014', NULL, 'SAHABUDDIN STORE', '01880054160', NULL, 'TALTOLA BAZAR', 'NANGOLKOT BAZAR', 'COD', 'pending', 'pending', 352.00, 0.00, 0.00, 0.00, 352.00, NULL, '2026-03-15 08:51:41.080088+00', '2026-03-15 08:51:41.080088+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (32, 'ORD-20260315-0015', NULL, 'al arafa trades', '01775621230', NULL, 'SOKTOLA', 'LAKSAM', 'COD', 'pending', 'pending', 660.00, 0.00, 0.00, 0.00, 660.00, NULL, '2026-03-15 09:28:37.228239+00', '2026-03-15 09:28:37.228239+00', NULL, NULL, NULL);
INSERT INTO public.orders VALUES (21, 'ORD-20260315-0004', NULL, 'AL ARAFA TRADERS', '01775621230', NULL, 'SOKTOLA', 'LAKSAM', 'COD', 'pending', 'cancelled', 110.00, 0.00, 0.00, 0.00, 110.00, NULL, '2026-03-15 07:57:12.770355+00', '2026-03-15 09:29:16.63704+00', NULL, '2026-03-15 09:29:16.63704+00', NULL);


--
-- TOC entry 3863 (class 0 OID 34936)
-- Dependencies: 248
-- Data for Name: product_reviews; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--



--
-- TOC entry 3861 (class 0 OID 34913)
-- Dependencies: 246
-- Data for Name: product_variations; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--



--
-- TOC entry 3859 (class 0 OID 34848)
-- Dependencies: 244
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.products VALUES (84, 'Lux Nourished Glow Soap 98G', '', 2, NULL, NULL, NULL, '1122367', 1, 18, '', 'public/images/products/lux-nourished-glow-soap.webp', '{}', 110.00, 110.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-21 08:47:11.398267+00', '2026-03-04 04:05:28.817687+00', 105.00, 10.00);
INSERT INTO public.products VALUES (15, 'হলুদ গুরা/turmeric powder', '', 2, NULL, NULL, NULL, '112246', 1, 1, '', 'public/images/products/-turmeric-powder.jpg', '{}', 550.00, 520.00, 1.00, 25.00, 25.00, 0.00, 'flat', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:28:08.418878+00', '2026-03-07 04:42:43.854703+00', 400.00, 5.00);
INSERT INTO public.products VALUES (78, 'LG Glass Single Gas Stove', '', 3, NULL, NULL, NULL, '11223654-78', 1, 23, '', 'public/images/products/lg-glass-single-gas-stove.jpeg', '{}', 2300.00, 2070.00, 10.00, 50.00, 0.00, 0.00, 'percent', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:21:45.415254+00', '2026-02-19 04:21:45.415254+00', 0.00, 1.00);
INSERT INTO public.products VALUES (119, 'Steel Majuni Scrubber -2 PC', 'Steel Majuni 1 PC Scrubber
High Quality Product (Brand Product)
Low Price than Market
Color: Silver
Easy to Use
Long Lasting', 2, NULL, NULL, NULL, 'PRD-26-2-119', 1, NULL, '', 'public/images/products/thumbprd-26-2-119.webp', NULL, 60.00, 60.00, 2.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-11 04:28:40.351897+00', '2026-03-11 04:28:40.38087+00', 10.00, 10.00);
INSERT INTO public.products VALUES (134, 'radhuni murgir masala', '', 2, NULL, NULL, NULL, 'PRD-26-2-134', 1, 2, '', 'public/images/products/thumbprd-26-2-134.jpg', NULL, 95.00, 95.00, 1.00, 10.00, 0.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-12 04:15:38.971704+00', '2026-03-12 04:15:38.98727+00', 90.00, 1.00);
INSERT INTO public.products VALUES (75, 'Pinium Burner', '', 3, NULL, NULL, NULL, '11223569', 1, 23, '', 'public/images/products/pinium-burner.jpeg', '{}', 280.00, 160.00, 10.00, 100.00, 0.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:15:17.452394+00', '2026-03-15 08:13:42.139357+00', 180.00, 1.00);
INSERT INTO public.products VALUES (11, 'পেঁয়াজ/Onion', '🧅 দেশী পেঁয়াজ
রান্নার আসল স্বাদ আসে ভালো মানের দেশী পেঁয়াজ থেকে।
ঝাঁঝালো ঘ্রাণ, ঘন স্বাদ আর প্রতিটি তরকারিতে বাড়তি মজা — এটাই আসল দেশী পেঁয়াজের বিশেষত্ব।
✔️ তাজা ও মানসম্মত
✔️ রান্নায় স্বাদ বাড়ায় কয়েক গুণ
✔️ ভর্তা, ভাজি, তরকারি—সবকিছুর জন্য উপযোগী
প্রতিদিনের রান্না হোক আরও সুস্বাদু ও পরিপূর্ণ।
🛒 আজই সংগ্রহ করুন Noor Super Mart থেকে।
সাশ্রয়ী দাম, নিশ্চিত মান, আর দ্রুত ডেলিভারি—আপনার দোরগোড়ায় 🚚✨', 2, NULL, NULL, NULL, '112242', 1, 1, '', 'public/images/products/onion.jpg', '{}', 60.00, 55.00, 1.00, 200.00, 20.00, 0.00, 'flat', 5.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:11:56.204528+00', '2026-02-28 07:53:17.685084+00', 0.00, 1.00);
INSERT INTO public.products VALUES (83, 'Lux Soap', '', 2, NULL, NULL, NULL, '1122366', 1, 20, '', 'public/images/products/lux-soap.jpg', '{}', 50.00, 50.00, 5.00, 100.00, 100.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-21 08:40:55.833407+00', '2026-02-23 05:44:44.951888+00', 0.00, 1.00);
INSERT INTO public.products VALUES (81, 'LG Double Glass Gas Stove', '', 3, NULL, NULL, NULL, '112236654', 1, 23, '', 'public/images/products/lg-double-glass-gas-stove.jpg', '{}', 4500.00, 4200.00, 1.00, 50.00, 50.00, 0.00, 'flat', 300.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:30:16.174721+00', '2026-02-23 07:14:27.564424+00', 0.00, 1.00);
INSERT INTO public.products VALUES (82, 'Digital Scale 40kg', '', 3, NULL, NULL, NULL, '1122365-82', 1, 23, '', 'public/images/products/digital-scale-40kg.jpg', '{}', 2700.00, 2550.00, 1.00, 50.00, 50.00, 0.00, 'flat', 150.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:33:55.243546+00', '2026-02-23 07:02:12.996247+00', 0.00, 1.00);
INSERT INTO public.products VALUES (76, '100 MM Pinium', '🔹 100mm Pinium
উচ্চমানের 100mm Pinium আপনার দৈনন্দিন কাজকে করে আরও সহজ, দ্রুত এবং কার্যকর। শক্তিশালী ও টেকসই নির্মাণের কারণে এটি দীর্ঘদিন ব্যবহারযোগ্য।
✨ বিশেষ বৈশিষ্ট্য:
প্রিমিয়াম মানের উপাদান
শক্তিশালী ও দীর্ঘস্থায়ী
সহজে ব্যবহারযোগ্য
নির্ভরযোগ্য পারফরম্যান্স
💡 কেন কিনবেন:
দৈনন্দিন কাজের জন্য সঠিক মানের টুল/উপকরণ, যা আপনাকে সময় ও শ্রম বাঁচাবে।
📦 অর্ডার করুন আজই এবং আপনার কাজকে আরও প্রফেশনাল বানান!', 3, NULL, NULL, NULL, '1122365-76', 1, 23, '', 'public/images/products/100-mm-pinium.jpeg', '{}', 150.00, 140.00, 1.00, 100.00, 100.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:16:49.251511+00', '2026-02-23 08:07:15.67516+00', 0.00, 1.00);
INSERT INTO public.products VALUES (77, 'Gas Saver', '', 3, NULL, NULL, NULL, '11223654-77', 1, 23, '', 'public/images/products/gas-saver.jpeg', '{}', 220.00, 200.00, 10.00, 1000.00, 1000.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:19:25.571029+00', '2026-02-23 07:16:59.163157+00', 0.00, 1.00);
INSERT INTO public.products VALUES (6, 'চিনিগুড়া চাল', 'চিনিগুড়া চাল
সুগন্ধি, ঝরঝরে আর নরম ভাতের জন্য সেরা পছন্দ — মানসম্মত চিনিগুড়া চাল।
পোলাও, কাচ্চি, বিরিয়ানি বা বিশেষ দিনের আয়োজন—সবকিছুতেই এনে দেয় আলাদা সুবাস আর স্বাদ।
✔️ সুগন্ধে ভরপুর
✔️ রান্নার পর ঝরঝরে থাকে
✔️ বিশেষ অনুষ্ঠান ও দৈনন্দিন ব্যবহারের জন্য উপযোগী
পরিবারের সাথে জমে উঠুক প্রতিটি খাবারের মুহূর্ত।
🛒 আজই সংগ্রহ করুন Noor Super Mart থেকে।
নিশ্চিত মান, সাশ্রয়ী দাম ও দ্রুত ডেলিভারি—সব একসাথে আপনার দোরগোড়ায় 🚚✨', 2, NULL, NULL, NULL, '112238', 1, 1, '', 'public/images/products/thumb112238.jpg', '{}', 180.00, 150.00, 1.00, 50.00, 50.00, 0.00, 'flat', 30.00, 0.00, 'exclusive', 0.00, 'static', true, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 06:43:31.170873+00', '2026-02-28 08:00:42.806325+00', 0.00, 1.00);
INSERT INTO public.products VALUES (1, 'Ramadaan Special', '', 1, NULL, NULL, NULL, '112233', 1, 6, '', 'public/images/products/ramadaan-special.jpg', '{}', 2850.00, 2500.00, 1.00, 10.00, 10.00, 0.00, 'flat', 350.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 05:46:26.405123+00', '2026-02-28 08:10:27.894307+00', 0.00, 1.00);
INSERT INTO public.products VALUES (16, 'মরিচ গুরা/Chilli powder', '', 2, NULL, NULL, NULL, '112247', 1, 1, '', 'public/images/products/-chilli-powder.jpg', '{}', 590.00, 550.00, 10.00, 20.00, 0.00, 0.00, 'flat', 40.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:33:08.703543+00', '2026-02-05 07:33:08.703543+00', 0.00, 1.00);
INSERT INTO public.products VALUES (17, 'মরিয়ম খেজুর/Dates', '', 2, NULL, NULL, NULL, '112248', 1, 1, '', 'public/images/products/-dates.jpg', '{}', 1800.00, 1600.00, 10.00, 10.00, 0.00, 0.00, 'flat', 200.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:37:03.565212+00', '2026-02-05 07:37:03.565212+00', 0.00, 1.00);
INSERT INTO public.products VALUES (19, 'চিড়া/Chira', '', 2, NULL, NULL, NULL, '112250', 1, 1, '', 'public/images/products/chira.jpg', '{}', 100.00, 90.00, 10.00, 100.00, 0.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:42:55.43568+00', '2026-02-05 07:42:55.43568+00', 0.00, 1.00);
INSERT INTO public.products VALUES (21, 'সরিষার তেল/mustard oil', '', 2, NULL, NULL, NULL, '112252-21', 1, 3, '', 'public/images/products/-mustard-oil.jpg', '{}', 350.00, 315.00, 10.00, 120.00, 0.00, 0.00, 'percent', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:51:40.466434+00', '2026-02-05 07:51:40.466434+00', 0.00, 1.00);
INSERT INTO public.products VALUES (24, 'দাব্বাস খেজুর', '', 2, NULL, NULL, NULL, '112252-24', 1, 1, '', 'public/images/products/-.png', '{}', 350.00, 320.00, 10.00, 20.00, 0.00, 0.00, 'flat', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 10:27:13.583991+00', '2026-02-05 10:28:13.732859+00', 0.00, 1.00);
INSERT INTO public.products VALUES (23, 'ডানো ফুল ক্রিম মিল্ক পাওডার', '', 2, NULL, NULL, NULL, '112254', 1, 8, '', 'public/images/products/----.webp', '{}', 480.00, 460.00, 10.00, 20.00, 0.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 08:04:46.392214+00', '2026-02-05 08:04:46.392214+00', 0.00, 1.00);
INSERT INTO public.products VALUES (25, 'chikpeas/ছোলা', '', 2, NULL, NULL, NULL, '112257', 1, 1, '', 'public/images/products/chikpeas.jpg', '{}', 110.00, 105.00, 10.00, 50.00, 0.00, 0.00, 'flat', 5.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 03:51:11.506687+00', '2026-02-07 03:51:11.506687+00', 0.00, 1.00);
INSERT INTO public.products VALUES (26, 'Gas stove glass single', '', 3, NULL, NULL, NULL, '112260', 1, 22, '', 'public/images/products/gas-stove-glass-single.jpg', '{}', 2250.00, 2150.00, 10.00, 20.00, 0.00, 0.00, 'flat', 100.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 04:02:46.331448+00', '2026-02-07 04:02:46.331448+00', 0.00, 1.00);
INSERT INTO public.products VALUES (27, 'Tea Bag', '', 2, NULL, NULL, NULL, '112261', 1, 11, '', 'public/images/products/tea-bag.jpg', '{}', 90.00, 90.00, 10.00, 50.00, 0.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 04:11:33.258691+00', '2026-02-07 04:11:33.258691+00', 0.00, 1.00);
INSERT INTO public.products VALUES (29, 'ধনিয়া /Coriander', '', 2, NULL, NULL, NULL, '112263', 1, 4, '', 'public/images/products/-coriander.jpg', '{}', 480.00, 450.00, 10.00, 100.00, 0.00, 0.00, 'flat', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 04:38:00.008498+00', '2026-02-07 04:38:00.008498+00', 0.00, 1.00);
INSERT INTO public.products VALUES (30, 'জিরা/Cumin', '', 2, NULL, NULL, NULL, '112264', 1, 1, '', 'public/images/products/cumin.jpg', '{}', 980.00, 950.00, 10.00, 100.00, 0.00, 0.00, 'flat', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 04:42:47.94981+00', '2026-02-07 04:44:24.192875+00', 0.00, 1.00);
INSERT INTO public.products VALUES (31, 'গরম মসলা/Garam Masala', '', 2, NULL, NULL, NULL, '112265', 1, 7, '', 'public/images/products/-garam-masala.jpg', '{}', 80.00, 50.00, 10.00, 100.00, 0.00, 0.00, 'percentage', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 05:11:22.808539+00', '2026-02-07 05:19:12.801763+00', 0.00, 1.00);
INSERT INTO public.products VALUES (32, 'দারুচিনি/Cinnamon', '', 2, NULL, NULL, NULL, '112266', 1, 2, '', 'public/images/products/cinnamon.jpg', '{}', 100.00, 80.00, 10.00, 100.00, 0.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 05:16:41.765053+00', '2026-02-07 05:16:41.765053+00', 0.00, 1.00);
INSERT INTO public.products VALUES (34, 'এলাচ/Cardamom', '', 2, NULL, NULL, NULL, '112268', 1, 7, '', 'public/images/products/cardamom.jpg', '{}', 150.00, 150.00, 10.00, 100.00, 0.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 07:44:30.659243+00', '2026-02-07 07:44:30.659243+00', 0.00, 1.00);
INSERT INTO public.products VALUES (35, 'তেজপাতা/Bay Leaf', '', 2, NULL, NULL, NULL, '1122690', 1, 7, '', 'public/images/products/bay-leaf.jpg', '{}', 30.00, 25.00, 10.00, 20.00, 0.00, 0.00, 'flat', 5.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 08:34:25.509604+00', '2026-02-07 08:34:25.509604+00', 0.00, 1.00);
INSERT INTO public.products VALUES (36, 'কফি/coffee', '', 1, NULL, NULL, NULL, '112270', 1, 10, '', 'public/images/products/coffee.jpg', '{}', 1000.00, 950.00, 10.00, 100.00, 0.00, 0.00, 'flat', 50.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 09:00:55.924135+00', '2026-02-07 09:00:55.924135+00', 0.00, 1.00);
INSERT INTO public.products VALUES (37, 'Nutty Biscuit', '', 2, NULL, NULL, NULL, '112272', 1, 15, '', 'public/images/products/nutty-biscuit.jpg', '{}', 50.00, 50.00, 10.00, 200.00, 0.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 09:54:38.26073+00', '2026-02-07 09:54:38.26073+00', 0.00, 1.00);
INSERT INTO public.products VALUES (45, 'চিনি/suger', '', 2, NULL, NULL, NULL, '112282', 1, 1, '', 'public/images/products/suger.jpg', '{}', 110.00, 105.00, 1.00, 100.00, 0.00, 0.00, 'flat', 5.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 05:24:09.541363+00', '2026-03-12 04:22:13.209487+00', 100.00, 1.00);
INSERT INTO public.products VALUES (43, 'কাবেরি মেহেদি/Kaberi mehedi', '', 2, NULL, NULL, NULL, '112280', 1, 23, '', 'public/images/products/-kaberi-mehedi.jpg', '{}', 70.00, 50.00, 10.00, 100.00, 0.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 04:15:26.74473+00', '2026-02-08 04:15:26.74473+00', 0.00, 1.00);
INSERT INTO public.products VALUES (44, 'ডিম/Egg', '', 2, NULL, NULL, NULL, '112281', 1, 23, '', 'public/images/products/egg.jpg', '{}', 10.00, 10.00, 100.00, 1000.00, 0.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 04:23:35.824201+00', '2026-02-08 04:23:35.824201+00', 0.00, 1.00);
INSERT INTO public.products VALUES (42, 'স্মার্ট কোন মেহেদি', '', 2, NULL, NULL, NULL, '112279', 1, 23, '', 'public/images/products/thumb112279.jpg', '{}', 70.00, 50.00, 1.00, 100.00, 100.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 04:11:41.307109+00', '2026-03-01 05:03:36.597612+00', 0.00, 1.00);
INSERT INTO public.products VALUES (39, 'চিপস', '', 2, NULL, NULL, NULL, '112274', 1, 21, '', 'public/images/products/thumb112274.jpg', '{}', 10.00, 10.00, 1.00, 500.00, 0.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 10:04:16.182258+00', '2026-02-28 02:36:21.088405+00', 0.00, 1.00);
INSERT INTO public.products VALUES (38, 'বারবিকিউ চানাচুর', '', 2, NULL, NULL, NULL, '112273', 1, 25, '', 'public/images/products/thumb112273.jpg', '{}', 50.00, 50.00, 1.00, 100.00, 100.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 09:58:50.987286+00', '2026-02-28 07:15:03.901875+00', 0.00, 1.00);
INSERT INTO public.products VALUES (22, 'বনফুল লাচ্চা সেমাই', '', 2, NULL, NULL, NULL, '112253', 1, 10, '', 'public/images/products/thumb112253.jpg', '{}', 50.00, 50.00, 1.00, 50.00, 50.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:54:07.837328+00', '2026-02-28 07:10:09.509584+00', 0.00, 1.00);
INSERT INTO public.products VALUES (20, 'খেসারি ডাল', '', 2, NULL, NULL, NULL, '112251', 1, 1, '', 'public/images/products/thumb112251.jpg', '{}', 120.00, 110.00, 1.00, 100.00, 100.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:46:30.962635+00', '2026-02-28 07:14:08.129777+00', 0.00, 1.00);
INSERT INTO public.products VALUES (28, 'গোল মরিচ', '', 2, NULL, NULL, NULL, '112262', 1, 4, '', 'public/images/products/thumb112262.jpg', '{}', 350.00, 300.00, 1.00, 100.00, 100.00, 0.00, 'flat', 50.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 04:31:28.932876+00', '2026-02-28 07:15:48.013721+00', 0.00, 1.00);
INSERT INTO public.products VALUES (46, 'IGT Gas sensor Regulator', '', 3, NULL, NULL, NULL, '112283', 1, 23, '', 'public/images/products/igt-gas-sensor-regulator.jpg', '{}', 1350.00, 1150.00, 10.00, 50.00, 0.00, 0.00, 'flat', 200.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 05:25:55.493441+00', '2026-02-08 05:25:55.493441+00', 0.00, 1.00);
INSERT INTO public.products VALUES (47, 'Fruit Combo', '', 4, NULL, NULL, NULL, '112284', 1, 12, '', 'public/images/products/fruit-combo.jpeg', '{}', 1100.00, 1000.00, 5.00, 20.00, 0.00, 0.00, 'flat', 100.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 09:07:22.44218+00', '2026-02-08 09:07:22.44218+00', 0.00, 1.00);
INSERT INTO public.products VALUES (48, 'আপেল/Apple', '', 4, NULL, NULL, NULL, '112285', 1, 1, '', 'public/images/products/apple.jpg', '{}', 300.00, 280.00, 10.00, 20.00, 0.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 09:10:51.429173+00', '2026-02-08 09:10:51.429173+00', 0.00, 1.00);
INSERT INTO public.products VALUES (112, 'Omera 12kg Cylinder', '', 2, 2, 1, 1, '1122', 1, 17, '', 'public/images/products/thumb1122.png', NULL, 1800.00, 1800.00, 1.00, 10.00, 10.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-06 15:10:12.283015+00', '2026-03-06 15:10:12.294215+00', 1800.00, 1.00);
INSERT INTO public.products VALUES (120, 'Glue & Tapes -3', '', 2, NULL, NULL, NULL, '1122560', 1, NULL, '', 'public/images/products/thumb1122560.jpg', NULL, 90.00, 90.00, 3.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-11 05:03:47.039846+00', '2026-03-11 05:03:47.050075+00', 30.00, 6.00);
INSERT INTO public.products VALUES (49, 'মালতা/Malta', '', 4, NULL, NULL, NULL, '112286', 1, 1, '', 'public/images/products/malta.jpg', '{}', 350.00, 340.00, 1.00, 100.00, 0.00, 0.00, 'flat', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 09:13:32.632812+00', '2026-03-12 04:10:39.712238+00', 310.00, 1.00);
INSERT INTO public.products VALUES (79, 'Pan Support', '', 3, NULL, NULL, NULL, '225863', 1, 23, '', 'public/images/products/pan-support.jpg', '{}', 150.00, 70.00, 100.00, 100.00, 0.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:24:35.085857+00', '2026-03-15 08:16:32.106685+00', 110.00, 1.00);
INSERT INTO public.products VALUES (80, 'Serbo Meter Regulator', '', 3, NULL, NULL, NULL, '45000', 1, 23, '', 'public/images/products/serbo-meter-regulator.jpeg', '{}', 450.00, 380.00, 100.00, 100.00, 0.00, 0.00, 'flat', 70.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:26:55.663955+00', '2026-02-19 04:26:55.663955+00', 0.00, 1.00);
INSERT INTO public.products VALUES (40, 'Bag', '', 2, NULL, NULL, NULL, '1458712', 1, 17, '', 'public/images/products/bag.webp', '{}', 2889.00, 2449.00, 10.00, 100.00, 0.00, 0.00, 'flat', 440.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 10:31:16.910894+00', '2026-02-07 10:31:16.910894+00', 0.00, 1.00);
INSERT INTO public.products VALUES (86, 'Bahrain Food item Mango Tang instant drink powder 500gm', '', 2, NULL, NULL, NULL, '14582684', 1, 8, '', 'public/images/products/bahrain-food-item-mango-tang-instant-drink-powder.png', '{}', 350.00, 400.00, 1.00, 100.00, 100.00, 0.00, '', 50.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-22 03:52:08.21865+00', '2026-02-23 05:08:16.64292+00', 0.00, 1.00);
INSERT INTO public.products VALUES (87, 'Tang Orange Drink Powder (India) 500gm', '', 2, NULL, NULL, NULL, '1122380', 1, 8, '', 'public/images/products/tang-orange-drink-powder-india.jpg', '{}', 320.00, 350.00, 1.00, 100.00, 100.00, 0.00, '', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-22 16:17:57.427164+00', '2026-02-23 05:22:34.158563+00', 0.00, 1.00);
INSERT INTO public.products VALUES (53, 'fresh apple', '🍎 আপেল – Noor Super Mart
সতেজ, রসালো এবং পুষ্টিকর আপেল আপনার স্বাস্থ্য ও স্বাদ দুটোই ধরে রাখে। প্রতিদিনের খাদ্যাভ্যাসে এটি যুক্ত করলে শরীর সুস্থ, শক্তিশালী এবং রোগপ্রতিরোধ ক্ষমতা বৃদ্ধি পায়।
✨ স্বাস্থ্য ও উপকারিতা:
১০০% সতেজ ও প্রাকৃতিক, সরাসরি আপনার কাছে পৌঁছে যায়
ভিটামিন C ও অ্যান্টিঅক্সিডেন্ট সমৃদ্ধ, যা শরীরের প্রতিরোধ ক্ষমতা বাড়ায়
ফাইবার সমৃদ্ধ, হজম শক্তি বৃদ্ধি করে এবং পাকস্থলীর স্বাস্থ্য বজায় রাখে
প্রাকৃতিকভাবে ওজন নিয়ন্ত্রণে সহায়ক
হৃদয় ও মস্তিষ্কের জন্য উপকারী, সার্বিক শক্তি বৃদ্ধি করে
💡 কেন Noor Super Mart থেকে কিনবেন:
সর্বোচ্চ মানের সতেজ আপেল, যাতে প্রতিটি কামড়ই খেতে সুস্বাদু
বাড়িতে ডেলিভারি সুবিধা, নিরাপদ ও ঝামেলাহীন
মূল্য সাশ্রয়ী এবং ক্রেতাদের জন্য নিয়মিত স্টক
📦 উপলব্ধ প্যাক সাইজ:
  কেজি  অনুযায়ী
🌱 আজই অর্ডার করুন এবং স্বাস্থ্যকর, সতেজ ও প্রাকৃতিক আপেল উপভোগ করুন Noor Super Mart থেকে! 🍎✨', 2, NULL, NULL, NULL, '112296', 1, 1, '', 'public/images/products/fresh-apple.jpg', '{}', 350.00, 300.00, 1.00, 100.00, 100.00, 0.00, 'flat', 50.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-09 09:05:39.704174+00', '2026-02-23 08:51:44.420015+00', 0.00, 1.00);
INSERT INTO public.products VALUES (62, 'তেলাপিয়া মাছ', '🐟 তেলাপিয়া মাছ – Noor Super Mart
সতেজ, পুষ্টিকর এবং স্বাদে টকস তেলাপিয়া মাছ আপনার দৈনন্দিন খাদ্যাভ্যাসকে করে আরও স্বাস্থ্যসম্মত। প্রোটিনে ভরপুর এই মাছ পেশী গঠন, শক্তি বৃদ্ধি এবং সার্বিক সুস্থতা বজায় রাখতে সাহায্য করে।
✨ সুবিধা ও উপকারিতা:
১০০% সতেজ ও হাইজিনিক, সরাসরি আপনার দরজায় ডেলিভারি
প্রোটিন এবং ওমেগা-৩ সমৃদ্ধ, যা হৃদয় ও মস্তিষ্কের জন্য উপকারী
হজম সহজ করে এবং পাচন শক্তি বৃদ্ধি করে
দীর্ঘমেয়াদি স্বাস্থ্য ও রোগ প্রতিরোধ ক্ষমতা বাড়াতে সহায়ক
সহজ অর্ডার, নিরাপদ পেমেন্ট ও প্যাক/কেজি অনুযায়ী ক্রয় সুবিধা
💡 কেন Noor Super Mart থেকে কিনবেন:
Noor Super Mart সরবরাহ করে সর্বোচ্চ মানের সতেজ মাছ, যা স্বাস্থ্যসম্মত ও স্বাদে উজ্জ্বল। প্রতিদিনের খাবারের জন্য এটি একটি সেরা পছন্দ।
📦 উপলব্ধ প্যাক সাইজ:
পিস/কেজি অনুযায়ী
🌱 আজই অর্ডার করুন এবং স্বাস্থ্যকর, সতেজ ও প্রোটিন সমৃদ্ধ তেলাপিয়া মাছ উপভোগ করুন Noor Super Mart থেকে! 🐟✨', 5, NULL, NULL, NULL, '112297-62', 1, 1, '', 'public/images/products/thumb112297-62.jpg', '{}', 200.00, 190.00, 1.00, 20.00, 20.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-09 10:24:47.504191+00', '2026-03-01 05:04:23.591323+00', 0.00, 1.00);
INSERT INTO public.products VALUES (88, 'Olympic Foodie Masala Noodles 8pcs pack', 'Olympic Industries Limited, is one of the longest running and most reputed manufacturing-based companies in Bangladesh, with a heritage of over 50 years and group profile including interests in Pharmaceuticals, Power, and Information Technology, among other FMCG. Starting with steel production dating back to 1950, Olympic has steadily diversified over the years into various consumer goods including biscuits, confectioneries, batteries, and ball pens, with over 40 brands and 53 SKU’s. The companyhas gotten to where it is today by staying true to its core beliefs, in providing high quality, innovative products which its consumers can rely on.

Brand: Olympic
Net weight: 496gm
Product type: Noodles
Origin: Bangladesh
Manufacturer: Olympic Industries Limited
1 Carton x 12 Packs
Marketed and Distributed by Olympic Industries Limited', 2, NULL, NULL, NULL, '1122549', 1, 24, '', 'public/images/products/olympic-foodie-masala-noodles-8pcs-pack.webp', NULL, 150.00, 170.00, 3.00, 100.00, 10.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-23 04:24:48.598346+00', '2026-03-04 04:04:02.313483+00', 140.00, 1.00);
INSERT INTO public.products VALUES (113, 'Kaveri No.1 Indian Cone Mehedi (12PCS PACK)', '', 2, NULL, NULL, NULL, '1122550', 1, NULL, '', 'public/images/products/thumb1122550.jpg', NULL, 360.00, 360.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-07 04:25:54.935802+00', '2026-03-07 04:25:54.941624+00', 250.00, 5.00);
INSERT INTO public.products VALUES (67, 'মাস কলাইর ডাল(Black Gram / Urad Dal)', '🌿 মাস কলাইর ডাল (Black Gram / Urad Dal) – স্বাস্থ্যকর ও পুষ্টিকর
মাস কলাইর ডাল প্রাকৃতিকভাবে প্রোটিন ও ফাইবারে সমৃদ্ধ, যা শরীরের জন্য শক্তি ও পুষ্টি যোগায়। এটি হজম শক্তি বৃদ্ধি, রোগ প্রতিরোধ ক্ষমতা ও সার্বিক স্বাস্থ্য উন্নত করতে সহায়তা করে।
✨ উপকারিতা:
প্রোটিন ও পুষ্টিতে সমৃদ্ধ, শক্তি বৃদ্ধি করে
হজম শক্তি ভালো রাখে, পাচন সহজ করে
রক্তচাপ ও রক্তে শর্করার মাত্রা নিয়ন্ত্রণে সহায়ক
হাড় ও দাতের স্বাস্থ্য সমর্থন করে
দীর্ঘস্থায়ী শক্তি ও রোগ প্রতিরোধ ক্ষমতা বাড়ায়
💡 কেন ব্যবহার করবেন:
প্রতিদিনের খাদ্যাভ্যাসে মাস কলাইর ডাল ব্যবহার করলে শরীর সুস্থ, শক্তিশালী এবং পুষ্টিকর হয়।
📦 প্যাক সাইজ:
২৫০g / ৫০০g / ১kg
আজই অর্ডার করুন এবং স্বাস্থ্যের জন্য প্রাকৃতিক শক্তি যোগ করুন 🌿✨', 2, NULL, NULL, NULL, '112301', 1, 1, '', 'public/images/products/thumb112301.jpg', '{}', 225.00, 225.00, 1.00, 20.00, 20.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-10 05:48:17.772579+00', '2026-03-07 04:30:50.998807+00', 180.00, 1.00);
INSERT INTO public.products VALUES (85, 'Sunsilk Lusciously Thick & Long Shampoo', '', 2, NULL, NULL, NULL, '1122368', 1, 14, '', 'public/images/products/sunsilk-lusciously-thick--long-shampoo.jpg', '{}', 550.00, 650.00, 1.00, 100.00, 100.00, 0.00, '', 100.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-21 08:52:13.471505+00', '2026-02-23 05:15:41.284402+00', 0.00, 1.00);
INSERT INTO public.products VALUES (90, 'Olympic foodie instant noodles masala 12pack 744gm', '', 2, NULL, NULL, NULL, '1122571', 1, NULL, '', 'public/images/products/olympic-foodie-instant-noodles-masala-12pack-744gm.webp', NULL, 210.00, 255.00, 1.00, 100.00, 100.00, 0.00, '', 45.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-23 05:38:57.790923+00', '2026-02-23 05:40:35.610472+00', 0.00, 1.00);
INSERT INTO public.products VALUES (73, 'SS Pan Support', '🔥 SS Pan Support
উচ্চমানের স্টেইনলেস স্টিল দিয়ে তৈরি এই SS Pan Support আপনার গ্যাস স্টোভকে করে আরও নিরাপদ ও স্থিতিশীল। রান্নার সময় হাঁড়ি, কড়াই বা প্যান যেন না নড়ে যায়—তার জন্য এটি শক্তভাবে ধরে রাখে।
✨ বিশেষ বৈশিষ্ট্য:
প্রিমিয়াম স্টেইনলেস স্টিল (SS) নির্মাণ
মরিচা প্রতিরোধী ও টেকসই
ভারী হাঁড়ি ও কড়াই সহজে ধরে রাখে
উচ্চ তাপ সহনশীল
সহজে পরিষ্কার করা যায়
🍳 উপযোগিতা:
ডাবল/সিঙ্গেল বার্নার গ্যাস স্টোভের জন্য উপযুক্ত।
দীর্ঘদিন নিরাপদ ও ঝামেলাহীন রান্নার জন্য আজই সংগ্রহ করুন 🔥✨', 3, NULL, NULL, NULL, '1122336', 1, 23, '', 'public/images/products/ss-pan-support.jpeg', '{}', 150.00, 120.00, 1.00, 50.00, 50.00, 0.00, 'flat', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:12:27.615654+00', '2026-03-15 07:56:08.729424+00', 110.00, 1.00);
INSERT INTO public.products VALUES (33, 'লবঙ্গ', '', 2, NULL, NULL, NULL, '112267', 1, 2, '', 'public/images/products/thumb112267.jpg', '{public/images/products/-gallery-1-gncev.webp,public/images/products/-gallery-2-hplj7.webp}', 110.00, 100.00, 1.00, 100.00, 0.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-07 05:28:51.606231+00', '2026-02-28 02:36:54.342663+00', 0.00, 1.00);
INSERT INTO public.products VALUES (152, 'Orignal OSAKA Pvc Tape -Pack of 12 (White, Black, Red)', '', 2, NULL, NULL, NULL, '1122563', 1, 17, '', 'public/images/products/thumb1122563.jpeg', NULL, 20.00, 0.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 08:29:42.229423+00', '2026-03-15 08:41:21.636796+00', 13.33, 1.00);
INSERT INTO public.products VALUES (121, 'Magnet Premium Quality Strong Super Glue 1pcs', 'Magnet Premium Quality Strong Super Glue 1pcs



Magnet Premium Quality Super Glue is a powerful adhesive that bonds materials in seconds. This glue accommodates a wide array of surfaces, from metals and wood to plastics. You can fix items in seconds with its fast-drying formula. It is also Heat- and water-resistant to provide strong and long-lasting repairs. Whether one is mending a broken item or working with crafts, this is the perfect glue for quick fixes and DIY projects, as it provides a really reliable bond.



Features and Benefits



Easy to use: it is easy to apply and dries out fast. Just a couple of drops press the surfaces together, and in no time, the glue sets, holding firm. Without extra tools or effort, one can easily use it. This makes it quite handy for anyone.



Sound Quality: Magnet Super Glue is a high-quality product. It forms a strong bond with different materials, including wood, metal, and plastic. You don''t have to be concerned about the bond weakening after some time. This strong glue makes your repairs last for quite an extended period.



Quick-drying: You don''t have to wait for ages for it to set. You apply the glue and press the pieces together. In just a couple of seconds, it would bond the surfaces. Therefore, this is fantastic when quick fixes are needed and also saves time when working on projects.



Strong Bond: The adhesive has a strong bond that is resistant and strong. It is suitable for gluing up broken objects or doing DIY projects, whether a chair or some craft; your things will be kept together very firmly with this super glue to keep repairs easy and quick.



Water and Heat Resistant: This glue is resistant to water and Heat. It ensures that the bond remains intact even under challenging conditions. Whether it''s exposed to Heat or moisture, the strong bond stays strong, offering long-lasting reliability for both indoor and outdoor use.



Magnet Premium Quality Strong Super Glue Price in Bangladesh

The Magnet Premium Quality Strong Super Glue is available at the best price in Bangladesh on your own Rokomari.com. Grab it now from us!

 
', 2, NULL, NULL, NULL, '1122552', 1, 17, '', 'public/images/products/thumb1122552.jpg', NULL, 20.00, 0.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-11 08:11:44.730033+00', '2026-03-15 08:45:16.406264+00', 11.00, 1.00);
INSERT INTO public.products VALUES (89, 'Dala Palm-ছড়া খেজুর 1kg', 'খেজুর শরীরের রোগ প্রতিরোধ ক্ষমতা বৃদ্ধির পাশাপাশি, ওজন নিয়ন্ত্রণ, হাড়ের গঠনে খেজুর নিরলস ভুমিকা পালন করে। খেজুরে অন্যান্য ফলের তুলনায় পটাশিয়াম, ম্যাংগানিজ, ম্যাগনেশিয়াম, আয়রন, কপারের পরিমাণবেশি থাকে। যার কারণে রমজান মাস ছাড়াও অনেকে খাদ্যতালিকায় খেজুর সবসময় রাখে। খেজুরে থাকা উচ্চ মাত্রার ভিটামিন ‘বি’ নার্ভকে শান্ত করে রক্তচাপ কমাতে সহায়তা করে।', 2, NULL, NULL, NULL, '1122572', 1, 1, '', 'public/images/products/dala-palm-chhara-khajur.jpg', NULL, 600.00, 650.00, 1.00, 50.00, 50.00, 0.00, '', 50.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-23 05:34:32.193953+00', '2026-03-02 07:28:50.361859+00', 0.00, 1.00);
INSERT INTO public.products VALUES (96, 'নাগাল খেজুর Nagal Dates | 1kg | সৌদি আরব', 'নাগাল খেজুর Nagal Dates | 500gm,1kg,3kg, 5kg বক্স | সৌদি আরব:
নাগাল খেজুর (Nagal Dates) সৌদি আরবের বিখ্যাত খেজুর। প্রাকৃতিকভাবে মিষ্টি স্বাদ ও স্বাস্থ্যকর পুষ্টি উপাদান সমৃদ্ধ এই খেজুর আপনার প্রতিদিনের খাবারকে আরও স্বাদযুক্ত ও পুষ্টিকর করে তোলে।

ব্যবহার:
✅ স্বাস্থ্যকর স্ন্যাকস হিসেবে
✅ পোলাও, বিরিয়ানি, ডেজার্ট ও বেকড পণ্য তৈরিতে
✅ সালাদ, মিল্কশেক ও স্মুদি-তে
✅ উপহার হিসেবে বা উৎসবের আয়োজনেও ব্যবহার করা যায়

স্বাস্থ্য উপকারিতা:
✅ প্রাকৃতিক অ্যান্টিঅক্সিডেন্ট সমৃদ্ধ
✅ হজম প্রক্রিয়ায় সহায়ক
✅ রোগ প্রতিরোধ ক্ষমতা বাড়ায়
✅ রক্তস্বল্পতা রোধে সহায়ক
✅ শক্তি বৃদ্ধি ও ক্লান্তি দূর করে

সৌদি আরব থেকে আমদানি করা এই নাগাল খেজুর বক্সে সহজলভ্য – 3kg ও 5kg বক্সে। আজই আপনার রান্নাঘর ও ডায়েটের অংশ বানান এই প্রাকৃতিক সুস্বাদু উপহারটি!

পাকা খেজুরের মধ্যে নাগাল খুবই সুস্বাদু ও বহুল জনপ্রিয়।

খেজুরের উপকারিতা–
১. খেজুরে রয়েছে প্রচুর ভিটামিন, খনিজ, ক্যালসিয়াম ও পটাশিয়া। খেজুরে থাকা অ্যান্টিঅক্সিডেন্ট শরীরে রোগ প্রতিরোধের ক্ষমতা বাড়ায়।
২. ফাইবারও মিলবে খেজুরে। তাই এই ফল ডায়েটে রাখতে পারেন নিশ্চিন্তে।
৩. প্রতিটি খেজুরে রয়েছে ২০ থেকে ২৫ মিলিগ্রাম ম্যাগনেসিয়াম, যা উচ্চ রক্তচাপ কমাতে সাহায্য করে।
৪. রক্তস্বল্পতায় ভোগা রোগীরা প্রতিদিন খেজুর খেতে পারেন। একজন সুস্থ মানুষের শরীরে যতটুকু আয়রন প্রয়োজন, তার প্রায় ১১ ভাগ পূরণ করে খেজুর।
৫. যারা চিনি খান না তারা খেজুর খেতে পারেন। চিনির বিকল্প খেজুরের রস ও গুড়।
৬. কোষ্ঠকাঠিন্যের সমস্যায় রাতে পানিতে খেজুর ভিজিয়ে রাখুন। পর দিন সকালে খেজুর ভেজানো পানি পান করুন। দূর হবে কোষ্ঠকাঠিন্য।
৭. খেজুরে থাকা নানা খনিজ হৃদস্পন্দনের হার ঠিক রাখতে সাহায্য করে।

কেন কিনবেন Noor Super Mart থেকে?
Noor Super Mart থেকে আপনি পাচ্ছেন সরাসরি আমদানিকৃত,আমরা আপনাকে দিচ্ছি সেরা দামে অরিজিনাল পণ্য, সেফ প্যাকেজিং ও দেশব্যাপী দ্রুত ডেলিভারি।

অর্ডার এবং সরবরাহ:
আমাদের ওয়েবসাইট থেকে সহজেই Shopping & Retail অর্ডার করতে পারবেন। নিরাপদ প্যাকেজিং এবং দ্রুত ডেলিভারির নিশ্চয়তা।', 2, NULL, NULL, NULL, '1122573', 1, 1, '', 'public/images/products/--nagal-dates--1kg-500gm---.jpg', NULL, 350.00, 400.00, 1.00, 100.00, 100.00, 0.00, '', 50.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-23 06:04:58.183074+00', '2026-03-02 07:29:43.767787+00', 0.00, 1.00);
INSERT INTO public.products VALUES (98, 'Kaveri Mehedi Indian-6pcs Pack', '🌿 কাবেরি মেহেদী (Kaveri Mehedi)
প্রাকৃতিক উপাদানে তৈরি কাবেরি মেহেদী আপনার হাতে ও চুলে দেবে গাঢ়, দীর্ঘস্থায়ী ও আকর্ষণীয় রঙ। সহজে ব্যবহারযোগ্য এই মেহেদী দ্রুত শুকায় এবং ত্বকে জ্বালাপোড়া করে না।
✨ বিশেষ বৈশিষ্ট্য:
প্রাকৃতিক উপাদানে প্রস্তুত
গাঢ় ও লং-লাস্টিং রঙ
ত্বকের জন্য নিরাপদ
সহজে প্রয়োগযোগ্য
হাত ও চুল—দুই কাজেই ব্যবহারযোগ্য
🎨 ব্যবহার:
হাতের ডিজাইন, ঈদ/বিয়ে/যেকোনো উৎসব কিংবা নিয়মিত হেয়ার কেয়ার—সব ক্ষেত্রেই উপযুক্ত।
📦 উপলব্ধ সাইজ:
২৫g / ৫০g / ১০০g (স্টক অনুযায়ী)
আপনার সৌন্দর্যে প্রাকৃতিক ছোঁয়া আনতে আজই অর্ডার করুন 🌿✨', 2, NULL, NULL, NULL, '1122574', 1, 11, '', 'public/images/products/kaveri-mehedi-indian-6pcs-pack.jpg', NULL, 200.00, 200.00, 0.05, 100.00, 50.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-23 07:12:38.965391+00', '2026-02-23 07:53:50.370857+00', 0.00, 1.00);
INSERT INTO public.products VALUES (148, 'CHANNEL SINGLE', '', 3, NULL, NULL, 1, 'PRD-26-3-148', 1, 13, '', 'public/images/products/thumbprd-26-3-148.jpg', NULL, 80.00, 30.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 08:12:49.076112+00', '2026-03-15 08:12:49.0839+00', 35.00, 1.00);
INSERT INTO public.products VALUES (149, 'Mondete Rainbow Pencils', '', 2, NULL, NULL, NULL, '1122562', 1, NULL, '', 'public/images/products/thumb1122562.jpeg', NULL, 120.00, 120.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 08:12:58.457584+00', '2026-03-15 08:12:58.462858+00', 70.00, 1.00);
INSERT INTO public.products VALUES (151, '22 MM MINI', '', 3, NULL, NULL, NULL, '12369', 1, 17, '', 'public/images/products/thumb12369.jpg', NULL, 300.00, 180.00, 1.00, 50.00, 50.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 08:19:05.454001+00', '2026-03-15 08:19:05.459361+00', 200.00, 1.00);
INSERT INTO public.products VALUES (72, 'LG SS Double Burner Gas Stove', '🔥 LG SS Double Burner Gas Stove
আপনার রান্নাঘরের জন্য শক্তপোক্ত, স্টাইলিশ এবং দীর্ঘস্থায়ী একটি গ্যাস চুলা। উচ্চমানের স্টেইনলেস স্টিল বডি দিয়ে তৈরি এই ডাবল বার্নার গ্যাস স্টোভ একসাথে দুইটি আইটেম দ্রুত ও সমান তাপে রান্না করতে সাহায্য করে।
✨ বিশেষ বৈশিষ্ট্য:
প্রিমিয়াম স্টেইনলেস স্টিল (SS) বডি
২টি শক্তিশালী বার্নার
দ্রুত ও সমান তাপ বিতরণ
কম গ্যাস খরচ
সহজে পরিষ্কার করা যায়
টেকসই ও মরিচা প্রতিরোধী
🍳 উপযোগিতা:
দৈনন্দিন রান্না, বড় পরিবার বা ছোট ব্যবসার রান্নাঘরের জন্য উপযুক্ত।
📦 প্যাকেজে যা থাকবে:
১টি ডাবল বার্নার গ্যাস স্টোভ
আপনার রান্নাকে সহজ ও দ্রুত করতে আজই অর্ডার করুন 🔥✨', 3, NULL, NULL, NULL, '112235', 1, 23, '', 'public/images/products/lg-ss-double-burner-gas-stove.jpeg', '{}', 2500.00, 2250.00, 1.00, 50.00, 50.00, 0.00, 'flat', 250.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:09:34.906361+00', '2026-02-23 08:01:30.644503+00', 0.00, 1.00);
INSERT INTO public.products VALUES (70, 'মেথি/Methi', '🌿 Methi (মেথি) – স্বাস্থ্যকর ও প্রাকৃতিক উপকারিতা
প্রাকৃতিক Methi (মেথি) শুধু রান্নার স্বাদ বাড়ায় না, বরং শরীরের জন্য অনেক উপকারও দেয়।
✨ উপকারিতা:
রক্তে শর্করার মাত্রা নিয়ন্ত্রণে সাহায্য করে
হজম শক্তি বৃদ্ধি করে, পাচন সুগম করে
চুল ও ত্বকের জন্য প্রাকৃতিক পুষ্টি যোগায়
ওজন নিয়ন্ত্রণে সহায়তা করে
হাড় ও হাড়ের ঘনত্ব বজায় রাখতে সাহায্য করে
💡 কেন ব্যবহার করবেন:
প্রতিদিনের খাদ্যাভ্যাসে মেথি যুক্ত করলে শরীর সুস্থ, শক্তিশালী এবং রোগ প্রতিরোধ ক্ষমতা বাড়ে।
📦 প্যাক সাইজ:
২৫g / ৫০g / ১০০g
আপনার স্বাস্থ্যকে শক্তিশালী করতে আজই অর্ডার করুন 🌿✨', 2, NULL, NULL, NULL, '112304', 1, 5, '', 'public/images/products/methi.jpg', '{}', 150.00, 150.00, 1.00, 100.00, 100.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-10 06:34:34.344978+00', '2026-02-23 08:11:16.450409+00', 0.00, 1.00);
INSERT INTO public.products VALUES (13, 'আলু/Potato', '🥔 তাজা আলু
প্রতিদিনের রান্নার সবচেয়ে প্রয়োজনীয় উপাদান — তাজা, পরিষ্কার আর ভালো মানের আলু।
ভাজি, ভর্তা, তরকারি বা ফ্রেঞ্চ ফ্রাই — সবকিছুতেই পারফেক্ট!
✔️ টাটকা ও মানসম্মত
✔️ দীর্ঘদিন ভালো থাকে
✔️ সব ধরনের রান্নায় উপযোগী
ঘরের রান্না হোক আরও সুস্বাদু ও পরিপূর্ণ।
🛒 আজই অর্ডার করুন Noor Super Mart থেকে।
সাশ্রয়ী দামে নিত্যপ্রয়োজনীয় সব পণ্য, দ্রুত ও নির্ভরযোগ্য সার্ভিসে আপনার দোরগোড়ায় 🚚✨', 2, NULL, NULL, NULL, '112244', 1, 1, '', 'public/images/products/potato.jpg', '{}', 25.00, 25.00, 1.00, 1000.00, 100.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:20:08.384836+00', '2026-02-28 07:47:59.705828+00', 0.00, 1.00);
INSERT INTO public.products VALUES (9, 'ছোলা/chikpeas', 'ছোলা
পুষ্টিকর, সুস্বাদু আর ঘরের নিত্যপ্রয়োজনীয় একটি উপাদান — ভালো মানের ছোলা।
ভুনা ছোলা, চানা চাট, হালিম বা ইফতারের নানা আইটেম—সবকিছুর জন্য পারফেক্ট।
✔️ প্রোটিন সমৃদ্ধ
✔️ পরিষ্কার ও মানসম্মত
✔️ ভিজিয়ে রান্না করলে নরম ও সুস্বাদু
পরিবারের জন্য স্বাস্থ্যকর ও মজাদার আয়োজন চাই?
তাহলে রাখুন ভালো মানের ছোলা ঘরেই।
🛒 আজই অর্ডার করুন Noor Super Mart থেকে।
নিশ্চিত মান, সাশ্রয়ী দাম এবং দ্রুত ডেলিভারি—সব একসাথে 🚚✨', 2, NULL, NULL, NULL, '112240', 1, 1, '', 'public/images/products/chikpeas.jpg', '{}', 120.00, 110.00, 1.00, 100.00, 10.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', true, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:01:37.823959+00', '2026-02-28 07:56:14.246687+00', 0.00, 1.00);
INSERT INTO public.products VALUES (68, 'জাইফল/Nutmeg', '🌿 জাইফল (Nutmeg) – প্রাকৃতিক স্বাস্থ্য সুবিধা
জাইফল শুধু রান্নার স্বাদ বাড়ায় না, এটি প্রাকৃতিকভাবে শরীরের জন্য নানা উপকারও দিতে পারে। গবেষণা দেখা গেছে জাইফলে কিছু শক্তিশালী উপাদান আছে যা শরীরকে স্বাস্থ্য সমর্থন করতে সাহায্য করে। �
WebMD +1
✨ প্রধান উপকারিতা:
• অ্যান্টিঅক্সিডেন্ট শক্তি: জাইফল অনেক অ্যান্টিঅক্সিডেন্ট প্রদান করে, যা মুক্ত কণার ক্ষতির বিরুদ্ধে প্রতিরোধে সাহায্য করে। �
• প্রদাহ কমায়: এর কিছু উপাদান শরীরের প্রদাহ কমাতে সাহায্য করতে পারে। �
• হজম সমর্থন: এটি পাচনের সহায়ক হিসেবে কাজ করে, পাচনতন্ত্রের স্বস্তি ও গ্যাসের সমস্যা হ্রাস করতে সাহায্য করতে পারে। �
• মাইক্রোবিয়াল প্রভাব: কিছু গবেষণায় দেখা গেছে জাইফলের অ্যান্টিব্যাকটেরিয়াল প্রভাব আছে, যা মুখের ব্যাকটেরিয়া ও দুর্গন্ধ কমাতে সাহায্য করতে পারে। �
• আরামদায়ক ঘুম: প্রচলিতভাবে ধূমপান বা অতিরিক্ত চিন্তা কমাতে ও ঘুম উন্নত করতে ব্যবহার করা হয়। �
Healthline
Healthline
MedicineNet
WebMD +1
India TV News
💡 দ্রষ্টব্য:
জাইফল ছোট পরিমাণে ব্যবহারে নিরাপদ, কিন্তু একেবারে বেশি মাত্রায় খাওয়া উচিত নয়, কারণ তা বিষক্রিয়া বা অপ্রত্যাশিত প্রভাব সৃষ্টি করতে পারে। �
WebMD
📦 উপলব্ধ প্যাক সাইজ:
২৫g / ৫০g / ১০০g
🌱 আজই সংগ্রহ করুন এবং প্রাকৃতিকভাবে আপনার স্বাস্থ্যকে শক্তিশালী করুন!', 2, NULL, NULL, NULL, '112302', 1, 4, '', 'public/images/products/nutmeg.jpg', '{}', 50.00, 50.00, 1.00, 100.00, 100.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-10 06:20:18.492204+00', '2026-02-23 08:16:52.443486+00', 0.00, 1.00);
INSERT INTO public.products VALUES (66, 'বুটের ডাল (Chana Dal / Split Bengal Gram)', ' 🌿 বুটের ডাল (Chana Dal / Split Bengal Gram) – পুষ্টিকর ও স্বাস্থ্যসম্মত
বুটের ডাল প্রোটিন, ফাইবার ও গুরুত্বপূর্ণ খনিজে সমৃদ্ধ, যা শরীরের শক্তি, হজম শক্তি এবং সার্বিক স্বাস্থ্য বজায় রাখতে সহায়তা করে।
✨ উপকারিতা:
প্রোটিন ও পুষ্টিতে সমৃদ্ধ, শরীরের শক্তি বৃদ্ধি করে
হজম সহজ করে, পাচনতন্ত্র সুগম রাখে
রক্তচাপ ও রক্তে শর্করার মাত্রা নিয়ন্ত্রণে সাহায্য করে
হাড় ও হাড়ের ঘনত্ব বজায় রাখতে সহায়তা করে
দীর্ঘস্থায়ী শক্তি এবং রোগ প্রতিরোধ ক্ষমতা বাড়ায়
💡 কেন ব্যবহার করবেন:
প্রতিদিনের খাদ্যাভ্যাসে বুটের ডাল যুক্ত করলে শরীর সুস্থ, শক্তিশালী এবং পুষ্টিকর হয়।
📦 উপলব্ধ প্যাক সাইজ:
২৫০g / ৫০০g / ১kg
আজই অর্ডার করুন এবং স্বাস্থ্যকর পুষ্টি উপভোগ করুন 🌿✨ – পুষ্টিকর ও স্বাস্থ্যসম্মত
বুটের ডাল প্রোটিন, ফাইবার ও গুরুত্বপূর্ণ খনিজে সমৃদ্ধ, যা শরীরের শক্তি, হজম শক্তি এবং সার্বিক স্বাস্থ্য বজায় রাখতে সহায়তা করে।
✨ উপকারিতা:
প্রোটিন ও পুষ্টিতে সমৃদ্ধ, শরীরের শক্তি বৃদ্ধি করে
হজম সহজ করে, পাচনতন্ত্র সুগম রাখে
রক্তচাপ ও রক্তে শর্করার মাত্রা নিয়ন্ত্রণে সাহায্য করে
হাড় ও হাড়ের ঘনত্ব বজায় রাখতে সহায়তা করে
দীর্ঘস্থায়ী শক্তি এবং রোগ প্রতিরোধ ক্ষমতা বাড়ায়
💡 কেন ব্যবহার করবেন:
প্রতিদিনের খাদ্যাভ্যাসে বুটের ডাল যুক্ত করলে শরীর সুস্থ, শক্তিশালী এবং পুষ্টিকর হয়।
📦 উপলব্ধ প্যাক সাইজ:
২৫০g / ৫০০g / ১kg
আজই অর্ডার করুন এবং স্বাস্থ্যকর পুষ্টি উপভোগ করুন 🌿✨', 2, NULL, NULL, NULL, '112298-66', 1, 1, '', 'public/images/products/--chana-dal--split-bengal-gram.jpg', '{}', 70.00, 70.00, 1.00, 50.00, 100.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-10 05:36:38.48077+00', '2026-02-23 08:22:45.093946+00', 0.00, 1.00);
INSERT INTO public.products VALUES (52, 'পেয়ারা/Guava', '🍐 পেয়ারা – Noor Super Mart
সতেজ, রসালো এবং পুষ্টিকর পেয়ারা স্বাস্থ্যের জন্য এক অসাধারণ উপহার। এটি প্রতিদিনের খাদ্যাভ্যাসে যুক্ত করলে শরীর সুস্থ, রোগ প্রতিরোধ ক্ষমতা বৃদ্ধি পায় এবং হজম শক্তি উন্নত হয়।
✨ স্বাস্থ্য ও উপকারিতা:
১০০% সতেজ ও প্রাকৃতিক, সরাসরি আপনার কাছে পৌঁছে যায়
ভিটামিন C, ফাইবার ও অ্যান্টিঅক্সিডেন্ট সমৃদ্ধ
হজম শক্তি বৃদ্ধি করে এবং পাকস্থলীর স্বাস্থ্য রক্ষা করে
রোগ প্রতিরোধ ক্ষমতা বাড়াতে সহায়ক
হৃদয় ও ত্বকের স্বাস্থ্যের জন্য উপকারী
ওজন নিয়ন্ত্রণ ও দীর্ঘস্থায়ী শক্তি বৃদ্ধিতে সহায়ক
💡 কেন Noor Super Mart থেকে কিনবেন:
সর্বোচ্চ মানের সতেজ পেয়ারা, খেতে সুস্বাদু ও রসালো
নিরাপদ অর্ডার ও দ্রুত বাড়িতে ডেলিভারি
মূল্য সাশ্রয়ী এবং নিয়মিত স্টক
📦 উপলব্ধ প্যাক সাইজ:
১ কেজি / প্যাক অনুযায়ী
🌱 আজই অর্ডার করুন এবং স্বাস্থ্যকর, সতেজ ও পুষ্টিকর পেয়ারা উপভোগ করুন Noor Super Mart থেকে! 🍐✨', 2, NULL, NULL, NULL, '112295', 1, 1, '', 'public/images/products/guava.jpg', '{}', 40.00, 30.00, 3.00, 100.00, 100.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-09 09:00:58.058993+00', '2026-02-23 08:53:52.021482+00', 0.00, 1.00);
INSERT INTO public.products VALUES (129, 'parasuit oil 200ml', '', 2, NULL, NULL, 1, 'PRD-26-2-129', 1, 4, '', 'public/images/products/thumbprd-26-2-129.jpg', NULL, 200.00, 200.00, 1.00, 10.00, 10.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-12 04:07:47.697844+00', '2026-03-12 04:07:47.713313+00', 195.00, 1.00);
INSERT INTO public.products VALUES (139, 'মিনিকেট চাল স্ট্যান্ডার্ড (৫ কেজি) |', 'মিনিকেট চাল
উদ্যোক্তাদের নিজস্ব হাতে তৈরী ভেঁজালমুক্ত
সঠিক ওজন
পরিমান- ৫ কেজি', 2, NULL, NULL, NULL, '1122561', 1, NULL, '', 'public/images/products/thumb1122561.jpg', NULL, 350.00, 350.00, 1.00, 50.00, 50.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 07:41:31.843146+00', '2026-03-15 07:41:31.85092+00', 350.00, 1.00);
INSERT INTO public.products VALUES (140, 'SINGLE IGNOITION SS CHULA', '', 3, NULL, NULL, 1, '12', 1, 17, '', 'public/images/products/thumb12.jpg', NULL, 110.00, 105.00, 1.00, 50.00, 50.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 07:46:25.527589+00', '2026-03-15 07:46:25.535438+00', 105.00, 1.00);
INSERT INTO public.products VALUES (71, '22mm Regulator', '🔹 22mm Gas Regulator
আপনার রান্নাঘরের নিরাপত্তা ও গ্যাসের স্থিতিশীল প্রবাহ নিশ্চিত করতে 22mm Gas Regulator–টি একটি নির্ভরযোগ্য সমাধান। দীর্ঘস্থায়ী ও মানসম্মত উপাদানে তৈরি, যা নিরাপদ ও ঝামেলাহীন রান্নার জন্য উপযুক্ত।
✨ বিশেষ বৈশিষ্ট্য:
২২মিমি স্ট্যান্ডার্ড ফিট
টেকসই ও দীর্ঘস্থায়ী
নিরাপদ গ্যাস প্রবাহ নিয়ন্ত্রণ
সহজে সংযুক্ত করা যায়
💡 কেন কিনবেন:
নিরাপদ, নির্ভরযোগ্য এবং দীর্ঘস্থায়ী গ্যাস কন্ট্রোল—আপনার রান্নার অভিজ্ঞতাকে করে ঝামেলাহীন ও নিরাপদ।
📦 আজই অর্ডার করুন এবং রান্না করুন নিশ্চিন্তে 🔥', 3, NULL, NULL, NULL, '112305', 1, 23, '', 'public/images/products/22mm-regulator.jpg', '{}', 350.00, 300.00, 1.00, 1000.00, 1000.00, 0.00, 'flat', 50.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-14 23:29:54.622226+00', '2026-03-15 08:17:06.213334+00', 210.00, 1.00);
INSERT INTO public.products VALUES (115, 'Kaveri No.1 Indian Cone Mehedi (12PCS PACK)', 'Experience the rich tradition and vibrant beauty of Indian mehndi art with the Kaveri No.1 Indian Cone Mehedi. This convenient 12-piece pack is perfect for personal use, professional artists, or for sharing the joy of intricate body art with friends and family.', 2, NULL, NULL, NULL, '1122551', 1, NULL, '', 'public/images/products/thumb1122551.jpg', NULL, 360.00, 360.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-07 04:34:52.528876+00', '2026-03-07 04:34:52.534052+00', 250.00, 3.00);
INSERT INTO public.products VALUES (51, 'ড্রাগন ফল', 'ড্রাগন ফল – Noor Super Mart
সতেজ, রসালো এবং রঙিন ড্রাগন ফল শুধু চোখে নয়, স্বাস্থ্যের জন্যও উপকারী। এটি প্রতিদিনের খাদ্যাভ্যাসে যুক্ত করলে শরীরের রোগ প্রতিরোধ ক্ষমতা বৃদ্ধি পায়, হজম শক্তি উন্নত হয় এবং শরীরকে দেয় প্রাকৃতিক শক্তি।
✨ স্বাস্থ্য ও উপকারিতা:
১০০% সতেজ ও প্রাকৃতিক, সরাসরি আপনার কাছে পৌঁছে যায়
ভিটামিন C, ফাইবার এবং অ্যান্টিঅক্সিডেন্টে সমৃদ্ধ, যা শরীরকে সুস্থ রাখে
হজম শক্তি বৃদ্ধি করে এবং পাকস্থলীর স্বাস্থ্য বজায় রাখে
রোগ প্রতিরোধ ক্ষমতা বাড়ায় এবং শরীরকে শক্তিশালী রাখে
হৃদয় ও ত্বকের স্বাস্থ্যের জন্য উপকারী
প্রাকৃতিকভাবে ওজন নিয়ন্ত্রণে সহায়ক
💡 কেন Noor Super Mart থেকে কিনবেন:
সর্বোচ্চ মানের সতেজ ড্রাগন ফল, রঙিন, রসালো ও সুস্বাদু
নিরাপদ অর্ডার ও দ্রুত বাড়িতে ডেলিভারি
মূল্য সাশ্রয়ী এবং নিয়মিত স্টক
📦 উপলব্ধ প্যাক সাইজ:
 ১ কেজি  অনুযায়ী
🌱 আজই অর্ডার করুন এবং স্বাস্থ্যকর, সতেজ ও পুষ্টিকর ড্রাগন ফল উপভোগ করুন Noor Super Mart থেকে! 🐉✨', 4, NULL, NULL, NULL, '112288', 1, 1, '', 'public/images/products/thumb112288.jpg', '{}', 450.00, 420.00, 1.00, 100.00, 100.00, 0.00, 'flat', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 09:21:21.76798+00', '2026-02-28 07:03:08.814073+00', 0.00, 1.00);
INSERT INTO public.products VALUES (56, 'আনারস/ Pineapple', '🍍 আনারস – Noor Super Mart
সতেজ, রসালো এবং ভিটামিন C সমৃদ্ধ আনারস স্বাদে মন মাতানো ও শরীরের জন্য উপকারী। প্রতিদিনের খাদ্যাভ্যাসে আনারস খেলে হজম শক্তি বৃদ্ধি, রোগ প্রতিরোধ ক্ষমতা এবং শরীরের শক্তি বাড়ে।
✨ সুবিধা ও উপকারিতা:
১০০% সতেজ ও রসালো
ভিটামিন C, অ্যান্টিঅক্সিডেন্ট সমৃদ্ধ
হজম সহজ করে ও পাচন শক্তি বৃদ্ধি করে
রোগ প্রতিরোধ ক্ষমতা ও সার্বিক স্বাস্থ্য উন্নত করে
সহজ অর্ডার, নিরাপদ পেমেন্ট ও বাড়িতে ডেলিভারি
📦 উপলব্ধ প্যাক সাইজ:
পিস / কেজি অনুযায়ী
🌱 আজই অর্ডার করুন এবং সতেজ, স্বাস্থ্যকর ও পুষ্টিকর আনারস উপভোগ করুন Noor Super Mart থেকে! 🍍✨', 2, NULL, NULL, NULL, '112297-56', 1, 13, '', 'public/images/products/-pineapple.jpg', '{}', 40.00, 30.00, 3.00, 100.00, 100.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-09 09:29:29.706798+00', '2026-02-23 08:48:51.386242+00', 0.00, 1.00);
INSERT INTO public.products VALUES (41, 'নেহা মেহিদি', '', 2, NULL, NULL, NULL, '112278', 1, 23, '', 'public/images/products/thumb112278.webp', '{}', 30.00, 30.00, 10.00, 100.00, 0.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 04:05:05.530129+00', '2026-02-27 01:09:16.942579+00', 0.00, 1.00);
INSERT INTO public.products VALUES (99, 'Fair & Lovely Advance Multi-Vitamin Fairness Expert Cream - 100g (U.A.E)', 'Fair and Lovely Advanced Multi Vitamin Cream’- is ideal for all skin types. Improved with Vitamin B3, Vitamin C, Vitamin E, Vitamin B6, the cream gives you moment reasonableness and shine and is protected to utilize. Its high-level sun assurance recipe safeguards your skin from hurtful UVA and UVB radiations. With standard use, dull fixes and spots on the skin begin getting lighter. Ingredients: Water, palmatic corrosive and stearic corrosive, niancinamide, glycerin, dimethicone, ethylhexyl methoxycinnamate,, butyl methoxydibenzoylmethane, titanium dioxide, sodium ascorbyl phosphate, isopropyl myristate, tocopheryl acetic acid derivation, allantoin, pyridoxine hydrocholride, cetyl liquor, aluminum hydroxide, phenoxyethanol, metthylparaben, propylparaben, potassium hydroxide, disodium EDTA, scent, alpha-isomethyl ionone, benzyl salicylate, butylphenyl methylpropional, citronellol, geraniol, hexyl cinnamal, coumarin, cinnamyl liquor, limonene, linalool.
Information:
**Brand Fair & Lovely
Made In  U.A.E.
Size 100gm**
Benefits: This Fair and Lovey isn''t only a cream, it''s your day by day master reasonableness. The Fair and Lovely multivitamin cream focuses on similar reasonableness issues as 5 Expert Fairness Solutions: 1. The Laser-live exactness chips away at skin marks; 2. Face cleaning that decreased sun tan; 3. Face strips that objectives bluntness; 4. cell reinforcements eases up dark circles; 5. Nutrient cover that eases up skin tone. This item is dermatologist tried.
How to Apply: Simply dab the cream over the face and next and delicately rub in. Utilize twice day by day on purified face for best outcomes. You may attempt to see that you won''t ever stress over glancing dull in the day, again. This item doesn''t contain any hurtful fixings or fade.
Highlights:
 Works like laser light to give best decency.
Dives deep inside to light up brown complexion cells.
Gives sparkling brilliant reasonableness.
 Enters the open pores and eliminates overabundance oil.

 Bundling may shift.
 Get sun ensured reasonableness.', 2, NULL, NULL, NULL, '1122580', 1, 2, '', 'public/images/products/fair--lovely-advance-multi-vitamin-fairness-expert-cream---100g-u.a.webp', NULL, 710.00, 700.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-24 04:14:37.370534+00', '2026-02-27 09:44:27.593011+00', 0.00, 1.00);
INSERT INTO public.products VALUES (63, 'গরুর মাংস', '🥩 গরুর মাংস – Noor Super Mart
শক্তিশালী প্রোটিনে ভরপুর, সতেজ ও স্বাস্থ্যকর গরুর মাংস পেশী গঠন, শক্তি বৃদ্ধি এবং সার্বিক স্বাস্থ্য সমর্থনে সাহায্য করে।
✨ সুবিধা ও উপকারিতা:
১০০% সতেজ ও হাইজিনিক
প্রোটিন, ভিটামিন B ও গুরুত্বপূর্ণ খনিজ সমৃদ্ধ
হজম সহজ ও শক্তি বৃদ্ধি
দ্রুত ডেলিভারি ও নিরাপদ পেমেন্ট
প্যাক/কেজি অনুযায়ী সহজ ক্রয়
📦 আজই অর্ডার করুন এবং প্রাকৃতিকভাবে স্বাস্থ্য ও শক্তি বৃদ্ধি করুন Noor Super Mart থেকে! 🥩✨', 5, NULL, NULL, NULL, '112298-63', 1, 1, '', 'public/images/products/thumb112298-63.jpg', '{}', 800.00, 750.00, 1.00, 100.00, 100.00, 0.00, 'flat', 50.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-09 10:29:17.217523+00', '2026-02-28 07:02:17.679316+00', 0.00, 1.00);
INSERT INTO public.products VALUES (150, 'NOGEL', '', 3, NULL, NULL, 1, '2563', 1, 17, '', 'public/images/products/thumb2563.jpg', NULL, 50.00, 14.00, 1.00, 500.00, 500.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 08:15:43.564527+00', '2026-03-15 08:15:43.575029+00', 20.00, 1.00);
INSERT INTO public.products VALUES (132, 'তীর আটা ২ কেজি', '', 2, NULL, NULL, NULL, 'PRD-26-2-132', 1, NULL, '', 'public/images/products/thumbprd-26-2-132.webp', NULL, 130.00, 130.00, 1.00, 10.00, 10.00, 0.00, 'flat', 5.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-12 04:09:31.633993+00', '2026-03-15 07:06:06.3978+00', 60.00, 10.00);
INSERT INTO public.products VALUES (141, 'DOUBLE IGNITION SS CHULA', '', 3, NULL, NULL, 1, '1232', 1, 13, '', 'public/images/products/thumb1232.jpg', NULL, 120.00, 120.00, 1.00, 120.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 07:47:31.82888+00', '2026-03-15 07:47:31.834438+00', 120.00, 1.00);
INSERT INTO public.products VALUES (65, 'মসুরির ডাল (Masoor Dal / Red Lentil)', '🌿 মসুরির ডাল (Masoor Dal / Red Lentil) – পুষ্টিকর ও স্বাস্থ্যসম্মত
মসুরির ডাল প্রোটিন, ফাইবার এবং গুরুত্বপূর্ণ খনিজে সমৃদ্ধ, যা শরীরকে শক্তি দেয়, হজম শক্তি বৃদ্ধি করে এবং সার্বিক স্বাস্থ্য বজায় রাখে।
✨ উপকারিতা:
প্রোটিন ও পুষ্টিতে সমৃদ্ধ, শরীরের শক্তি বৃদ্ধি করে
হজম শক্তি বৃদ্ধি করে, পাচন সুগম করে
রক্তচাপ নিয়ন্ত্রণে সাহায্য করে
হাড় ও দাতের স্বাস্থ্য বজায় রাখতে সহায়তা করে
রোগ প্রতিরোধ ক্ষমতা ও সার্বিক সুস্থতা বাড়ায়
💡 কেন ব্যবহার করবেন:
প্রতিদিনের খাদ্যাভ্যাসে মসুরির ডাল যুক্ত করলে শরীর সুস্থ, শক্তিশালী এবং পুষ্টিকর হয়।
📦 উপলব্ধ প্যাক সাইজ:
২৫০g / ৫০০g / ১kg
আজই অর্ডার করুন এবং স্বাস্থ্যকর পুষ্টি উপভোগ করুন 🌿✨', 2, NULL, NULL, NULL, '112297-65', 1, 1, '', 'public/images/products/--masoor-dal--red-lentil.jpg', '{}', 180.00, 180.00, 1.00, 50.00, 50.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-10 05:29:19.651525+00', '2026-02-26 05:27:05.539371+00', 0.00, 1.00);
INSERT INTO public.products VALUES (50, 'আনার', 'ড্রাগন ফল – Noor Super Mart
সতেজ, রসালো এবং রঙিন ড্রাগন ফল শুধু চোখে নয়, স্বাস্থ্যের জন্যও উপকারী। এটি প্রতিদিনের খাদ্যাভ্যাসে যুক্ত করলে শরীরের রোগ প্রতিরোধ ক্ষমতা বৃদ্ধি পায়, হজম শক্তি উন্নত হয় এবং শরীরকে দেয় প্রাকৃতিক শক্তি।
✨ স্বাস্থ্য ও উপকারিতা:
১০০% সতেজ ও প্রাকৃতিক, সরাসরি আপনার কাছে পৌঁছে যায়
ভিটামিন C, ফাইবার এবং অ্যান্টিঅক্সিডেন্টে সমৃদ্ধ, যা শরীরকে সুস্থ রাখে
হজম শক্তি বৃদ্ধি করে এবং পাকস্থলীর স্বাস্থ্য বজায় রাখে
রোগ প্রতিরোধ ক্ষমতা বাড়ায় এবং শরীরকে শক্তিশালী রাখে
হৃদয় ও ত্বকের স্বাস্থ্যের জন্য উপকারী
প্রাকৃতিকভাবে ওজন নিয়ন্ত্রণে সহায়ক
💡 কেন Noor Super Mart থেকে কিনবেন:
সর্বোচ্চ মানের সতেজ ড্রাগন ফল, রঙিন, রসালো ও সুস্বাদু
নিরাপদ অর্ডার ও দ্রুত বাড়িতে ডেলিভারি
মূল্য সাশ্রয়ী এবং নিয়মিত স্টক
📦 উপলব্ধ প্যাক সাইজ:
 ১ কেজি  অনুযায়ী
🌱 আজই অর্ডার করুন এবং স্বাস্থ্যকর, সতেজ ও পুষ্টিকর ড্রাগন ফল উপভোগ করুন Noor Super Mart থেকে! 🐉✨', 4, NULL, NULL, NULL, '112287', 1, 1, '', 'public/images/products/thumb112287.jpg', '{}', 450.00, 430.00, 1.00, 200.00, 200.00, 0.00, 'flat', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-08 09:15:58.931718+00', '2026-02-28 02:33:41.229251+00', 0.00, 1.00);
INSERT INTO public.products VALUES (116, 'Joya Sanitary Napkin Belt System-8’s Pack', 'Joya Belt is the most affordable and best belt system sanitary napkin available in the market. It is the first scented sanitary napkin in Bangladesh. Joya Belt ensures highest quality with maximum absorption. The pad length is 490 mm (unfolded). ', 2, NULL, NULL, NULL, '1122556', 1, NULL, 'sanitary napkin ', 'public/images/products/thumb1122556.png', NULL, 60.00, 60.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-08 07:55:59.919133+00', '2026-03-08 07:55:59.932302+00', 30.00, 15.00);
INSERT INTO public.products VALUES (12, 'আদা/Ginger', '🌿 চিনা Original আদা
ঝাঁঝালো ঘ্রাণ আর তাজা স্বাদের আসল China Original আদা — রান্নার স্বাদ বাড়াতে একদম পারফেক্ট।
মাংস, মাছ, সবজি বা চায়ের সাথে — প্রতিটি রেসিপিতে এনে দেয় আলাদা ফ্লেভার।
✔️ ঘ্রাণে ও স্বাদে সমৃদ্ধ
✔️ রান্নায় বাড়তি ঝাঁজ ও ফ্রেশনেস
✔️ আদা-চায়ের জন্য উপযুক্ত
প্রতিদিনের রান্নায় চাই আসল স্বাদ?
তাহলে বেছে নিন মানসম্মত আদা।
🛒 এখনই অর্ডার করুন Noor Super Mart থেকে।
নিশ্চিত মান, সাশ্রয়ী দাম, আর দ্রুত ডেলিভারি — সব একসাথে! 🚚✨', 2, NULL, NULL, NULL, '112243', 1, 1, '', 'public/images/products/ginger.webp', '{}', 200.00, 165.00, 1.00, 50.00, 50.00, 0.00, 'flat', 35.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:15:46.060248+00', '2026-02-28 07:51:09.355137+00', 0.00, 1.00);
INSERT INTO public.products VALUES (10, 'বেসন/gram flour', 'বেসন (ছোলার গুঁড়া)
মচমচে পকোড়া, বেগুনি, পেঁয়াজু কিংবা ঘরের নানান মজাদার আইটেম—সবকিছুর জন্য চাই ভালো মানের বেসন।
তাজা, পরিষ্কার ও ঝরঝরে বেসনেই আসে আসল স্বাদ আর সুন্দর টেক্সচার।
✔️ ভাজার জন্য পারফেক্ট
✔️ নরম ও মসৃণ ব্যাটার তৈরি হয়
✔️ রমজান, বিকেলের নাস্তা বা বিশেষ আয়োজনের জন্য উপযোগী
রান্নায় আনুন বাড়তি মজা আর খাঁটি স্বাদ।
🛒 আজই নিয়ে নিন আপনার প্রয়োজনীয় পরিমাণ
Noor Super Mart থেকে — সাশ্রয়ী দাম, নিশ্চিত মান ও দ্রুত ডেলিভারি 🚚✨', 2, NULL, NULL, NULL, '112241', 1, 1, '', 'public/images/products/gram-flour.jpg', '{}', 120.00, 110.00, 1.00, 100.00, 20.00, 0.00, 'percentage', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:09:00.130338+00', '2026-02-28 07:54:38.618169+00', 0.00, 1.00);
INSERT INTO public.products VALUES (143, 'SINGLE AUTO IGNITION GD', '', 3, NULL, NULL, NULL, 'PRD-26-3-143', 1, 13, '', 'public/images/products/thumbprd-26-3-143.jpg', NULL, 140.00, 140.00, 1.00, 50.00, 50.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 07:49:50.36941+00', '2026-03-15 07:49:50.376085+00', 140.00, 1.00);
INSERT INTO public.products VALUES (7, 'Fundools ইনস্ট্যান্ট নুডলস', '🍜 Fundools ইনস্ট্যান্ট নুডলস
ব্যস্ত দিনের দ্রুত সমাধান — সুস্বাদু Fundools ইনস্ট্যান্ট নুডলস!
মাত্র কয়েক মিনিটেই তৈরি করুন গরম, মজাদার আর পেট ভরানো এক প্লেট নুডলস।
✔️ দ্রুত রান্না
✔️ বাচ্চা থেকে বড় সবার পছন্দ
✔️ হালকা ক্ষুধায় পারফেক্ট স্ন্যাকস
অফিস, পড়াশোনা বা সন্ধ্যার আড্ডা—সব সময়েই জমে যায় এক প্লেট নুডলস।
🛒 আজই সংগ্রহ করুন Noor Super Mart থেকে।
সহজ অর্ডার, সাশ্রয়ী দাম আর দ্রুত ডেলিভারি—আপনার হাতের নাগালে 🚚✨', 2, NULL, NULL, NULL, '112239', 1, 24, '', 'public/images/products/fundools--.jpg', '{}', 170.00, 155.00, 1.00, 100.00, 100.00, 0.00, 'flat', 15.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 06:55:41.634358+00', '2026-02-28 07:58:33.341945+00', 0.00, 1.00);
INSERT INTO public.products VALUES (5, 'Salt /  লবন', 'লবণ
রান্নার স্বাদ ঠিক রাখতে ভালো মানের লবণ অপরিহার্য।
আমাদের লবণ বিশুদ্ধ, ঝরঝরে আর প্রতিদিনের রান্নার জন্য সম্পূর্ণ নিরাপদ।
এক কাপ ভাত, তেলমশলা মিশ্রিত তরকারি বা দই-চাটনি—সবকিছুর স্বাদ আসে আরও সুন্দর, যখন ব্যবহার করা হয় মানসম্মত লবণ।
✔️ বিশুদ্ধ ও মানসম্মত
✔️ স্বাদের সাথে কোনো ঝামেলা নয়
✔️ প্রতিদিনের রান্নায় ব্যবহার উপযোগী
✔️ ছোট ছোট রান্নায়ও স্বাদ সমানভাবে ছড়িয়ে দেয়
একটি ছোট উপাদান, কিন্তু খাবারের স্বাদে বিশাল পার্থক্য।
পরিবারের প্রতিটি খাবারে আনুন নির্ভরযোগ্য স্বাদ ও পরিপূর্ণতা।
🛒 আজই সংগ্রহ করুন Noor Super Mart থেকে।
নিশ্চিত মান, সাশ্রয়ী দাম এবং দ্রুত ডেলিভারি—সব একসাথে আপনার দোরগোড়ায় 🚚✨', 2, NULL, NULL, NULL, '112237', 1, 1, '', 'public/images/products/salt.webp', '{}', 42.00, 40.74, 1.00, 100.00, 100.00, 0.00, '', 3.00, 0.00, 'exclusive', 0.00, 'static', true, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 06:36:44.480948+00', '2026-02-28 08:03:26.819197+00', 0.00, 1.00);
INSERT INTO public.products VALUES (4, 'সয়াবিন তেল ২লিটার', 'সয়াবিন তেল
সুস্থ রান্না ও পরিবারের প্রতিদিনের পুষ্টির জন্য প্রয়োজন ভালো মানের সয়াবিন তেল।
আমাদের সয়াবিন তেল বিশুদ্ধ, হালকা স্বাদের এবং উচ্চ মানের পুষ্টি সংরক্ষণ করে।
✔️ স্বাস্থ্যসম্মত রান্নার জন্য উপযোগী
✔️ ঘরে ব্যবহার উপযোগী, ফ্রায়িং, ভাজি বা তরকারিতে পারফেক্ট
✔️ দীর্ঘ সময়ের জন্য তাজা ও গন্ধহীন
পরিবারের স্বাস্থ্যের কথা ভেবে প্রতিদিনের রান্নায় ব্যবহার করুন বিশুদ্ধ সয়াবিন তেল।
সুস্বাদু খাবারের স্বাদ বাড়াতে আর স্বাস্থ্যও রক্ষা করতে একদম পারফেক্ট।
🛒 আজই অর্ডার করুন Noor Super Mart থেকে।
নিশ্চিত মান, সাশ্রয়ী দাম এবং দ্রুত ডেলিভারি—সরাসরি আপনার দোরগোড়ায় 🚚✨', 2, NULL, NULL, NULL, '112236', 1, 19, '', 'public/images/products/thumb112236.jpg', '{}', 400.00, 400.00, 1.00, 50.00, 50.00, 0.00, 'percentage', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 06:24:38.097684+00', '2026-02-28 08:05:51.260827+00', 0.00, 1.00);
INSERT INTO public.products VALUES (104, 'POND''S Dream Flower Fragrant Talc Powder with Vitamin B3', 'Freshness and fragrance

Safe and natural

Brand

POND''S

Item Form

Powder

Skin Type

All

Product Benefits

Fresh

Recommended Uses For Product

Talc

Material Feature

Natural

Scent

Fragrant

Item Weight

100 Grams

Active Ingredients

benzyl alcohol,citronellol,eugenol,geraniol,limonene,linalool,menthol,niacinamide,vitamin b3

Number of Items

1

', 2, NULL, NULL, NULL, '1122967', 1, 2, '', 'public/images/products/thumb1122967.jpg', NULL, 220.00, 280.00, 1.00, 50.00, 50.00, 0.00, 'flat', 21.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-28 08:36:38.516099+00', '2026-02-28 08:36:38.522553+00', 0.00, 1.00);
INSERT INTO public.products VALUES (144, 'WATER TRAY CHULA', '', 3, NULL, NULL, 1, '12365', 1, 17, '', 'public/images/products/thumb12365.jpg', NULL, 60.00, 45.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 07:50:52.889382+00', '2026-03-15 07:50:52.894826+00', 45.00, 1.00);
INSERT INTO public.products VALUES (64, 'মুরগীর মাংস(Chicken)', '🍗 মুরগি (Chicken) – প্রোটিন সমৃদ্ধ স্বাস্থ্যকর খাবার
মুরগি প্রাকৃতিকভাবে প্রোটিনে ভরপুর, যা শরীরের পেশী গঠন, শক্তি বৃদ্ধি এবং সার্বিক স্বাস্থ্য সমর্থনে সাহায্য করে।
✨ উপকারিতা:
প্রোটিন সমৃদ্ধ, পেশী ও শক্তি বৃদ্ধিতে সহায়ক
ভিটামিন B ও খনিজে সমৃদ্ধ, শরীরের সঠিক ফাংশন বজায় রাখে
হজম সহজ করে এবং পাচন শক্তি বৃদ্ধি করে
রোগ প্রতিরোধ ক্ষমতা বৃদ্ধি করতে সাহায্য করে
ওজন নিয়ন্ত্রণে সহায়ক
💡 কেন ব্যবহার করবেন:
প্রতিদিনের খাদ্যাভ্যাসে মুরগি যুক্ত করলে শরীর সুস্থ, শক্তিশালী ও পুষ্টিকর থাকে।



✨ কেন Noor Super Mart থেকে কিনবেন:
১০০% সতেজ ও হাইজিনিক
মানসম্মত ও প্রিমিয়াম প্রোটিন সমৃদ্ধ মুরগি
সরাসরি ডেলিভারি, সময় ও শ্রম বাঁচায়
মূল্য সাশ্রয়ী ও রেগুলার স্টক
ক্রেতাদের জন্য সহজ অর্ডার ও নিরাপদ পেমেন্ট সিস্টেম
💡 আপনার সুবিধা:
প্রতিদিনের প্রয়োজন অনুযায়ী প্যাক/কেজি ক্রয় করতে পারবেন
স্বাস্থ্য ও পুষ্টির দিক থেকে নিশ্চিন্ত
দ্রুত ডেলিভারি, বাড়িতে পৌঁছে যাবে ফ্রেশ
📦 উপলব্ধ সাইজ/প্যাক:
পিস / কেজি অনুযায়ী
আজই অর্ডার করুন এবং স্বাস্থ্যকর ও সতেজ মুরগি উপভোগ করুন Noor Super Mart থেকে!
', 5, NULL, NULL, NULL, '112299', 1, 1, '', 'public/images/products/-chicken.jpg', '{}', 250.00, 250.00, 1.00, 100.00, 100.00, 0.00, 'flat', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-09 10:31:22.968172+00', '2026-03-12 04:11:23.390891+00', 240.00, 1.00);
INSERT INTO public.products VALUES (3, 'রসুন-Garlic', 'চায়না রসুন
রান্নার স্বাদ ও ঘ্রাণ বাড়াতে চায়না রসুনের জুড়ি নেই।
তাজা, ঝরঝরে এবং ঘ্রাণে সমৃদ্ধ চায়না রসুন প্রতিটি রেসিপিতে এনে দেয় নিখুঁত স্বাদ।
✔️ রান্নায় স্বাদ বাড়ায়
✔️ ঘন ও তাজা গুঁড়া/ক্লোভ
✔️ ভাজি, তরকারি, মাংস বা মাছের সাথে পারফেক্ট মিল
প্রতিদিনের রান্নায় স্বাদ ও সুগন্ধের জন্য বেছে নিন মানসম্মত চায়না রসুন।
রান্না হোক সুস্বাদু, ঘরও ভরে উঠুক সুবাসে।
🛒 আজই সংগ্রহ করুন Noor Super Mart থেকে।
বিশ্বাসযোগ্য মান, সাশ্রয়ী দাম এবং দ্রুত ডেলিভারি—সরাসরি আপনার দোরগোড়ায় 🚚✨', 2, NULL, NULL, NULL, '১১২২৩৫', 1, 1, '', 'public/images/products/-garlic.jpeg', '{}', 180.00, 165.00, 1.00, 50.00, 50.00, 0.00, 'flat', 15.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 06:20:16.760816+00', '2026-02-28 08:07:13.652959+00', 0.00, 1.00);
INSERT INTO public.products VALUES (2, 'মুড়ি', 'মুড়ি
হালকা, ক্রিস্পি আর স্বাদে ভরপুর মুড়ি—সকালের নাস্তা, বিকেলের হালকা স্ন্যাকস বা চা-সাথে perfect companion।
প্রতিটি মুখে আসে খুশির ছোট্ট মুহূর্ত, আর পরিবারের সবাইকে খুশি করে।
✔️ টাটকা ও ক্রিস্পি
✔️ হালকা ও সহজে হজমযোগ্য
✔️ চা, জুস বা সরাসরি খাবার হিসেবে উপযোগী
স্বাদ, স্বস্তি আর খুশি একসাথে আনুন ঘরে।
🛒 আজই অর্ডার করুন Noor Super Mart থেকে।
সাশ্রয়ী দাম, নিশ্চিত মান এবং দ্রুত ডেলিভারি—আপনার দোরগোড়ায় 🚚✨', 2, NULL, NULL, NULL, '112234', 1, 1, '', 'public/images/products/thumb112234.webp', '{public/images/products/gallery1122341jyuuv.webp,public/images/products/gallery1122342ew687.webp}', 80.00, 70.00, 1.00, 100.00, 50.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 06:02:41.317414+00', '2026-02-28 08:09:55.039928+00', 0.00, 1.00);
INSERT INTO public.products VALUES (100, 'Glow & Lovely Advanced Multi Vitamin Cream 47 gm', '✨ Glow & Lovely Advanced Multi Vitamin Cream
আপনার ত্বকের যত্নের জন্য একদম প্রিমিয়াম সমাধান — Glow & Lovely Advanced Multi Vitamin Cream।
ভিটামিন সমৃদ্ধ এই ক্রিমটি নিয়মিত ব্যবহারে ত্বককে রাখে নরম, মসৃণ, উজ্জ্বল এবং হালকা।
প্রতিদিনের দূষণ, সূর্যের আলো এবং ক্লান্তি থেকে ত্বককে রক্ষা করতে এটি দেয় অতিরিক্ত পুষ্টি ও তাজা অনুভূতি।
✔️ ত্বক নরম ও মসৃণ করে — স্পর্শেই বোঝা যায় মান
✔️ ভিটামিন সমৃদ্ধ — ত্বকের স্বাভাবিক স্বাস্থ্য ও উজ্জ্বলতা বৃদ্ধি করে
✔️ হালকা ও ত্বকে দ্রুত শোষিত — তেল ঝাপসা করে না
✔️ প্রতিদিনের ব্যবহার উপযোগী — সকাল-সন্ধ্যা যেকোনো সময়ে
নিজের স্বাভাবিক উজ্জ্বলতা ফিরিয়ে আনুন, ত্বককে দিন সতেজ ও রোদমাখা ভাব।
স্কিনকে ভালো রাখার জন্য আর অপেক্ষা না করে আজই ব্যবহার শুরু করুন।
🛒 Noor Super Mart থেকে আজই সংগ্রহ করুন।
সাশ্রয়ী দাম, মান নিশ্চিত এবং দ্রুত ডেলিভারি—আপনার দোরগোড়ায় 🚚✨', 2, NULL, NULL, NULL, '11229864', 1, NULL, '', 'public/images/products/thumbprd-26-2-100.webp', NULL, 165.00, 170.00, 1.00, 100.00, 50.00, 0.00, '', 5.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-28 08:19:26.282271+00', '2026-02-28 08:21:27.180862+00', 0.00, 1.00);
INSERT INTO public.products VALUES (101, 'fair & lovely winter fairness cream 80gm', 'Expert fairness with advanced multivitamins
Enriched with 24 hour active moisturizers
Ideal for winter
Sun protection with uva/uvb defense
Gives soft , supple fair skin', 2, NULL, NULL, NULL, '1122965', 1, NULL, '', 'public/images/products/thumb1122965.webp', NULL, 180.00, 200.00, 1.00, 50.00, 50.00, 0.00, '', 20.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-28 08:26:51.951842+00', '2026-02-28 08:26:51.960788+00', 0.00, 1.00);
INSERT INTO public.products VALUES (14, 'বাহরাইন অরেঞ্জ/আমের ট্যাং', '🍊 বাহরাইন অরেঞ্জ
রসালো, মিষ্টি আর হালকা টক স্বাদের একদম ফ্রেশ বাহরাইন অরেঞ্জ।
প্রতিটি কামড়ে পাবেন প্রাকৃতিক রস আর সতেজতার ছোঁয়া।
✔️ ভিটামিন C সমৃদ্ধ
✔️ জুস করার জন্য পারফেক্ট
✔️ বাচ্চা থেকে বড় সবার পছন্দ
সকালের নাশতা হোক বা বিকেলের জুস — বাহরাইন অরেঞ্জ মানেই ফ্রেশ এনার্জি।
🛒 এখনই অর্ডার করুন Noor Super Mart থেকে।
ঘরে বসেই পেয়ে যান টাটকা ফল, দ্রুত ডেলিভারিতে! 🚚✨
🥭 আমের টেং
গরমে এক গ্লাস ঠান্ডা আমের টেং — মুড ফ্রেশ করার সেরা উপায়!
ঘন আমের স্বাদ, মিষ্টি সুবাস আর তৃপ্তির একদম পারফেক্ট কম্বিনেশন।
✔️ সহজে তৈরি
✔️ বাচ্চাদের প্রিয়
✔️ অতিথি আপ্যায়নে ঝামেলাহীন সমাধান
এক চুমুকেই পাবেন আসল আমের স্বাদের মজা!
🥭 আজই নিয়ে নিন আপনার প্রয়োজনীয় পরিমাণ
শুধুমাত্র Noor Super Mart-এ — সহজ অর্ডার, নিশ্চিন্ত সার্ভিস।', 2, NULL, NULL, NULL, '112245', 1, 9, '', 'public/images/products/thumb112245.jpg', '{public/images/products/gallery1122451rri4h.webp}', 1800.00, 1950.00, 1.00, 20.00, 20.00, 0.00, 'flat', 150.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-05 07:24:25.146056+00', '2026-03-01 07:23:01.488658+00', 0.00, 1.00);
INSERT INTO public.products VALUES (133, 'surf exel', '', 2, NULL, NULL, NULL, 'PRD-26-2-133', 1, 8, '', 'public/images/products/thumbprd-26-2-133.jpg', NULL, 120.00, 120.00, 1.00, 10.00, 0.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-12 04:13:28.791123+00', '2026-03-12 04:13:28.806612+00', 118.00, 1.00);
INSERT INTO public.products VALUES (118, '20 g Cricket Plastic Ball', 'Color : Neon Green
Box Contains : 6 Pieces
Material : Plastic
Packaging Type : Box
Weight : 20 g
Country of Origin : Made in India', 2, NULL, NULL, NULL, 'PRD-26-2-118', 1, 13, '', 'public/images/products/thumbprd-26-2-118.jpg', NULL, 40.00, 0.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-10 08:20:37.251166+00', '2026-03-15 08:38:02.32249+00', 20.00, 1.00);
INSERT INTO public.products VALUES (123, 'FRESH EXERCISE BOOK ( KHATA ) PER PCS', '', 2, NULL, NULL, NULL, '1122553', 1, NULL, '', 'public/images/products/thumb1122553.jpg', NULL, 35.00, 35.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-11 08:17:28.232482+00', '2026-03-11 08:17:28.239289+00', 25.00, 1.00);
INSERT INTO public.products VALUES (69, 'জয়ত্রী', '🌿 জয়ত্রী (Joytree / Mace) – প্রাকৃতিক স্বাস্থ্যকর উপকারিতা
প্রাকৃতিক জয়ত্রী শুধু খাবারে সুগন্ধ ও স্বাদ বাড়ায় না, বরং শরীরের জন্য নানা উপকারও দেয়।
✨ উপকারিতা:
হজম শক্তি বৃদ্ধি করে, পাকস্থলীর স্বাস্থ্য উন্নত করে
রক্তচাপ নিয়ন্ত্রণে সাহায্য করে
প্রদাহ কমাতে সহায়তা করে
দেহের রোগ প্রতিরোধ ক্ষমতা বাড়ায়
চুল ও ত্বকের জন্য প্রাকৃতিক পুষ্টি যোগায়
💡 কেন ব্যবহার করবেন:
প্রতিদিনের খাদ্যাভ্যাসে জয়ত্রী যুক্ত করলে শরীর সুস্থ, হজম শক্তি ভালো এবং রোগপ্রতিরোধ ক্ষমতা বাড়ে।
📦 প্যাক সাইজ:
২৫g / ৫০g / ১০০g
আজই অর্ডার করুন এবং স্বাস্থ্যকে প্রাকৃতিকভাবে শক্তিশালী করুন 🌿✨', 2, NULL, NULL, NULL, '112303', 1, 16, '', 'public/images/products/thumb112303.jpg', '{}', 210.00, 200.00, 1.00, 100.00, 100.00, 0.00, 'flat', 10.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-10 06:29:21.702793+00', '2026-03-15 07:07:01.723132+00', 180.00, 1.00);
INSERT INTO public.products VALUES (127, 'ABC COTTON BUD', 'Yiwu Tango Household Products Co., Ltd. supplies high-quality ABC COTTON BUD to global B2B buyers, wholesalers, and distributors. Based in China , the company specializes in manufacturing and exporting premium Cotton Bud that meet international beauty industry standards.

The price of ABC COTTON BUD from China varies according to order volume and destination market. Yiwu Tango Household Products Co., Ltd. offers competitive and flexible pricing for different Cotton Bud specifications, ensuring cost-effective solutions for bulk importers and international distributors.', 2, NULL, NULL, NULL, '1122555', 1, NULL, '', 'public/images/products/thumb1122555.jpg', NULL, 140.00, 140.00, 12.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-11 08:30:21.51777+00', '2026-03-11 08:31:33.814466+00', 84.00, 100.00);
INSERT INTO public.products VALUES (142, 'DOUBLE GD AUTO', '', 3, NULL, NULL, 1, '1234', 1, 17, '', 'public/images/products/thumb1234.jpg', NULL, 170.00, 170.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 07:48:39.463111+00', '2026-03-15 07:48:39.46891+00', 170.00, 1.00);
INSERT INTO public.products VALUES (145, 'ELBO', '', 3, NULL, NULL, 1, '1258', 1, 17, '', 'public/images/products/thumb1258.jpg', NULL, 35.00, 30.00, 1.00, 1000.00, 1000.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 08:03:35.27203+00', '2026-03-15 08:03:35.277265+00', 30.00, 1.00);
INSERT INTO public.products VALUES (74, 'Honeycomb Burner', '🔥 Honeycomb Burner
আপনার রান্নাঘরের জন্য উন্নতমানের Honeycomb Burner—দ্রুত, সমান তাপে এবং নিরাপদ রান্নার নিশ্চয়তা দেয়। এর বিশেষ হানিকম কাঠামো তাপকে সমানভাবে ছড়িয়ে দেয়, যা রান্নাকে আরও কার্যকর ও জ্বালানি সাশ্রয়ী করে।
✨ বিশেষ বৈশিষ্ট্য:
হানিকম স্টাইলের ডিজাইন
সমান তাপ বিতরণ
দ্রুত রান্না ও কম গ্যাস খরচ
টেকসই ও মরিচা প্রতিরোধী
সহজে পরিষ্কার করা যায়
🍳 উপযোগিতা:
ডাবল বা সিঙ্গেল বার্নার গ্যাস স্টোভে ব্যবহারযোগ্য, দৈনন্দিন রান্না বা বড় পরিবারের জন্য উপযুক্ত।
আপনার রান্নার অভিজ্ঞতাকে আরও উন্নত করতে আজই সংগ্রহ করুন 🔥✨', 3, NULL, NULL, NULL, '1122356', 1, 23, '', 'public/images/products/honeycomb-burner.jpeg', '{}', 250.00, 220.00, 1.00, 100.00, 50.00, 0.00, 'flat', 30.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-02-19 04:13:54.119099+00', '2026-03-15 08:04:17.140182+00', 150.00, 1.00);
INSERT INTO public.products VALUES (147, 'CALM', '', 3, NULL, NULL, 1, '1236', 1, 2, '', 'public/images/products/thumb1236.jpg', NULL, 500.00, 380.00, 1.00, 20.00, 20.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-15 08:05:49.505182+00', '2026-03-15 08:05:49.511374+00', 380.00, 1.00);
INSERT INTO public.products VALUES (126, 'রাইজ ব্লেড[rise bleet]', '', 2, NULL, NULL, NULL, '1122554', 1, NULL, '', 'public/images/products/thumb1122554.jpg', NULL, 20.00, 20.00, 1.00, 100.00, 100.00, 0.00, '', 0.00, 0.00, 'exclusive', 0.00, 'static', false, '{}', 0, 0.00, 0, 0, 0, 0, 0, '2026-03-11 08:24:25.303765+00', '2026-03-15 08:49:27.153654+00', 220.00, 1.00);


--
-- TOC entry 3835 (class 0 OID 34137)
-- Dependencies: 220
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--



--
-- TOC entry 3839 (class 0 OID 34175)
-- Dependencies: 224
-- Data for Name: sub_categories; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.sub_categories VALUES (2, 2, 'cilynder', 1, true, '2026-03-06 15:08:23.918909+00', '2026-03-06 15:08:23.918909+00');


--
-- TOC entry 3841 (class 0 OID 34194)
-- Dependencies: 226
-- Data for Name: sub_sub_categories; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.sub_sub_categories VALUES (1, 2, 'omera', 1, true, '2026-03-06 15:08:45.318389+00', '2026-03-06 15:08:45.318389+00');


--
-- TOC entry 3855 (class 0 OID 34628)
-- Dependencies: 240
-- Data for Name: units; Type: TABLE DATA; Schema: public; Owner: super_shop_dev_user
--

INSERT INTO public.units VALUES (1, '1kg', '1kg', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (2, '100', '100', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (3, '1Ltr', '1Ltr', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (4, '250gm', '250gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (5, '250', '250', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (6, '20 bags', '20 bags', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (7, '25gm', '25gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (8, '500gm', '500gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (9, '2kg', '2kg', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (10, '200gm', '200gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (11, '1box', '1box', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (12, '1 Dala', '1 Dala', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (13, '1pes', '1pes', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (14, '360ml', '360ml', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (15, '175gm', '175gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (16, '50gm', '50gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (17, '1 Pcs', '1 Pcs', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (18, '98gm', '98gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (19, '2Ltr', '2Ltr', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (20, '75gm,', '75gm,', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (21, '10gm', '10gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (22, '1pcs', '1pcs', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (23, '1psc', '1psc', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (24, '496gm', '496gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (25, '150gm', '150gm', '2026-02-22 20:03:57.929166+00', '2026-02-22 20:03:57.929166+00');
INSERT INTO public.units VALUES (51, '12', 'kg', '2026-03-14 15:43:27.101768+00', '2026-03-14 15:43:27.101768+00');


--
-- TOC entry 3896 (class 0 OID 0)
-- Dependencies: 241
-- Name: attributes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.attributes_id_seq', 1, false);


--
-- TOC entry 3897 (class 0 OID 0)
-- Dependencies: 237
-- Name: branches_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.branches_id_seq', 1, true);


--
-- TOC entry 3898 (class 0 OID 0)
-- Dependencies: 249
-- Name: brands_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.brands_id_seq', 1, true);


--
-- TOC entry 3899 (class 0 OID 0)
-- Dependencies: 221
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.categories_id_seq', 5, true);


--
-- TOC entry 3900 (class 0 OID 0)
-- Dependencies: 229
-- Name: customers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.customers_id_seq', 9, true);


--
-- TOC entry 3901 (class 0 OID 0)
-- Dependencies: 217
-- Name: employees_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.employees_id_seq', 1, false);


--
-- TOC entry 3902 (class 0 OID 0)
-- Dependencies: 235
-- Name: hero_sections_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.hero_sections_id_seq', 1, false);


--
-- TOC entry 3903 (class 0 OID 0)
-- Dependencies: 233
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.order_items_id_seq', 79, true);


--
-- TOC entry 3904 (class 0 OID 0)
-- Dependencies: 231
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.orders_id_seq', 32, true);


--
-- TOC entry 3905 (class 0 OID 0)
-- Dependencies: 247
-- Name: product_reviews_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.product_reviews_id_seq', 1, false);


--
-- TOC entry 3906 (class 0 OID 0)
-- Dependencies: 245
-- Name: product_variations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.product_variations_id_seq', 1, false);


--
-- TOC entry 3907 (class 0 OID 0)
-- Dependencies: 227
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.products_id_seq', 87, true);


--
-- TOC entry 3908 (class 0 OID 0)
-- Dependencies: 243
-- Name: products_id_seq1; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.products_id_seq1', 152, true);


--
-- TOC entry 3909 (class 0 OID 0)
-- Dependencies: 219
-- Name: refresh_tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.refresh_tokens_id_seq', 1, false);


--
-- TOC entry 3910 (class 0 OID 0)
-- Dependencies: 223
-- Name: sub_categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.sub_categories_id_seq', 2, true);


--
-- TOC entry 3911 (class 0 OID 0)
-- Dependencies: 225
-- Name: sub_sub_categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.sub_sub_categories_id_seq', 1, true);


--
-- TOC entry 3912 (class 0 OID 0)
-- Dependencies: 239
-- Name: units_id_seq; Type: SEQUENCE SET; Schema: public; Owner: super_shop_dev_user
--

SELECT pg_catalog.setval('public.units_id_seq', 51, true);


--
-- TOC entry 3634 (class 2606 OID 34846)
-- Name: attributes attributes_name_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_name_key UNIQUE (name);


--
-- TOC entry 3636 (class 2606 OID 34844)
-- Name: attributes attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_pkey PRIMARY KEY (id);


--
-- TOC entry 3626 (class 2606 OID 34485)
-- Name: branches branches_name_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_name_key UNIQUE (name);


--
-- TOC entry 3628 (class 2606 OID 34483)
-- Name: branches branches_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_pkey PRIMARY KEY (id);


--
-- TOC entry 3674 (class 2606 OID 34984)
-- Name: brands brands_name_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.brands
    ADD CONSTRAINT brands_name_key UNIQUE (name);


--
-- TOC entry 3676 (class 2606 OID 34982)
-- Name: brands brands_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.brands
    ADD CONSTRAINT brands_pkey PRIMARY KEY (id);


--
-- TOC entry 3580 (class 2606 OID 34173)
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- TOC entry 3582 (class 2606 OID 34171)
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- TOC entry 3603 (class 2606 OID 34350)
-- Name: customers customers_email_unique; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_email_unique UNIQUE (email);


--
-- TOC entry 3605 (class 2606 OID 34348)
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- TOC entry 3566 (class 2606 OID 34133)
-- Name: employees employees_email_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_email_key UNIQUE (email);


--
-- TOC entry 3568 (class 2606 OID 34131)
-- Name: employees employees_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);


--
-- TOC entry 3623 (class 2606 OID 34467)
-- Name: hero_sections hero_sections_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.hero_sections
    ADD CONSTRAINT hero_sections_pkey PRIMARY KEY (id);


--
-- TOC entry 3621 (class 2606 OID 34421)
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- TOC entry 3615 (class 2606 OID 34401)
-- Name: orders orders_order_number_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_order_number_key UNIQUE (order_number);


--
-- TOC entry 3617 (class 2606 OID 34399)
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- TOC entry 3670 (class 2606 OID 34948)
-- Name: product_reviews product_reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT product_reviews_pkey PRIMARY KEY (id);


--
-- TOC entry 3672 (class 2606 OID 34950)
-- Name: product_reviews product_reviews_product_id_customer_id_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT product_reviews_product_id_customer_id_key UNIQUE (product_id, customer_id);


--
-- TOC entry 3661 (class 2606 OID 34925)
-- Name: product_variations product_variations_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_variations
    ADD CONSTRAINT product_variations_pkey PRIMARY KEY (id);


--
-- TOC entry 3663 (class 2606 OID 34927)
-- Name: product_variations product_variations_sku_product_id_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_variations
    ADD CONSTRAINT product_variations_sku_product_id_key UNIQUE (sku, product_id);


--
-- TOC entry 3601 (class 2606 OID 34322)
-- Name: legacy_products products_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.legacy_products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- TOC entry 3655 (class 2606 OID 34882)
-- Name: products products_pkey1; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey1 PRIMARY KEY (id);


--
-- TOC entry 3657 (class 2606 OID 34884)
-- Name: products products_sku_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_sku_key UNIQUE (sku);


--
-- TOC entry 3574 (class 2606 OID 34146)
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 3576 (class 2606 OID 34148)
-- Name: refresh_tokens refresh_tokens_token_hash_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_token_hash_key UNIQUE (token_hash);


--
-- TOC entry 3585 (class 2606 OID 34186)
-- Name: sub_categories sub_categories_name_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.sub_categories
    ADD CONSTRAINT sub_categories_name_key UNIQUE (name);


--
-- TOC entry 3587 (class 2606 OID 34184)
-- Name: sub_categories sub_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.sub_categories
    ADD CONSTRAINT sub_categories_pkey PRIMARY KEY (id);


--
-- TOC entry 3590 (class 2606 OID 34205)
-- Name: sub_sub_categories sub_sub_categories_name_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.sub_sub_categories
    ADD CONSTRAINT sub_sub_categories_name_key UNIQUE (name);


--
-- TOC entry 3592 (class 2606 OID 34203)
-- Name: sub_sub_categories sub_sub_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.sub_sub_categories
    ADD CONSTRAINT sub_sub_categories_pkey PRIMARY KEY (id);


--
-- TOC entry 3630 (class 2606 OID 34637)
-- Name: units units_name_key; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_name_key UNIQUE (name);


--
-- TOC entry 3632 (class 2606 OID 34635)
-- Name: units units_pkey; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_pkey PRIMARY KEY (id);


--
-- TOC entry 3578 (class 2606 OID 34150)
-- Name: refresh_tokens uq_employee_active_token; Type: CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT uq_employee_active_token UNIQUE (employee_id, revoked_at);


--
-- TOC entry 3606 (class 1259 OID 34351)
-- Name: idx_customers_phone; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_customers_phone ON public.customers USING btree (phone);


--
-- TOC entry 3607 (class 1259 OID 34352)
-- Name: idx_customers_referral_code; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_customers_referral_code ON public.customers USING btree (referral_code);


--
-- TOC entry 3569 (class 1259 OID 34134)
-- Name: idx_employees_email; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_employees_email ON public.employees USING btree (email);


--
-- TOC entry 3570 (class 1259 OID 34135)
-- Name: idx_employees_uuid; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_employees_uuid ON public.employees USING btree (uuid);


--
-- TOC entry 3593 (class 1259 OID 34326)
-- Name: idx_legacy_products_brand_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_legacy_products_brand_id ON public.legacy_products USING btree (brand_id);


--
-- TOC entry 3594 (class 1259 OID 34325)
-- Name: idx_legacy_products_category_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_legacy_products_category_id ON public.legacy_products USING btree (category_id);


--
-- TOC entry 3595 (class 1259 OID 34324)
-- Name: idx_legacy_products_name; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_legacy_products_name ON public.legacy_products USING btree (name);


--
-- TOC entry 3596 (class 1259 OID 34328)
-- Name: idx_legacy_products_published; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_legacy_products_published ON public.legacy_products USING btree (published);


--
-- TOC entry 3597 (class 1259 OID 34323)
-- Name: idx_legacy_products_slug; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_legacy_products_slug ON public.legacy_products USING btree (slug);


--
-- TOC entry 3598 (class 1259 OID 34327)
-- Name: idx_legacy_products_status; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_legacy_products_status ON public.legacy_products USING btree (status);


--
-- TOC entry 3599 (class 1259 OID 34329)
-- Name: idx_legacy_products_user_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_legacy_products_user_id ON public.legacy_products USING btree (user_id);


--
-- TOC entry 3618 (class 1259 OID 34438)
-- Name: idx_order_items_order_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_order_items_order_id ON public.order_items USING btree (order_id);


--
-- TOC entry 3619 (class 1259 OID 34439)
-- Name: idx_order_items_product_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_order_items_product_id ON public.order_items USING btree (product_id);


--
-- TOC entry 3608 (class 1259 OID 34435)
-- Name: idx_orders_created_at; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_orders_created_at ON public.orders USING btree (created_at DESC);


--
-- TOC entry 3609 (class 1259 OID 34432)
-- Name: idx_orders_customer_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_orders_customer_id ON public.orders USING btree (customer_id);


--
-- TOC entry 3610 (class 1259 OID 34437)
-- Name: idx_orders_customer_mobile; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_orders_customer_mobile ON public.orders USING btree (customer_mobile);


--
-- TOC entry 3611 (class 1259 OID 34436)
-- Name: idx_orders_order_number; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_orders_order_number ON public.orders USING btree (order_number);


--
-- TOC entry 3612 (class 1259 OID 34433)
-- Name: idx_orders_order_status; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_orders_order_status ON public.orders USING btree (order_status);


--
-- TOC entry 3613 (class 1259 OID 34434)
-- Name: idx_orders_payment_status; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_orders_payment_status ON public.orders USING btree (payment_status);


--
-- TOC entry 3637 (class 1259 OID 34903)
-- Name: idx_products_avg_rating; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_avg_rating ON public.products USING btree (avg_rating);


--
-- TOC entry 3638 (class 1259 OID 34898)
-- Name: idx_products_brand; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_brand ON public.products USING btree (brand_id);


--
-- TOC entry 3639 (class 1259 OID 34895)
-- Name: idx_products_category; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_category ON public.products USING btree (category_id);


--
-- TOC entry 3640 (class 1259 OID 34905)
-- Name: idx_products_five_star_count; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_five_star_count ON public.products USING btree (five_star_count);


--
-- TOC entry 3641 (class 1259 OID 34906)
-- Name: idx_products_four_star_count; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_four_star_count ON public.products USING btree (four_star_count);


--
-- TOC entry 3642 (class 1259 OID 34909)
-- Name: idx_products_one_star_count; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_one_star_count ON public.products USING btree (one_star_count);


--
-- TOC entry 3643 (class 1259 OID 34901)
-- Name: idx_products_price; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_price ON public.products USING btree (retail_price);


--
-- TOC entry 3644 (class 1259 OID 34900)
-- Name: idx_products_stock; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_stock ON public.products USING btree (current_stock_qty);


--
-- TOC entry 3645 (class 1259 OID 34896)
-- Name: idx_products_sub_category; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_sub_category ON public.products USING btree (sub_category_id);


--
-- TOC entry 3646 (class 1259 OID 34897)
-- Name: idx_products_sub_subcategory; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_sub_subcategory ON public.products USING btree (sub_sub_category_id);


--
-- TOC entry 3647 (class 1259 OID 34910)
-- Name: idx_products_tags; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_tags ON public.products USING gin (to_tsvector('simple'::regconfig, tags));


--
-- TOC entry 3648 (class 1259 OID 34907)
-- Name: idx_products_three_star_count; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_three_star_count ON public.products USING btree (three_star_count);


--
-- TOC entry 3649 (class 1259 OID 34904)
-- Name: idx_products_total_reviews; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_total_reviews ON public.products USING btree (total_reviews);


--
-- TOC entry 3650 (class 1259 OID 34902)
-- Name: idx_products_total_sold; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_total_sold ON public.products USING btree (total_sold);


--
-- TOC entry 3651 (class 1259 OID 34908)
-- Name: idx_products_two_star_count; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_two_star_count ON public.products USING btree (two_star_count);


--
-- TOC entry 3652 (class 1259 OID 34899)
-- Name: idx_products_unit_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_unit_id ON public.products USING btree (unit_id);


--
-- TOC entry 3653 (class 1259 OID 34911)
-- Name: idx_products_variation_attributes; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_products_variation_attributes ON public.products USING gin (variation_attributes);


--
-- TOC entry 3571 (class 1259 OID 34156)
-- Name: idx_refresh_tokens_employee_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_refresh_tokens_employee_id ON public.refresh_tokens USING btree (employee_id);


--
-- TOC entry 3572 (class 1259 OID 34157)
-- Name: idx_refresh_tokens_token_hash; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_refresh_tokens_token_hash ON public.refresh_tokens USING btree (token_hash);


--
-- TOC entry 3664 (class 1259 OID 34965)
-- Name: idx_reviews_customer; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_reviews_customer ON public.product_reviews USING btree (customer_id);


--
-- TOC entry 3665 (class 1259 OID 34961)
-- Name: idx_reviews_product; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_reviews_product ON public.product_reviews USING btree (product_id);


--
-- TOC entry 3666 (class 1259 OID 34962)
-- Name: idx_reviews_product_created; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_reviews_product_created ON public.product_reviews USING btree (product_id, created_at DESC);


--
-- TOC entry 3667 (class 1259 OID 34964)
-- Name: idx_reviews_rating; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_reviews_rating ON public.product_reviews USING btree (rating);


--
-- TOC entry 3668 (class 1259 OID 34963)
-- Name: idx_reviews_status; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_reviews_status ON public.product_reviews USING btree (status);


--
-- TOC entry 3583 (class 1259 OID 34192)
-- Name: idx_sub_cat_parent; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_sub_cat_parent ON public.sub_categories USING btree (category_id);


--
-- TOC entry 3588 (class 1259 OID 34211)
-- Name: idx_sub_sub_cat_parent; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_sub_sub_cat_parent ON public.sub_sub_categories USING btree (sub_category_id);


--
-- TOC entry 3658 (class 1259 OID 34934)
-- Name: idx_variations_attributes; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_variations_attributes ON public.product_variations USING gin (variation_attributes);


--
-- TOC entry 3659 (class 1259 OID 34933)
-- Name: idx_variations_product_id; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE INDEX idx_variations_product_id ON public.product_variations USING btree (product_id);


--
-- TOC entry 3624 (class 1259 OID 34468)
-- Name: only_one_row; Type: INDEX; Schema: public; Owner: super_shop_dev_user
--

CREATE UNIQUE INDEX only_one_row ON public.hero_sections USING btree ((1));


--
-- TOC entry 3688 (class 2620 OID 34441)
-- Name: orders trigger_orders_updated_at; Type: TRIGGER; Schema: public; Owner: super_shop_dev_user
--

CREATE TRIGGER trigger_orders_updated_at BEFORE UPDATE ON public.orders FOR EACH ROW EXECUTE FUNCTION public.update_orders_updated_at();


--
-- TOC entry 3681 (class 2606 OID 34422)
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- TOC entry 3682 (class 2606 OID 34427)
-- Name: order_items order_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.legacy_products(id) ON DELETE SET NULL;


--
-- TOC entry 3680 (class 2606 OID 34402)
-- Name: orders orders_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE SET NULL;


--
-- TOC entry 3686 (class 2606 OID 34956)
-- Name: product_reviews product_reviews_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT product_reviews_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- TOC entry 3687 (class 2606 OID 34951)
-- Name: product_reviews product_reviews_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT product_reviews_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- TOC entry 3685 (class 2606 OID 34928)
-- Name: product_variations product_variations_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.product_variations
    ADD CONSTRAINT product_variations_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- TOC entry 3683 (class 2606 OID 34885)
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- TOC entry 3684 (class 2606 OID 34890)
-- Name: products products_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id);


--
-- TOC entry 3677 (class 2606 OID 34151)
-- Name: refresh_tokens refresh_tokens_employee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_employee_id_fkey FOREIGN KEY (employee_id) REFERENCES public.employees(id) ON DELETE CASCADE;


--
-- TOC entry 3678 (class 2606 OID 34187)
-- Name: sub_categories sub_categories_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.sub_categories
    ADD CONSTRAINT sub_categories_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- TOC entry 3679 (class 2606 OID 34206)
-- Name: sub_sub_categories sub_sub_categories_sub_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: super_shop_dev_user
--

ALTER TABLE ONLY public.sub_sub_categories
    ADD CONSTRAINT sub_sub_categories_sub_category_id_fkey FOREIGN KEY (sub_category_id) REFERENCES public.sub_categories(id) ON DELETE CASCADE;


--
-- TOC entry 3872 (class 0 OID 0)
-- Dependencies: 7
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: super_shop_dev_user
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;


-- Completed on 2026-03-15 17:13:22

--
-- PostgreSQL database dump complete
--

