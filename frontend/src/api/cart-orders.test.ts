import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { addToCart, clearCart, getCart, removeFromCart } from "./cart";
import { ApiError } from "./client";
import { createOrder, getOrder, listOrders } from "./orders";

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
    headers: init?.headers as Record<string, string>,
  };
}

async function captureError(promise: Promise<unknown>): Promise<ApiError> {
  try {
    await promise;
  } catch (error) {
    if (error instanceof ApiError) {
      return error;
    }
    throw error;
  }
  throw new Error("expected the request to fail");
}

describe("cart", () => {
  it("reads the cart with the bearer token", async () => {
    const cart = { items: [], total_items: 0, success: true };
    fetchMock.mockResolvedValue(jsonResponse(cart));

    const result = await getCart("jwt-value");

    expect(result).toEqual(cart);
    expect(lastCall().url).toBe("https://api.test/api/cart");
    expect(lastCall().headers.Authorization).toBe("Bearer jwt-value");
  });

  it("converts the camelCase input to the API's snake_case body", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ success: true, message: "Product added to cart" }),
    );

    await addToCart("jwt-value", { productId: 7, quantity: 2 });

    expect(lastCall().init?.method).toBe("POST");
    expect(lastCall().body).toEqual({ product_id: 7, quantity: 2 });
  });

  it("removes one product by query parameter", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ success: true }));

    await removeFromCart("jwt-value", 3);

    expect(lastCall().init?.method).toBe("DELETE");
    expect(lastCall().url).toBe("https://api.test/api/cart?product_id=3");
  });

  it("clears the whole cart", async () => {
    fetchMock.mockResolvedValue(jsonResponse({ success: true }));

    await clearCart("jwt-value");

    expect(lastCall().init?.method).toBe("DELETE");
    expect(lastCall().url).toBe("https://api.test/api/cart/clear");
  });
});

describe("orders", () => {
  it("posts only the items and returns the created order", async () => {
    const order = {
      id: 1,
      order_number: "ORD-1",
      total_amount: 699.99,
      total_coins_used: 700,
      status: "pending",
      items: [],
    };
    fetchMock.mockResolvedValue(jsonResponse(order, 201));

    const result = await createOrder("jwt-value", {
      items: [{ product_id: 1, quantity: 1 }],
    });

    expect(result).toEqual(order);
    expect(lastCall().url).toBe("https://api.test/api/orders");
    expect(lastCall().body).toEqual({
      items: [{ product_id: 1, quantity: 1 }],
    });
    expect(lastCall().headers.Authorization).toBe("Bearer jwt-value");
  });

  it("surfaces 402 with the server's message", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ error: "insufficient coins" }, 402),
    );

    const error = await captureError(
      createOrder("jwt-value", { items: [{ product_id: 3, quantity: 1 }] }),
    );

    expect(error.status).toBe(402);
    expect(error.message).toBe("insufficient coins");
  });

  it("surfaces 409 for stock", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ error: "insufficient stock" }, 409),
    );

    const error = await captureError(
      createOrder("jwt-value", { items: [{ product_id: 1, quantity: 999 }] }),
    );

    expect(error.status).toBe(409);
    expect(error.message).toBe("insufficient stock");
  });

  it("reads a single order and the user's order list", async () => {
    fetchMock.mockImplementation(async () => jsonResponse([]));

    await getOrder("jwt-value", 12);
    expect(lastCall().url).toBe("https://api.test/api/orders/12");

    await listOrders("jwt-value", "user-uuid");
    expect(lastCall().url).toBe("https://api.test/api/users/user-uuid/orders");
    expect(lastCall().headers.Authorization).toBe("Bearer jwt-value");
  });
});
