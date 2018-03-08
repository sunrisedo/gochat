# YLWS BACKGROUND API
websocket链接: wss://ttt.tdex.com/chat

[一、聊天室](#1)  
&nbsp; &nbsp; [心跳握手](#1.1)  
&nbsp; &nbsp; [加入频道](#1.2)  
&nbsp; &nbsp; [退出频道](#1.3)  
&nbsp; &nbsp; [发送消息](#1.4)  


---
<h3 id='1.1'>心跳握手</h3>
WebSocket API 支持双向心跳，无论是 Server 还是 Client 都可以发起 ping message，对方返回 pong message。

注：返回的数据里面的 "pong" 的值为收到的 "ping" 的值 注：WebSocket Client 和 WebSocket Server 建立连接之后，WebSocket Server 每隔 5s（这个频率可能会变化） 会向 WebSocket Client 发起一次心跳，WebSocket Client 忽略心跳5次后，WebSocket Server 将会主动断开连接。

请求参数：

参数名|类型|说明
---|---|---
ping|string|Unix时间戳
pong|string|Unix时间戳
返回参数: 

参数名|类型|说明
---|---|---
ping|string|Unix时间戳
pong|string|Unix时间戳


请求示例:
```
ping=1520414598
```

返回示例：
```
{
    "pong":1520473202,
}
```

---
<h3 id='1.2'>加入频道</h3>

请求参数：

参数名|类型|说明
---|---|---
id|string|请求ID(时间戳)
task|string|消息类型
uid|string|用户ID
sub|string|频道ID


返回参数:

参数名|类型|说明
---|---|---
data|object|数据对象
id|string|请求ID(时间戳)
result|string|消息类型
chats|array|该频道信息列表


请求示例:
```
id=1520474198&uid=123131&sub=roomworld&task=join
```
返回示例:
```
{
    "task":"join",
    "data":{
        "id":1520474198,
        "result":"success"
        "chats"[
            {
                "sub":"roomservice",
                "uid":112312,
                "time":1520474198,
                "msg":"hello world"
            },
            {
                "sub":"roomservice",
                "uid":112312,
                "time":1520474198,
                "msg":"hello world"
            },
        ]
    }
}
```

加入频道后每当频道用户列表 userupdate 有更新时，client 会收到数据，示例：

参数名|类型|说明
---|---|---
task|string|消息类型
data|object|数据对象
sub|string|频道ID
uid|string|用户ID
action|string|信息(join,exit)
uids|array|该频道用户列表

```
{
    "task":"userupdate",
    "data":{
        "sub":"roomworld",
        "uid":112312,
        "action":"join",
        "uids":[112312]
    }
}
```

加入频道后每当频道用户 msg 有更新时，client 会收到数据，示例：

参数名|类型|说明
---|---|---
task|string|消息类型
data|object|数据对象
sub|string|频道ID
uid|string|用户ID
time|int64|时间戳
msg|string|信息

```
{
    "task":"msg",
    "data":{
        "sub":"roomservice",
        "uid":112312,
        "time":1520474198,
        "msg":"hello world"
    }
}
```
---
<h3 id='1.3'>退出频道</h3>
请求参数：

参数名|类型|说明
---|---|---
id|string|请求ID(时间戳)
task|string|消息类型
uid|string|用户ID
sub|string|频道ID


返回参数:

参数名|类型|说明
---|---|---
data|object|数据对象
id|string|请求ID(时间戳)
result|string|消息类型



请求示例:
```
id=1520474198&uid=123131&sub=roomworld&task=exit
```
返回示例:
```
{
    "task":"exit",
    "data":{
        "id":1520474198,
        "result":"success"
    }
}
```
---
<h3 id='1.4'>发送消息</h3>

请求参数：

参数名|类型|说明
---|---|---
task|string|消息类型
uid|string|用户ID
sub|string|频道ID
msg|string|消息


返回参数: 

参数名|类型|说明
---|---|---
task|string|消息类型
data|object|数据对象
sub|string|频道ID
uid|string|用户ID
time|int64|时间戳
msg|string|信息


请求示例:
```
msg=asdasd&sub=roomservice&task=msg
```

```
{
    "task":"msg",
    "data":{
        "sub":"roomservice",
        "uid":112312,
        "time":1520474198,
        "msg":"hello world"
    }
}
```

