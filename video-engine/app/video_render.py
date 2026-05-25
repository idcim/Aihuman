import subprocess
from dataclasses import dataclass
from pathlib import Path

from app.schemas import RenderRequest

FONT_FILE = "/usr/share/fonts/TTF/DejaVuSans.ttf"


@dataclass(frozen=True)
class RenderedVideo:
    command: list[str]
    output_path: Path
    preview_duration_seconds: int


def render_preview(request: RenderRequest) -> RenderedVideo:
    task_id = _safe_id(request.task_id or "preview")
    output_dir = Path("/tmp/aihuman-renders")
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / f"{task_id}.mp4"

    # Keep the MVP preview fast while proving the end-to-end clipping pipeline.
    preview_duration = min(max(request.duration_seconds, 5), 15)
    command = [
        "ffmpeg",
        "-y",
        "-f",
        "lavfi",
        "-i",
        f"testsrc2=size=720x1280:rate=30:duration={preview_duration}",
        "-f",
        "lavfi",
        "-i",
        f"sine=frequency=880:duration={preview_duration}",
        "-vf",
        "drawbox=x=0:y=0:w=iw:h=150:color=black@0.55:t=fill,"
        f"drawtext=fontfile={FONT_FILE}:text='AI HUMAN VIDEO FACTORY':fontcolor=white:fontsize=36:x=40:y=44,"
        "drawbox=x=0:y=1130:w=iw:h=150:color=black@0.55:t=fill,"
        f"drawtext=fontfile={FONT_FILE}:text='AUTO CLIP MVP PREVIEW':fontcolor=white:fontsize=32:x=40:y=1184",
        "-c:v",
        "libx264",
        "-preset",
        "ultrafast",
        "-pix_fmt",
        "yuv420p",
        "-c:a",
        "aac",
        "-shortest",
        str(output_path),
    ]
    try:
        subprocess.run(command, check=True, capture_output=True, text=True)
    except subprocess.CalledProcessError as error:
        stderr = (error.stderr or "").strip()
        raise RuntimeError(f"ffmpeg render failed: {stderr[-1200:]}") from error
    return RenderedVideo(
        command=command,
        output_path=output_path,
        preview_duration_seconds=preview_duration,
    )


def _safe_id(value: str) -> str:
    return "".join(char for char in value if char.isalnum() or char in {"-", "_"})[:80] or "preview"
