package com.techelevator.controller;

import com.techelevator.authentication.AuthProvider;
import com.techelevator.authentication.UnauthorizedException;
import com.techelevator.model.Favorites;
import com.techelevator.model.Moderator;
import com.techelevator.model.Post;
import com.techelevator.model.Subnoddit;
import com.techelevator.model.SubnodditDao;
import com.techelevator.model.User;
import com.techelevator.model.UserDao;
import com.techelevator.model.Vote;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

/**
 * ApiController
 */
@CrossOrigin
@RestController
@RequestMapping("/api")
public class ApiController {
	
	private SubnodditDao subnodditDao;
	private UserDao userDao;

    @Autowired
    private AuthProvider authProvider;
    
    @Autowired
    private ApiController(SubnodditDao subnodditDao, UserDao userDao) {
    	this.subnodditDao = subnodditDao;
    	this.userDao = userDao;
    }

    @RequestMapping(path = "/", method = RequestMethod.GET)
    public String authorizedOnly() throws UnauthorizedException {
        /*
        You can lock down which roles are allowed by checking
        if the current user has a role.
        
        In this example, if the user does not have the admin role
        we send back an unauthorized error.
        */
        if (!authProvider.userHasRole(new String[] { "admin" })) {
            throw new UnauthorizedException();
        }
        return "Success";
    }
    
    @PostMapping("/subnoddits/create")
    public void createSubnoddit(@RequestBody Subnoddit subnoddit) {
    	String[] splitString = subnoddit.getSubnodditName().split(" ");
    	String subnodditString = "";
		for (int i=0; i < splitString.length; i++) {
			if (i < splitString.length-1) {
				subnodditString += splitString[i] + "_";
			} else {
				subnodditString += splitString[i];
			}
		}
		subnoddit.setSubnodditName(subnodditString);
    	long userId = userDao.getUserId(subnoddit.getUsername());
    	subnodditDao.createSubnoddit(userId, subnoddit);
    	Subnoddit newSubnoddit = subnodditDao.getSubnodditId(subnoddit.getSubnodditName());
    	subnodditDao.createModerator(userId, newSubnoddit.getSubnodditId());
    }
    
    @PutMapping("/subnoddits/update") 
    public Subnoddit updateSubnoddit(@RequestBody Subnoddit subnoddit) {
    	return subnodditDao.updateSubnoddit(subnoddit);
    }
    
    @GetMapping("/public/subnoddits")
	public List<Subnoddit> getAllSubnoddits() {
		return subnodditDao.getAllSubnoddits();
	}
    
    @GetMapping("/public/subnoddits/active")
    public List<Subnoddit> getActiveSubnoddits() {
    	return subnodditDao.getActiveSubnoddits();
    }
    
    @GetMapping("/public/subnoddits/{subnodditName}")
    public Subnoddit getSubnodditByName(@PathVariable String subnodditName) {
    	return subnodditDao.getSubnodditByName(subnodditName);
    }
    
    @GetMapping("/public/{subnodditName}/{postId}") 
    public Post getPostById(@PathVariable String subnodditName, @PathVariable long postId) {
    	Post post = subnodditDao.getPostById(postId);
    	return post;
    }
    
    @PostMapping("/{subnodditName}/{postId}/createreply")
    public void createCommentReply(@RequestBody Post post) {
    	User user = userDao.getUserByUsername(post.getUsername());
    	post.setUserId(user.getId());
    	if (post.getTopLevelId() == post.getParentPostId()) {
    		subnodditDao.createPostReply(post);
    	}
    	else {
    		subnodditDao.createCommentReply(post);
    	}
    }
    
    @GetMapping("/public/{subnodditName}/{postId}/replies") 
    public List<Post> getRepliesForPost(@PathVariable String subnodditName, @PathVariable long postId) {
    	return subnodditDao.getRepliesForPost(postId);
    }
    
    @GetMapping("/public/allposts/{subnodditName}")
    public List<Post> getAllPostsForSubnoddit(@PathVariable String subnodditName) {
    	return subnodditDao.getPostsForSubnoddit(subnodditName);
    }
    
    @GetMapping("/public/allpostspopular/{subnodditName}")
    public List<Post> getAllPostsForSubnodditPopular(@PathVariable String subnodditName) {
    	return subnodditDao.getPostsForSubnodditPopular(subnodditName);
    }
    
    @GetMapping("/public/subnoddits/search/{searchTerm}")
    public List<Subnoddit> getSubnodditBySearch(@PathVariable String searchTerm) {
    	return subnodditDao.getSubnodditsBySearch(searchTerm);
    }
    
    @DeleteMapping("/subnoddits/delete/{subnodditName}")
    public void deleteSubnoddit(@PathVariable String subnodditName) {
    	subnodditDao.deleteSubnoddit(subnodditName);
    }
    
    @PostMapping("/post/create")
    public Post createPost(@RequestBody Post post) {
    	Long userId = userDao.getUserId(post.getUsername());
    	return subnodditDao.createPost(post, userId);
    }
    
