// @vitest-environment jsdom

import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, lastCall, mockApi } from "../../api/testing";
import type { Product } from "../../api/types";
import { ProductsPage } from "./ProductsPage";

const PRODUCTS: Product[] = [
  {
    id: 1,
    category_id: 1,
    name: "Gaming Laptop G15",
    description: "A laptop for games.",
    price: 1299.99,
    stock_quantity: 4,
    image_url: "",
    average_rating: 4.5,
    total_comments: 10,
    created_at: "2026-11-01T00:00:00Z",
    updated_at: "2026-11-01T00:00:00Z",
  },
  {
    id: 2,
    category_id: 1,
    name: "Wireless Earbuds Pro",
    description: "Earbuds.",
    price: 149.99,
    stock_quantity: 20,
    image_url: "",
    average_rating: 4.2,
    total_comments: 3,
    created_at: "2026-11-01T00:00:00Z",
    updated_at: "2026-11-01T00:00:00Z",
  },
];

mockApi();

// The server does the filtering; the stub only has to answer a request that
// looks like the one the page sends.
function serve() {
  fetchMock.mockImplementation((input) => {
    const url = new URL(String(input));
    const q = (url.searchParams.get("q") ?? "").toLowerCase();
    const found = PRODUCTS.filter((product) =>
      product.name.toLowerCase().includes(q),
    );
    return Promise.resolve(jsonResponse(found));
  });
}

function renderPage() {
  render(
    <MantineProvider>
      <QueryClientProvider client={new QueryClient()}>
        <MemoryRouter initialEntries={["/products"]}>
          <Routes>
            <Route path="/products" element={<ProductsPage />} />
            <Route path="/products/:id" element={<p>the product</p>} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

describe("ProductsPage", () => {
  it("lists the catalogue without sending a search", async () => {
    serve();
    renderPage();

    expect(await screen.findByText("Gaming Laptop G15")).toBeTruthy();
    expect(screen.getByText("Wireless Earbuds Pro")).toBeTruthy();
    expect(screen.getByText("2 products")).toBeTruthy();
    expect(lastCall().url).toBe(
      "https://api.test/api/products?page=1&limit=100",
    );
  });

  it("sends the query once typing pauses", async () => {
    serve();
    renderPage();
    const user = userEvent.setup();

    await user.type(
      screen.getByRole("textbox", { name: "Search products" }),
      "laptop",
    );

    await waitFor(
      () => {
        expect(lastCall().url).toContain("q=laptop");
      },
      { timeout: 3000 },
    );

    expect(await screen.findByText("Gaming Laptop G15")).toBeTruthy();
    expect(screen.queryByText("Wireless Earbuds Pro")).toBeNull();
    expect(screen.getByText("1 products")).toBeTruthy();
  });

  it("says so when nothing matches", async () => {
    serve();
    renderPage();
    const user = userEvent.setup();

    await user.type(
      screen.getByRole("textbox", { name: "Search products" }),
      "zzz",
    );

    expect(await screen.findByText('No products match "zzz".')).toBeTruthy();
    // The field's own clear button and the one in the empty state share a
    // label, because they do the same thing.
    expect(
      screen.getAllByRole("button", { name: "Clear search" }),
    ).toHaveLength(2);
  });
});
