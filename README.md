# 说明

浏览器中运行的JS程序无法修改某些Header，因此JS程序将通过此代理程序来获得这部分的Header的内容。

代理程序通过将受限Header替换为对应的非受限Header来绕过浏览器的限制。

具体替换的对应关系是：

- Cookie  -> X-Cookie
- Referer -> X-Referer
- Origin  -> X-Origin
