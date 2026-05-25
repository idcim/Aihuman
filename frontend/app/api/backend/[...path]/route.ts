import { NextRequest, NextResponse } from "next/server";

const backendURL = process.env.BACKEND_API_URL ?? "http://localhost:18080";

type RouteContext = {
  params: Promise<{ path: string[] }>;
};

export async function GET(_request: NextRequest, context: RouteContext) {
  return proxy(context);
}

export async function POST(request: NextRequest, context: RouteContext) {
  return proxy(context, await request.text());
}

async function proxy(context: RouteContext, body?: string) {
  const { path } = await context.params;
  const target = `${backendURL}/api/v1/${path.join("/")}`;
  const response = await fetch(target, {
    method: body === undefined ? "GET" : "POST",
    headers: body === undefined ? undefined : { "Content-Type": "application/json" },
    body,
    cache: "no-store",
  });

  const text = await response.text();
  return new NextResponse(text, {
    status: response.status,
    headers: {
      "Content-Type": response.headers.get("Content-Type") ?? "application/json",
    },
  });
}

