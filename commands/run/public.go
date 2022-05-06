package run

import (
	"os/user"
	"strconv"
)

func findUser(uid int) string {
	user, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return "-"
	}
	return user.Username
}
