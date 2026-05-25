# API Contract

All public API responses use:

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## Create Video Task

`POST /api/v1/video/create`

```json
{
  "product_name": "AI coffee coupon",
  "script": "Introduce a weekend coupon in a warm tone.",
  "avatar_image_url": "https://example.com/avatar.png",
  "duration_seconds": 60
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task id",
    "product_name": "AI coffee coupon",
    "script": "Introduce a weekend coupon in a warm tone.",
    "avatar_image_url": "https://example.com/avatar.png",
    "duration_seconds": 60,
    "status": "queued",
    "created_at": "2026-05-26T00:00:00Z",
    "updated_at": "2026-05-26T00:00:00Z"
  }
}
```

## Get Video Task

`GET /api/v1/video/{id}`

Returns the same task payload.

## AI Service Script Draft

`POST /v1/script`

```json
{
  "product_name": "AI coffee coupon",
  "selling_points": ["fast setup", "local-store friendly"],
  "tone": "warm",
  "duration_seconds": 60
}
```

## Video Engine Render Plan

`POST /v1/render/plan`

```json
{
  "avatar_image_path": "/assets/avatar.png",
  "voiceover_path": "/assets/voice.wav",
  "output_path": "/output/video.mp4",
  "duration_seconds": 60
}
```

