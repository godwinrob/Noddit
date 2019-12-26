package com.techelevator.model;

import java.util.List;

public interface SubnodditDao {
	
	public List<Subnoddit> getAllSubnoddits();
	
	public List<Subnoddit> getActiveSubnoddits();
	
	public List<Subnoddit> getSubnodditsBySearch(String search);
	
	public Subnoddit updateSubnoddit(Subnoddit subnoddit);
	
	public Subnoddit getSubnodditByName(String subnodditName);

	public Subnoddit getSubnodditId(String subnodditName);
	
	public List<Post> getAllPosts();
	
	public List<Post> getMostRecentPosts();
	
	public List<Post> getRepliesForPost(long postId);
	
	public List<Post> getPostsByUserId(long userId);
	
	public List<Post> getPostsForSubnoddit(String subnodditName);
	
	public List<Post> getPostsForSubnodditPopular(String subnodditName);
	
	public List<Post> getTodaysPopularPosts();	
	
	public Post createPost(Post post, long userId);
	
	public Post updatePost(long postId, Post post);
	
	public Post getPostById(long postId);
	
	public List<Favorites> getFavoriesForUser(long userId);
	
	public Favorites createFavoritePost(Favorites favorite);
	
	public Favorites createFavoriteSubnoddit(Favorites favorite);
	
	public List<Moderator> getModeratorsForSubnoddit(long subnodditId);
	
	public List<Vote> getVotesForPost(long postId);
	
	public void createSubnoddit(long id, Subnoddit subnoddit);	
	
	public void deleteSubnoddit(String subnodditName);
	
	public void deletePost(long postId);

	public void createPostReply(Post post);
	
	public void deleteFavoriteSubnoddit(Favorites favorite);
	
	public void deleteFavoritePost(Favorites favorite);

	public void newUserDefaultFavorite(long userId);

	public void createModerator(long userId, long subnodditId);
	
	public void createCommentReply(Post post);
	
	public void addVoteForPost(Vote vote, long userId);
	
}
