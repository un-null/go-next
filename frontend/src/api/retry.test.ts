import { describe, expect, it } from "vitest";

import { ApiError } from "./client";
import { retryUnlessClientError } from "./retry";

describe("retryUnlessClientError", () => {
  it("does not retry a 402, the server has already answered", () => {
    expect(
      retryUnlessClientError(0, new ApiError(402, "insufficient coins", null)),
    ).toBe(false);
  });

  it("does not retry a 409 or a 404", () => {
    expect(
      retryUnlessClientError(0, new ApiError(409, "insufficient stock", null)),
    ).toBe(false);
    expect(
      retryUnlessClientError(0, new ApiError(404, "not found", null)),
    ).toBe(false);
  });

  it("retries a server error once", () => {
    expect(retryUnlessClientError(0, new ApiError(500, "internal", null))).toBe(
      true,
    );
    expect(retryUnlessClientError(1, new ApiError(500, "internal", null))).toBe(
      false,
    );
  });

  it("retries a request that never reached the server", () => {
    expect(
      retryUnlessClientError(0, new ApiError(0, "network failed", null)),
    ).toBe(true);
  });

  it("retries an error that is not an ApiError", () => {
    expect(retryUnlessClientError(0, new Error("boom"))).toBe(true);
  });
});
