import { describe, expect, it } from "vitest";

import type { CartItem, Product } from "../api/types";
import { cartTotals } from "./cartTotals";

function product(id: number, price: number): Product {
  return {
    id,
    category_id: 1,
    name: `Product ${id}`,
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

function item(productId: number, quantity: number): CartItem {
  return {
    id: productId,
    user_id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
    product_id: productId,
    quantity,
    created_at: "2026-11-01T00:00:00Z",
    updated_at: "2026-11-01T00:00:00Z",
  };
}

describe("cartTotals", () => {
  it("returns zero for an empty cart", () => {
    expect(cartTotals([], [product(1, 10)])).toEqual({
      totalCents: 0,
      coins: 0,
      missing: 0,
    });
  });

  it("multiplies each line by its quantity", () => {
    const totals = cartTotals(
      [item(1, 2), item(2, 1)],
      [product(1, 19.99), product(2, 5)],
    );

    expect(totals.totalCents).toBe(4498);
    expect(totals.coins).toBe(45);
  });

  it("rounds coins up, the way the server charges", () => {
    expect(cartTotals([item(1, 1)], [product(1, 49.99)]).coins).toBe(50);
    expect(cartTotals([item(1, 1)], [product(1, 0.01)]).coins).toBe(1);
  });

  it("skips a product that is no longer in the catalogue", () => {
    const totals = cartTotals([item(1, 1), item(99, 3)], [product(1, 10)]);

    expect(totals.totalCents).toBe(1000);
    expect(totals.missing).toBe(1);
  });

  it("counts nothing while the catalogue is still loading", () => {
    const totals = cartTotals([item(1, 2)], undefined);

    expect(totals.totalCents).toBe(0);
    expect(totals.missing).toBe(1);
  });
});
