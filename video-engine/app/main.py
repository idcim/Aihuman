from fastapi import FastAPI

from app.ffmpeg_plan import build_talking_head_plan
from app.schemas import RenderPlanResponse, RenderRequest, RenderResponse
from app.storage import ObjectStorage
from app.video_render import render_preview

app = FastAPI(title="AI Human Video Engine", version="0.1.0")
storage = ObjectStorage()


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "video-engine"}


@app.post("/v1/render/plan", response_model=RenderPlanResponse)
def create_render_plan(request: RenderRequest) -> RenderPlanResponse:
    plan = build_talking_head_plan(
        avatar_image_path=request.avatar_image_path,
        voiceover_path=request.voiceover_path,
        output_path=request.output_path,
        duration_seconds=request.duration_seconds,
    )
    return RenderPlanResponse(command=plan.command, output_path=plan.output_path)


@app.post("/v1/render", response_model=RenderResponse)
def render_video(request: RenderRequest) -> RenderResponse:
    rendered = render_preview(request)
    object_key = f"renders/{request.task_id or 'preview'}.mp4"
    result_url = storage.upload_file(rendered.output_path, object_key)
    return RenderResponse(
        command=rendered.command,
        output_path=str(rendered.output_path),
        object_key=object_key,
        result_url=result_url,
        preview_duration_seconds=rendered.preview_duration_seconds,
    )