    @PutMapping("/post/update/{postId}")
    public Post updatePost(@PathVariable int postId, @RequestBody Post post) {
    	return subnodditDao.updatePost(postId, post);
    }
    
    @PostMapping("/post/vote")
    public void addVoteForPost(@RequestBody Vote vote) {
    	Long userId = userDao.getUserId(vote.getUsername());
    	subnodditDao.addVoteForPost(vote, userId);
    }
    
    @GetMapping("/public/allposts")
    public List<Post> getAllPosts() {
    	return subnodditDao.getAllPosts();
    }
    
    @GetMapping("/public/popularposts")
    public List<Post> getPopularPosts() {
    	return subnodditDao.getTodaysPopularPosts();
    }
    
    @GetMapping("/post/votes/{postId}")
    public List<Vote> getVotesForPost(@PathVariable long postId) {
    	return subnodditDao.getVotesForPost(postId);
    }
    
    @GetMapping("/public/post/user/{username}")
    public List<Post> getPostsByUserId(@PathVariable String username) {
    	long id = userDao.getUserId(username);
    	return subnodditDao.getPostsByUserId(id);
    }
    
    @GetMapping("/public/recentposts")
    public List<Post> getMostRecentPosts() {
    	return subnodditDao.getMostRecentPosts();
    }
    
    @DeleteMapping("/post/delete/{postId}")
    public void deletePost(@PathVariable int postId) {
    	subnodditDao.deletePost(postId);
    }
    
    @PostMapping("/favorites/create/post")
    public Favorites createFavoritePost(@RequestBody Favorites favorite) {
    	favorite.setUserId(userDao.getUserId(favorite.getUsername()));
    	return subnodditDao.createFavoritePost(favorite);
    }
    
    @PostMapping("/favorites/create/subnoddit")
    public Favorites createFavoriteSubnoddit(@RequestBody Favorites favorite) {
    	favorite.setUserId(userDao.getUserId(favorite.getUsername()));
    	Subnoddit newSubnoddit = subnodditDao.getSubnodditId(favorite.getSubnodditName());
    	favorite.setSubnodditId(newSubnoddit.getSubnodditId());
    	return subnodditDao.createFavoriteSubnoddit(favorite);
    }  
    
    @GetMapping("/favorites/{username}")
    public List<Favorites> getFavoriesForUser(@PathVariable String username) {
    	long id = userDao.getUserId(username);
    	return subnodditDao.getFavoriesForUser(id);
    }
    
    @DeleteMapping("/favorites/delete/post/{postId}")
    public void deleteFavoritePost(@RequestBody Favorites favorite) {
    	subnodditDao.deleteFavoritePost(favorite);
    }
    
    @DeleteMapping("/favorites/subnoddit/{subnodditId}")
    public void deleteFavoriteSubnoddit(@PathVariable long subnodditId) {
    	Favorites favorite = new Favorites();
    	favorite.setSubnodditId(subnodditId);
    	favorite.setUserId(userDao.getUserId(authProvider.getCurrentUser().getUsername()));
    	subnodditDao.deleteFavoriteSubnoddit(favorite);
    }
    
    @GetMapping("/user/{username}")
    public User getInfoForUser(@PathVariable String username) {
    	return userDao.getUserByUsername(username);
    }
    
    @PutMapping("/user/update/email/{username}")
    public void updateUserEmail(@PathVariable String username, @RequestBody String emailAddress) {
    	long userId = userDao.getUserId(username);
    	userDao.updateUserEmail(emailAddress, userId);
    }
    
    @PutMapping("/user/update/username/{username}")
    public void updateUsername(@PathVariable String username, @RequestBody User user) {
    	long userId = userDao.getUserId(username);
    	user.setId(userId);
    	userDao.updateUsername(user);
    }
    
    @PutMapping("/user/update/name/{username}")
    public User updateUserRealName(@PathVariable String username, @RequestBody User user) {
    	return userDao.updateUserRealName(username, user);
    }
    
    @PutMapping("/user/update/avatar/{username}")
    public User updateUserAvatar(@PathVariable String username, @RequestBody User user) {
    	return userDao.updateUserAvatar(username, user);
    }
    
    @DeleteMapping("/user/delete/{username}")
    public void deleteUser(@PathVariable String username) throws UnauthorizedException {
    	if (!authProvider.userHasRole(new String[] { "admin", "super_admin" })) {
            throw new UnauthorizedException();
        }
    	userDao.deleteUser(username);
    }
    
    @GetMapping("/public/moderators/{subnodditName}")
    public List<Moderator> getModeratorsForSubnoddit(@PathVariable String subnodditName) {
    	Subnoddit newSubnoddit = subnodditDao.getSubnodditId(subnodditName);
    	long id = newSubnoddit.getSubnodditId();
    	return subnodditDao.getModeratorsForSubnoddit(id);
    }

}