from pydantic import BaseModel, Field


class ScriptRequest(BaseModel):
    product_name: str = Field(..., min_length=1)
    selling_points: list[str] = Field(default_factory=list)
    tone: str = "warm"
    duration_seconds: int = Field(default=60, ge=15, le=180)


class ScriptScene(BaseModel):
    order: int
    title: str
    voiceover: str
    duration_seconds: int


class ScriptResponse(BaseModel):
    title: str
    scenes: list[ScriptScene]
    total_duration_seconds: int

