# YLWS BACKGROUND API
websocket链接: ws://192.168.2.68:10001/chat 

<!-- [一、相册标签管理](#1)  
&nbsp; &nbsp; [标签添加](#1.1)  
&nbsp; &nbsp; [标签删除](#1.2)  
&nbsp; &nbsp; [标签列表](#1.3)  
&nbsp; &nbsp; [标签绑定](#1.4)   -->

---
<h3 id='1.1'>加入频道</h3>

请求参数：

参数名|类型|说明
---|---|---
task|string|消息类型
uid|string|用户ID
sub|string|频道ID


返回参数:

参数名|类型|说明
---|---|---
task|string|消息类型
data|obj|数据

请求示例:
```
uid=123131&sub=roomworld&task=join
```
返回示例：
```
{
    "task":"join",
    "data":{}
}
```
---
<h3 id='1.2'>退出频道</h3>

请求参数：

参数名|类型|说明
---|---|---
uid|string|用户ID
sub|string|频道ID
type|string|消息类型


返回参数: 

参数名|类型|说明
---|---|---
task|string|消息类型
data|obj|数据

请求示例:
```
uid=123131&sub=roomworld&task=exit
```
返回示例：
```
{
    "task":"join",
    "data":{}
}
```
---
<h3 id='1.3'>发送消息</h3>

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
data|obj|数据


请求示例:
```
msg=asdasd&sub=roomservice&task=msg
```

返回示例：
```
{
    "task":"join",
    "data":{}
}
```

---
<h3 id='1.4'>连接心跳</h3>

请求参数：

参数名|类型|说明
---|---|---
task|string|消息类型
time|string|Unix时间戳

返回参数: 

参数名|类型|说明
---|---|---
task|string|消息类型
data|obj|数据


请求示例:
```
time=1520414598&task=heart
```

返回示例：
```
{
    "task":"join",
    "data":{}
}
```