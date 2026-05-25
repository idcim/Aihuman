# AI 自动剪辑与数字人视频工厂

这是一个本地可运行的 AI 自动剪辑与数字人口播视频工厂 MVP。当前版本已经具备中文后台、任务编排、Redis 队列、FFmpeg 预览成片、MinIO/OSS 成果存储和成果链接查看能力。

## 已支持能力

- 中文运营后台：创建视频任务、查看任务状态、查看渲染命令和成果链接。
- 后端编排：Go + Gin API，PostgreSQL 持久化，Redis Stream 任务队列。
- AI 文案服务：FastAPI 生成中文口播脚本草稿。
- 视频引擎：FastAPI 调用 FFmpeg 生成竖屏 MP4 预览视频。
- OSS 存储：通过 S3 兼容协议上传到 MinIO，并返回可访问的成果 URL。
- Docker Compose：一条命令启动前端、后端、AI 服务、视频引擎、PostgreSQL、Redis、Qdrant、MinIO。

## 目录结构

```text
frontend/         Next.js 中文运营后台
backend/          Go + Gin 主 API、任务编排和 worker
ai-service/       Python + FastAPI 文案服务
video-engine/     Python + FastAPI FFmpeg 渲染服务
docs/             架构、API 和路线文档
docker-compose.yml
.env.example
```

## 快速启动

1. 复制环境变量文件：

```powershell
Copy-Item .env.example .env
```

2. 启动完整本地环境：

```powershell
docker compose up --build
```

3. 打开中文后台：

```text
http://localhost:3000
```

4. 后端服务入口：

```text
http://localhost:18080
```

5. MinIO 控制台：

```text
http://localhost:19001
```

默认账号密码来自 `.env.example`：

```text
MINIO_ROOT_USER=aihuman
MINIO_ROOT_PASSWORD=aihuman_dev_password
```

## 服务端口

| 服务 | 端口 |
| --- | --- |
| frontend | 3000 |
| backend | 18080 -> 8080 |
| ai-service | 8001 |
| video-engine | 8002 |
| PostgreSQL | 15432 -> 5432 |
| Redis | 16379 -> 6379 |
| Qdrant | 16333 -> 6333 |
| MinIO API | 19000 -> 9000 |
| MinIO Console | 19001 -> 9001 |

## 开发说明

- API 统一返回 `{ "code": 0, "message": "success", "data": {} }`。
- 后台通过 `frontend/app/api` 代理访问后端和 AI 服务，浏览器无需直接跨域访问容器服务。
- 新建视频任务后，后端会写入 PostgreSQL，并推送到 Redis Stream。
- worker 消费任务后调用 video-engine 的 `/v1/render`，生成 MP4 并上传到 OSS。
- 成果文件默认存储在 MinIO bucket `aihuman-assets` 下的 `renders/` 前缀。
