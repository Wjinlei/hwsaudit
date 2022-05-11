package public

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Wjinlei/golib/os/cmd"
	"github.com/Wjinlei/hwsaudit/global"
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

func IsMatchUser(uid int, user string) bool {
	user = strings.TrimSpace(user)
	/* Match user */
	if strings.Contains(user, "-") && user != "-" {
		if global.FindUser(uid) != strings.Trim(user, "-") {
			return true
		}
	} else {
		if global.FindUser(uid) == user {
			return true
		}
	}
	return false
}

func IsMatchMode(fileMode os.FileMode, mode string) bool {
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

	for i, m := range strings.Split(strings.TrimSpace(mode), "") {
		/* Mode length only 3 */
		if i > 2 {
			break
		}

		if ok := contains(table[perm[i]], m); !ok {
			return false
		}
	}

	return true
}

func IsMatchAcl(filePath string, facl string) (string, bool) {
	cmdReader, err := cmd.New().Shell("getfacl -c -s -p " + filePath + " |grep -E :.+:")
	if err != nil {
		return "", false
	}

	data, err := ioutil.ReadAll(cmdReader)
	if err != nil {
		return "", false
	}

	out := strings.TrimSpace(string(data))
	if out == "" {
		return "", false
	}

	for _, item := range strings.Split(facl, ",") {
		rule := strings.Split(item, ":")
		user := strings.TrimSpace(strings.ReplaceAll(rule[0], "*", ""))
		if user != "" {
			if !strings.Contains(out, ":"+user+":") {
				continue
			}
		}

		mode := "*"
		if len(rule) > 1 {
			mode = rule[1]
		}

		for _, line := range strings.Split(out, "\n") {
			if user == "" {
				lineMode := strings.Split(line, ":")[2]
				if ok := contains(lineMode, strings.TrimSpace(mode)); ok {
					return strings.ReplaceAll(out, "\n", ","), true
				}
			} else {
				if strings.Contains(line, ":"+user+":") {
					lineMode := strings.Split(line, ":")[2]
					if ok := contains(lineMode, strings.TrimSpace(mode)); ok {
						return strings.ReplaceAll(out, "\n", ","), true
					}
				}
			}
		}
	}
	return "", false
}

func contains(a string, b string) bool {
	switch b {
	case "0", "1", "2", "3", "4", "5", "6", "7":
		for _, bit := range strings.Split(table[b], "") {
			if bit == "-" {
				continue
			}
			if !strings.Contains(a, bit) {
				return false
			}
		}
		break
	case "*", "":
		break
	default:
		return false
	}
	return true
}

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
