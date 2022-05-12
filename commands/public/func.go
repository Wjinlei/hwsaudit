package public

/*
#cgo CFLAGS: -I./clib
#cgo LDFLAGS: -L${SRCDIR}/clib -lmyacl -lacl
#include <stdlib.h>
#include "myacl.h"
*/
import "C"

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

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

func IsMatchAcl(path string, facl string) (string, bool) {
	// Get file acl
	cPath := C.CString(path)
	cFacl := C.getfacl(cPath) // cgo
	o := C.GoString(cFacl)
	C.free(unsafe.Pointer(cPath))
	C.free(unsafe.Pointer(cFacl))

	o = strings.TrimSpace(o)
	if o == "" {
		return "", false
	}

	for _, item := range strings.Split(facl, ",") {
		rule := strings.Split(item, ":")
		user := strings.TrimSpace(strings.ReplaceAll(rule[0], "*", ""))
		if user != "" {
			if !strings.Contains(o, ":"+user+":") {
				continue
			}
		}

		mode := "*"
		if len(rule) > 1 {
			mode = rule[1]
		}

		for _, line := range strings.Split(o, "\n") {
			if user == "" {
				lineMode := strings.Split(line, ":")[2]
				if ok := contains(lineMode, strings.TrimSpace(mode)); ok {
					return strings.ReplaceAll(o, "\n", ","), true
				}
			} else {
				if strings.Contains(line, ":"+user+":") {
					lineMode := strings.Split(line, ":")[2]
					if ok := contains(lineMode, strings.TrimSpace(mode)); ok {
						return strings.ReplaceAll(o, "\n", ","), true
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
