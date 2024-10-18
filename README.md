## Moments 极简朋友圈

docker-compose.yaml 配置文件如下:

```yaml
services:
  moments:
    image: jkjoy/moments:latest
    environment:
      JWT_KEY: "BbYS93dHHfIC1cQR8rI6"
      WEBHOOK_URL: "https://open.feishu.cn/open-apis/bot/v2/hook/*" #飞书webhook 
      SITE_URL: "https://www.moments.cn" #访问地址
      QQ_WEBHOOK_URL: "https://http.asbid.cn" #QQ机器人的API
      QQ_USER_ID: "123456" #接收消息的QQ号码
    ports:
      - "3000:3000"
    volumes:
      - ./data:/app/data
      - /etc/localtime:/etc/localtime:ro
      - /etc/timezone:/etc/timezone:ro
```

mod版本增加了webhook评论通知

没有修改原版的数据库,使用系统变量读取,

`WEBHOOK_URL`为你使用的webhook地址, 可以是飞书webhook, 也可以是其他的.
`SITE_URL`为你的moments的访问地址, 可以是域名,也可以是ip地址.用来拼接memo的访问地址

`QQ_WEBHOOK_URL`为你使用的QQ机器人的API地址,需要自行部署,或者使用公共服务
`QQ_USER_ID`为你接收消息的QQ号码

修改了`fancybox`为`viewimage`





















```md
v0.2.1使用了golang作为服务端重写,目前已经基本实现了0.2.0版本的大部分功能.包体积更小了.

1. 增加了多用户模式,后台可以自由开启是否运行注册多用户.
2. 支持在Linux/MacOS/Windows平台双击本地启动.
3. 标签的定义,以#号开头,空格/空行结尾的中间的部分会被认为是标签.
4. 完善了tag标签的选择,在memo发言的输入框里点击右键可以选择标签来插入.
5. 支持了完整的markdown,但是目前样式只适配了常用的几个标签,更多的待接下来完善.
6. 默认用户名密码`admin/a123456`,登陆后后台可以修改.

[更多说明](https://discussion.mblog.club/post/pto2hqoFzDKzZMpvoPZKYuP)

[交流TG群](https://t.me/simple_moments)

[交流论坛](https://discussion.mblog.club/)

[0.2.0的README](https://github.com/kingwrcy/moments/blob/master/README.md)


#### v0.2.5发布说明 2024-08-14

1. 增加代码内容/发言内容强制换行.
2. 尝试减小代码高亮的引入文件大小,加快首页打开速度
3. 修复首次加入标签时异常的问题

#### v0.2.4发布说明 2024-08-09

1. 发言输入框不再支持#号开头的内容识别为标签了,标签改为单独一列,右键输入标签继续可以.
2. 首页默认不自动加载下一页了,后台增加`是否自动加载下一页`的开关,需要的可以手动开启.
3. 增加代码块的支持,支持语法高亮,使用方式是3个`符号之后跟上代码的语言即可.
4. 修复传第二张图片会把前一张图片删除的bug.
5. 尝试修复ios环境下safari浏览图片超宽的bug.
6. 增加`回到顶部`按钮,pc和手机模式下都有.
7. 修复登出按钮在pc无法看到的bug.
8. 代码块增加一键复制按钮.
```