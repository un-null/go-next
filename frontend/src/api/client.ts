export interface RequestOptions {
  method?: "GET" | "POST" | "PATCH" | "DELETE";
  body?: unknown;
  token?: string;
  signal?: AbortSignal;
}

export class ApiError extends Error {
  readonly status: number;
  readonly body: unknown;

  constructor(status: number, message: string, body: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.body = body;
  }
}

export async function request<T>(
  path: string,
  options: RequestOptions = {},
): Promise<T> {
  const { method = "GET", body, token, signal } = options;

  const headers: Record<string, string> = { Accept: "application/json" };
  if (body !== undefined) {
    headers["Content-Type"] = "application/json";
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const url = `${baseUrl()}${path}`;

  let response: Response;
  try {
    response = await fetch(url, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
      signal,
    });
  } catch (cause) {
    if (cause instanceof DOMException && cause.name === "AbortError") {
      throw cause;
    }
    throw new ApiError(0, "network request failed", cause);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const payload = parseJson(await response.text());

  if (!response.ok) {
    throw new ApiError(
      response.status,
      errorMessage(payload, response.status),
      payload,
    );
  }

  return payload as T;
}

function baseUrl(): string {
  const url = import.meta.env.VITE_API_URL;
  if (!url) {
    throw new Error("VITE_API_URL is not set");
  }
  return url.replace(/\/$/, "");
}

function parseJson(text: string): unknown {
  if (!text) {
    return undefined;
  }
  try {
    return JSON.parse(text);
  } catch {
    return undefined;
  }
}

function errorMessage(payload: unknown, status: number): string {
  if (typeof payload === "object" && payload !== null) {
    const record = payload as Record<string, unknown>;
    for (const key of ["error", "message"]) {
      const value = record[key];
      if (typeof value === "string" && value) {
        return value;
      }
    }
  }
  return `HTTP ${status}`;
}
