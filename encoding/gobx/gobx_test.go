package gobx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
	Age  int
}

func TestEncode(t *testing.T) {
	a := assert.New(t)
	user := User{Name: "chen", Age: 2}
	data, err := Encode(user)
	a.NoError(err)

	var load User
	err = Decode(data, &load)
	a.NoError(err)
	fmt.Println(load)
	a.Equal(user, load)
}
