// @vitest-environment jsdom

import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router";
import { describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, mockApi } from "../../api/testing";
import type { Category, Product } from "../../api/types";
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

const BOOK: Product = {
  id: 3,
  category_id: 2,
  name: "Go Programming Essentials",
  description: "A book.",
  price: 39.99,
  stock_quantity: 12,
  image_url: "",
  average_rating: 4.8,
  total_comments: 6,
  created_at: "2026-11-01T00:00:00Z",
  updated_at: "2026-11-01T00:00:00Z",
};

const CATEGORIES: Category[] = [
  { id: 1, name: "Electronics" },
  { id: 2, name: "Books" },
];

mockApi();

function serve(catalog: Product[] = PRODUCTS) {
  fetchMock.mockImplementation((input) => {
    const url = new URL(String(input));

    if (url.pathname === "/api/categories") {
      return Promise.resolve(jsonResponse(CATEGORIES));
    }

    const q = (url.searchParams.get("q") ?? "").toLowerCase();
    const category = Number(url.searchParams.get("category") ?? 0);
    const found = catalog.filter(
      (product) =>
        product.name.toLowerCase().includes(q) &&
        (category === 0 || product.category_id === category),
    );
    return Promise.resolve(jsonResponse(found));
  });
}

function renderPage(path = "/products") {
  render(
    <MantineProvider>
      <QueryClientProvider client={new QueryClient()}>
        <MemoryRouter initialEntries={[path]}>
          <Routes>
            <Route path="/products" element={<ProductsPage />} />
            <Route path="/products/:id" element={<p>the product</p>} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

function lastProductUrl(): string {
  const call = fetchMock.mock.calls
    .filter(([input]) => String(input).includes("/api/products?"))
    .at(-1);
  return call ? String(call[0]) : "";
}

describe("ProductsPage", () => {
  it("lists the catalogue without sending a search", async () => {
    serve();
    renderPage();

    expect(await screen.findByText("Gaming Laptop G15")).toBeTruthy();
    expect(screen.getByText("Wireless Earbuds Pro")).toBeTruthy();
    expect(screen.getByText("2 products")).toBeTruthy();
    expect(lastProductUrl()).toBe(
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
        expect(lastProductUrl()).toContain("q=laptop");
      },
      { timeout: 3000 },
    );

    expect(await screen.findByText("Gaming Laptop G15")).toBeTruthy();
    expect(screen.queryByText("Wireless Earbuds Pro")).toBeNull();
    expect(screen.getByText("1 product")).toBeTruthy();
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
    expect(
      screen.getAllByRole("button", { name: "Clear search" }),
    ).toHaveLength(1);
    expect(screen.getByRole("button", { name: "Clear filters" })).toBeTruthy();
  });

  it("narrows the catalogue by category", async () => {
    serve([...PRODUCTS, BOOK]);
    renderPage();
    const user = userEvent.setup();

    expect(await screen.findByText("Go Programming Essentials")).toBeTruthy();

    await user.selectOptions(screen.getByLabelText("Category"), "2");

    await waitFor(() => {
      expect(lastProductUrl()).toContain("category=2");
    });

    expect(await screen.findByText("Go Programming Essentials")).toBeTruthy();
    expect(screen.queryByText("Gaming Laptop G15")).toBeNull();
    expect(screen.getByText("1 product")).toBeTruthy();
  });

  it("ignores a category that does not exist", async () => {
    serve([...PRODUCTS, BOOK]);
    renderPage("/products?category=99");

    expect(await screen.findByText("Gaming Laptop G15")).toBeTruthy();
    expect(screen.queryByText("No products in this category.")).toBeNull();
    expect(lastProductUrl()).not.toContain("category=");
  });
});
