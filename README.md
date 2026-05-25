# AI Human Video Factory

AI 自动剪辑与数字人口播视频工厂。当前仓库先聚焦 Phase 1 MVP：用户输入商品、文案和人物图片，系统生成 60 秒口播短视频任务。

## MVP Scope

- 文案生成与脚本标准化
- 数字人口播任务编排
- FFmpeg 视频渲染服务契约
- 统一任务状态 API
- Docker Compose 本地开发环境

暂不优先开发 3D 数字人、实时直播、复杂时间线编辑器、自动发布和视频矩阵。

## Repository Layout

```text
backend/          Go + Gin main API and task orchestration
ai-service/       Python + FastAPI script and storyboard service
video-engine/     Python + FastAPI FFmpeg render planning service
docs/             Architecture, API, and MVP roadmap
docker-compose.yml
.env.example
```

## Quick Start

1. Copy `.env.example` to `.env` and adjust values if needed.
2. Start local infrastructure and services:

```powershell
docker compose up --build
```

3. Create a video task:

```powershell
Invoke-RestMethod -Method Post http://localhost:18080/api/v1/video/create `
  -ContentType 'application/json' `
  -Body '{"product_name":"AI coffee coupon","script":"Introduce a weekend coupon in a warm tone.","avatar_image_url":"https://example.com/avatar.png","duration_seconds":60}'
```

4. Check health:

```powershell
Invoke-RestMethod http://localhost:18080/healthz
Invoke-RestMethod http://localhost:8001/healthz
Invoke-RestMethod http://localhost:8002/healthz
```

## Service Ports

| Service | Port |
| --- | --- |
| backend | 18080 -> 8080 |
| ai-service | 8001 |
| video-engine | 8002 |
| PostgreSQL | 15432 -> 5432 |
| Redis | 16379 -> 6379 |
| Qdrant | 16333 -> 6333 |
| MinIO API | 19000 -> 9000 |
| MinIO Console | 19001 -> 9001 |

## Development Notes

- API responses use the unified `{ "code": 0, "message": "success", "data": {} }` envelope.
- Backend controllers should stay thin; orchestration belongs in services.
- AI and video processing stay outside the API layer behind HTTP contracts.
- No secrets should be hard-coded. Use environment variables.
