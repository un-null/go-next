import { QueryClient } from "@tanstack/react-query";
import { describe, expect, it } from "vitest";

import { productKeys, productQuery, productsQuery } from "./queries";
import { fetchMock, jsonResponse, lastCall, mockApi } from "./testing";

mockApi();

describe("productKeys", () => {
  it("keeps every product query under one prefix", () => {
    expect(productKeys.all).toEqual(["products"]);
    expect(productKeys.list()).toEqual(["products", "list"]);
    expect(productKeys.detail(3)).toEqual(["products", "detail", 3]);
  });

  it("lets one invalidation cover the list and every detail", () => {
    for (const key of [productKeys.list(), productKeys.detail(7)]) {
      expect(key.slice(0, productKeys.all.length)).toEqual([
        ...productKeys.all,
      ]);
    }
  });
});

describe("product queries", () => {
  it("fetches the whole catalogue in one request", async () => {
    fetchMock.mockResolvedValue(jsonResponse([]));

    const products = await new QueryClient().fetchQuery(productsQuery());

    expect(products).toEqual([]);
    expect(lastCall().url).toBe(
      "https://api.test/api/products?page=1&limit=100",
    );
  });

  it("fetches a single product by id", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ id: 3, name: "Gaming Laptop G15" }),
    );

    await new QueryClient().fetchQuery(productQuery(3));

    expect(lastCall().url).toBe("https://api.test/api/products/3");
  });
});
