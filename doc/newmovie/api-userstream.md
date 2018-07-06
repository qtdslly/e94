# 用户自定义视频流  

## 用户新增视频流

#### 方法
`POST` `/cms/user/stream/add`

#### 功能
用户新增视频流

#### 请求内容

```
{
	"installation_id":1807042120121435,
	"title":"江苏卫视5",
	"url":"http://www.ezhantao5.com"
}

```

#### 返回结果
```
{
    "err_code": 0
}
```

## 用户修改视频流

#### 方法
`POST` `/cms/user/stream/update`

#### 功能
用户修改视频流

#### 请求内容

```
{
	"user_stream_id":2,
	"title":"湖北卫视",
	"url":"http://www.ezhantao.com"
}
```

#### 返回结果
```
{
    "err_code": 0
}
```



## 用户删除视频流

#### 方法
`POST` `/cms/user/stream/delete`

#### 功能
用户删除视频流

#### 请求内容

```
{
	"user_stream_id":2
}
```

#### 返回结果
```
{
    "err_code": 0
}
```



## 获取获取自定义视频流列表

#### 方法
`GET` `/cms/user/stream/list`

#### 功能
获取获取自定义视频流列表

#### 参数
| 参数   | 类型 | 必选 |  说明 |
|:--- | --- | ---  | --- |
|installation_id|int|是|类型|
|offset|int|偏移量|id|
|limit|string|是|返回数量|

#### 返回结果
```
{
    "count": 7,
    "data": [
        {
            "id": 3,
            "title": "湖南卫视",
            "url": "http://www.baidu.com"
        },
        {
            "id": 4,
            "title": "湖北卫视",
            "url": "http://www.ezhantao.com"
        },
        {
            "id": 5,
            "title": "江苏卫视",
            "url": "http://www.ezhantao1.com"
        }
    ],
    "err_code": 0,
    "has_more": true
}
```