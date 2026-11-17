import type { CartItem, Product } from "../api/types";

export interface CartTotals {
  totalCents: number;
  coins: number;
  missing: number;
}

export function cartTotals(
  items: CartItem[],
  products: Product[] | undefined,
): CartTotals {
  const prices = new Map(
    products?.map((product) => [product.id, product.price]),
  );

  let totalCents = 0;
  let missing = 0;

  for (const item of items) {
    const price = prices.get(item.product_id);
    if (price === undefined) {
      missing += 1;
      continue;
    }

    totalCents += Math.round(price * 100) * item.quantity;
  }

  return { totalCents, coins: Math.ceil(totalCents / 100), missing };
}
