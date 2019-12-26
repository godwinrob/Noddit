package com.techelevator.model;

import java.sql.Timestamp;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.List;

import javax.sql.DataSource;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.support.rowset.SqlRowSet;
import org.springframework.stereotype.Component;

import com.techelevator.authentication.PasswordHasher;

@Component
public class JdbcSubnodditDao implements SubnodditDao {

	private JdbcTemplate jdbcTemplate;

	@Autowired
	public JdbcSubnodditDao(DataSource dataSource, PasswordHasher passwordHasher) {
		this.jdbcTemplate = new JdbcTemplate(dataSource);
	}

	@Override
	public List<Subnoddit> getAllSubnoddits() {
		List<Subnoddit> subnodditList = new ArrayList<>();
		String sql = "SELECT * FROM subnoddits;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql);

		while (results.next()) {
			Subnoddit subnoddit = mapResultsToSubnoddit(results);
			subnodditList.add(subnoddit);
		}

		return subnodditList;
	}

	@Override
	public List<Post> getAllPosts() {
		List<Post> posts = new ArrayList<>();
		String sql = "SELECT * FROM posts;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql);

		while (results.next()) {
			Post post = mapResultsToPost(results);
			posts.add(post);
		}

		return posts;
	}

	@Override
	public List<Post> getMostRecentPosts() {
		List<Post> posts = new ArrayList<>();
		String sql = "SELECT  posts.post_id, posts.parent_post_id, posts.sn_id, subnoddits.sn_name, posts.user_id, users.username, posts.title, posts.body, posts.image_address, posts.created_date, posts.post_score, posts.top_level_id "
				+ "FROM posts JOIN users ON users.id = posts.user_id "
				+ "JOIN subnoddits ON subnoddits.sn_id = posts.sn_id " + "WHERE parent_post_id is null "
				+ "ORDER BY post_id DESC LIMIT 5;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql);

		while (results.next()) {
			Post post = mapResultsToPost(results);
			posts.add(post);
		}

		return posts;
	}
	
	@Override
	public List<Post> getTodaysPopularPosts() {
		List<Post> posts = new ArrayList<>();
		Instant instant = Instant.now().minus(24, ChronoUnit.HOURS);
		Timestamp timestamp = Timestamp.from(instant); 
		String sql = "SELECT posts.post_id, posts.parent_post_id, posts.sn_id, subnoddits.sn_name, posts.user_id, users.username, posts.title, posts.body, posts.image_address, posts.created_date, posts.post_score, posts.top_level_id "
				+ "FROM posts JOIN users ON users.id = posts.user_id "
				+ "JOIN subnoddits ON subnoddits.sn_id = posts.sn_id WHERE parent_post_id is null "
				+ "AND created_date >= ? ORDER BY post_score DESC, created_date DESC LIMIT 10;";
		
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, timestamp);
		
		while (results.next()) {
			Post post = mapResultsToPost(results);
			posts.add(post);
		}
		
