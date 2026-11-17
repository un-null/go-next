import { queryOptions } from "@tanstack/react-query";

import { getCart } from "./cart";
import { getProduct, listProducts } from "./catalog";

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
