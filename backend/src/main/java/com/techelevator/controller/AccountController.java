package com.techelevator.controller;

import javax.validation.Valid;

import com.techelevator.authentication.AuthProvider;
import com.techelevator.authentication.JwtTokenHandler;
import com.techelevator.authentication.RegistrationResult;
import com.techelevator.authentication.UnauthorizedException;
import com.techelevator.model.SubnodditDao;
import com.techelevator.model.User;
import com.techelevator.model.UserDao;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.ObjectError;
import org.springframework.web.bind.annotation.*;

/**
 * AccountController
 */
@CrossOrigin
@RestController
public class AccountController {
	
	private SubnodditDao subnodditDao;
	private UserDao userDao;
	
	@Autowired
    private AccountController(SubnodditDao subnodditDao, UserDao userDao) {
    	this.subnodditDao = subnodditDao;
    	this.userDao = userDao;
    }
	
    @Autowired
    private AuthProvider auth;

    @Autowired
    private JwtTokenHandler tokenHandler;

    @PostMapping("/login")
    public String login(@RequestBody User user) throws UnauthorizedException {
        if(auth.signIn(user.getUsername(), user.getPassword())) {
            User currentUser = auth.getCurrentUser();
            return tokenHandler.createToken(user.getUsername(), currentUser.getRole());
        } else {
            throw new UnauthorizedException();
        }
    }



    @PostMapping("/register")
    public RegistrationResult register(@Valid @RequestBody User user, BindingResult result) {
    	user.setUsername(user.getUsername().toLowerCase());
    	RegistrationResult registrationResult = new RegistrationResult();
    	if(result.hasErrors()) {
            for(ObjectError error : result.getAllErrors()) {
                registrationResult.addError(error.getDefaultMessage());
            }
        }
    	else {
    		auth.register(user.getUsername(), user.getPassword(), user.getRole());
    		long id = userDao.getUserId(user.getUsername());
    		subnodditDao.newUserDefaultFavorite(id);
    		registrationResult.setSuccess(true);
    	}
    	return registrationResult;
    }

}