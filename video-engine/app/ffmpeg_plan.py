from dataclasses import dataclass


@dataclass(frozen=True)
class RenderPlan:
    command: list[str]
    output_path: str


def build_talking_head_plan(
    avatar_image_path: str,
    voiceover_path: str,
    output_path: str,
    duration_seconds: int,
) -> RenderPlan:
    if duration_seconds < 15 or duration_seconds > 180:
        raise ValueError("duration_seconds must be between 15 and 180")

    command = [
        "ffmpeg",
        "-y",
        "-loop",
        "1",
        "-i",
        avatar_image_path,
        "-i",
        voiceover_path,
        "-t",
        str(duration_seconds),
        "-vf",
        "scale=1080:1920:force_original_aspect_ratio=increase,crop=1080:1920",
        "-c:v",
        "libx264",
        "-preset",
        "veryfast",
        "-c:a",
        "aac",
        "-shortest",
        output_path,
    ]
    return RenderPlan(command=command, output_path=output_path)

