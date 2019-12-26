package com.techelevator.model;

import java.util.ArrayList;
import java.util.Base64;
import java.util.List;

import javax.sql.DataSource;

import com.techelevator.authentication.PasswordHasher;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.support.rowset.SqlRowSet;
import org.springframework.stereotype.Component;

@Component
public class JdbcUserDao implements UserDao {

    private JdbcTemplate jdbcTemplate;
    private PasswordHasher passwordHasher;
    private SubnodditDao subnodditDao;

    /**
     * Create a new user dao with the supplied data source and the password hasher
     * that will salt and hash all the passwords for users.
     *
     * @param dataSource an SQL data source
     * @param passwordHasher an object to salt and hash passwords
     */
    @Autowired
    public JdbcUserDao(DataSource dataSource, PasswordHasher passwordHasher) {
        this.jdbcTemplate = new JdbcTemplate(dataSource);
        this.passwordHasher = passwordHasher;
    }

    /**
     * Save a new user to the database. The password that is passed in will be
     * salted and hashed before being saved. The original password is never
     * stored in the system. We will never have any idea what it is!
     *
     * @param userName the user name to give the new user
     * @param password the user's password
     * @param role the user's role
     * @return the new user
     */
    @Override
    public User saveUser(String userName, String password, String role) {
        byte[] salt = passwordHasher.generateRandomSalt();
        String hashedPassword = passwordHasher.computeHash(password, salt);
        String saltString = new String(Base64.getEncoder().encode(salt));
        long newId = jdbcTemplate.queryForObject(
                "INSERT INTO users(username, password, salt, role) VALUES (?, ?, ?, ?) RETURNING id", Long.class,
                userName, hashedPassword, saltString, role);
        
        User newUser = new User();
        newUser.setId(newId);
        newUser.setUsername(userName);
        newUser.setRole(role);

        return newUser;
    }

    @Override
    public void changePassword(User user, String newPassword) {
        byte[] salt = passwordHasher.generateRandomSalt();
        String hashedPassword = passwordHasher.computeHash(newPassword, salt);
        String saltString = new String(Base64.getEncoder().encode(salt));

        jdbcTemplate.update("UPDATE users SET password=?, salt=? WHERE id=?", hashedPassword, saltString, user.getId());
    }

    /**
     * Look for a user with the given username and password. Since we don't
     * know the password, we will have to get the user's salt from the database,
     * hash the password, and compare that against the hash in the database.
     *
     * @param userName the user name of the user we are checking
     * @param password the password of the user we are checking
     * @return true if the user is found and their password matches
     */
    @Override
    public User getValidUserWithPassword(String userName, String password) {
        String sqlSearchForUser = "SELECT * FROM users WHERE UPPER(username) = ?";

        SqlRowSet results = jdbcTemplate.queryForRowSet(sqlSearchForUser, userName.toUpperCase());
        if (results.next()) {
            String storedSalt = results.getString("salt");
            String storedPassword = results.getString("password");
            String hashedPassword = passwordHasher.computeHash(password, Base64.getDecoder().decode(storedSalt));
            if (storedPassword.equals(hashedPassword)) {
                return mapResultToUser(results);
            } else {
                return null;
            }
        } else {
            return null;
        }
    }

    /**
     * Get all of the users from the database.
     * @return a List of user objects
     */
    @Override
    public List<User> getAllUsers() {
        List<User> users = new ArrayList<User>();
        String sqlSelectAllUsers = "SELECT id, username, role, avatar_address, first_name, last_name, email_address, join_date FROM users";
        SqlRowSet results = jdbcTemplate.queryForRowSet(sqlSelectAllUsers);

        while (results.next()) {
            User user = mapResultToUser(results);
            users.add(user);
        }

        return users;
    }
    
    @Override
    public User getUserByUsername(String username) {
        String sqlSelectUserByUsername = "SELECT id, username, role, avatar_address, first_name, last_name, email_address, join_date FROM users WHERE username = ?";
        SqlRowSet results = jdbcTemplate.queryForRowSet(sqlSelectUserByUsername, username);

        if (results.next()) {
            return mapResultToUser(results);
        } else {
            return null;
        }
    }
    
    @Override
    public long getUserId(String username) {
    	String sql = "SELECT id FROM users WHERE username = ?;";
    	
    	SqlRowSet results = jdbcTemplate.queryForRowSet(sql, username);
    	
    	long test = 0;
    	
    	if(results.next()) {
    		return results.getLong("id");
    	}
    	
    	return test;
    	
    }

	@Override
	public void updateUserEmail(String emailAddress, long userId) {
		String sql = "UPDATE users SET email_address = ? WHERE id = ?;";
		jdbcTemplate.update(sql, emailAddress, userId);
	}
	
	@Override
	public User updateUserRealName(String username, User user) {
		String sql = "UPDATE users SET first_name = ?, last_name = ? WHERE username = ?;";
		jdbcTemplate.update(sql, user.getFirstName(), user.getLastName(), username);
		return user;
	}
	
	@Override
	public void updateUsername(User user) {
		String sql = "UPDATE users SET username = ? WHERE id = ?;";
		jdbcTemplate.update(sql, user.getNewUsername(), user.getId());
	}
	
	@Override
	public User updateUserAvatar(String username, User user) {
		String sql = "UPDATE users SET avatar_address = ? WHERE username = ?;";
		jdbcTemplate.update(sql, user.getAvatarAddress(), username);
		return user;
	}
	
	@Override
	public void deleteUser(String username) {
		String sql = "DELETE FROM users WHERE username = ?;";
		jdbcTemplate.update(sql);	
	}

    private User mapResultToUser(SqlRowSet results) {
        User user = new User();
        user.setId(results.getLong("id"));
        user.setUsername(results.getString("username"));
        user.setRole(results.getString("role"));
        user.setAvatarAddress(results.getString("avatar_address"));
        user.setFirstName(results.getString("first_name"));
        user.setLastName(results.getString("last_name"));
        user.setEmailAddress(results.getString("email_address"));
        user.setJoinDate(results.getDate("join_date"));
        return user;
    }
    

}
