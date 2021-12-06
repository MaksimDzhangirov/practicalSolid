package goodCode2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

func GetUserFile(id uint) (io.Reader, error) {
	filename := fmt.Sprintf("user_%d.json", id)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func GetUserHTTP(id uint) (io.Reader, error) {
	uri := fmt.Sprintf("http://some-api.com/users/%d", id)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func GetDummyUser(userJSON UserJSON) (io.Reader, error) {
	data, err := json.Marshal(userJSON)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func GetUser(reader io.Reader) (*User, error) {
	data, err := ioutil.ReadAll(reader)
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
