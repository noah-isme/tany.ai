import { cookies } from "next/headers";
import { NextResponse, type NextRequest } from "next/server";

import { ACCESS_TOKEN_COOKIE } from "@/lib/auth";
import { resolveApiUrl } from "@/lib/env";

export const dynamic = "force-dynamic";

export async function POST(request: NextRequest) {
  const formData = await request.formData();
  const file = formData.get("file");
  if (!(file instanceof File)) {
    return NextResponse.json({ error: { message: "file field is required" } }, { status: 400 });
  }

  const proxyForm = new FormData();
  proxyForm.append("file", file);

  const cookieStore = await cookies();
  const token = cookieStore.get(ACCESS_TOKEN_COOKIE)?.value ?? "";

  let response: Response;
  try {
    response = await fetch(resolveApiUrl("/api/admin/uploads"), {
      method: "POST",
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      body: proxyForm,
    });
  } catch {
    return NextResponse.json(
      { error: { message: "gagal terhubung ke server upload" } },
      { status: 502 },
    );
  }

  const contentType = response.headers.get("content-type") ?? "";
  const isJSON = contentType.includes("application/json");
  const payload = isJSON ? await response.json() : await response.text();

  if (!response.ok) {
    if (isJSON && payload && typeof payload === "object" && "error" in payload) {
      return NextResponse.json(payload, { status: response.status });
    }
    const message = typeof payload === "string" && payload ? payload : "unggahan gagal";
    return NextResponse.json({ error: { message } }, { status: response.status });
  }

  if (isJSON) {
    return NextResponse.json(payload, { status: response.status });
  }

  return NextResponse.json({ data: payload }, { status: response.status });
}
