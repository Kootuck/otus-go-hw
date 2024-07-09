package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

type users [100_000]User

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	users, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(users, domain)
}

func getUsers(r io.Reader) (u users, err error) {
	scanner := bufio.NewScanner(r)
	var i int
	for scanner.Scan() {
		var user User
		if er := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			err = fmt.Errorf("umarshalling error: %w", er)
			return
		}
		u[i] = user
		i++
	}
	if scanner.Err() != nil {
		err = fmt.Errorf("read file error: %w", err)
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		email := strings.ToLower(user.Email)
		if idx := strings.LastIndex(email, "@"); idx > 0 {
			userDomain := email[idx+1:]
			if strings.HasSuffix(userDomain, domain) {
				result[userDomain]++
			}
		}
	}

	return result, nil
}
