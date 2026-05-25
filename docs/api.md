# API 契约

所有公开 API 响应统一使用：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 后端入口

`GET /`

返回轻量 HTML 服务入口页，包含健康检查和主要路由说明。

## 中文运营后台

前端后台由 `frontend` 服务提供：

```text
http://localhost:3000
```

后台可创建任务、查看任务列表、查看任务详情、打开 OSS 成果链接。

## 创建视频任务

`POST /api/v1/video/create`

```json
{
  "product_name": "新品咖啡券",
  "script": "用温暖、有活力的语气介绍周末优惠。",
  "avatar_image_url": "https://example.com/avatar.png",
  "duration_seconds": 60
}
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task id",
    "product_name": "新品咖啡券",
    "script": "用温暖、有活力的语气介绍周末优惠。",
    "avatar_image_url": "https://example.com/avatar.png",
    "duration_seconds": 60,
    "status": "queued",
    "created_at": "2026-05-26T00:00:00Z",
    "updated_at": "2026-05-26T00:00:00Z"
  }
}
```

## 查询视频任务

`GET /api/v1/video/{id}`

成功完成后会返回渲染命令和 OSS 成果信息：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task id",
    "status": "succeeded",
    "render_command": ["ffmpeg", "..."],
    "output_url": "http://localhost:19000/aihuman-assets/renders/task-id.mp4",
    "object_key": "renders/task-id.mp4",
    "preview_seconds": 15,
    "created_at": "2026-05-26T00:00:00Z",
    "updated_at": "2026-05-26T00:00:10Z"
  }
}
```

## 任务列表

`GET /api/v1/video`

按创建时间倒序返回任务记录。

任务状态：

- `queued`：已进入 Redis Stream 队列。
- `running`：worker 正在处理。
- `succeeded`：已生成预览视频并上传 OSS。
- `failed`：处理失败，`error_message` 会包含失败原因。

## AI 文案服务

`POST /v1/script`

```json
{
  "product_name": "新品咖啡券",
  "selling_points": ["快速领取", "周末可用"],
  "tone": "温暖",
  "duration_seconds": 60
}
```

## 视频渲染服务

### 生成渲染计划

`POST /v1/render/plan`

```json
{
  "avatar_image_path": "/assets/avatar.png",
  "voiceover_path": "/assets/voice.wav",
  "output_path": "/output/video.mp4",
  "duration_seconds": 60
}
```

### 生成并上传预览视频

`POST /v1/render`

```json
{
  "task_id": "task id",
  "product_name": "新品咖啡券",
  "script": "用温暖、有活力的语气介绍周末优惠。",
  "avatar_image_path": "https://example.com/avatar.png",
  "voiceover_path": "/tmp/task.wav",
  "output_path": "/output/task.mp4",
  "duration_seconds": 60
}
```
