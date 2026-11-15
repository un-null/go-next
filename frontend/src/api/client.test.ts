import { describe, expect, it, vi } from "vitest";

import { ApiError, request } from "./client";
import { captureError, fetchMock, jsonResponse, mockApi } from "./testing";

mockApi("https://api.test/");

describe("request", () => {
  it("joins the base URL and sends method, body and JSON header", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ id: 1 }));

    await request("/api/orders", { method: "POST", body: { items: [] } });

    expect(fetchMock).toHaveBeenCalledTimes(1);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe("https://api.test/api/orders");
    expect(init?.method).toBe("POST");
    expect(init?.body).toBe(JSON.stringify({ items: [] }));
    expect(init?.headers).toMatchObject({
      Accept: "application/json",
      "Content-Type": "application/json",
    });
  });

  it("reports a missing base URL instead of a network failure", async () => {
    vi.stubEnv("VITE_API_URL", "");

    const error = await request("/api/products").catch(
      (cause: unknown) => cause,
    );

    expect(error).toBeInstanceOf(Error);
    expect(error).not.toBeInstanceOf(ApiError);
    expect((error as Error).message).toBe("VITE_API_URL is not set");
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("only attaches the bearer token when one is given", async () => {
    fetchMock.mockImplementation(async () => jsonResponse({ ok: true }));

    await request("/api/products");
    await request("/api/cart", { token: "jwt-value" });

    expect(fetchMock.mock.calls[0][1]?.headers).not.toHaveProperty(
      "Authorization",
    );
    expect(fetchMock.mock.calls[1][1]?.headers).toMatchObject({
      Authorization: "Bearer jwt-value",
    });
  });

  it("parses the Echo message key", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ message: "Invalid credentials" }, 401),
    );

    const error = await captureError(request("/api/login"));

    expect(error.status).toBe(401);
    expect(error.message).toBe("Invalid credentials");
    expect(error.body).toEqual({ message: "Invalid credentials" });
  });

  it("parses the handler error key", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ error: "insufficient stock" }, 409),
    );

    const error = await captureError(request("/api/orders"));

    expect(error.status).toBe(409);
    expect(error.message).toBe("insufficient stock");
  });

  it("falls back to the status when the body is not JSON", async () => {
    fetchMock.mockResolvedValue(
      new Response("<html>oops</html>", { status: 500 }),
    );

    const error = await captureError(request("/api/products"));

    expect(error.status).toBe(500);
    expect(error.message).toBe("HTTP 500");
    expect(error.body).toBeUndefined();
  });

  it("returns undefined for 204", async () => {
    fetchMock.mockResolvedValue(new Response(null, { status: 204 }));

    await expect(
      request("/api/cart/clear", { method: "DELETE" }),
    ).resolves.toBe(undefined);
  });

  it("reports a failed connection as status 0", async () => {
    fetchMock.mockRejectedValue(new TypeError("Failed to fetch"));

    const error = await captureError(request("/api/products"));

    expect(error.status).toBe(0);
    expect(error.message).toBe("network request failed");
  });

  it("lets an abort through untouched", async () => {
    fetchMock.mockRejectedValue(new DOMException("aborted", "AbortError"));

    await expect(request("/api/products")).rejects.toBeInstanceOf(DOMException);
  });
});
