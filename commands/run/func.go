package run

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Wjinlei/golib/os/cmd"
)

var table map[string]string

func init() {
	/* Mode table */
	table = map[string]string{
		"0": "!!!",
		"1": "--x",
		"2": "-w-",
		"3": "-wx",
		"4": "r--",
		"5": "r-x",
		"6": "rw-",
		"7": "rwx",
	}
}

func isMatchUser(uid int) bool {
	/* Match user */
	if strings.Contains(opt.user, "-") && opt.user != "-" {
		if findUser(uid) != strings.Trim(opt.user, "-") {
			return true
		}
	} else {
		if findUser(uid) == opt.user {
			return true
		}
	}
	return false
}

func isMatchMode(fileMode os.FileMode) bool {
	/* Get Mode().Perm() []string */
	perm := strings.Split(fmt.Sprintf("%#o", fileMode.Perm()), "")

	/* Get normal perm mode */
	if len(perm) > 3 {
		perm = perm[1:4]
	}

	/* Fix mode < 3 */
	for len(perm) < 3 {
		perm = append(perm, "0")
	}

	for i, mode := range strings.Split(opt.fileMode, "") {
		/* Mode length only 3 */
		if i > 2 {
			break
		}

		switch mode {
		case "0", "1", "2", "3", "4", "5", "6", "7":
			for _, bit := range strings.Split(table[mode], "") {
				if bit == "-" {
					continue
				}
				if !strings.Contains(table[perm[i]], bit) {
					return false
				}
			}
			break
		case "*":
			continue
		default:
			return false
		}
	}

	return true
}

func isMatchAcl(filePath string) (string, bool) {
	out, err := cmd.New().Shell(fmt.Sprintf("getfacl -c -s -p %s |grep -E :.+:", filePath))
	if err != nil {
		return "", false
	}
	data, err := ioutil.ReadAll(out)
	if err != nil {
		return "", false
	}

	str := string(data)
	str = strings.TrimSpace(str)

	if str == "" {
		return "", false
	}

	for _, rule := range strings.Split(opt.fileAcl, ",") {
		rule = strings.Trim(rule, "*")
		rule = strings.ReplaceAll(rule, "0", "---")
		rule = strings.ReplaceAll(rule, "1", table["1"])
		rule = strings.ReplaceAll(rule, "2", table["2"])
		rule = strings.ReplaceAll(rule, "3", table["3"])
		rule = strings.ReplaceAll(rule, "4", table["4"])
		rule = strings.ReplaceAll(rule, "5", table["5"])
		rule = strings.ReplaceAll(rule, "6", table["6"])
		rule = strings.ReplaceAll(rule, "7", table["7"])

		if !strings.HasPrefix(rule, ":") {
			rule = ":" + rule
		}

		if !strings.Contains(str, rule) {
			return "", false
		}
	}

	return strings.ReplaceAll(str, "\n", ","), true
}
