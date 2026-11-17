import { request } from "./client";
import type { CartResponse } from "./types";

export interface SetCartQuantityInput {
  productId: number;
  quantity: number;
}

export function getCart(token: string): Promise<CartResponse> {
  return request<CartResponse>("/api/cart", { token });
}

export function setCartQuantity(
  token: string,
  input: SetCartQuantityInput,
): Promise<void> {
  return request<void>("/api/cart", {
    method: "POST",
    token,
    body: { product_id: input.productId, quantity: input.quantity },
  });
}

export function removeFromCart(
  token: string,
  productId: number,
): Promise<void> {
  const query = new URLSearchParams({ product_id: String(productId) });
  return request<void>(`/api/cart?${query}`, { method: "DELETE", token });
}

export function clearCart(token: string): Promise<void> {
  return request<void>("/api/cart/clear", { method: "DELETE", token });
}
