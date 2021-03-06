## 项目：Vocabulary

voc用来平日记忆英文单词用，把需要查询到单词用translate命令查询，查询结果在屏幕显示的同时，还会追加到本地文件~/.word/vocabulary.txt中。 执行 voc --store 导入生词和查询结果到本地文件数据库中。voc 命令复习单词，复习重复的频率暂时按照斐波那契数列有天数为单位重复。


* 操作环境 Gentoo/Linux 终端

* 依赖的程序 translate-shell, 其配置文件如下：

~/.translate-shell/init.trans
```
{
    :show-original  true
    :indent         2
    :hl             "en"
    :tl             ["zh-CN"]
    :user-agent     "Mozilla/5.0 (X11; Linux x86_64; rv:33.0) Gecko/20100101 Firefox/63.0"
    :no-ansi        true
}
```

### 安装方法

#### 源码安装
```
$ sudo emerge -avq translate-shell
$ go get github.com/laowalter/vocabulary
$ cd ${GOPATH}/src/github.com/laowalter/vocabulary/
$ go build
```

#### 直接下载可执行文件


1. 安装translate-shell
2. 直接下载 voc 和 translate 两个文件。


### 使用方法


1. 首次使用需要初始化 
```
voc --init
```

2. 日常使用 

	1. 查询  
	```
	$ translate Ural Mountains
	```

	2. 存入数据库
	```
	$ voc --store
	```
	3. 列出当前未导入数据库的单词清单
	```
	$ voc --list
	```
	4. 列出数据库中的单词总数
	```
	$ voc --total
	```
    5. 复习当日和当日前应该复习但没有复习的单词
	```
	$ voc
	```
		1. 如果觉得记住了，就点击 p ,本单词将进入下一周期记忆, 本日不再显示；
		2. 如果点击m，可以修改本单词，比如原来查询单词时用到是复数;
		3. 如果点击d，则本记录永久从单词数据库记录中删除； 
		4. 如果点击r，把当前单词重置（reset）为新单词；
		5. space 显示翻译。

3. 其他情况
	* 如果单词同时出现在vocabulary.txt和数据库中，则重置单词为新单词。

### Todo

	增加一个配置文件用于个性化设置自己的记忆曲线。
