-- Seed data from original Noddit application
-- This includes sample users, subnoddits, posts, comments, and votes

-- Note: Passwords are already hashed with PBKDF2
-- All sample user passwords are their usernames (e.g., user 'rgodwin' password is 'rgodwin')

-- Insert sample users
INSERT INTO users (id, username, password, salt, role, avatar_address, first_name, last_name, email_address, join_date) VALUES
(1, 'rgodwin', 'LR2la7Lh7VWlf26I94NSng==', 'd+M5scpkXcjNVEc86Ogsy9LOaYgKPaLwaSu4a+6fd8hSmuV00+KA7Oks0jSpxqDxgcTAxMuIhw0zR4QQ4u/keI6ewye+yGLOO6TikvJ5OzHL8CoGmbdyaxFiIIpYoV3nUTMMNuabu5A0DFNnSw9Bfw+6p/avKIOuGIV0UTEdl5g=', 'super_admin', 'https://i.imgur.com/Au7tMN8.jpg', 'Rob', 'Goblin', 'rob@robgoblin.com', '2019-12-09'),
(2, 'csamad', 'LtDiqdc4kgqaHRTW7fL+vQ==', 't3Aw4VGyaWcaqsqt2nINmgt/HZegjO/QcVBmMb9hWZn2l6MfWVSsajMTDD8Z7h+grLRmPIBtt6Gk2KR4DD2/6AVwPBLsjA4yj9Zw+7KXIVhP3GS5LU5J4pM9yQ4EwiB3RruBdjKZ0zUQI7Ef+pFDcwpOF6JFSuVsI67/15xWyZA=', 'super_admin', NULL, NULL, NULL, NULL, NULL),
(3, 'eknutson', 'wuzqd9amUC/UYyF7fOML4Q==', 'YQxeIMEC8ww2dbZH51h2Zlq87n2cdq7BRFyDt7tW6pEI0pknGOlasolljlOCKOVbOU3KiWhkLBDx9kYvGfTYjX2ABO12aBya2dglN6bBUi5Lw1+6RttUTNMkmuU/AX0kGA2ACnSzf9+kDs2TVMoMbI93/av9mIq4xHAuNeG2awc=', 'super_admin', NULL, NULL, NULL, NULL, NULL),
(4, 'jminihan', 'dAz8u6SGzCMw69Q9l9ZuDg==', '5NUKnNW9n4u28RYcW5GlD7dRTJFLO10CrzWQ068jvPhTjh0LguiZTC7bQ/gZmQMjhVrdXszGCgz16iivutPbNLRBRx/e3d7r5LRMLIni+XS113F31Vl9rgskJBJljYWr6O13tDlsrIbeWis38rc8HP1Sos9I/U7jGmnH2Bq2Wac=', 'super_admin', NULL, NULL, NULL, NULL, NULL),
(5, 'test', 'JpKoqORob92om5JDo9zk9w==', 'jR78v9MHNX721r+KGFcOMDn1Yiq/BAyzMph57Cmr2IxR7EhFbS27xxVMJKLhUJNCy7/PVgHk596mNxWDYAq9dzzVoVq6jRhbW6Zv4+uNe2bv50ZA1DSHq/v9apeF7YRhm9Km2KyKaYdmuqaJtEZ2ZFhRmn5/KgNOGZfZbr+pAdg=', 'user', NULL, NULL, NULL, NULL, NULL),
(6, 'asd1', 'Ox4VXrivKmyQV4zytFt4jQ==', 'ezHmf1EncbnQ7eWPVTOVRKD3VgXAJpYkMQuRkJQchNPcMexIvJ36UZFZpIK63QUdRMjYIv1nkn43EJYEvXW1AQAkP/EZQTDPhZxDENbEfPI8SxpkO9eJVK2Ydx/EyzoQKvgJVG4OVV6vKSwUbJmC2U9s2h1k0eL66byKc5EbLoA=', 'user', NULL, NULL, NULL, NULL, NULL),
(7, 'david', 'OrV3rOe2/2FoFYMxpB+EYA==', 'CpdfwaDXO5g2JzhKDCNfDd4v8Vz2JuQaEDMJKJPaHuzCYE2GEdu9k4CP514TkMEr9+PkK6n3h93olfT71FuDRRa2n/4gFLQksicgTt9vqG0XLbaRjXL2P7+ApKBXCCbLg9iyDqxOYXDYzafn5nXsAnPPcVgjbcRdp9kSQ2tTnlY=', 'user', NULL, NULL, NULL, NULL, NULL)
ON CONFLICT (username) DO NOTHING;

-- Update user ID sequence
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users), true);

-- Insert subnoddits (communities)
INSERT INTO subnoddits (sn_id, sn_name, sn_description) VALUES
(1, 'Cats', 'A home for cats and cat accessories'),
(2, 'Dogs', 'Dogs are not as cool as cats, but sometimes we like them too'),
(3, 'Harold', 'Our inspiration to hiding the pain in our lives'),
(4, 'Politics', 'Don''t post here'),
(5, 'Gardening', 'Plant stuff in the ground and hope it grows'),
(9, 'test_subnoddit', 'this is a test'),
(10, 'star_wars', 'All the Jedi things!')
ON CONFLICT (sn_name) DO NOTHING;

-- Update subnoddit ID sequence
SELECT setval('subnoddits_sn_id_seq', (SELECT MAX(sn_id) FROM subnoddits), true);

