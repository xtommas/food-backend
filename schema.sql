--
-- PostgreSQL database dump
--

\restrict hbrbTLgkdHQkn1irlZ75aDDpKP877SAHmaxpotkzd3VG4ldPSILh85nfz0bioVc

-- Dumped from database version 17.10
-- Dumped by pg_dump version 17.10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: citext; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


--
-- Name: set_updated_at(); Type: FUNCTION; Schema: public; Owner: dockerfood
--

CREATE FUNCTION public.set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.set_updated_at() OWNER TO dockerfood;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: dishes; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.dishes (
    id bigint NOT NULL,
    restaurant_id bigint NOT NULL,
    name text NOT NULL,
    price bigint NOT NULL,
    description text NOT NULL,
    categories text[] NOT NULL,
    photo text,
    available boolean DEFAULT true NOT NULL,
    updated_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    CONSTRAINT categories_length_check CHECK (((array_length(categories, 1) >= 1) AND (array_length(categories, 1) <= 5))),
    CONSTRAINT dishes_price_check CHECK (((price)::double precision >= (0)::double precision))
);


ALTER TABLE public.dishes OWNER TO dockerfood;

--
-- Name: dishes_id_seq; Type: SEQUENCE; Schema: public; Owner: dockerfood
--

CREATE SEQUENCE public.dishes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dishes_id_seq OWNER TO dockerfood;

--
-- Name: dishes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dockerfood
--

ALTER SEQUENCE public.dishes_id_seq OWNED BY public.dishes.id;


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.order_items (
    id bigint NOT NULL,
    order_id bigint NOT NULL,
    dish_id bigint NOT NULL,
    quantity integer NOT NULL,
    subtotal bigint NOT NULL,
    dish_name text NOT NULL,
    unit_price bigint NOT NULL
);


ALTER TABLE public.order_items OWNER TO dockerfood;

--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: dockerfood
--

CREATE SEQUENCE public.order_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_items_id_seq OWNER TO dockerfood;

--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dockerfood
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    restaurant_id bigint NOT NULL,
    total bigint NOT NULL,
    address text NOT NULL,
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    status text NOT NULL,
    updated_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    CONSTRAINT orders_status_check CHECK ((status = ANY (ARRAY['pending'::text, 'confirmed'::text, 'preparing'::text, 'ready'::text, 'delivered'::text, 'cancelled'::text])))
);


ALTER TABLE public.orders OWNER TO dockerfood;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: dockerfood
--

CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO dockerfood;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dockerfood
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.permissions (
    id bigint NOT NULL,
    code text NOT NULL
);


ALTER TABLE public.permissions OWNER TO dockerfood;

--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: dockerfood
--

CREATE SEQUENCE public.permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.permissions_id_seq OWNER TO dockerfood;

--
-- Name: permissions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dockerfood
--

ALTER SEQUENCE public.permissions_id_seq OWNED BY public.permissions.id;


--
-- Name: restaurant_staff; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.restaurant_staff (
    user_id bigint NOT NULL,
    restaurant_id bigint NOT NULL,
    role text DEFAULT 'staff'::text NOT NULL
);


ALTER TABLE public.restaurant_staff OWNER TO dockerfood;

--
-- Name: restaurants; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.restaurants (
    id bigint NOT NULL,
    name text NOT NULL,
    photo text,
    address text NOT NULL,
    city text NOT NULL,
    state text,
    province text,
    country text NOT NULL,
    latitude numeric(9,6),
    longitude numeric(9,6),
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    version integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.restaurants OWNER TO dockerfood;

--
-- Name: restaurants_id_seq; Type: SEQUENCE; Schema: public; Owner: dockerfood
--

CREATE SEQUENCE public.restaurants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.restaurants_id_seq OWNER TO dockerfood;

--
-- Name: restaurants_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dockerfood
--

ALTER SEQUENCE public.restaurants_id_seq OWNED BY public.restaurants.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO dockerfood;

--
-- Name: users; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    photo text,
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    name text NOT NULL,
    email public.citext NOT NULL,
    password_hash bytea NOT NULL,
    activated boolean NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    role text NOT NULL
);


