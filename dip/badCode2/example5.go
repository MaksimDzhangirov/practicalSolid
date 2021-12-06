package badCode2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type User struct {
	// какие-то поля
}

type UserJSON struct {
	// какие-то поля
}

func (j UserJSON) ToUser() *User {
	return &User{
		// какие-то поля
	}
}

func GetUser(id uint) (*User, error) {
	filename := fmt.Sprintf("user_%d.json", id)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var user UserJSON
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return user.ToUser(), nil
}
