from app.schemas import ScriptRequest, ScriptResponse, ScriptScene


def generate_script(request: ScriptRequest) -> ScriptResponse:
    points = request.selling_points or ["clear value", "simple purchase path", "limited-time offer"]
    per_scene = max(5, request.duration_seconds // 3)
    scenes = [
        ScriptScene(
            order=1,
            title="hook",
            voiceover=f"Looking for {request.product_name}? Here is the key reason to try it today.",
            duration_seconds=per_scene,
        ),
        ScriptScene(
            order=2,
            title="value",
            voiceover=f"{request.product_name} helps with {', '.join(points[:3])}.",
            duration_seconds=per_scene,
        ),
        ScriptScene(
            order=3,
            title="call_to_action",
            voiceover=f"Choose {request.product_name} now and make the next step easier.",
            duration_seconds=request.duration_seconds - per_scene * 2,
        ),
    ]
    return ScriptResponse(
        title=f"{request.product_name} short video script",
        scenes=scenes,
        total_duration_seconds=request.duration_seconds,
    )

