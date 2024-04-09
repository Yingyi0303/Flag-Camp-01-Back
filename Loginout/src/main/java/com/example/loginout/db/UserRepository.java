package com.example.loginout.db;

import com.example.loginout.db.entity.UserEntity;
import org.springframework.data.repository.CrudRepository;
import org.springframework.data.jdbc.repository.query.Query;
import org.springframework.data.jdbc.repository.query.Modifying;
import java.util.List;


public interface UserRepository extends CrudRepository<UserEntity, Long> {


    List<UserEntity> findByLastName(String lastName);


    List<UserEntity> findByFirstName(String firstName);


    UserEntity findByUsername(String username);


    @Modifying
    @Query("UPDATE users SET first_name = :firstName, last_name = :lastName WHERE username = :username")
    void updateNameByUsername(String username, String firstName, String lastName);
}
