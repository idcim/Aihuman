import { NextResponse } from "next/server";

const services = [
  ["\u540e\u7aef\u670d\u52a1", process.env.BACKEND_API_URL ?? "http://localhost:18080", "/healthz"],
  ["AI \u670d\u52a1", process.env.AI_SERVICE_URL ?? "http://localhost:8001", "/healthz"],
  ["\u89c6\u9891\u5f15\u64ce", process.env.VIDEO_ENGINE_URL ?? "http://localhost:8002", "/healthz"],
] as const;

export async function GET() {
  const statuses = await Promise.all(
    services.map(async ([name, baseURL, path]) => {
      const url = `${baseURL}${path}`;
      try {
        const response = await fetch(url, { cache: "no-store" });
        return {
          name,
          ok: response.ok,
          url,
          detail: response.ok ? "\u5728\u7ebf" : `HTTP ${response.status}`,
        };
      } catch (error) {
        return {
          name,
          ok: false,
          url,
          detail: error instanceof Error ? error.message : "\u4e0d\u53ef\u8bbf\u95ee",
        };
      }
    }),
  );

  return NextResponse.json({ services: statuses });
}
