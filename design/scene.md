# scene代表了游戏的具体场景。本文档针对游戏具体场景定义通用接口。
# 请求参数：auth_credential，登录接口返回。请求本文档中的任何接口，都需要传递这个参数。

## 查询scene详情
### url path: scene/detail

### 请求
#### 请求参数scene：指定scene，格式要求非空字符串，只能包括字母和数字，另外支持'.','-','_'这三个特殊字符
#### auth_credential

### 响应
#### 响应的http body是个json。举例就是{"scene":"round1","layout":[{"slot1":{"visible":"Y","value":"1USD"}]}



## 揭晓指定场景的某个slot
### url path: scene/slot/reveal

### 请求
#### 请求参数scene：指定scene，格式要求非空字符串，只能包括字母和数字，另外支持'.','-','_'这三个特殊字符
#### 请求参数slot：指定slot，格式要求非空字符串，只能包括字母和数字，另外支持'.','-','_'这三个特殊字符
#### auth_credential

### 响应
#### 响应的http body是个json。举例就是{"scene":"round1","layout":[{"slot1":{"visible":"Y","value":"1USD"}]}


### 处理逻辑
#### 每次请求，如果成功处理，消耗用户的一个credit。如果用户credit为0，则操作失败。
#### 处理失败需要返回错误。



