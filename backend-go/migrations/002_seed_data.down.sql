-- Rollback seed data

DELETE FROM favorites;
DELETE FROM post_votes;
DELETE FROM mod;
DELETE FROM posts;
DELETE FROM subnoddits;
DELETE FROM users WHERE id <= 10;

-- Reset sequences
SELECT setval('users_id_seq', 1, false);
SELECT setval('subnoddits_sn_id_seq', 1, false);
SELECT setval('posts_post_id_seq', 1, false);
