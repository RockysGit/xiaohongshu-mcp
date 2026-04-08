---
name: xhs-api
description: >
  小红书 API 技能。通过 HTTP API 调用小红书服务，支持登录管理、内容搜索、Feed 获取等功能。
  包含：检查登录状态、获取登录二维码、搜索内容、获取用户主页、获取 Feed 详情、获取 Feed 列表。
---

# 小红书 API 技能

通过 HTTP API 调用小红书 MCP 服务，实现内容查询和账号管理功能。

## 前置条件

- 小红书 MCP 服务已启动，默认地址：`http://10.2.248.59:18060`
- 首次使用需要扫码登录

## 功能列表

| 功能 | API 端点 | 方法   | 描述 |
|------|----------|------|------|
| 检查登录状态 | `/api/v1/login/status` | GET  | 检查当前是否已登录 |
| 获取登录二维码 | `/api/v1/login/qrcode` | GET  | 获取扫码登录的二维码 |
| 搜索内容 | `/api/v1/feeds/search` | GET  | 根据关键词搜索笔记 |
| 获取用户主页 | `/api/v1/user/profile` | POST | 获取指定用户的主页信息 |
| 获取 Feed 详情 | `/api/v1/feeds/detail` | POST | 获取笔记详情和评论 |
| 获取 Feed 列表 | `/api/v1/feeds/list` | GET  | 获取推荐 Feed 列表 |

---

## 工具 1: 检查登录状态

检查当前用户的登录状态，确认是否需要重新登录。

### 请求

```bash
curl -X GET "http://10.2.248.59:18060/api/v1/login/status"
```

### 响应示例

```json
{
  "success": true,
  "data": {
    "is_logged_in": true,
    "username": "用户名"
  },
  "message": "检查登录状态成功"
}
```

### 使用场景

- 在执行其他操作前，先检查登录状态
- 如果 `is_logged_in` 为 `false`，需要调用"获取登录二维码"进行登录

---

## 工具 2: 获取登录二维码

获取登录二维码，用于用户扫码登录。

### 请求

```bash
curl -X GET "http://10.2.248.59:18060/api/v1/login/qrcode"
```

### 响应示例

```json
{
  "success": true,
  "data": {
    "timeout": "300",
    "is_logged_in": false,
    "img": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
  },
  "message": "获取登录二维码成功"
}
```

### 响应字段说明

| 字段 | 类型 | 描述 |
|------|------|------|
| `timeout` | string | 二维码过期时间（秒） |
| `is_logged_in` | boolean | 当前是否已登录 |
| `img` | string | Base64 编码的二维码图片 |

### 使用流程

1. 调用此接口获取二维码
2. 将 `img` 字段的 Base64 图片展示给用户
3. 用户使用小红书 APP 扫码
4. 扫码成功后，调用"检查登录状态"确认登录成功

---

## 工具 3: 搜索内容

根据关键词搜索小红书笔记，支持多种筛选条件。

### 请求

```bash
curl -X POST "http://10.2.248.59:18060/api/v1/feeds/search" \
  -H "Content-Type: application/json" \
  -d '{
    "keyword": "搜索关键词",
    "filters": {
      "sort_by": "综合",
      "note_type": "不限",
      "publish_time": "不限"
    }
  }'
```

### 请求参数

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| `keyword` | string | 是 | 搜索关键词 |
| `filters.sort_by` | string | 否 | 排序方式：`综合`(默认) / `最新` / `最多点赞` / `最多评论` / `最多收藏` |
| `filters.note_type` | string | 否 | 笔记类型：`不限`(默认) / `视频` / `图文` |
| `filters.publish_time` | string | 否 | 发布时间：`不限`(默认) / `一天内` / `一周内` / `半年内` |
| `filters.search_scope` | string | 否 | 搜索范围：`不限`(默认) / `已看过` / `未看过` / `已关注` |
| `filters.location` | string | 否 | 位置距离：`不限`(默认) / `同城` / `附近` |

### 响应示例

```json
{
  "success": true,
  "data": {
    "feeds": [
      {
        "xsecToken": "security_token_value",
        "id": "feed_id_1",
        "modelType": "note",
        "noteCard": {
          "type": "normal",
          "displayTitle": "笔记标题",
          "user": {
            "userId": "user_id_1",
            "nickname": "用户昵称",
            "avatar": "https://example.com/avatar.jpg"
          },
          "interactInfo": {
            "liked": false,
            "likedCount": "100",
            "collected": false,
            "collectedCount": "50",
            "commentCount": "30"
          },
          "cover": {
            "url": "https://example.com/cover.jpg"
          }
        }
      }
    ],
    "count": 10
  },
  "message": "搜索Feeds成功"
}
```

