from app.schemas import ScriptRequest, ScriptResponse, ScriptScene


def generate_script(request: ScriptRequest) -> ScriptResponse:
    points = request.selling_points or ["价值清晰", "下单路径简单", "限时优惠"]
    per_scene = max(5, request.duration_seconds // 3)
    scenes = [
        ScriptScene(
            order=1,
            title="开场钩子",
            voiceover=f"如果你正在关注{request.product_name}，这条短视频先告诉你最值得马上了解的一点。",
            duration_seconds=per_scene,
        ),
        ScriptScene(
            order=2,
            title="核心价值",
            voiceover=f"{request.product_name}的亮点是：{'、'.join(points[:3])}。用简单直接的方式，让用户快速理解为什么值得选择。",
            duration_seconds=per_scene,
        ),
        ScriptScene(
            order=3,
            title="行动引导",
            voiceover=f"现在就选择{request.product_name}，把下一步变得更轻松、更明确。",
            duration_seconds=request.duration_seconds - per_scene * 2,
        ),
    ]
    return ScriptResponse(
        title=f"{request.product_name}口播短视频文案",
        scenes=scenes,
        total_duration_seconds=request.duration_seconds,
    )
