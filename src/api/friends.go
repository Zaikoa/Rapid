package database

import (
	"database/sql"

	custom "github.com/Zaikoa/rapid/src/handling"
)

/*
Handles adding a friend to the users friend list
*/
func AddFriend(friendCode string, id int) (bool, error) {
	if id == 0 {
		return false, custom.NewError("User must be logged in to use this method")
	}

	var to_friend_id int
	query := `SELECT id FROM users WHERE friend_code=$1`
	err := conn.QueryRow(query, friendCode).Scan(&to_friend_id)

	if err != nil || to_friend_id == 0 { // id has to be 0 if it does not exist since sql does not support null primary ids
		return false, custom.NewError("User you are trying to add does not exist")
	}

	if AreMutualFriends(id, to_friend_id) {
		return false, custom.NewError("User is already friends with the specified user")
	}

	query = `INSERT INTO friends (user_one, user_two) VALUES ($1, $2)`
	_, err = conn.Exec(query, id, to_friend_id)
	if err != nil {
		return false, err
	}
	return true, nil

}

/*
Removes a friend from users friend list
*/
func DeleteFriend(id int, username string) (bool, error) {
	if id == 0 {
		return false, custom.NewError("User must be logged in to use this method")
	}

	query := `DELETE FROM friends WHERE (user_one=$1 AND user_two=$2) OR (user_one=$2 AND user_two=$1)`
	result, err := GetUserID(username)
	if err != nil {
		return false, err
	}
	_, err = conn.Exec(query, id, result)
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
Checks if user one is friends with user two
*/
func IsFriend(user_one_id int, user_two_id int) bool {
	var temp int
	query := `
	SELECT id 
	FROM friends
	WHERE friends.user_one=$1 AND friends.user_two=$2`
	err := conn.QueryRow(query, user_one_id, user_two_id).Scan(&temp)
	if err == sql.ErrNoRows || temp == 0 { // primary id cannot be null so we must check if its 0 instead
		return false
	}
	return true
}

// AreMutualFriends determines if two users are mutual friends
func AreMutualFriends(user_one_id int, user_two_id int) bool {
	return IsFriend(user_one_id, user_two_id) || IsFriend(user_two_id, user_one_id)
}

type Friend struct {
	Name       string
	FriendCode string
}

/*
Retrieves a list of friends for a given user
*/
func GetFriendsList(id int) ([]Friend, error) {
	query := `
	SELECT nickname, friend_code
	FROM (
		SELECT users.nickname, users.friend_code
		FROM users
		JOIN friends ON users.id = friends.user_two
		WHERE friends.user_one = $1
		UNION
		SELECT users.nickname, users.friend_code
		FROM users
		JOIN friends ON users.id = friends.user_one
		WHERE friends.user_two = $1
	) AS combined_data
	GROUP BY nickname, friend_code`

	rows, err := conn.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendsList []Friend
	for rows.Next() {
		var friend Friend
		if err := rows.Scan(&friend.Name, &friend.FriendCode); err != nil {
			return nil, err
		}
		friendsList = append(friendsList, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return friendsList, nil
}