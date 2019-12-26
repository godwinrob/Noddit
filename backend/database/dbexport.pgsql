--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5
-- Dumped by pg_dump version 11.5

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

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: favorites; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.favorites (
    user_id bigint,
    sn_id bigint,
    post_id bigint
);


ALTER TABLE public.favorites OWNER TO postgres;

--
-- Name: mod; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mod (
    sn_id bigint NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE public.mod OWNER TO postgres;

--
-- Name: post_votes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.post_votes (
    post_id bigint NOT NULL,
    user_id bigint NOT NULL,
    vote character varying(8) NOT NULL
);


ALTER TABLE public.post_votes OWNER TO postgres;

--
-- Name: posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.posts (
    post_id integer NOT NULL,
    parent_post_id bigint,
    sn_id bigint NOT NULL,
    user_id bigint NOT NULL,
    title character varying(100) NOT NULL,
    body character varying(2000) NOT NULL,
    image_address character varying,
    post_score bigint,
    top_level_id bigint,
    created_date timestamp(6) without time zone NOT NULL
);


ALTER TABLE public.posts OWNER TO postgres;

--
-- Name: posts_post_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.posts_post_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.posts_post_id_seq OWNER TO postgres;

--
-- Name: posts_post_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.posts_post_id_seq OWNED BY public.posts.post_id;


--
-- Name: subnoddits; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.subnoddits (
    sn_id integer NOT NULL,
    sn_name character varying(30) NOT NULL,
    sn_description character varying(200) NOT NULL
);


ALTER TABLE public.subnoddits OWNER TO postgres;

--
-- Name: subnoddits_sn_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.subnoddits_sn_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.subnoddits_sn_id_seq OWNER TO postgres;

--
-- Name: subnoddits_sn_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.subnoddits_sn_id_seq OWNED BY public.subnoddits.sn_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(32) NOT NULL,
    salt character varying(256) NOT NULL,
    role character varying(255) DEFAULT 'user'::character varying NOT NULL,
    avatar_address character varying(200),
    first_name character varying(20),
    last_name character varying(20),
    email_address character varying(30),
    join_date date
);


ALTER TABLE public.users OWNER TO postgres;

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


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: posts post_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts ALTER COLUMN post_id SET DEFAULT nextval('public.posts_post_id_seq'::regclass);


