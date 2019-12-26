package com.techelevator.model;

public class Favorites {
	
	private long userId;
	private long subnodditId;
	private long postId;
	private String subnodditName;
	private String username;
	
	public long getUserId() {
		return userId;
	}
	public void setUserId(long userId) {
		this.userId = userId;
	}
	public long getSubnodditId() {
		return subnodditId;
	}
	public void setSubnodditId(long subnodditId) {
		this.subnodditId = subnodditId;
	}
	public long getPostId() {
		return postId;
	}
	public void setPostId(long postId) {
		this.postId = postId;
	}
	public String getSubnodditName() {
		return subnodditName;
	}
	public void setSubnodditName(String subnodditName) {
		this.subnodditName = subnodditName;
	}
	public String getUsername() {
		return username;
	}
	public void setUsername(String username) {
		this.username = username;
	}

}
