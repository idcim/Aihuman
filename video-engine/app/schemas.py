from pydantic import BaseModel, Field


class RenderRequest(BaseModel):
    avatar_image_path: str = Field(..., min_length=1)
    voiceover_path: str = Field(..., min_length=1)
    output_path: str = Field(..., min_length=1)
    duration_seconds: int = Field(default=60, ge=15, le=180)
    task_id: str | None = None
    product_name: str | None = None
    script: str | None = None


class RenderPlanResponse(BaseModel):
    command: list[str]
    output_path: str


class RenderResponse(BaseModel):
    command: list[str]
    output_path: str
    object_key: str
    result_url: str
    preview_duration_seconds: int
