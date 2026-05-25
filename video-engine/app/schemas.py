from pydantic import BaseModel, Field


class RenderRequest(BaseModel):
    avatar_image_path: str = Field(..., min_length=1)
    voiceover_path: str = Field(..., min_length=1)
    output_path: str = Field(..., min_length=1)
    duration_seconds: int = Field(default=60, ge=15, le=180)


class RenderPlanResponse(BaseModel):
    command: list[str]
    output_path: str

