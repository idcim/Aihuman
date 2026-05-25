import { NextRequest, NextResponse } from "next/server";

const aiServiceURL = process.env.AI_SERVICE_URL ?? "http://localhost:8001";

export async function POST(request: NextRequest) {
  const response = await fetch(`${aiServiceURL}/v1/script`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: await request.text(),
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

