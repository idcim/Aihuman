from fastapi import FastAPI

from app.ffmpeg_plan import build_talking_head_plan
from app.schemas import RenderPlanResponse, RenderRequest

app = FastAPI(title="AI Human Video Engine", version="0.1.0")


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

