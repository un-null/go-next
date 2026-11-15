import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { login, signup } from "./auth";
import { getProduct, listCategories, listProducts } from "./catalog";

const fetchMock = vi.fn<typeof fetch>();

beforeEach(() => {
  vi.stubEnv("VITE_API_URL", "https://api.test");
  fetchMock.mockReset();
  vi.stubGlobal("fetch", fetchMock);
});

afterEach(() => {
  vi.unstubAllEnvs();
  vi.unstubAllGlobals();
});

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

function lastCall() {
  const [url, init] = fetchMock.mock.calls[fetchMock.mock.calls.length - 1];
  return {
    url,
    init,
    body: init?.body ? JSON.parse(String(init.body)) : undefined,
  };
}

describe("auth", () => {
  it("posts credentials and returns the token envelope", async () => {
    const envelope = { token: "jwt-value", user: { id: "u1", coins: 1000 } };
    fetchMock.mockResolvedValue(jsonResponse(envelope));

    const result = await login({
      email: "a@example.com",
      password: "password123",
    });

    expect(result).toEqual(envelope);
    expect(lastCall().url).toBe("https://api.test/api/login");
    expect(lastCall().init?.method).toBe("POST");
    expect(lastCall().body).toEqual({
      email: "a@example.com",
      password: "password123",
    });
  });

  it("unwraps the user from the signup envelope", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse(
        {
          message: "User created successfully",
          user: { id: "u1", coins: 1000 },
        },
        201,
      ),
    );

    const user = await signup({
      name: "A",
      email: "a@example.com",
      password: "password123",
    });

    expect(user).toEqual({ id: "u1", coins: 1000 });
    expect(lastCall().url).toBe("https://api.test/api/signup");
    expect(lastCall().body).toEqual({
      name: "A",
      email: "a@example.com",
      password: "password123",
    });
  });
});

describe("catalog", () => {
  it("builds the pagination query and returns the bare array", async () => {
    const products = [{ id: 3 }, { id: 2 }];
    fetchMock.mockResolvedValue(jsonResponse(products));

    const result = await listProducts({ page: 2, limit: 20 });

    expect(result).toEqual(products);
    expect(lastCall().url).toBe(
      "https://api.test/api/products?page=2&limit=20",
    );
    expect(lastCall().body).toBeUndefined();
  });

  it("defaults to the server's page 1 and limit 10", async () => {
    fetchMock.mockResolvedValue(jsonResponse([]));

    await listProducts();

    expect(lastCall().url).toBe(
      "https://api.test/api/products?page=1&limit=10",
    );
  });

  it("addresses a single product by id", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ id: 3 }));

    await getProduct(3);

    expect(lastCall().url).toBe("https://api.test/api/products/3");
  });

  it("returns categories as they arrive", async () => {
    const categories = [{ id: 1, name: "Electronics" }];
    fetchMock.mockResolvedValue(jsonResponse(categories));

    const result = await listCategories();

    expect(result).toEqual(categories);
    expect(lastCall().url).toBe("https://api.test/api/categories");
  });

  it("sends no Authorization header on public endpoints", async () => {
    fetchMock.mockResolvedValue(jsonResponse([]));

    await listCategories();

    expect(lastCall().init?.headers).not.toHaveProperty("Authorization");
  });
});
