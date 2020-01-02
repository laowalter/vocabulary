## 项目：Vocabulary

这个项目用来平日记忆英文单词用，把需要查询到单词用tranlate命令查询，查询结果在屏幕显示的同时，还会追加到本地文件voc.txt中。每天固定时间有exportdb程序导入生词和查询结果到本地文件数据库中。在由dailyReview每日弹出复习单词，复习重复的频率和方法有dailyReview执行。

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

