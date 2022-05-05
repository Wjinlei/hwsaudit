package run

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func findUser(uid int) string {
    user, err := user.LookupId(strconv.Itoa(uid))
    if err != nil {
        return "-"
    }
    return user.Username
}

// GetAbsPath 返回传入路径的绝对路径,如果绝对路径获取失败则返回原路径
func GetAbsPath(filePath string) string {
    if filepath.IsAbs(filePath) {
        return filePath
    }
    if strings.HasSuffix(filePath, ".") || strings.HasSuffix(filePath, "..") {
        filePath = filePath + "/"
    }
    filePathAbs, err := filepath.Abs(filepath.Dir(filePath))
    if err != nil {
        return filePath
    }
    name := filepath.Base(filePath)
    if name == "." || name == ".." {
        return filepath.FromSlash(filePathAbs)
    }
    return filepath.FromSlash(fmt.Sprintf("%s/%s", filePathAbs, filepath.Base(filePath)))
}
