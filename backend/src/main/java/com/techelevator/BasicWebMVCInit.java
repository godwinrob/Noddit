/*
 * After adding this file, no web.xml or springmvc-servlet.xml is necessary.
 * If Eclipse complains about a missing web.xml, add this section to the pom.xml file:
    <properties>
    	<failOnMissingWebXml>false</failOnMissingWebXml>
    </properties>
 */
package com.techelevator;

import java.util.ArrayList;
import java.util.List;

import javax.servlet.ServletContext;
import javax.servlet.ServletRegistration;

import org.apache.commons.dbcp2.BasicDataSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.WebApplicationInitializer;
import org.springframework.web.context.support.AnnotationConfigWebApplicationContext;
import org.springframework.web.servlet.DispatcherServlet;
import org.springframework.web.servlet.ViewResolver;
import org.springframework.web.servlet.config.annotation.DefaultServletHandlerConfigurer;
import org.springframework.web.servlet.config.annotation.EnableWebMvc;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.view.InternalResourceViewResolver;

import com.techelevator.authentication.JwtAuthInterceptor;

import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;



@Configuration
@EnableWebMvc
@ComponentScan
public class BasicWebMVCInit implements WebMvcConfigurer, WebApplicationInitializer{
	
	@Override
	public void onStartup(ServletContext servletContext) {			
		AnnotationConfigWebApplicationContext dispatcherContext = new AnnotationConfigWebApplicationContext();
		dispatcherContext.register(BasicWebMVCInit.class);
			
		ServletRegistration.Dynamic dispatcher = servletContext.addServlet("dispatcher", new DispatcherServlet(dispatcherContext));
		dispatcher.addMapping("/");
	}
	
	//This method enables pass through of requests not handled by the dispatcher servlet (for static assets, for example).
	@Override
	public void configureDefaultServletHandling(DefaultServletHandlerConfigurer configurer) {
		configurer.enable();
	}
	
	@Bean
	public ViewResolver viewResolver() {
		InternalResourceViewResolver resolver = new InternalResourceViewResolver();	
		resolver.setPrefix("/WEB-INF/jsp/");
		resolver.setSuffix(".jsp");		
		return resolver;
	}
	

	@Bean
	public BasicDataSource dataSource() {
		BasicDataSource dataSource = new BasicDataSource();
		dataSource.setDriverClassName("org.postgresql.Driver");
		
		String dbUrl = System.getenv("JDBC_DATABASE_URL");
		if(dbUrl != null) {
			dataSource.setUrl(dbUrl);
		} else {
			dataSource.setUrl("jdbc:postgresql://localhost:5432/msgboard");
			dataSource.setUsername("postgres");
			dataSource.setPassword("postgres1");
		}
		return dataSource;
	}
	
	@Bean
	public JwtAuthInterceptor interceptor() {
		List<String> exceptions = new ArrayList<String>();
		exceptions.add("/register");
		exceptions.add("/login");
		exceptions.add("/");
		exceptions.add("/api/subnoddits");
		exceptions.add("/api/subnoddits/active");
		exceptions.add("/api/subnoddits/{subnodditName}");
		exceptions.add("/api/{subnodditName}");
		exceptions.add("/api/{subnodditName}/{postId}");
		exceptions.add("/api/{subnodditName}/{postId}/replies");
		exceptions.add("/api/allposts/{subnodditName}");
		exceptions.add("/api/allposts");
		exceptions.add("/api/popularposts");
		exceptions.add("/api/recentposts");
		exceptions.add("/api/user/{username}");
		exceptions.add("/api/subnoddits/active");
		exceptions.add("/api/post/votes/{postId}");
		return new JwtAuthInterceptor(exceptions);
	}
	
	@Override
	public void addInterceptors(InterceptorRegistry registry) {
		registry.addInterceptor(interceptor());
	}

}