ALTER TABLE public.users OWNER TO dockerfood;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: dockerfood
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO dockerfood;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dockerfood
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users_permissions; Type: TABLE; Schema: public; Owner: dockerfood
--

CREATE TABLE public.users_permissions (
    user_id bigint NOT NULL,
    permission_id bigint NOT NULL
);


ALTER TABLE public.users_permissions OWNER TO dockerfood;

--
-- Name: dishes id; Type: DEFAULT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.dishes ALTER COLUMN id SET DEFAULT nextval('public.dishes_id_seq'::regclass);


--
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: permissions id; Type: DEFAULT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);


--
-- Name: restaurants id; Type: DEFAULT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.restaurants ALTER COLUMN id SET DEFAULT nextval('public.restaurants_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: dishes dishes_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.dishes
    ADD CONSTRAINT dishes_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: restaurant_staff restaurant_staff_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.restaurant_staff
    ADD CONSTRAINT restaurant_staff_pkey PRIMARY KEY (user_id, restaurant_id);


--
-- Name: restaurants restaurants_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.restaurants
    ADD CONSTRAINT restaurants_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users_permissions users_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT users_permissions_pkey PRIMARY KEY (user_id, permission_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: dishes_categories_idx; Type: INDEX; Schema: public; Owner: dockerfood
--

CREATE INDEX dishes_categories_idx ON public.dishes USING gin (categories);


--
-- Name: dishes_names_idx; Type: INDEX; Schema: public; Owner: dockerfood
--

CREATE INDEX dishes_names_idx ON public.dishes USING gin (to_tsvector('simple'::regconfig, name));


--
-- Name: dishes_restaurant_id_idx; Type: INDEX; Schema: public; Owner: dockerfood
--

CREATE INDEX dishes_restaurant_id_idx ON public.dishes USING btree (restaurant_id);


--
-- Name: order_items_order_id_idx; Type: INDEX; Schema: public; Owner: dockerfood
--

CREATE INDEX order_items_order_id_idx ON public.order_items USING btree (order_id);


--
-- Name: orders_restaurant_id_idx; Type: INDEX; Schema: public; Owner: dockerfood
--

CREATE INDEX orders_restaurant_id_idx ON public.orders USING btree (restaurant_id);


--
-- Name: orders_status_idx; Type: INDEX; Schema: public; Owner: dockerfood
--

CREATE INDEX orders_status_idx ON public.orders USING btree (status);


--
-- Name: orders_user_id_idx; Type: INDEX; Schema: public; Owner: dockerfood
--

CREATE INDEX orders_user_id_idx ON public.orders USING btree (user_id);


--
-- Name: dishes dishes_set_updated_at; Type: TRIGGER; Schema: public; Owner: dockerfood
--

CREATE TRIGGER dishes_set_updated_at BEFORE UPDATE ON public.dishes FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: orders orders_set_updated_at; Type: TRIGGER; Schema: public; Owner: dockerfood
--

CREATE TRIGGER orders_set_updated_at BEFORE UPDATE ON public.orders FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();


--
-- Name: dishes dishes_restaurant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.dishes
    ADD CONSTRAINT dishes_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_dish_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_dish_id_fkey FOREIGN KEY (dish_id) REFERENCES public.dishes(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: orders orders_restaurant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE;


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: restaurant_staff restaurant_staff_restaurant_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.restaurant_staff
    ADD CONSTRAINT restaurant_staff_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES public.restaurants(id) ON DELETE CASCADE;


--
-- Name: restaurant_staff restaurant_staff_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.restaurant_staff
    ADD CONSTRAINT restaurant_staff_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users_permissions users_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT users_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: users_permissions users_permissions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dockerfood
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT users_permissions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict hbrbTLgkdHQkn1irlZ75aDDpKP877SAHmaxpotkzd3VG4ldPSILh85nfz0bioVc

