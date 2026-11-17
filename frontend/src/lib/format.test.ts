import { describe, expect, it } from "vitest";

import { formatDate, formatPrice, thumbnailUrl } from "./format";

describe("formatPrice", () => {
  it("shows cents and whole dollars the same way", () => {
    expect(formatPrice(699.99)).toBe("$699.99");
    expect(formatPrice(149.99)).toBe("$149.99");
    expect(formatPrice(50)).toBe("$50.00");
  });
});

describe("formatDate", () => {
  it("shows a readable date", () => {
    expect(formatDate("2026-11-17T12:00:00Z")).toBe("Nov 17, 2026");
  });

  it("returns the input when it cannot be parsed", () => {
    expect(formatDate("not a date")).toBe("not a date");
  });
});

describe("thumbnailUrl", () => {
  it("asks Unsplash for a resized copy", () => {
    const url = thumbnailUrl(
      "https://images.unsplash.com/photo-1510557880182-3d4d3cba35a5",
      600,
    );

    expect(url).toContain("w=600");
    expect(url).toContain("auto=format");
    expect(url).toContain("q=70");
  });

  it("keeps an existing query string", () => {
    const url = thumbnailUrl("https://example.com/a.jpg?crop=faces");

    expect(url).toContain("crop=faces");
    expect(url).toContain("w=600");
  });

  it("returns the input when it is not an absolute URL", () => {
    expect(thumbnailUrl("not a url")).toBe("not a url");
  });

  it("returns an empty string when there is no image", () => {
    expect(thumbnailUrl("")).toBe("");
  });
});
