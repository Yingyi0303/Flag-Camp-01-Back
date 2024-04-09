package com.example.loginout;

import com.example.loginout.db.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.web.SecurityFilterChain;


@Configuration
@EnableWebSecurity
public class WebSecurityConfig {

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    // Configure the SecurityFilterChain
    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
                //.csrf().disable()
                .authorizeRequests(authorizeRequests ->
                        authorizeRequests
                                .requestMatchers("/register").permitAll() // Allow registration without authentication
                                .anyRequest().authenticated() // All other requests require authentication
                )
                .formLogin(form -> form
                        .loginPage("/login").permitAll() // Customize the login page URL if needed
                        .defaultSuccessUrl("/home", true) // Redirect after successful login
                )
                .logout(logout -> logout
                        .logoutSuccessUrl("/login?logout") // Customize the logout success URL if needed
                );
        return http.build();
    }

    @Bean
    public UserDetailsService userDetailsService(final PasswordEncoder passwordEncoder) {
        return new UserDetailsService() {
            @Autowired
            private UserRepository userRepository;

            @Override
            public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
                com.example.loginout.db.entity.UserEntity userEntity = userRepository.findByUsername(username);
                if (userEntity == null) {
                    throw new UsernameNotFoundException("User '" + username + "' not found");
                }

                return User.builder()
                        .username(userEntity.username())
                        // In a real application, you should also consider user roles and permissions
                        .password(passwordEncoder.encode(userEntity.password()))
                        .roles("USER") // This should be dynamic based on your application's roles
                        .build();
            }
        };
    }
}


