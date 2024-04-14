package api

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/kshwedha/geodata-go/src/common/db"
)

func OneQueryExec() {
	// var count int
	// query := "select count(*) from users where $1 = $2 limit 1 offset 1;"
	// err = conn.QueryRow(query, key, val).Scan(&count)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return false
	// }
	// fmt.Println(count)
	// return count > 0
	return
}

func DoesExists(key string, val string) bool {
	conn, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	query := fmt.Sprintf("select count(%s) from users where %s='%s' limit 1;", key, key, val)
	results := db.ExecPsqlRows(conn, query)
	for results.Next() {
		var count int
		err := results.Scan(&count)
		if err != nil {
			return true
		}
		if count > 0 {
			return false
		}
	}
	return true
}

func RegisterUser(username string, password string, email string) error {
	hash := md5.New()
	hash.Write([]byte(password))
	password = fmt.Sprintf("%x", hash.Sum(nil))
	query := fmt.Sprintf("insert into users (username, password, email) values ('%s', '%s', '%s');", username, password, email)
	conn, err := db.InitDB()
	if err != nil {
		return err
	}
	results, err := db.ExecPsqlResult(conn, query)
	if err != nil {
		return err
	}
	if results == 1 {
		return nil
	}
	return fmt.Errorf("!! could not register user")
}

func LoginUser(loginid string, password string) (string, error) {
	if loginid == "" || password == "" {
		return "", fmt.Errorf("username or password cannot be empty")
	}
	if DoesExists("username", loginid) && DoesExists("email", loginid) {
		return "", fmt.Errorf("user does not exist")
	}
	hash := md5.New()
	hash.Write([]byte(password))
	password = fmt.Sprintf("%x", hash.Sum(nil))
	query := fmt.Sprintf("select count(*) from users where username='%s' and password='%s' limit 1;", loginid, password)
	conn, err := db.InitDB()
	if err != nil {
		return "", err
	}
	results := db.ExecPsqlRows(conn, query)
	for results.Next() {
		var count int
		err := results.Scan(&count)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		if count > 0 {
			token, err := CreateSession(loginid)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
			return token, nil
		}
	}
	return "!! ", fmt.Errorf("invalid username or password")
}

func getUserID(loginid string) (string, error) {
	query := fmt.Sprintf("select id from users where username='%s' or email='%s' limit 1;", loginid, loginid)
	conn, err := db.InitDB()
	if err != nil {
		return "0", err
	}
	results := db.ExecPsqlRows(conn, query)
	for results.Next() {
		var id string
		err := results.Scan(&id)
		if err != nil {
			return "0", err
		}
		return id, nil
	}
	return "0", fmt.Errorf("user does not exist")
}

func checkExistingSession(userid string) (bool, string) {
	query := fmt.Sprintf("select session_token from sessions where user_id='%s' and expiration_time > CURRENT_TIMESTAMP;", userid)
	conn, err := db.InitDB()
	if err != nil {
		return false, ""
	}
	results := db.ExecPsqlRows(conn, query)
	for results.Next() {
		var token string
		err := results.Scan(&token)
		if err != nil {
			return false, ""
		}
		if len(token) > 0 {
			return true, token
		}
	}
	return false, ""
}

func CreateSession(loginid string) (string, error) {
	userid, _ := getUserID(loginid)
	if userid == "0" {
		return "", fmt.Errorf("user does not exist")
	}
	exists, etoken := checkExistingSession(userid)
	if exists {
		return etoken, nil
	}
	token, err := generateSessionToken(128)
	if err != nil {
		return "", err
	}
	query := fmt.Sprintf("insert into sessions (user_id, session_token, expiration_time) values ('%s', '%s', CURRENT_TIMESTAMP + INTERVAL '7 days');", userid, token)
	conn, err := db.InitDB()
	if err != nil {
		return "", err
	}
	_, err = db.ExecPsqlResult(conn, query)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return token, nil
}

func generateSessionToken(length int) (string, error) {
	// Create a byte slice to store random bytes
	token := make([]byte, length/2)

	// Read random bytes into the byte slice
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to a hexadecimal string
	return hex.EncodeToString(token), nil
}

func SaveFile(file []byte, file_name string) error {
	query := fmt.Sprintf("insert into files (filename, file_content) values ('%s', '%s');", file_name, file)
	conn, err := db.InitDB()
	if err != nil {
		return err
	}
	_, err = db.ExecPsqlResult(conn, query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
