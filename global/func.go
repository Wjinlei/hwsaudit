package global

import (
	"os/user"
	"strconv"
)

func FindUser(uid int) string {
	user, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return "-"
	}
	return user.Username
}
