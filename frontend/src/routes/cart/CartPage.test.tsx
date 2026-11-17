// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router";
import { beforeEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, lastCall, mockApi } from "../../api/testing";
import type {
  CartItem,
  CartResponse,
  Order,
  Product,
  User,
} from "../../api/types";
import { getSession, setSession } from "../../lib/session";
import { CartPage } from "./CartPage";

mockApi();

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-01T00:00:00Z",
  updated_at: "2026-11-01T00:00:00Z",
};

function product(id: number, name: string, price: number): Product {
  return {
    id,
    category_id: 1,
    name,
    description: "",
    price,
    stock_quantity: 10,
    image_url: "",
    average_rating: 4,
    total_comments: 0,
    created_at: "2026-11-01T00:00:00Z",
    updated_at: "2026-11-01T00:00:00Z",
  };
}

function line(productId: number, quantity: number): CartItem {
  return {
    id: productId,
    user_id: USER.id,
    product_id: productId,
    quantity,
    created_at: "2026-11-01T00:00:00Z",
    updated_at: "2026-11-01T00:00:00Z",
  };
}

function cart(items: CartItem[]): CartResponse {
  return { items, total_items: items.length, success: true };
}

const PRODUCTS = [product(1, "Monitor", 399.99), product(2, "Earbuds", 149.99)];

const ORDER: Order = {
  id: 7,
  user_id: USER.id,
  order_number: "ORD-7",
  total_amount: 799.98,
  total_coins_used: 800,
  status: "completed",
  created_at: "2026-11-17T00:00:00Z",
  updated_at: "2026-11-17T00:00:00Z",
};

function serve(
  cartResponse: CartResponse,
  order: () => Response = () => jsonResponse(ORDER, 201),
) {
  fetchMock.mockImplementation(async (input, init) => {
    const url = String(input);
    if (url.includes("/api/products")) {
      return jsonResponse(PRODUCTS);
    }
    if (url.includes("/api/orders")) {
      return order();
    }
    if (url.includes("/api/cart")) {
      return jsonResponse(cartResponse);
    }
    throw new Error(`unexpected request: ${url} ${init?.method}`);
  });
}

function renderCart() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <MantineProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={["/cart"]}>
          <Routes>
            <Route path="/cart" element={<CartPage />} />
            <Route path="/products" element={<p>the catalogue</p>} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

beforeEach(() => {
  setSession({ token: "jwt-value", user: USER });
});

describe("CartPage", () => {
  it("says the cart is empty and offers the catalogue", async () => {
    serve(cart([]));
    renderCart();

    expect(await screen.findByText("Your cart is empty.")).toBeTruthy();
    expect(
      screen
        .getByRole("link", { name: "Browse products" })
        .getAttribute("href"),
    ).toBe("/products");
  });

  it("shows each line with its own total and the estimated total", async () => {
    serve(cart([line(1, 2)]));
    renderCart();

    expect(await screen.findByText("Monitor")).toBeTruthy();
    expect(screen.getByText("$399.99 each")).toBeTruthy();
    // Line total and estimated total are the same number here, so both match.
    expect(screen.getAllByText("$799.98")).toHaveLength(2);
    expect(screen.getByText("800")).toBeTruthy();
  });

  it("sends the new quantity, not a difference", async () => {
    serve(cart([line(1, 2)]));
    renderCart();

    const user = userEvent.setup();
    await user.click(
      await screen.findByRole("button", { name: "Add one Monitor" }),
    );

    await waitFor(() => {
      expect(lastCall().init?.method).toBe("POST");
      expect(lastCall().body).toEqual({ product_id: 1, quantity: 3 });
    });
    expect(lastCall().headers.Authorization).toBe("Bearer jwt-value");
    expect(screen.getByText("3")).toBeTruthy();
  });

  it("removes the line when the last one is taken away", async () => {
    serve(cart([line(1, 1)]));
    renderCart();

    const user = userEvent.setup();
    await user.click(
      await screen.findByRole("button", { name: "Remove Monitor" }),
    );

    await waitFor(() => {
      expect(lastCall().init?.method).toBe("DELETE");
      expect(lastCall().url).toBe("https://api.test/api/cart?product_id=1");
    });
    expect(await screen.findByText("Your cart is empty.")).toBeTruthy();
  });

  it("shows a product that left the catalogue as a removable line", async () => {
    serve(cart([line(99, 1)]));
    renderCart();

    expect(
      await screen.findByText("This item is no longer for sale."),
    ).toBeTruthy();
    expect(
      screen.getByText("One item is no longer for sale and is not counted."),
    ).toBeTruthy();
  });

  it("shows the server's message when the cart cannot be read", async () => {
    serve(cart([line(1, 1)]));
    fetchMock.mockImplementation(async (input) => {
      const url = String(input);
      if (url.includes("/api/products")) {
        return jsonResponse(PRODUCTS);
      }
      return jsonResponse({ message: "cart is unavailable" }, 500);
    });
    renderCart();

    expect(await screen.findByText("cart is unavailable")).toBeTruthy();
    expect(screen.getByRole("button", { name: "Try again" })).toBeTruthy();
  });

  it("sends every line to the server and keeps the server's numbers", async () => {
    serve(cart([line(1, 2)]));
    renderCart();

    const user = userEvent.setup();
    await user.click(
      await screen.findByRole("button", { name: "Place order" }),
    );

    // A successful order invalidates the cart, so the last call is that GET.
    const orderCall = () => {
      const call = fetchMock.mock.calls.find(([input]) =>
        String(input).includes("/api/orders"),
      );
      if (!call) {
        throw new Error("no order request");
      }
      return { init: call[1], body: JSON.parse(String(call[1]?.body)) };
    };

    await waitFor(() =>
      expect(orderCall().body).toEqual({
        items: [{ product_id: 1, quantity: 2 }],
      }),
    );
    expect(
      (orderCall().init?.headers as Record<string, string>).Authorization,
    ).toBe("Bearer jwt-value");

    expect(await screen.findByText("Order placed")).toBeTruthy();
    expect(screen.getByText("ORD-7")).toBeTruthy();
    expect(screen.getByText("800 coins")).toBeTruthy();
    expect(
      screen.getByRole("link", { name: "View order" }).getAttribute("href"),
    ).toBe("/orders/7");
    expect(getSession()?.user.coins).toBe(200);
  });

  it("shows the server's reason when the order is refused", async () => {
    serve(cart([line(1, 2)]), () =>
      jsonResponse({ message: "insufficient coins" }, 402),
    );
    renderCart();

    const user = userEvent.setup();
    await user.click(
      await screen.findByRole("button", { name: "Place order" }),
    );

    expect(await screen.findByText("insufficient coins")).toBeTruthy();
    expect(screen.getByText("Estimated total")).toBeTruthy();
  });

  it("ends the session when the order comes back 401", async () => {
    serve(cart([line(1, 1)]), () =>
      jsonResponse({ message: "Invalid or expired token" }, 401),
    );
    renderCart();

    const user = userEvent.setup();
    await user.click(
      await screen.findByRole("button", { name: "Place order" }),
    );

    await waitFor(() => expect(getSession()).toBeNull());
  });
});
