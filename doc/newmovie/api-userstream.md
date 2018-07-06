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



## 获取全局广告

#### 方法
`GET` `/cms/v1.0/ad/config`

#### 功能
获取全局广告,客户端定时拉取更新,方便控量及时生效

#### 参数
| 参数   | 类型 | 必选 |  说明 |
|:--- | --- | ---  | --- |
|app_key|string|是|app key|
|app_version|string|是|app版本|
|os_type|int|是|系统类型|
|channel|string|是|渠道|

#### 返回结果

```
{
  "err_code": 0,
  "data": {
    "providers": [
      {
        "id": 6,
        "provider": 4,
        "name": "SSP",
        "app_id": "",
        "app_key": "47630da927784f9e96853e7d048306b8"
      },
      {
        "id": 3,
        "provider": 1,
        "name": "百度Android",
        "app_id": "b893b417",
        "app_key": ""
      },
      {
        "id": 1,
        "provider": 2,
        "name": "广点通Android",
        "app_id": "1105693811",
        "app_key": ""
      },
      {
        "id": 7,
        "provider": 5,
        "name": "推啊",
        "app_id": "22079",
        "app_key": ""
      }
    ],
    "configs": [
      {
        "page": "stream",
        "type": "iva",
        "ad": {
          "id": 12,
          "name": "播放中banner",
          "units": [
            {
              "id": 26,
              "name": "SSP播放中banner",
              "type": 1,
              "provider_id": 6,
              "placement_id": "23be51fa188a45cb919170b90cb0bd60"
            },
            {
              "id": 27,
              "name": "百度Android播放中banner",
              "type": 1,
              "provider_id": 3,
              "placement_id": "3105082"
            },
            {
              "id": 28,
              "name": "广点通Android播放中banner",
              "type": 1,
              "provider_id": 1,
              "placement_id": "3090912579495415"
            },
            {
              "id": 31,
              "name": "直播插入汽车广告",
              "type": 3,
              "provider_id": 1,
              "iva": {
                "id": 1,
                "app_id": 1,
                "name": "汽车",
                "image_url": "http://pic133.nipic.com/file/20170615/13057068_204253899038_2.jpg",
                "jump_type": 1,
                "jump_url": "http://advod.starschinalive.com/ad/Public/Uploads/58df126086d05.mp4",
                "show_url": "",
                "click_url": "",
                "width": 0,
                "height": 0,
                "gravity": 0,
                "padding": 0,
                "style": 0,
                "created_at": "2017-06-23T16:38:43+08:00",
                "updated_at": "2017-06-23T16:38:43+08:00"
              }
            },
            {
              "id": 32,
              "name": "直播插入图片广告",
              "type": 3,
              "provider_id": 2,
              "iva": {
                "id": 2,
                "app_id": 1,
                "name": "测试图片",
                "image_url": "http://pic58.nipic.com/file/20150112/12299514_224005339000_2.jpg",
                "jump_type": 0,
                "jump_url": "",
                "show_url": "",
                "click_url": "",
                "width": 0,
                "height": 0,
                "gravity": 0,
                "padding": 0,
                "style": 0,
                "created_at": "2017-06-23T16:38:43+08:00",
                "updated_at": "2017-06-23T16:38:43+08:00"
              }
            }
          ]
        }
      },
      {
        "page": "player",
        "type": "feed",
        "ad": {
          "id": 8,
          "name": "前贴",
          "units": [
            {
              "id": 14,
              "name": "广点通Android前贴",
              "type": 1,
              "provider_id": 1,
              "placement_id": "4020613952354881"
            },
            {
              "id": 15,
              "name": "百度Android前贴",
              "type": 1,
              "provider_id": 3,
              "placement_id": "3146786"
            },
            {
              "id": 16,
              "name": "SSP Android前贴",
              "type": 1,
              "provider_id": 6,
              "placement_id": "fc27880a8f8a47dba165345f619ae8eb"
            }
          ]
        }
      },
      {
        "page": "app",
        "type": "loading",
        "ad": {
          "id": 1,
          "name": "开屏",
          "units": [
            {
              "id": 1,
              "name": "SSP开屏",
              "type": 1,
              "provider_id": 6,
              "placement_id": "b081790a82aa4492b9a736cbb9793eb1"
            },
            {
              "id": 2,
              "name": "广点通Android开屏",
              "type": 1,
              "provider_id": 1,
              "placement_id": "1060216559290464"
            },
            {
              "id": 3,
              "name": "百度Android开屏",
              "type": 1,
              "provider_id": 3,
              "placement_id": "3105100"
            }
          ]
        }
      },
      {
        "page": "app",
        "type": "rebootloading",
        "ad": {
          "id": 2,
          "name": "重启开屏",
          "units": [
            {
              "id": 4,
              "name": "SSP重启开屏",
              "type": 1,
              "provider_id": 6,
              "placement_id": "54c0c18454ff4e5fba23292e58b5f908"
            }
          ]
        }
      },
      {
        "page": "player",
        "type": "loading",
        "ad": {
          "id": 4,
          "name": "焦点图信息流",
          "units": [
            {
              "id": 8,
              "name": "广点通焦点图信息流",
              "type": 1,
              "provider_id": 1,
              "placement_id": "5090826190751756"
            },
            {
              "id": 9,
              "name": "百度焦点图信息流",
              "type": 1,
              "provider_id": 3,
              "placement_id": "3530011"
            },
            {
              "id": 10,
              "name": "SSP焦点图信息流",
              "type": 1,
              "provider_id": 6,
              "placement_id": "05d30a6710c54a21a25f127eb471896f"
            },
            {
              "id": 34,
              "name": "loading图片",
              "type": 2,
              "provider_id": 2,
              "custom": {
                "id": 2,
                "app_id": 1,
                "name": "测试",
                "content_type": 1,
                "content_url": "http://pic58.nipic.com/file/20150112/12299514_224005339000_2.jpg",
                "jump_type": 0,
                "jump_url": "",
                "duration": 5,
                "cancelable": false,
                "action": 0,
                "app_name": "",
                "bundle_id": "",
                "hint": "",
                "description": "",
                "show_url": "",
                "click_url": "",
                "created_at": "2017-06-26T09:54:20+08:00",
                "updated_at": "2017-06-26T09:54:20+08:00"
              }
            },
            {
              "id": 35,
              "name": "loading视频",
              "type": 2,
              "provider_id": 3,
              "custom": {
                "id": 3,
                "app_id": 1,
                "name": "测试",
                "content_type": 2,
                "content_url": "http://advod.starschinalive.com/ad/Public/Uploads/58df126086d05.mp4",
                "jump_type": 1,
                "jump_url": "http://clickc.admaster.com.cn/c/a82755,b1579049,c3887,i0,m101,8a2,8b1,h",
                "duration": 15,
                "cancelable": false,
                "action": 0,
                "app_name": "",
                "bundle_id": "",
                "hint": "",
                "description": "",
                "show_url": "",
                "click_url": "",
                "created_at": "2017-06-26T09:54:20+08:00",
                "updated_at": "2017-06-26T09:54:20+08:00"
              }
            }
          ]
        }
      },
      {
        "page": "app",
        "type": "exit",
        "ad": {
          "id": 5,
          "name": "信息流",
          "units": [
            {
              "id": 11,
              "name": "SSP 信息流",
              "type": 1,
              "provider_id": 6,
              "placement_id": "13264f5094d94adf809dcc14e349d8f9"
            }
          ]
        }
      },
      {
        "page": "home",
        "type": "overlay",
        "ad": {
          "id": 6,
          "name": "首页悬浮",
          "units": [
            {
              "id": 12,
              "name": "推啊首页悬浮",
              "type": 1,
              "provider_id": 7,
              "placement_id": "478"
            }
          ]
        }
      }
    ]
  }
}

```