--
-- Name: subnoddits sn_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subnoddits ALTER COLUMN sn_id SET DEFAULT nextval('public.subnoddits_sn_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: favorites; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.favorites (user_id, sn_id, post_id) FROM stdin;
1	\N	2
4	1	\N
5	1	\N
1	2	\N
1	1	\N
1	3	\N
6	1	\N
7	1	\N
7	\N	186
7	5	\N
7	10	\N
2	3	\N
2	2	\N
\.


--
-- Data for Name: mod; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.mod (sn_id, user_id) FROM stdin;
6	2
9	2
10	7
\.


--
-- Data for Name: post_votes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.post_votes (post_id, user_id, vote) FROM stdin;
61	4	upvote
61	1	upvote
20	1	upvote
18	1	downvote
6	1	upvote
19	1	upvote
41	1	upvote
40	1	downvote
17	1	upvote
1	1	upvote
16	1	upvote
17	2	upvote
61	2	upvote
17	2	downvote
172	2	upvote
186	7	upvote
192	7	upvote
192	2	downvote
40	2	upvote
6	2	upvote
41	2	downvote
197	2	upvote
\.


--
-- Data for Name: posts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.posts (post_id, parent_post_id, sn_id, user_id, title, body, image_address, post_score, top_level_id, created_date) FROM stdin;
11	\N	3	2	More Harold	Harold	https://brobible.files.wordpress.com/2019/11/istock-153696622.jpg	1	\N	2019-12-10 10:00:00
68	61	3	2	post-reply	test	\N	1	61	2019-12-13 16:56:37.134
69	20	1	2	post-reply	test	\N	1	20	2019-12-13 16:57:04.845
19	6	1	2	test comment	I agree!	\N	2	2	2019-12-12 10:00:00
16	\N	3	3	I am the Big Sad.	Big sad!	https://assets.change.org/photos/9/io/nc/faiONCulPSIOaSe-800x450-noPad.jpg?1525707064	2	\N	2019-12-10 10:00:00
76	68	3	1	comment-reply	fuck this thread	\N	1	61	2019-12-16 10:20:57.069
40	6	1	2	comment-reply	Parent level test	\N	1	2	2019-12-13 13:12:14.766
116	78	3	2	comment-reply		\N	1	77	2019-12-18 10:18:47.712
17	\N	1	2	This is a cat	CAT CAT	https://www.rd.com/wp-content/uploads/2019/11/cat-10-e1573844975155-1024x692.jpg	0	\N	2019-12-10 10:00:00
78	77	3	2	post-reply	Amazing!	\N	1	77	2019-12-17 13:08:41.869
79	78	3	2	comment-reply		\N	1	77	2019-12-18 09:10:02.505
80	78	3	2	comment-reply		\N	1	77	2019-12-18 09:10:10.71
81	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:48.665
82	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:55.867
83	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:56.067
84	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:56.247
85	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:56.443
86	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:56.583
87	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:56.711
88	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:58.287
89	80	3	2	comment-reply		\N	1	77	2019-12-18 09:47:58.427
90	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:06.901
91	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:08.499
92	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:08.621
93	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:09.178
94	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:09.32
95	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:09.452
96	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:10.169
97	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:10.285
98	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:12.356
99	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:12.52
100	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:13.098
101	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:13.23
102	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:13.943
103	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:14.122
104	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:14.274
105	80	3	2	comment-reply		\N	1	77	2019-12-18 09:48:14.721
106	105	3	2	comment-reply		\N	1	77	2019-12-18 09:49:28.656
107	105	3	2	comment-reply		\N	1	77	2019-12-18 09:49:58.093
108	105	3	2	comment-reply		\N	1	77	2019-12-18 09:50:01.877
109	105	3	2	comment-reply		\N	1	77	2019-12-18 09:50:27.937
110	104	3	2	comment-reply		\N	1	77	2019-12-18 09:50:29.608
111	104	3	2	comment-reply		\N	1	77	2019-12-18 09:50:31.799
112	105	3	2	comment-reply		\N	1	77	2019-12-18 09:50:33.831
113	78	3	2	comment-reply		\N	1	77	2019-12-18 10:03:14.275
114	78	3	2	comment-reply		\N	1	77	2019-12-18 10:14:12.531
115	78	3	2	comment-reply		\N	1	77	2019-12-18 10:16:08.534
117	78	3	2	comment-reply		\N	1	77	2019-12-18 10:19:32.08
118	116	3	2	comment-reply		\N	1	77	2019-12-18 10:25:07.401
119	116	3	2	comment-reply		\N	1	77	2019-12-18 10:25:13.321
120	116	3	2	comment-reply		\N	1	77	2019-12-18 10:26:44.988
121	116	3	2	comment-reply		\N	1	77	2019-12-18 10:26:47.574
122	116	3	2	comment-reply		\N	1	77	2019-12-18 10:26:57.391
123	116	3	2	comment-reply		\N	1	77	2019-12-18 10:27:01.677
124	116	3	2	comment-reply		\N	1	77	2019-12-18 10:27:13.594
125	116	3	2	comment-reply		\N	1	77	2019-12-18 10:27:19.105
126	77	3	2	post-reply		\N	1	77	2019-12-18 10:30:46.961
127	77	3	2	post-reply		\N	1	77	2019-12-18 10:30:47.999
6	2	1	1	Wow	your cat looks so cool	\N	4	2	2019-12-10 10:00:00
41	19	1	2	comment-reply	Nested test	\N	1	2	2019-12-13 13:12:26.474
2	\N	1	1	This is Evie	She is a great kitty	https://imgur.com/nOQfW6Z.jpg	1	\N	2019-12-10 10:00:00
4	\N	2	1	Dogs are bad	Here is a cat instead	https://imgur.com/erIPkF5.jpg	1	\N	2019-12-10 10:00:00
5	\N	5	1	I grew this	Look at this sweet plant that I grew	https://imgur.com/gb8kBNe.jpg	1	\N	2019-12-10 10:00:00
18	\N	2	2	Dog making cash	Tiny dog BIG money	https://i.imgur.com/LRoLTlK.jpg	99	\N	2019-12-11 11:00:00
1	\N	3	1	Hide the Pain with Harold	He can do it an so can you	https://imgur.com/sM8MtJn.jpg	2	\N	2019-12-10 10:00:00
61	\N	3	2	More Harold	Even more Harold!	https://i.imgur.com/caYzJPI.jpg	2	\N	2019-12-13 15:42:39.403
128	77	3	2	post-reply		\N	1	77	2019-12-18 10:30:48.455
129	77	3	2	post-reply		\N	1	77	2019-12-18 10:30:48.587
130	116	3	2	comment-reply		\N	1	77	2019-12-18 10:30:49.559
131	77	3	2	post-reply	tttt	\N	1	77	2019-12-18 10:31:12.633
132	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:26.476
133	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:40.318
134	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:40.705
135	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:40.873
136	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:41.005
137	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:41.157
138	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:41.289
139	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:41.453
140	77	3	2	post-reply		\N	1	77	2019-12-18 10:31:41.581
154	68	3	2	comment-reply		\N	1	61	2019-12-18 10:46:38.709
155	61	3	2	post-reply		\N	1	61	2019-12-18 10:46:40.461
156	68	3	2	comment-reply		\N	1	61	2019-12-18 10:46:45.627
157	61	3	2	post-reply		\N	1	61	2019-12-18 10:46:50.841
158	68	3	2	comment-reply		\N	1	61	2019-12-18 10:46:52.405
141	77	3	2	post-reply		\N	1	77	2019-12-18 10:33:28.711
142	77	3	2	post-reply		\N	1	77	2019-12-18 10:33:30.999
143	116	3	2	comment-reply		\N	1	77	2019-12-18 10:33:32.579
144	116	3	2	comment-reply		\N	1	77	2019-12-18 10:33:39.138
145	77	3	2	post-reply		\N	1	77	2019-12-18 10:33:40.073
146	116	3	2	comment-reply		\N	1	77	2019-12-18 10:33:50.122
147	77	3	2	post-reply		\N	1	77	2019-12-18 10:33:51.236
148	116	3	2	comment-reply		\N	1	77	2019-12-18 10:35:40.663
149	77	3	2	post-reply		\N	1	77	2019-12-18 10:35:41.745
150	77	3	2	post-reply		\N	1	77	2019-12-18 10:35:49.713
151	77	3	2	post-reply		\N	1	77	2019-12-18 10:37:24.288
152	77	3	2	post-reply		\N	1	77	2019-12-18 10:38:17.729
153	116	3	2	comment-reply		\N	1	77	2019-12-18 10:38:20.074
159	153	3	2	comment-reply		\N	1	77	2019-12-18 10:49:44.21
160	77	3	2	post-reply		\N	1	77	2019-12-18 10:51:44.786
161	77	3	2	post-reply		\N	1	77	2019-12-18 10:51:46.292
162	116	3	2	comment-reply		\N	1	77	2019-12-18 10:57:30.684
163	78	3	2	comment-reply		\N	1	77	2019-12-18 10:57:32.321
164	77	3	2	post-reply		\N	1	77	2019-12-18 10:58:18.116
165	116	3	2	comment-reply		\N	1	77	2019-12-18 11:01:42.364
166	116	3	2	comment-reply		\N	1	77	2019-12-18 11:01:44.082
167	116	3	2	comment-reply		\N	1	77	2019-12-18 11:01:45.973
168	116	3	2	comment-reply		\N	1	77	2019-12-18 11:02:03.227
169	168	3	2	comment-reply		\N	1	77	2019-12-18 11:02:54.773
170	169	3	2	comment-reply		\N	1	77	2019-12-18 11:03:15.111
171	116	3	2	comment-reply		\N	1	77	2019-12-18 11:03:18.367
172	116	3	2	comment-reply	test	\N	2	77	2019-12-18 11:04:46.899
173	116	3	2	comment-reply		\N	1	77	2019-12-18 11:13:11.362
174	78	3	2	comment-reply		\N	1	77	2019-12-18 11:13:16.324
175	77	3	2	post-reply	f	\N	1	77	2019-12-18 11:17:52.31
176	77	3	2	post-reply		\N	1	77	2019-12-18 11:25:05.99
177	77	3	2	post-reply		\N	1	77	2019-12-18 11:25:07.705
178	11	3	2	post-reply	test	\N	1	11	2019-12-18 11:26:52.424
179	178	3	2	comment-reply	t	\N	1	11	2019-12-18 11:27:03.275
180	77	3	2	post-reply		\N	1	77	2019-12-18 11:28:33.175
181	179	3	2	comment-reply	TEST	\N	1	11	2019-12-18 11:45:17.098
182	11	3	2	post-reply		\N	1	11	2019-12-18 11:45:46.184
183	181	3	2	comment-reply		\N	1	11	2019-12-18 11:47:55.614
184	181	3	2	comment-reply		\N	1	11	2019-12-18 11:47:58.737
185	11	3	2	post-reply		\N	1	11	2019-12-18 11:57:04.892
193	192	10	7	post-reply		\N	1	192	2019-12-18 13:41:25.692
194	186	1	7	post-reply	I love this cat!	\N	1	186	2019-12-18 13:41:38.991
195	194	1	2	comment-reply	You're wrong!	\N	1	186	2019-12-18 13:42:36.508
196	192	10	2	post-reply		\N	1	192	2019-12-18 13:43:21.68
20	\N	1	3	Cats suck, Sheep are better	SHEEEEEEEEP	http://upload.wikimedia.org/wikipedia/commons/c/c4/Lleyn_sheep.jpg	2	\N	2019-12-12 15:59:43.811
187	\N	3	2	HIDE THE PAIN	AAAAAAAAHHHHH	https://cdn0.vox-cdn.com/thumbor/kUCRDpV6jltsX_Hy6m2aCu5Flm4=/0x529:1267x1374/1310x873/cdn0.vox-cdn.com/uploads/chorus_image/image/49580341/h1.0.0.jpg	1	\N	2019-12-18 12:49:26.747
186	\N	1	2	Chonky	chonky chonker	https://static.boredpanda.com/blog/wp-content/uploads/2019/10/cinderblock-fat-cat-workout-1-5db6a2874218f__700.jpg	2	\N	2019-12-18 12:48:12.765
192	\N	10	7	SO EXCITED!!!!!!111!!!	spoilers n stuff	https://cnet3.cbsistatic.com/img/oYUaiXftBwdqY2bBsVKlwGKvMCM=/1200x675/2019/08/26/1612c4aa-32ab-48af-8532-17936f310691/rey-new-double-lightsaber.jpg	1	\N	2019-12-18 13:38:32.134
247	186	1	2	post-reply	test	\N	1	186	2019-12-19 13:34:43.855
248	\N	3	2	Pillow	sweet pillow	https://i.pinimg.com/originals/d0/09/25/d009254378a88cada108c9593d8dbcb1.jpg	1	\N	2019-12-19 13:36:03.401
250	186	1	2	post-reply	test	\N	1	186	2019-12-19 13:40:32.706
251	186	1	2	post-reply	test	\N	1	186	2019-12-19 13:40:35.618
252	186	1	2	post-reply	TEEEST	\N	1	186	2019-12-19 13:42:32.186
253	186	1	2	post-reply	test	\N	1	186	2019-12-19 13:48:02.508
254	194	1	2	comment-reply	TEST	\N	1	186	2019-12-19 13:48:13.583
255	250	1	2	comment-reply	AAAAAAAAAAA	\N	1	186	2019-12-19 13:49:36.705
256	186	1	2	post-reply	BBBBBBB	\N	1	186	2019-12-19 13:49:45.493
257	186	1	2	post-reply	CCCCC	\N	1	186	2019-12-19 13:50:30.746
258	186	1	2	post-reply	DDDDDD	\N	1	186	2019-12-19 13:53:41.311
259	197	10	2	post-reply	TEST	\N	1	197	2019-12-19 13:53:51.512
262	\N	10	2	Dew It	DEW IT	https://cdn.deseretnews.com/images/article/hires/1836806/1836806.jpg	1	\N	2019-12-19 13:55:40.324
263	262	10	2	post-reply	JUST DEW IT	\N	1	262	2019-12-19 13:55:57.447
264	263	10	2	comment-reply	DEWING IT	\N	1	262	2019-12-19 13:56:04.033
197	\N	10	2	Han shot first	You all know it	https://s.hdnux.com/photos/01/06/74/54/18598686/7/920x920.jpg	2	\N	2019-12-18 14:51:31.069
267	\N	9	2	this is a test!	TESTING	https://i.imgur.com/MySsgmm.jpg	1	\N	2019-12-19 14:10:00.901
268	267	9	2	post-reply	test 	\N	1	267	2019-12-19 14:10:15.12
\.


--
-- Data for Name: subnoddits; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.subnoddits (sn_id, sn_name, sn_description) FROM stdin;
1	Cats	A home for cats and cat accessories
2	Dogs	Dogs are not as cool as cats, but sometimes we like them too
3	Harold	Our inspiration to hiding the pain in our lives
4	Politics	Don't post here
5	Gardening	Plant stuff in the ground and hope it grows
9	test_subnoddit	this is a test
10	star_wars	All the Jedi things!
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, password, salt, role, avatar_address, first_name, last_name, email_address, join_date) FROM stdin;
1	rgodwin	LR2la7Lh7VWlf26I94NSng==	d+M5scpkXcjNVEc86Ogsy9LOaYgKPaLwaSu4a+6fd8hSmuV00+KA7Oks0jSpxqDxgcTAxMuIhw0zR4QQ4u/keI6ewye+yGLOO6TikvJ5OzHL8CoGmbdyaxFiIIpYoV3nUTMMNuabu5A0DFNnSw9Bfw+6p/avKIOuGIV0UTEdl5g=	super_admin	https://i.imgur.com/Au7tMN8.jpg	Rob	Goblin	rob@robgoblin.com	2019-12-09
5	test	JpKoqORob92om5JDo9zk9w==	jR78v9MHNX721r+KGFcOMDn1Yiq/BAyzMph57Cmr2IxR7EhFbS27xxVMJKLhUJNCy7/PVgHk596mNxWDYAq9dzzVoVq6jRhbW6Zv4+uNe2bv50ZA1DSHq/v9apeF7YRhm9Km2KyKaYdmuqaJtEZ2ZFhRmn5/KgNOGZfZbr+pAdg=	user	\N	\N	\N	\N	\N
6	asd1	Ox4VXrivKmyQV4zytFt4jQ==	ezHmf1EncbnQ7eWPVTOVRKD3VgXAJpYkMQuRkJQchNPcMexIvJ36UZFZpIK63QUdRMjYIv1nkn43EJYEvXW1AQAkP/EZQTDPhZxDENbEfPI8SxpkO9eJVK2Ydx/EyzoQKvgJVG4OVV6vKSwUbJmC2U9s2h1k0eL66byKc5EbLoA=	user	\N	\N	\N	\N	\N
7	david	OrV3rOe2/2FoFYMxpB+EYA==	CpdfwaDXO5g2JzhKDCNfDd4v8Vz2JuQaEDMJKJPaHuzCYE2GEdu9k4CP514TkMEr9+PkK6n3h93olfT71FuDRRa2n/4gFLQksicgTt9vqG0XLbaRjXL2P7+ApKBXCCbLg9iyDqxOYXDYzafn5nXsAnPPcVgjbcRdp9kSQ2tTnlY=	user	\N	\N	\N	\N	\N
3	eknutson	wuzqd9amUC/UYyF7fOML4Q==	YQxeIMEC8ww2dbZH51h2Zlq87n2cdq7BRFyDt7tW6pEI0pknGOlasolljlOCKOVbOU3KiWhkLBDx9kYvGfTYjX2ABO12aBya2dglN6bBUi5Lw1+6RttUTNMkmuU/AX0kGA2ACnSzf9+kDs2TVMoMbI93/av9mIq4xHAuNeG2awc=	super_admin	\N	\N	\N	\N	\N
4	jminihan	dAz8u6SGzCMw69Q9l9ZuDg==	5NUKnNW9n4u28RYcW5GlD7dRTJFLO10CrzWQ068jvPhTjh0LguiZTC7bQ/gZmQMjhVrdXszGCgz16iivutPbNLRBRx/e3d7r5LRMLIni+XS113F31Vl9rgskJBJljYWr6O13tDlsrIbeWis38rc8HP1Sos9I/U7jGmnH2Bq2Wac=	super_admin	\N	\N	\N	\N	\N
2	csamad	LtDiqdc4kgqaHRTW7fL+vQ==	t3Aw4VGyaWcaqsqt2nINmgt/HZegjO/QcVBmMb9hWZn2l6MfWVSsajMTDD8Z7h+grLRmPIBtt6Gk2KR4DD2/6AVwPBLsjA4yj9Zw+7KXIVhP3GS5LU5J4pM9yQ4EwiB3RruBdjKZ0zUQI7Ef+pFDcwpOF6JFSuVsI67/15xWyZA=	super_admin	\N	\N	\N	\N	\N
\.


