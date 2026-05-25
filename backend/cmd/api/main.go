package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/idcim/aihuman/backend/internal/api"
	"github.com/idcim/aihuman/backend/internal/config"
	"github.com/idcim/aihuman/backend/internal/video"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	store, closeStore := configureStore(context.Background(), cfg.DatabaseURL, logger)
	defer closeStore()
	service := video.NewService(store, logger)

	queue, closeQueue := configureQueue(context.Background(), cfg.RedisAddr, cfg.RedisStream, logger)
	defer closeQueue()
	if queue != nil {
		service.UseQueue(queue)
		if cfg.WorkerEnabled {
			worker := video.NewWorker(service, queue, cfg.VideoEngineURL, logger)
			worker.Start(context.Background())
		}
	}

	handler := video.NewHandler(service)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(logger))

	router.GET("/healthz", func(c *gin.Context) {
		api.Success(c, gin.H{"status": "ok", "service": "backend"})
	})
	router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(homePageHTMLUTF8))
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("/video/create", handler.Create)
		v1.GET("/video", handler.List)
		v1.GET("/video/:id", handler.Get)
	}

	addr := ":" + cfg.Port
	logger.Info("backend listening", "addr", addr)
	if err := router.Run(addr); err != nil && err != http.ErrServerClosed {
		logger.Error("backend stopped", "error", err)
		os.Exit(1)
	}
}

func configureStore(ctx context.Context, databaseURL string, logger *slog.Logger) (video.Store, func()) {
	if databaseURL == "" {
		logger.Info("using memory task store")
		return video.NewMemoryStore(), func() {}
	}

	store, err := video.NewPostgresStore(ctx, databaseURL)
	if err != nil {
		logger.Error("failed to initialize postgres task store", "error", err)
		os.Exit(1)
	}

	logger.Info("using postgres task store")
	return store, func() {
		if err := store.Close(); err != nil {
			logger.Warn("failed to close postgres task store", "error", err)
		}
	}
}

func configureQueue(ctx context.Context, redisAddr string, redisStream string, logger *slog.Logger) (video.Queue, func()) {
	if redisAddr == "" {
		logger.Info("video queue disabled")
		return nil, func() {}
	}

	queue, err := video.NewRedisQueue(ctx, redisAddr, redisStream, logger)
	if err != nil {
		logger.Error("failed to initialize redis queue", "error", err)
		os.Exit(1)
	}

	return queue, func() {
		if err := queue.Close(); err != nil {
			logger.Warn("failed to close redis queue", "error", err)
		}
	}
}

const homePageHTMLUTF8 = "<!doctype html>" +
	"<html lang=\"zh-CN\">" +
	"<head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">" +
	"<title>AI \u6570\u5b57\u4eba\u89c6\u9891\u5de5\u5382</title>" +
	"<style>body{margin:0;min-height:100vh;background:#f7f8fb;color:#151923;font-family:ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif}main{width:min(960px,calc(100vw - 32px));margin:0 auto;padding:40px 0}h1{font-size:40px;line-height:1.1;margin:0 0 12px}p{color:#667085;line-height:1.7}.status{display:inline-block;margin:18px 0;padding:10px 14px;border:1px solid #99d8cf;border-radius:8px;background:#ecfdf9;color:#115e59;font-weight:700}.card{background:#fff;border:1px solid #d9dee8;border-radius:8px;padding:18px;margin-top:16px}a{color:#115e59;font-weight:700;text-decoration:none}code{display:block;overflow:auto;background:#101828;color:#fff;border-radius:8px;padding:14px;margin-top:10px}</style>" +
	"</head><body><main>" +
	"<h1>AI \u6570\u5b57\u4eba\u89c6\u9891\u5de5\u5382</h1>" +
	"<p>\u7b2c\u4e00\u9636\u6bb5 MVP \u540e\u7aef\u670d\u52a1\u5df2\u8fd0\u884c\uff0c\u5e76\u5df2\u63a5\u5165 PostgreSQL \u548c Redis Stream \u4efb\u52a1\u961f\u5217\u3002</p>" +
	"<div class=\"status\">\u540e\u7aef\u5728\u7ebf</div>" +
	"<section class=\"card\"><h2>\u540e\u53f0\u64cd\u4f5c\u754c\u9762</h2><p><a href=\"http://localhost:3000\">http://localhost:3000</a></p></section>" +
	"<section class=\"card\"><h2>\u5065\u5eb7\u68c0\u67e5</h2><p><a href=\"/healthz\">/healthz</a></p></section>" +
	"<section class=\"card\"><h2>\u53ef\u7528 API</h2><code>POST /api/v1/video/create\nGET  /api/v1/video\nGET  /api/v1/video/:id</code></section>" +
	"</main></body></html>"

const homePageHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>AI 数字人视频工厂</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f7f8fb;
      --panel: #ffffff;
      --text: #151923;
      --muted: #667085;
      --line: #d9dee8;
      --accent: #0f766e;
      --accent-strong: #115e59;
      --code: #101828;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      background: var(--bg);
      color: var(--text);
      font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    main {
      width: min(1120px, calc(100vw - 32px));
      margin: 0 auto;
      padding: 40px 0;
    }
    header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 24px;
      padding-bottom: 28px;
      border-bottom: 1px solid var(--line);
    }
    h1 {
      margin: 0 0 10px;
      font-size: clamp(30px, 4vw, 48px);
      line-height: 1.08;
      letter-spacing: 0;
    }
    p { margin: 0; color: var(--muted); line-height: 1.65; }
    .status {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      flex: 0 0 auto;
      padding: 10px 14px;
      border: 1px solid #99d8cf;
      border-radius: 8px;
      background: #ecfdf9;
      color: var(--accent-strong);
      font-weight: 700;
      white-space: nowrap;
    }
    .dot {
      width: 9px;
      height: 9px;
      border-radius: 50%;
      background: var(--accent);
    }
    section {
      margin-top: 28px;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 16px;
    }
    .card {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      padding: 18px;
    }
    h2 {
      margin: 0 0 12px;
      font-size: 19px;
      letter-spacing: 0;
    }
    h3 {
      margin: 0 0 8px;
      font-size: 16px;
      letter-spacing: 0;
    }
    a {
      color: var(--accent-strong);
      font-weight: 700;
      text-decoration: none;
    }
    a:hover { text-decoration: underline; }
    code {
      display: block;
      overflow-x: auto;
      margin-top: 12px;
      padding: 14px;
      border-radius: 8px;
      background: var(--code);
      color: #f9fafb;
      font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", monospace;
      font-size: 13px;
      line-height: 1.55;
      white-space: pre;
    }
    .routes {
      display: grid;
      grid-template-columns: 160px 1fr;
      gap: 10px 14px;
      margin-top: 10px;
    }
    .method {
      color: var(--accent-strong);
      font-weight: 800;
      font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", monospace;
    }
    .path {
      color: var(--text);
      font-family: ui-monospace, SFMono-Regular, Consolas, "Liberation Mono", monospace;
      overflow-wrap: anywhere;
    }
    @media (max-width: 780px) {
      main { padding: 28px 0; }
      header { display: block; }
      .status { margin-top: 18px; }
      .grid { grid-template-columns: 1fr; }
      .routes { grid-template-columns: 1fr; }
    }
  </style>
</head>
<body>
  <main>
    <header>
      <div>
        <h1>AI 数字人视频工厂</h1>
        <p>第一阶段 MVP 后端服务已运行。当前支持视频任务 API、AI 文案服务契约和视频渲染计划契约。</p>
      </div>
      <div class="status"><span class="dot"></span>后端在线</div>
    </header>

    <section class="grid">
      <div class="card">
        <h2>后端 API</h2>
        <p>创建和查看短视频生成任务。</p>
        <p><a href="/healthz">打开健康检查</a></p>
      </div>
      <div class="card">
        <h2>AI 服务</h2>
        <p>文案和分镜服务运行在 8001 端口。</p>
        <p><a href="http://localhost:8001/healthz">打开 AI 健康检查</a></p>
      </div>
      <div class="card">
        <h2>视频引擎</h2>
        <p>渲染计划服务运行在 8002 端口。</p>
        <p><a href="http://localhost:8002/healthz">打开视频引擎健康检查</a></p>
      </div>
    </section>

    <section class="card">
      <h2>可用接口</h2>
      <div class="routes">
        <div class="method">GET</div><div class="path">/healthz</div>
        <div class="method">POST</div><div class="path">/api/v1/video/create</div>
        <div class="method">GET</div><div class="path">/api/v1/video</div>
        <div class="method">GET</div><div class="path">/api/v1/video/:id</div>
      </div>
    </section>

    <section class="card">
      <h2>创建测试任务</h2>
      <p>在 PowerShell 中运行这条命令，可以创建一个排队中的 60 秒数字人口播任务。</p>
      <code>Invoke-RestMethod -Method Post http://localhost:18080/api/v1/video/create -ContentType 'application/json' -Body '{"product_name":"周末咖啡团购券","script":"介绍周末咖啡团购券的核心优惠，并用温暖自然的语气引导用户下单。","avatar_image_url":"https://example.com/avatar.png","duration_seconds":60}'</code>
    </section>
  </main>
</body>
</html>`

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logger.Info(
			"http_request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"request_id", c.GetHeader("X-Request-Id"),
		)
	}
}
