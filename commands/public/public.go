package public

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/Wjinlei/golib"
	"github.com/Wjinlei/hwsaudit/global"
)

func WalkDir(save bool, root string, target string, user string, mode string, s bool, t bool, acl string) error {
	i := 1
	jsonFile := "home.json"
	textFile := "/tmp/home.txt"

	golib.Delete(textFile)
	golib.Delete(jsonFile)
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		stat, err := d.Info()
		if err != nil {
			return nil
		}

		/* Exclude */
		if stat.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		ok := true
		facl := ""
		result := Result{}
		sSys := stat.Sys().(*syscall.Stat_t)

		/* Match scan target */
		if target == "file" && d.IsDir() {
			return nil
		}
		if target == "dir" && !d.IsDir() {
			return nil
		}

		/* Match user */
		if user != "" && user != "*" {
			ok = IsMatchUser(int(sSys.Uid), user)
			if !ok {
				return nil
			}
		}

		/* Match mode */
		if mode != "" && mode != "*" {
			ok = IsMatchMode(stat.Mode(), mode)
			if !ok {
				return nil
			}
		}

		/* Match sUid || sGid */
		if s {
			if stat.Mode()&os.ModeSetuid != 0 || stat.Mode()&os.ModeSetgid != 0 {
				ok = true
			} else {
				return nil
			}
		}

		/* Match sticky */
		if t {
			if stat.Mode()&os.ModeSticky != 0 {
				ok = true
			} else {
				return nil
			}
		}

		/* Match acl */
		if acl != "" {
			facl, ok = IsMatchAcl(path, acl)
			if !ok {
				return nil
			}
		}

		/* Check ok */
		if ok {
			result.Id = i
			result.Name = d.Name()
			result.Path = golib.GetAbs(path)
			result.User = global.FindUser(int(sSys.Uid))
			result.Mode = stat.Mode().String()
			result.Facl = facl

			jsonResult, _ := json.Marshal(result)
			if save {
				golib.FileWrite(jsonFile, string(jsonResult)+"\n", golib.FileAppend)
				golib.FileWrite(textFile, fmt.Sprintf(
					"<file>%d{*}%s{*}%s{*}%s{*}%s{*}%s</file>\n",
					result.Id, result.Name, result.Path, result.User, result.Mode, result.Facl),
					golib.FileAppend)
			} else {
				fmt.Println(string(jsonResult))
			}
			i = i + 1
		}
		return nil
	})
	return nil
}
