export type ApiEnvelope<T> = {
  code: number;
  message: string;
  data: T;
};

export type VideoTask = {
  id: string;
  product_name: string;
  script: string;
  avatar_image_url: string;
  duration_seconds: number;
  status: "queued" | "running" | "succeeded" | "failed";
  render_command?: string[];
  output_url?: string;
  object_key?: string;
  preview_seconds?: number;
  error_message?: string;
  created_at: string;
  updated_at: string;
};

export type ScriptScene = {
  order: number;
  title: string;
  voiceover: string;
  duration_seconds: number;
};

export type ScriptResponse = {
  title: string;
  scenes: ScriptScene[];
  total_duration_seconds: number;
};

export type HealthStatus = {
  name: string;
  ok: boolean;
  url: string;
  detail: string;
};

export async function requestJSON<T>(url: string, init?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || `Request failed with ${response.status}`);
  }

  return response.json() as Promise<T>;
}