### 重要字段说明

| 字段 | 描述 |
|------|------|
| `xsecToken` | 安全令牌，调用详情接口时需要 |
| `id` | Feed ID，用于获取详情 |
| `noteCard.displayTitle` | 笔记标题 |
| `noteCard.user` | 作者信息 |
| `noteCard.interactInfo` | 互动数据（点赞、收藏、评论数） |

---

## 工具 4: 获取用户主页

获取指定用户的主页信息，包括基本信息、互动数据和发布的笔记列表。

### 请求

```bash
curl -X POST "http://10.2.248.59:18060/api/v1/user/profile" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "用户ID",
    "xsec_token": "安全令牌"
  }'
```

### 请求参数

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| `user_id` | string | 是 | 用户 ID（从搜索结果或 Feed 中获取） |
| `xsec_token` | string | 是 | 安全令牌（从搜索结果或 Feed 中获取） |

### 响应示例

```json
{
  "success": true,
  "data": {
    "data": {
      "userBasicInfo": {
        "nickname": "用户昵称",
        "desc": "用户个人描述",
        "redId": "xiaohongshu_id",
        "gender": 1,
        "ipLocation": "浙江",
        "images": "https://example.com/avatar.jpg",
        "imageb": "https://example.com/background.jpg"
      },
      "interactions": [
        {"type": "follows", "name": "关注", "count": "1000"},
        {"type": "fans", "name": "粉丝", "count": "5000"},
        {"type": "interaction", "name": "获赞与收藏", "count": "10000"}
      ],
      "feeds": [
        {
          "xsecToken": "token",
          "id": "feed_id",
          "noteCard": {
            "displayTitle": "笔记标题"
          }
        }
      ]
    }
  },
  "message": "获取用户主页成功"
}
```

### 重要字段说明

| 字段 | 描述 |
|------|------|
| `userBasicInfo.nickname` | 用户昵称 |
| `userBasicInfo.redId` | 小红书号 |
| `userBasicInfo.gender` | 性别（1: 男, 2: 女, 0: 未知） |
| `interactions` | 关注数、粉丝数、获赞与收藏数 |
| `feeds` | 用户发布的笔记列表 |

---

## 工具 5: 获取 Feed 详情

获取指定笔记的详细信息，包括正文内容、图片列表和评论。

### 请求

```bash
curl -X POST "http://10.2.248.59:18060/api/v1/feeds/detail" \
  -H "Content-Type: application/json" \
  -d '{
    "feed_id": "笔记ID",
    "xsec_token": "安全令牌",
    "load_all_comments": false,
    "comment_config": {
      "click_more_replies": true,
      "max_comment_items": 50,
      "scroll_speed": "normal"
    }
  }'
```

### 请求参数

| 参数 | 类型 | 必填 | 描述 |
|------|------|------|------|
| `feed_id` | string | 是 | Feed ID |
| `xsec_token` | string | 是 | 安全令牌 |
| `load_all_comments` | boolean | 否 | 是否加载全部评论，默认 false |
| `comment_config.click_more_replies` | boolean | 否 | 是否展开更多回复 |
| `comment_config.max_replies_threshold` | int | 否 | 回复数量阈值，超过则跳过 |
| `comment_config.max_comment_items` | int | 否 | 最大加载评论数，0 表示全部 |
| `comment_config.scroll_speed` | string | 否 | 滚动速度：`slow` / `normal` / `fast` |

### 响应示例

```json
{
  "success": true,
  "data": {
    "feed_id": "64f1a2b3c4d5e6f7a8b9c0d1",
    "data": {
      "note": {
        "noteId": "64f1a2b3c4d5e6f7a8b9c0d1",
        "title": "笔记标题",
        "desc": "笔记详细内容描述",
        "type": "normal",
        "time": 1702195200000,
        "ipLocation": "浙江",
        "user": {
          "userId": "user_id_123",
          "nickname": "作者昵称",
          "avatar": "https://example.com/avatar.jpg"
        },
        "interactInfo": {
          "liked": false,
          "likedCount": "100",
          "collected": false,
          "collectedCount": "80",
          "commentCount": "50"
        },
        "imageList": [
          {
            "width": 1080,
            "height": 1440,
            "urlDefault": "https://example.com/image1.jpg"
          }
        ]
      },
      "comments": {
        "list": [
          {
            "id": "comment_id_1",
            "content": "评论内容",
            "likeCount": "10",
            "createTime": 1702195200000,
            "userInfo": {
              "userId": "commenter_id",
              "nickname": "评论者昵称"
            },
            "subComments": []
          }
        ],
        "hasMore": true
      }
    }
  },
  "message": "获取Feed详情成功"
}
```