		return posts;
	}

	@Override
	public List<Post> getRepliesForPost(long postId) {
		List<Post> posts = new ArrayList<>();
		String sql = "SELECT  posts.post_id, posts.parent_post_id, posts.sn_id, subnoddits.sn_name, posts.user_id, "
				+ " users.username, posts.title, posts.body, posts.image_address, posts.created_date, posts.post_score, posts.top_level_id "
				+ "FROM posts " + "JOIN users ON users.id = posts.user_id "
				+ "JOIN subnoddits ON subnoddits.sn_id = posts.sn_id "
				+ "WHERE top_level_id = ? ORDER BY post_score DESC;";

		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, postId);

		while (results.next()) {
			Post post = mapResultsToPost(results);
			posts.add(post);
		}

		return posts;
	}

	@Override
	public Post getPostById(long postId) {
		String sql = "SELECT  posts.post_id, posts.parent_post_id, posts.sn_id, subnoddits.sn_name, posts.user_id, users.username, posts.title, posts.body, posts.image_address, posts.created_date, posts.post_score, posts.top_level_id "
				+ "FROM posts " + "JOIN users ON users.id = posts.user_id "
				+ "JOIN subnoddits ON subnoddits.sn_id = posts.sn_id " + "WHERE post_id = ?;";

		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, postId);

		Post post = new Post();

		if (results.next()) {
			post = mapResultsToPost(results);
		}

		return post;
	}
	
	@Override
	public List<Post> getPostsByUserId(long userId) {
		List<Post> posts = new ArrayList<>();
		String sql = "SELECT  posts.post_id, posts.parent_post_id, posts.sn_id, subnoddits.sn_name, posts.user_id, users.username, posts.title, posts.body, posts.image_address, posts.created_date, posts.post_score, posts.top_level_id "
				+ "FROM posts " + "JOIN users ON users.id = posts.user_id "
				+ "JOIN subnoddits ON subnoddits.sn_id = posts.sn_id "
				+ "WHERE user_id = ? AND top_level_id IS NULL AND parent_post_id IS NULL ORDER BY post_id DESC;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, userId);

		while (results.next()) {
			Post post = mapResultsToPost(results);
			posts.add(post);
		}

		return posts;
	}


	@Override
	public List<Post> getPostsForSubnoddit(String subnodditName) {
		List<Post> posts = new ArrayList<>();
		String sql = "SELECT  posts.post_id, posts.parent_post_id, posts.sn_id, subnoddits.sn_name, posts.user_id, users.username, posts.title, posts.body, posts.image_address, posts.created_date, posts.post_score, posts.top_level_id "
				+ "FROM posts " + "JOIN users ON users.id = posts.user_id "
				+ "JOIN subnoddits ON subnoddits.sn_id = posts.sn_id "
				+ "WHERE sn_name = ? AND parent_post_id is null ORDER BY post_id DESC;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, subnodditName);

		while (results.next()) {
			Post post = mapResultsToPost(results);
			posts.add(post);
		}

		return posts;
	}
	
	@Override
	public List<Post> getPostsForSubnodditPopular(String subnodditName) {
		List<Post> posts = new ArrayList<>();
		String sql = "SELECT  posts.post_id, posts.parent_post_id, posts.sn_id, subnoddits.sn_name, posts.user_id, users.username, posts.title, posts.body, posts.image_address, posts.created_date, posts.post_score, posts.top_level_id "
				+ "FROM posts " + "JOIN users ON users.id = posts.user_id "
				+ "JOIN subnoddits ON subnoddits.sn_id = posts.sn_id "
				+ "WHERE sn_name = ? AND parent_post_id is null ORDER BY post_score DESC;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, subnodditName);

		while (results.next()) {
			Post post = mapResultsToPost(results);
			posts.add(post);
		}

		return posts;
	}

	@Override
	public Subnoddit getSubnodditId(String subnodditName) {
		String sql = "SELECT * FROM subnoddits WHERE sn_name = ?;"; // TODO fix select *
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, subnodditName);

		Subnoddit subnoddit = new Subnoddit();
		if (results.next()) {
			subnoddit = mapResultsToSubnoddit(results);
		}

		return subnoddit;
	}

	@Override
	public List<Favorites> getFavoriesForUser(long userId) {
		List<Favorites> favorites = new ArrayList<>();
		String sql = "SELECT * FROM favorites JOIN subnoddits ON subnoddits.sn_id = favorites.sn_id WHERE user_id = ? ORDER BY sn_name;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, userId);

		while (results.next()) {
			Favorites favorite = mapResultsToFavorites(results);
			favorites.add(favorite);
		}

		return favorites;
	}

	@Override
	public void createSubnoddit(long id, Subnoddit subnoddit) {
		String sql = "INSERT INTO subnoddits (sn_name, sn_description) VALUES (?, ?);";
		jdbcTemplate.update(sql, subnoddit.getSubnodditName(), subnoddit.getSubnodditDescription());
	}
	
	@Override
	public Subnoddit getSubnodditByName(String subnodditName) {
		Subnoddit subnoddit = new Subnoddit();
		String sql = "SELECT * FROM subnoddits WHERE sn_name = ?;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, subnodditName);
		
		if (results.next()) {
			subnoddit = mapResultsToSubnoddit(results);
		}
		
		return subnoddit;
	}
	
	@Override
	public List<Subnoddit> getSubnodditsBySearch(String search) {
		List<Subnoddit> searchResults = new ArrayList<>();
		String likeSearch = "%" + search + "%";
		String sql = "SELECT sn_id, sn_name, sn_description FROM subnoddits WHERE sn_name ILIKE ? OR sn_description ILIKE ?;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, likeSearch, likeSearch);
		
		while (results.next()) {
			Subnoddit subnoddit = mapResultsToSubnoddit(results);
			searchResults.add(subnoddit);
		}
		return searchResults;
	}

	@Override
	public Subnoddit updateSubnoddit(Subnoddit subnoddit) {
		String sql = "UPDATE subnoddits SET sn_description = ? WHERE sn_name = ?;";
		jdbcTemplate.update(sql, subnoddit.getSubnodditDescription(), subnoddit.getSubnodditName());
		return subnoddit;
	}

	@Override
	public void deleteSubnoddit(String subnodditName) {
		String sql = "DELETE FROM subnoddits WHERE sn_name = ?;";
		jdbcTemplate.update(sql, subnodditName);
	}

	@Override
	public Post createPost(Post post, long userId) {
		String sql = "INSERT INTO posts (sn_id, user_id, title, body, image_address, created_date, post_score) VALUES (?, ?, ?, ?, ?, ?, ?);";
		Timestamp timestamp = new Timestamp(System.currentTimeMillis());
		int postId = jdbcTemplate.update(sql, post.getSubnodditId(), userId, post.getTitle(), post.getBody(),
				post.getImageAddress(), timestamp, 1);
		post.setPostId((long) postId);
		return post;
	}

	@Override
	public void createPostReply(Post post) {
		String sql = "INSERT INTO posts (parent_post_id, sn_id, user_id, title, body, image_address, created_date, top_level_id, post_score) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);";
		Timestamp timestamp = new Timestamp(System.currentTimeMillis());
		jdbcTemplate.update(sql, post.getParentPostId(), post.getSubnodditId(), post.getUserId(),
				"post-reply", post.getBody(), null, timestamp, post.getTopLevelId(), 1);
	}
	
	@Override
	public void createCommentReply(Post post) {
		String sql = "INSERT INTO posts (parent_post_id, sn_id, user_id, title, body, image_address, created_date, top_level_id, post_score) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);";
		Timestamp timestamp = new Timestamp(System.currentTimeMillis());
		jdbcTemplate.update(sql, post.getParentPostId(), post.getSubnodditId(), post.getUserId(),
				"comment-reply", post.getBody(), null, timestamp, post.getTopLevelId(), 1);
	}

	@Override
	public Post updatePost(long postId, Post post) {
		String sql = "UPDATE posts SET body = ?, image_address = ? WHERE post_id = ?;";
		jdbcTemplate.update(sql, post.getBody(), post.getImageAddress(), postId);
		return post;
	}

	@Override
	public void deletePost(long postId) {
		String votes = "DELETE FROM post_votes WHERE post_id = ?;";
		String sql = "DELETE FROM posts WHERE post_id = ?;";
		jdbcTemplate.update(votes, postId);
		jdbcTemplate.update(sql, postId);
		
	}

	@Override
	public Favorites createFavoritePost(Favorites favorite) {
		String sql = "INSERT INTO favorites (user_id, post_id) SELECT ?, ? WHERE NOT EXISTS ("
				+ "SELECT user_id, sn_id FROM favorites WHERE user_id = ? AND post_id = ?);";
		jdbcTemplate.update(sql, favorite.getUserId(), favorite.getPostId(), favorite.getUserId(), favorite.getPostId());
		return favorite;
	}

	@Override
	public Favorites createFavoriteSubnoddit(Favorites favorite) {
		String sql = "INSERT INTO favorites (user_id, sn_id) SELECT ?, ? WHERE NOT EXISTS ("
				+ "SELECT user_id, sn_id FROM favorites WHERE user_id = ? AND sn_id = ?);";
		jdbcTemplate.update(sql, favorite.getUserId(), favorite.getSubnodditId(), favorite.getUserId(), favorite.getSubnodditId());
		return favorite;
	}

	@Override
	public void newUserDefaultFavorite(long userId) {
		String sql = "INSERT INTO favorites (user_id, sn_id) VALUES (?, 1);";
		jdbcTemplate.update(sql, userId);
	}

	@Override
	public void deleteFavoriteSubnoddit(Favorites favorite) {
		String sql = "DELETE FROM favorites WHERE sn_id = ? AND user_id = ?;";
		jdbcTemplate.update(sql, favorite.getSubnodditId(), favorite.getUserId());
	}

	@Override
	public void deleteFavoritePost(Favorites favorite) {
		String sql = "DELETE FROM favorites WHERE post_id = ? AND user_id = ?;";
		jdbcTemplate.update(sql, favorite.getPostId(), favorite.getUserId());
	}
	
	@Override
	public void createModerator(long userId, long subnodditId) {
		String sql = "INSERT INTO mod (user_id, sn_id) VALUES (?, ?);";
		jdbcTemplate.update(sql, userId, subnodditId);
	}
	
	@Override
	public List<Moderator> getModeratorsForSubnoddit(long subnodditId) {
		List<Moderator> mods = new ArrayList<>();
		String sql = "SELECT mod.sn_id, mod.user_id, users.username "
				+ "FROM mod JOIN users on users.id = mod.user_id "
				+ "WHERE sn_id = ?;";
		
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, subnodditId);
		
		while (results.next()) {
			Moderator mod = mapResultsToModerator(results);
			mods.add(mod);
		}
		
		return mods;
	}
	
	@Override
	public List<Vote> getVotesForPost(long postId) {
		List<Vote> votes = new ArrayList<>();
		String sql = "SELECT post_id, user_id, users.username, vote FROM post_votes "
				+ "JOIN users ON users.id = post_votes.user_id "
				+ "WHERE post_id = ?;";
		
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql, postId);
		
		while (results.next()) {
			Vote vote = mapResultsToVote(results);
			votes.add(vote);
		}
		
		return votes;
	}

	@Override
	public void addVoteForPost(Vote vote, long userId) {
		String voteSql = "";
		
		if (vote.getVote().contains("upvote")) {
			voteSql = "UPDATE posts SET post_score = post_score + 1 WHERE post_id = ? "
					+ "AND post_id NOT IN (SELECT post_id FROM post_votes WHERE post_id = ? AND user_id = ?);";
		} else {
			voteSql = "UPDATE posts SET post_score = post_score - 1 WHERE post_id = ? "
					+ "AND post_id NOT IN (SELECT post_id FROM post_votes WHERE post_id = ? AND user_id = ?);";
		}
		jdbcTemplate.update(voteSql, vote.getPostId(), vote.getPostId(), userId);
		
		String sql = "INSERT INTO post_votes (post_id, user_id, vote) VALUES (?, ?, ?);";
		jdbcTemplate.update(sql, vote.getPostId(), userId, vote.getVote());
	}
	
	@Override
	public List<Subnoddit> getActiveSubnoddits() {
		List<Subnoddit> activeSubnoddits = new ArrayList<>();
		String sql = "SELECT DISTINCT subnoddits.sn_id, sn_name, sn_description, max(post_id) "
				+ "FROM subnoddits JOIN posts ON posts.sn_id = subnoddits.sn_id "
				+ "GROUP BY subnoddits.sn_id, sn_name ORDER BY max(post_id) DESC LIMIT 5;";
		SqlRowSet results = jdbcTemplate.queryForRowSet(sql);
		
		while (results.next()) {
			Subnoddit subnoddit = mapResultsToActiveSubnoddit(results);
			activeSubnoddits.add(subnoddit);
		}
		return activeSubnoddits;
	}
	

	private Subnoddit mapResultsToSubnoddit(SqlRowSet results) {
		Subnoddit subnoddit = new Subnoddit();
		subnoddit.setSubnodditId(results.getLong("sn_id"));
		subnoddit.setSubnodditName(results.getString("sn_name"));
		subnoddit.setSubnodditDescription(results.getString("sn_description"));
		return subnoddit;
	}
	
	private Subnoddit mapResultsToActiveSubnoddit(SqlRowSet results) {
		Subnoddit subnoddit = new Subnoddit();
		subnoddit.setSubnodditId(results.getLong("sn_id"));
		subnoddit.setSubnodditName(results.getString("sn_name"));
		subnoddit.setSubnodditDescription(results.getString("sn_description"));
		subnoddit.setPostId(results.getLong("max"));
		return subnoddit;
	}

	private Post mapResultsToPost(SqlRowSet results) {
		Post post = new Post();
		post.setPostId(results.getLong("post_id"));
		post.setParentPostId(results.getLong("parent_post_id")); // TODO not all posts will have this
		post.setSubnodditId(results.getLong("sn_id"));
		post.setSubnodditName(results.getString("sn_name"));
		post.setUserId(results.getLong("user_id"));
		post.setUsername(results.getString("username"));
		post.setTitle(results.getString("title"));
		post.setBody(results.getString("body"));
		post.setImageAddress(results.getString("image_address")); // TODO not all posts will have this
		post.setCreatedDate(results.getTimestamp("created_date"));
		post.setPostScore(results.getLong("post_score"));
		post.setTopLevelId(results.getLong("top_level_id"));

		return post;
	}

	private Favorites mapResultsToFavorites(SqlRowSet results) {
		Favorites favorites = new Favorites();
		favorites.setUserId(results.getLong("user_id"));
		favorites.setSubnodditId(results.getLong("sn_id"));
		favorites.setPostId(results.getLong("post_id"));
		favorites.setSubnodditName(results.getString("sn_name"));
		return favorites;
	}

	private Moderator mapResultsToModerator(SqlRowSet results) {
		Moderator moderator = new Moderator();
		moderator.setUserId(results.getLong("user_id"));
		moderator.setSubnodditId(results.getLong("sn_id"));
		moderator.setUsername(results.getString("username"));
		return moderator;
	}
	
	private Vote mapResultsToVote(SqlRowSet results) {
		Vote vote = new Vote();
		vote.setPostId(results.getLong("post_id"));
		vote.setUserId(results.getLong("user_id"));
		vote.setVote(results.getString("vote"));
		vote.setUsername(results.getString("username"));
		return vote;
	}

}