--
-- Name: posts_post_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.posts_post_id_seq', 269, true);


--
-- Name: subnoddits_sn_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.subnoddits_sn_id_seq', 10, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 10, true);


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (post_id);


--
-- Name: subnoddits sn_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subnoddits
    ADD CONSTRAINT sn_name UNIQUE (sn_name) INCLUDE (sn_name);


--
-- Name: mod subnoddit user link; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mod
    ADD CONSTRAINT "subnoddit user link" PRIMARY KEY (sn_id, user_id);


--
-- Name: subnoddits subnoddits_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subnoddits
    ADD CONSTRAINT subnoddits_pkey PRIMARY KEY (sn_id);


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
-- Name: favorites post_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT post_id FOREIGN KEY (post_id) REFERENCES public.posts(post_id);


--
-- Name: post_votes post_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post_votes
    ADD CONSTRAINT post_id FOREIGN KEY (post_id) REFERENCES public.posts(post_id) NOT VALID;


--
-- Name: posts sn_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT sn_id FOREIGN KEY (sn_id) REFERENCES public.subnoddits(sn_id);


--
-- Name: favorites sn_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT sn_id FOREIGN KEY (sn_id) REFERENCES public.subnoddits(sn_id);


--
-- Name: posts user_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: favorites user_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: post_votes user_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.post_votes
    ADD CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES public.users(id) NOT VALID;


--
-- PostgreSQL database dump complete
--

