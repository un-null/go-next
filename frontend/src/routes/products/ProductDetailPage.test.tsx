// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, mockApi } from "../../api/testing";
import type { CartItem, CartResponse, Product, User } from "../../api/types";
import { clearSession, setSession } from "../../lib/session";
import { ProductDetailPage } from "./ProductDetailPage";

mockApi();

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-01T00:00:00Z",
  updated_at: "2026-11-01T00:00:00Z",
};

const PRODUCT: Product = {
  id: 1,
  category_id: 1,
  name: "Monitor",
  description: "27-inch 4K",
  price: 399.99,
  stock_quantity: 10,
  image_url: "",
  average_rating: 4.5,
  total_comments: 2,
  created_at: "2026-11-01T00:00:00Z",
  updated_at: "2026-11-01T00:00:00Z",
};

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

function serve(cartResponse: CartResponse, post: () => Response) {
  fetchMock.mockImplementation(async (input, init) => {
    const url = String(input);

    if (url.includes("/api/products/")) {
      return jsonResponse(PRODUCT);
    }
    if (url.includes("/api/cart")) {
      return (init?.method ?? "GET") === "POST"
        ? post()
        : jsonResponse(cartResponse);
    }

    throw new Error(`unexpected request: ${url}`);
  });
}

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <MantineProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={["/products/1"]}>
          <Routes>
            <Route path="/products/:id" element={<ProductDetailPage />} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

afterEach(() => {
  act(() => clearSession());
});

describe("ProductDetailPage", () => {
  it("sends a guest to sign in and back", async () => {
    serve(cart([]), () => jsonResponse({ success: true }, 201));
    renderPage();

    const link = await screen.findByRole("link", {
      name: "Sign in to add to cart",
    });

    expect(link.getAttribute("href")).toBe("/login?next=%2Fproducts%2F1");
  });

  it("adds one to the quantity the cart already holds", async () => {
    setSession({ token: "jwt-value", user: USER });
    serve(cart([line(1, 2)]), () =>
      jsonResponse({ success: true, message: "added" }, 201),
    );

    renderPage();

    expect(await screen.findByText("2 in your cart.")).toBeTruthy();
    await userEvent.click(screen.getByRole("button", { name: "Add to cart" }));

    await waitFor(() =>
      expect(
        fetchMock.mock.calls.some(([, init]) => init?.method === "POST"),
      ).toBe(true),
    );

    // The invalidation refetches the cart, so the last call is a GET.
    const [, post] =
      fetchMock.mock.calls.find(([, init]) => init?.method === "POST") ?? [];

    expect(JSON.parse(String(post?.body))).toEqual({
      product_id: 1,
      quantity: 3,
    });
    expect((post?.headers as Record<string, string>).Authorization).toBe(
      "Bearer jwt-value",
    );
  });

  it("shows the server message when the cart refuses", async () => {
    setSession({ token: "jwt-value", user: USER });
    serve(cart([line(1, 1)]), () =>
      jsonResponse({ message: "insufficient stock" }, 409),
    );

    renderPage();

    await screen.findByText("1 in your cart.");
    await userEvent.click(screen.getByRole("button", { name: "Add to cart" }));

    expect(await screen.findByText("insufficient stock")).toBeTruthy();
  });
});
