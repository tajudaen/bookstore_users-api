package users

import (
	"fmt"

	"github.com/tajud99n/bookstore_users-api/database/mysql/users"
	"github.com/tajud99n/bookstore_users-api/logger"
	"github.com/tajud99n/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser       = "INSERT INTO users(firstname, lastname, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, firstname, lastname, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser       = "UPDATE users SET firstname=?, lastname=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users where id=?;"
	queryFindUserByStatus = "SELECT id, firstname, lastname, email, date_created, status FROM users WHERE status=?;"
)

func (user *User) Save() *errors.RestErr {
	stmt, err := users.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("Server Error")
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return errors.NewInternalServerError("Server Error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerError("Server Error")
	}
	user.Id = userId
	return nil
}

func (user *User) Get() *errors.RestErr {
	stmt, err := users.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("Server Error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error when trying to prepare get user user by id", getErr)
		return errors.NewInternalServerError("Server Error")
	}

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("Server Error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError("Server Error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("Server Error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return errors.NewInternalServerError("Server Error")
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, errors.NewInternalServerError("Server Error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status ", err)
		return nil, errors.NewInternalServerError("Server Error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan user row into user struct", err)
			return nil, errors.NewInternalServerError("Server Error")
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}
