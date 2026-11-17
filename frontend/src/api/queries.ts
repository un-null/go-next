import { queryOptions } from "@tanstack/react-query";

import { getCart } from "./cart";
import { getProduct, listProducts } from "./catalog";
import { getOrder, listOrders } from "./orders";

export const productKeys = {
  all: ["products"] as const,
  list: () => [...productKeys.all, "list"] as const,
  detail: (id: number) => [...productKeys.all, "detail", id] as const,
};

export const productsQuery = () =>
  queryOptions({
    queryKey: productKeys.list(),
    queryFn: () => listProducts({ limit: 100 }),
  });

export const productQuery = (id: number) =>
  queryOptions({
    queryKey: productKeys.detail(id),
    queryFn: () => getProduct(id),
  });

export const cartKeys = {
  all: ["cart"] as const,
};

export const cartQuery = (token: string) =>
  queryOptions({
    queryKey: cartKeys.all,
    queryFn: () => getCart(token),
  });

export const orderKeys = {
  all: ["orders"] as const,
  list: (userId: string) => [...orderKeys.all, "list", userId] as const,
  detail: (id: number) => [...orderKeys.all, "detail", id] as const,
};

export const ordersQuery = (token: string, userId: string) =>
  queryOptions({
    queryKey: orderKeys.list(userId),
    queryFn: () => listOrders(token, userId),
  });

// The list response has no items; only this one does. Same key shape as the
// products, so a detail page can be opened by id without knowing the user.
export const orderQuery = (token: string, id: number) =>
  queryOptions({
    queryKey: orderKeys.detail(id),
    queryFn: () => getOrder(token, id),
  });
