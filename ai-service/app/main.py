from fastapi import FastAPI

from app.schemas import ScriptRequest, ScriptResponse
from app.services.script_generator import generate_script

app = FastAPI(title="AI Human AI Service", version="0.1.0")


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "service": "ai-service"}


@app.post("/v1/script", response_model=ScriptResponse)
def create_script(request: ScriptRequest) -> ScriptResponse:
    return generate_script(request)

