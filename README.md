### 介绍
Go语言开发的Linux下的权限审计工具，功能还在开发中，欢迎PR。

### 使用示例
```sh
./hwsaudit run -u "-root" -m "777"      # 查找当前目录下,不属于root用户,并且权限是777的文件
./hwsaudit run -u "-root" -m "**2"      # 查找当前目录下,不属于root用户,并且其他用户有写权限的文件
./hwsaudit run -u "*" -m "**4"          # 查找当前目录下,属于任意用户,并且其他用户有读权限的文件
./hwsaudit run -u "w" -m "**6"          # 查找当前目录下,属于w用户,并且其他用户有读写权限的文件
./hwsaudit run -u "w" -m "**2" -s       # 查找当前目录下,属于w用户,并且其他用户有写权限的文件,并且拥有SetUid 或者 SetGid 的文件
./hwsaudit run -u "w" -m "**2" -t       # 查找当前目录下,属于w用户,并且其他用户有写权限的文件,并且拥有粘贴位Sticky的文件
./hwsaudit run -u "root" -a "*:2"       # 查找当前目录下,属于root用户,并且任意用户拥有acl写权限的文件
./hwsaudit run -u "root" -a "w:6"       # 查找当前目录下,属于root用户,并且w用户拥有acl读写权限的文件
./hwsaudit run -C "dir" -m "777"        # 查找当前目录下,目录权限是777的目录
./hwsaudit run -d "/wwwroot" -u "www"   # 查找/wwwroot目录下,属于www用户的文件
```

### 选项默认值
如果某位不写，则为默认值
- `-C 默认值 "file"`
- `-d 默认值 "./"`
- `-u 默认值 "*" 匹配所有`
- `-m 默认值 "*" 匹配所有`
- `-a 默认值 "*" 这里只有一个*代表不检查, *:* 才代表匹配所有`
- `-s 默认值 false`
- `-t 默认值 false`

### 获取帮助
- `./hwsaudit` or `./hwsaudit help [command]`
- 你还可以将输出结果扔给`fzf`,以供选择
