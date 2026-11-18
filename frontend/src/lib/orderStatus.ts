import type { OrderStatus } from "../api/types";

export const STATUS_COLOR: Record<OrderStatus, string> = {
  completed: "teal",
  pending: "yellow",
  cancelled: "gray",
};
