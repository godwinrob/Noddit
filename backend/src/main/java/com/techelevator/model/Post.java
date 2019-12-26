package com.techelevator.model;

import java.sql.Timestamp;
import java.time.LocalDateTime;
import java.util.Date;

import org.springframework.format.annotation.DateTimeFormat;

import com.fasterxml.jackson.annotation.JsonFormat;

public class Post {
	
	private long postId;
	private long parentPostId;
	private long SubnodditId;
	private String subnodditName;
	private long userId;
	private String username;
	private String title;
	private String imageAddress;
	private String body;
	private long topLevelId;
	
	@JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "MM-dd-yyyy' 'HH:mm:ss")
	private Timestamp createdDate;
	
	private long postScore;
	
	public long getPostId() {
		return postId;
	}
	public void setPostId(long postId) {
		this.postId = postId;
	}
	public long getParentPostId() {
		return parentPostId;
	}
	public void setParentPostId(long parentPostId) {
		this.parentPostId = parentPostId;
	}
	public long getSubnodditId() {
		return SubnodditId;
	}
	public void setSubnodditId(long subnodditId) {
		SubnodditId = subnodditId;
	}
	public long getUserId() {
		return userId;
	}
	public void setUserId(long userId) {
		this.userId = userId;
	}
	public String getTitle() {
		return title;
	}
	public void setTitle(String title) {
		this.title = title;
	}
	public String getImageAddress() {
		return imageAddress;
	}
	public void setImageAddress(String imageAddress) {
		this.imageAddress = imageAddress;
	}
	public String getBody() {
		return body;
	}
	public void setBody(String body) {
		this.body = body;
	}
	public Timestamp getCreatedDate() {
		return createdDate;
	}
	public void setCreatedDate(Timestamp createdDate) {
		this.createdDate = createdDate;
	}
	public String getUsername() {
		return username;
	}
	public void setUsername(String username) {
		this.username = username;
	}
	public String getSubnodditName() {
		return subnodditName;
	}
	public void setSubnodditName(String subnodditName) {
		this.subnodditName = subnodditName;
	}
	public long getPostScore() {
		return postScore;
	}
	public void setPostScore(long postScore) {
		this.postScore = postScore;
	}
	public long getTopLevelId() {
		return topLevelId;
	}
	public void setTopLevelId(long topLevelId) {
		this.topLevelId = topLevelId;
	}

}