### 重要字段说明

| 字段 | 描述 |
|------|------|
| `note.title` | 笔记标题 |
| `note.desc` | 笔记正文内容 |
| `note.time` | 发布时间戳（毫秒） |
| `note.imageList` | 图片列表 |
| `comments.list` | 评论列表 |
| `comments.hasMore` | 是否有更多评论 |

---

## 工具 6: 获取 Feed 列表

获取推荐的 Feed 列表。

### 请求

```bash
curl -X GET "http://10.2.248.59:18060/api/v1/feeds/list"
```

### 响应示例

```json
{
  "success": true,
  "data": {
    "feeds": [
      {
        "xsecToken": "security_token_value",
        "id": "feed_id_1",
        "modelType": "note",
        "noteCard": {
          "type": "normal",
          "displayTitle": "笔记标题",
          "user": {
            "userId": "user_id_1",
            "nickname": "用户昵称",
            "avatar": "https://example.com/avatar.jpg"
          },
          "interactInfo": {
            "liked": false,
            "likedCount": "100",
            "collected": false,
            "collectedCount": "50",
            "commentCount": "30"
          },
          "cover": {
            "url": "https://example.com/cover.jpg"
          },
          "video": {
            "capa": {
              "duration": 60
            }
          }
        }
      }
    ],
    "count": 10
  },
  "message": "获取Feeds列表成功"
}
```

### 重要字段说明

| 字段 | 描述 |
|------|------|
| `xsecToken` | 安全令牌，调用详情接口时需要 |
| `id` | Feed ID |
| `noteCard.video` | 视频信息（仅视频笔记有此字段） |
| `noteCard.video.capa.duration` | 视频时长（秒） |

---

## 典型使用流程

### 流程 1: 首次使用登录

```
1. 检查登录状态 → is_logged_in: false
2. 获取登录二维码 → 展示二维码给用户
3. 用户扫码登录
4. 检查登录状态 → is_logged_in: true
```

### 流程 2: 搜索并查看笔记详情

```
1. 搜索内容（关键词）→ 获取 feeds 列表
2. 从结果中选择感兴趣的笔记，记录 id 和 xsecToken
3. 获取 Feed 详情（feed_id, xsec_token）→ 获取完整内容和评论
```

### 流程 3: 查看用户主页

```
1. 从搜索结果或 Feed 详情中获取 user_id 和 xsecToken
2. 获取用户主页（user_id, xsec_token）→ 获取用户信息和笔记列表
```

---

## 错误处理

所有 API 在发生错误时返回统一格式：

```json
{
  "error": "错误消息",
  "code": "ERROR_CODE",
  "details": "详细错误信息"
}
```

### 常见错误代码

| 错误代码 | 描述 | 处理建议 |
|----------|------|----------|
| `STATUS_CHECK_FAILED` | 检查登录状态失败 | 检查服务是否正常运行 |
| `MISSING_KEYWORD` | 搜索缺少关键词 | 确保提供 keyword 参数 |
| `SEARCH_FEEDS_FAILED` | 搜索失败 | 检查登录状态，重试 |
| `GET_FEED_DETAIL_FAILED` | 获取详情失败 | 检查 feed_id 和 xsec_token 是否正确 |
| `GET_USER_PROFILE_FAILED` | 获取用户主页失败 | 检查 user_id 和 xsec_token 是否正确 |
| `LIST_FEEDS_FAILED` | 获取列表失败 | 检查登录状态 |

---

## 注意事项

1. **登录状态**: 大部分 API 需要有效的登录状态，建议先检查登录
2. **安全令牌**: `xsec_token` 是必需参数，从搜索结果或 Feed 列表中获取
3. **请求频率**: 避免过于频繁的请求，建议间隔 1-2 秒
4. **评论加载**: 加载全部评论可能耗时较长，建议使用 `comment_config` 控制
