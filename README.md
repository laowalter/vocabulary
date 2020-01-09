## 项目：Vocabulary

voc用来平日记忆英文单词用，把需要查询到单词用translate命令查询，查询结果在屏幕显示的同时，还会追加到本地文件~/.word/vocabulary.txt中。 执行 voc --store 导入生词和查询结果到本地文件数据库中。voc 命令复习单词，复习重复的频率暂时按照斐波那契数列有天数为单位重复。


* 操作环境 Gentoo/Linux 终端

* 依赖的程序 translate-shell, 其配置文件如下：

~/.translate-shell/init.trans
{
    :show-original  true
    :indent         2
    :hl             "en"
    :tl             ["zh-CN"]
    :user-agent     "Mozilla/5.0 (X11; Linux x86_64; rv:33.0) Gecko/20100101 Firefox/63.0"
    :no-ansi        true
}


### 使用方法：


1. 首次使用 voc --init

2. 日常使用 

	1. 查询  
	```
	translate Ural Mountains
	```

	2. 存入数据库
	```
	voc --store
	```
    3. 复习
	```
	voc
	```
		1. 如果觉得记住了，就点击 p ,进入下一轮记忆；
		2. space 显示翻译。