-- Insert sample posts (with original 2019 image URLs and recent timestamps)
INSERT INTO posts (post_id, parent_post_id, sn_id, user_id, title, body, image_address, post_score, top_level_id, created_date) VALUES
(1, NULL, 3, 1, 'Hide the Pain with Harold', 'He can do it an so can you', 'https://imgur.com/sM8MtJn.jpg', 2, NULL, NOW() - INTERVAL '12 hours'),
(2, NULL, 1, 1, 'This is Evie', 'She is a great kitty', 'https://imgur.com/nOQfW6Z.jpg', 1, NULL, NOW() - INTERVAL '11 hours'),
(4, NULL, 2, 1, 'Dogs are bad', 'Here is a cat instead', 'https://imgur.com/erIPkF5.jpg', 1, NULL, NOW() - INTERVAL '10 hours'),
(5, NULL, 5, 1, 'I grew this', 'Look at this sweet plant that I grew', 'https://imgur.com/gb8kBNe.jpg', 1, NULL, NOW() - INTERVAL '9 hours'),
(6, 2, 1, 1, 'Wow', 'your cat looks so cool', 'https://i.imgur.com/fqd9uUj.jpg', 4, 2, NOW() - INTERVAL '8 hours'),
(11, NULL, 3, 2, 'More Harold', 'Harold', 'https://brobible.files.wordpress.com/2019/11/istock-153696622.jpg', 1, NULL, NOW() - INTERVAL '7 hours'),
(16, NULL, 3, 3, 'I am the Big Sad.', 'Big sad!', 'https://assets.change.org/photos/9/io/nc/faiONCulPSIOaSe-800x450-noPad.jpg?1525707064', 2, NULL, NOW() - INTERVAL '6 hours'),
(17, NULL, 1, 2, 'This is a cat', 'CAT CAT', 'https://www.rd.com/wp-content/uploads/2019/11/cat-10-e1573844975155-1024x692.jpg', 0, NULL, NOW() - INTERVAL '5 hours'),
(18, NULL, 2, 2, 'Dog making cash', 'Tiny dog BIG money', 'https://images.dog.ceo/breeds/shiba/shiba-11.jpg', 99, NULL, NOW() - INTERVAL '4 hours'),
(19, 6, 1, 2, 'test comment', 'I agree!', NULL, 2, 2, NOW() - INTERVAL '3 hours'),
(20, NULL, 1, 3, 'Cats suck, Sheep are better', 'SHEEEEEEEEP', 'http://upload.wikimedia.org/wikipedia/commons/c/c4/Lleyn_sheep.jpg', 2, NULL, NOW() - INTERVAL '2 hours'),
(186, NULL, 1, 2, 'Chonky', 'chonky chonker', 'https://static.boredpanda.com/blog/wp-content/uploads/2019/10/cinderblock-fat-cat-workout-1-5db6a2874218f__700.jpg', 2, NULL, NOW() - INTERVAL '90 minutes'),
(187, NULL, 3, 2, 'HIDE THE PAIN', 'AAAAAAAAHHHHH', 'https://cdn0.vox-cdn.com/thumbor/kUCRDpV6jltsX_Hy6m2aCu5Flm4=/0x529:1267x1374/1310x873/cdn0.vox-cdn.com/uploads/chorus_image/image/49580341/h1.0.0.jpg', 1, NULL, NOW() - INTERVAL '60 minutes'),
(192, NULL, 10, 7, 'SO EXCITED!!!!!!111!!!', 'spoilers n stuff', 'https://cnet3.cbsistatic.com/img/oYUaiXftBwdqY2bBsVKlwGKvMCM=/1200x675/2019/08/26/1612c4aa-32ab-48af-8532-17936f310691/rey-new-double-lightsaber.jpg', 1, NULL, NOW() - INTERVAL '30 minutes'),
(197, NULL, 10, 2, 'Han shot first', 'You all know it', 'https://s.hdnux.com/photos/01/06/74/54/18598686/7/920x920.jpg', 2, NULL, NOW() - INTERVAL '15 minutes')
ON CONFLICT (post_id) DO NOTHING;

-- Update posts ID sequence
SELECT setval('posts_post_id_seq', (SELECT MAX(post_id) FROM posts), true);

-- Insert moderators (only for existing subnoddits)
INSERT INTO mod (sn_id, user_id) VALUES
(9, 2),
(10, 7)
ON CONFLICT DO NOTHING;

-- Insert sample votes (only for existing posts)
INSERT INTO post_votes (post_id, user_id, vote) VALUES
(20, 1, 'upvote'),
(18, 1, 'downvote'),
(6, 1, 'upvote'),
(19, 1, 'upvote'),
(17, 1, 'upvote'),
(1, 1, 'upvote'),
(16, 1, 'upvote'),
(17, 2, 'upvote'),
(186, 7, 'upvote'),
(192, 7, 'upvote'),
(192, 2, 'downvote'),
(6, 2, 'upvote'),
(197, 2, 'upvote')
ON CONFLICT DO NOTHING;

-- Insert favorites
INSERT INTO favorites (user_id, sn_id, post_id) VALUES
(1, NULL, 2),
(4, 1, NULL),
(5, 1, NULL),
(1, 2, NULL),
(1, 1, NULL),
(1, 3, NULL),
(6, 1, NULL),
(7, 1, NULL),
(7, NULL, 186),
(7, 5, NULL),
(7, 10, NULL),
(2, 3, NULL),
(2, 2, NULL);

-- Grant permissions (if needed)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;
