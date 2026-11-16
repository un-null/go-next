import { describe, expect, it } from "vitest";

import { safeNext, withNext } from "./redirect";

describe("safeNext", () => {
  it("keeps the path the guard produced", () => {
    expect(safeNext("/cart")).toBe("/cart");
    expect(safeNext("/cart?coupon=x")).toBe("/cart?coupon=x");
  });

  it("falls back when the URL carries no destination", () => {
    expect(safeNext(null)).toBe("/products");
    expect(safeNext("")).toBe("/products");
  });

  it("refuses to leave the site", () => {
    expect(safeNext("//evil.example/cart")).toBe("/products");
    expect(safeNext("https://evil.example/cart")).toBe("/products");
  });
});

describe("withNext", () => {
  it("leaves an ordinary link alone", () => {
    expect(withNext("/signup", "/products")).toBe("/signup");
  });

  it("carries the destination through the switch", () => {
    expect(withNext("/signup", "/cart?coupon=x")).toBe(
      "/signup?next=%2Fcart%3Fcoupon%3Dx",
    );
  });
});
