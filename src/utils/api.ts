import { resolveCredentials, type Resolved } from "./config";

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string,
    public retryAfter?: number,
    public details?: unknown,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

export type Client = {
  resolved: Resolved;
  get: <T = unknown>(path: string, query?: Record<string, string>) => Promise<T>;
  post: <T = unknown>(path: string, body?: unknown) => Promise<T>;
  patch: <T = unknown>(path: string, body?: unknown) => Promise<T>;
  del: <T = unknown>(path: string) => Promise<T>;
};

const DEFAULT_TIMEOUT_MS = 60_000;

/** Encode a single URL path segment; reject blank values. */
export function pathSegment(raw: string): string {
  const value = raw?.trim() ?? "";
  if (!value) throw new Error("A non-empty path value is required.");
  return encodeURIComponent(value);
}

function publicPath(baseURL: string, path: string): string {
  const p = path.startsWith("/") ? path : `/${path}`;
  return `${baseURL}/api/v1/public${p}`;
}

function parseRetryAfter(res: Response, data: any): number | undefined {
  const header = res.headers.get("retry-after");
  if (header) {
    const n = Number(header);
    if (Number.isFinite(n) && n >= 0) return n;
  }
  const body = data?.error?.retry_after ?? data?.retry_after;
  if (typeof body === "number" && Number.isFinite(body)) return body;
  return undefined;
}

async function request<T>(
  resolved: Resolved,
  method: string,
  path: string,
  body?: unknown,
  query?: Record<string, string>,
  timeoutMs = DEFAULT_TIMEOUT_MS,
): Promise<T> {
  const url = new URL(publicPath(resolved.baseURL, path));
  if (query) {
    for (const [k, v] of Object.entries(query)) url.searchParams.set(k, v);
  }

  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeoutMs);

  let res: Response;
  try {
    res = await fetch(url, {
      method,
      headers: {
        Authorization: `Bearer ${resolved.apiKey}`,
        Accept: "application/json",
        "User-Agent": "praxicraft-assess-cli",
        ...(body !== undefined ? { "Content-Type": "application/json" } : {}),
      },
      body: body !== undefined ? JSON.stringify(body) : undefined,
      signal: controller.signal,
    });
  } catch (e: any) {
    if (e?.name === "AbortError") {
      throw new ApiError(`Request timed out after ${Math.round(timeoutMs / 1000)}s`, 0, "TIMEOUT");
    }
    const code = e?.code ?? e?.cause?.code;
    if (code === "ECONNREFUSED" || code === "ENOTFOUND" || code === "EAI_AGAIN") {
      throw new ApiError("Couldn't reach the Assess API. Check your network and base URL.", 0, code);
    }
    if (code === "ECONNRESET" || /socket hang up/i.test(String(e?.message))) {
      throw new ApiError("Connection reset. Try again in a moment.", 0, "ECONNRESET");
    }
    throw e;
  } finally {
    clearTimeout(timer);
  }

  const text = await res.text();
  let data: any = null;
  try {
    data = text ? JSON.parse(text) : null;
  } catch {
    data = text ? { raw: text.slice(0, 500) } : null;
  }

  if (!res.ok) {
    const err = data?.error;
    const message =
      err?.message ||
      (typeof data?.raw === "string" && data.raw
        ? `HTTP ${res.status}: ${data.raw.slice(0, 160)}`
        : `HTTP ${res.status}`);
    throw new ApiError(
      message,
      res.status,
      err?.code,
      parseRetryAfter(res, data),
      err?.details ?? err?.required_plan,
    );
  }

  return data as T;
}

export function createClient(opts?: {
  profile?: string;
  apiKey?: string;
  baseURL?: string;
}): Client {
  const resolved = resolveCredentials(opts);
  return {
    resolved,
    get: (path, query) => request(resolved, "GET", path, undefined, query),
    post: (path, body) => request(resolved, "POST", path, body),
    patch: (path, body) => request(resolved, "PATCH", path, body),
    del: (path) => request(resolved, "DELETE", path),
  };
}

export async function assistantChat(
  resolved: Resolved,
  payload: {
    messages: Array<Record<string, unknown>>;
    tools?: Array<Record<string, unknown>>;
    max_tokens?: number;
  },
): Promise<{ message: Record<string, any>; finish_reason: string }> {
  return request(resolved, "POST", "/assistant/chat/", payload);
}
