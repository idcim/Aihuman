CREATE TABLE IF NOT EXISTS video_tasks (
  id TEXT PRIMARY KEY,
  product_name TEXT NOT NULL,
  script TEXT NOT NULL,
  avatar_image_url TEXT NOT NULL,
  duration_seconds INTEGER NOT NULL CHECK (duration_seconds BETWEEN 15 AND 180),
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_video_tasks_created_at ON video_tasks (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_video_tasks_status ON video_tasks (status);

