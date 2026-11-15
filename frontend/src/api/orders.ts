import { request } from "./client";
import type { CreateOrderInput, Order } from "./types";

// The server decides the user, order number, amount, coins and status; the
// client only sends what it wants to buy.
export function createOrder(
  token: string,
  input: CreateOrderInput,
): Promise<Order> {
  return request<Order>("/api/orders", { method: "POST", token, body: input });
}

export function getOrder(token: string, id: number): Promise<Order> {
  return request<Order>(`/api/orders/${id}`, { token });
}

export function listOrders(token: string, userId: string): Promise<Order[]> {
  return request<Order[]>(`/api/users/${userId}/orders`, { token });
}
