import {
  ActionIcon,
  Alert,
  Anchor,
  Button,
  Flex,
  Group,
  Image,
  Paper,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router";

import { removeFromCart, setCartQuantity } from "../../api/cart";
import { ApiError } from "../../api/client";
import { createOrder } from "../../api/orders";
import {
  cartKeys,
  cartQuery,
  productsQuery,
  userKeys,
} from "../../api/queries";
import type {
  CartItem,
  CartResponse,
  Order,
  Product,
  User,
} from "../../api/types";
import { cartTotals } from "../../lib/cartTotals";
import { formatPrice, thumbnailUrl } from "../../lib/format";
import {
  clearSession,
  type Session,
  setSession,
  useSession,
} from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";
import { useUser } from "../../lib/useUser";

const PLACEHOLDER = `data:image/svg+xml,${encodeURIComponent(
  '<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200"><rect width="200" height="200" fill="#e9ecef"/></svg>',
)}`;

const SKELETONS = [0, 1, 2];

export function CartPage() {
  const session = useSession();

  return session ? <Cart session={session} /> : null;
}

function Cart({ session }: { session: Session }) {
  const queryClient = useQueryClient();
  const user = useUser();
  const cart = useQuery(cartQuery(session.token));
  const products = useQuery(productsQuery());
  const isSlow = useSlowHint(cart.isPending || products.isPending);

  const change = useMutation({
    mutationFn: ({
      productId,
      quantity,
    }: {
      productId: number;
      quantity: number;
    }) =>
      quantity > 0
        ? setCartQuantity(session.token, { productId, quantity })
        : removeFromCart(session.token, productId),

    onMutate: async ({ productId, quantity }) => {
      await queryClient.cancelQueries({ queryKey: cartKeys.all });
      const previous = queryClient.getQueryData<CartResponse>(cartKeys.all);
      queryClient.setQueryData<CartResponse>(cartKeys.all, (current) =>
        withQuantity(current, productId, quantity),
      );
      return { previous };
    },
    onError: (_error, _input, context) => {
      if (context?.previous) {
        queryClient.setQueryData(cartKeys.all, context.previous);
      }
    },
  });

  const items = cart.data?.items ?? [];
  const byId = new Map(
    (products.data ?? []).map((product) => [product.id, product]),
  );
  const totals = cartTotals(items, products.data);

  const checkout = useMutation({
    mutationFn: () =>
      createOrder(session.token, {
        items: items.map((item) => ({
          product_id: item.product_id,
          quantity: item.quantity,
        })),
      }),

    onSuccess: (order) => {
      const coins = session.user.coins - order.total_coins_used;

      setSession({ ...session, user: { ...session.user, coins } });
      queryClient.setQueryData<User>(
        userKeys.detail(session.user.id),
        (current) => (current ? { ...current, coins } : current),
      );
      queryClient.invalidateQueries({ queryKey: cartKeys.all });
    },

    onError: (error) => {
      if (error instanceof ApiError && error.status === 401) {
        clearSession();
      }
    },
  });

  const isPending = cart.isPending || products.isPending;
  const isError = cart.isError || products.isError;
  const error = cart.error ?? products.error;

  if (checkout.isSuccess) {
    return <OrderPlaced order={checkout.data} />;
  }

  if (isError) {
    return (
      <Stack gap="lg">
        <Title order={1} size="h2">
          Cart
        </Title>
        <Alert color="red" variant="light" title="Could not load your cart">
          <Text size="sm">{error?.message}</Text>
          <Button
            mt="sm"
            size="xs"
            variant="light"
            color="red"
            onClick={() => {
              cart.refetch();
              products.refetch();
            }}
            loading={cart.isFetching || products.isFetching}
          >
            Try again
          </Button>
        </Alert>
      </Stack>
    );
  }

  return (
    <Stack gap="lg">
      <Title order={1} size="h2">
        Cart
      </Title>

      {isSlow && (
        <Alert color="yellow" variant="light" title="Waking the API up">
          The free instance sleeps after fifteen minutes. This can take about
          twenty seconds.
        </Alert>
      )}

      {isPending && (
        <Stack gap="sm">
          {SKELETONS.map((key) => (
            <Skeleton key={key} height={80} radius="sm" />
          ))}
        </Stack>
      )}

      {!isPending && items.length === 0 && (
        <Stack gap="xs" align="flex-start">
          <Text c="dimmed">Your cart is empty.</Text>
          <Anchor component={Link} to="/products">
            Browse products
          </Anchor>
        </Stack>
      )}

      {items.length > 0 && (
        <>
          <Stack gap="sm">
            {items.map((item) => {
              const product = byId.get(item.product_id);

              return product ? (
                <CartLine
                  key={item.id}
                  item={item}
                  product={product}
                  onChange={(quantity) =>
                    change.mutate({ productId: item.product_id, quantity })
                  }
                />
              ) : (
                <MissingLine
                  key={item.id}
                  onRemove={() =>
                    change.mutate({ productId: item.product_id, quantity: 0 })
                  }
                />
              );
            })}
          </Stack>

          <Paper withBorder p="md" radius="sm">
            <Stack gap="xs">
              <Group justify="space-between">
                <Text size="sm">Estimated total</Text>
                <Text ff="monospace" fw={600}>
                  {formatPrice(totals.totalCents / 100)}
                </Text>
              </Group>
              <Group justify="space-between">
                <Text size="sm">Coins</Text>
                <Text ff="monospace">{totals.coins.toLocaleString()}</Text>
              </Group>
              <Group justify="space-between">
                <Text size="sm" c="dimmed">
                  Your balance
                </Text>
                <Text size="sm" c="dimmed">
                  {user?.coins.toLocaleString()}
                </Text>
              </Group>
              <Text size="xs" c="dimmed">
                Coins are whole numbers, rounded up. The server sends the exact
                amount when you check out.
              </Text>
              {totals.missing > 0 && (
                <Text size="xs" c="red.8">
                  {totals.missing === 1
                    ? "One item is no longer for sale and is not counted."
                    : `${totals.missing} items are no longer for sale and are not counted.`}
                </Text>
              )}

              {checkout.isError && (
                <Alert
                  color="red"
                  variant="light"
                  title="Could not place the order"
                >
                  <Text size="sm">{checkout.error.message}</Text>
                </Alert>
              )}

              <Group justify="flex-end" mt="xs">
                <Button
                  onClick={() => checkout.mutate()}
                  loading={checkout.isPending}
                  disabled={totals.missing > 0}
                >
                  Place order
                </Button>
              </Group>
            </Stack>
          </Paper>
        </>
      )}
    </Stack>
  );
}

interface CartLineProps {
  item: CartItem;
  product: Product;
  onChange: (quantity: number) => void;
}

function CartLine({ item, product, onChange }: CartLineProps) {
  const removing = item.quantity === 1;

  return (
    <Paper withBorder p="sm" radius="sm">
      <Flex
        direction={{ base: "column", sm: "row" }}
        gap="sm"
        align={{ base: "stretch", sm: "center" }}
      >
        <Group gap="sm" wrap="nowrap" style={{ flex: 1, minWidth: 0 }}>
          <Image
            src={thumbnailUrl(product.image_url, 200)}
            fallbackSrc={PLACEHOLDER}
            alt=""
            w={56}
            h={56}
            radius="sm"
            fit="cover"
          />

          <Stack gap={2} style={{ minWidth: 0 }}>
            <Anchor
              component={Link}
              to={`/products/${product.id}`}
              size="sm"
              fw={600}
              lineClamp={1}
            >
              {product.name}
            </Anchor>
            <Text size="xs" c="dimmed">
              {formatPrice(product.price)} each
            </Text>
          </Stack>
        </Group>

        <Group gap="sm" wrap="nowrap" justify="space-between">
          <Group gap={4} wrap="nowrap">
            <ActionIcon
              variant="default"
              size="sm"
              aria-label={
                removing
                  ? `Remove ${product.name}`
                  : `Remove one ${product.name}`
              }
              onClick={() => onChange(item.quantity - 1)}
            >
              −
            </ActionIcon>
            <Text size="sm" ff="monospace" w={24} ta="center">
              {item.quantity}
            </Text>
            <ActionIcon
              variant="default"
              size="sm"
              aria-label={`Add one ${product.name}`}
              onClick={() => onChange(item.quantity + 1)}
            >
              +
            </ActionIcon>
          </Group>

          <Text size="sm" ff="monospace" fw={600} w={72} ta="right">
            {formatPrice(product.price * item.quantity)}
          </Text>
        </Group>
      </Flex>
    </Paper>
  );
}

function MissingLine({ onRemove }: { onRemove: () => void }) {
  return (
    <Paper withBorder p="sm" radius="sm">
      <Group justify="space-between" wrap="nowrap">
        <Text size="sm" c="dimmed">
          This item is no longer for sale.
        </Text>
        <Button
          variant="subtle"
          color="gray"
          size="compact-sm"
          onClick={onRemove}
        >
          Remove
        </Button>
      </Group>
    </Paper>
  );
}

function withQuantity(
  cart: CartResponse | undefined,
  productId: number,
  quantity: number,
): CartResponse | undefined {
  if (!cart) {
    return cart;
  }

  const items =
    quantity > 0
      ? cart.items.map((item) =>
          item.product_id === productId ? { ...item, quantity } : item,
        )
      : cart.items.filter((item) => item.product_id !== productId);

  return { ...cart, items, total_items: items.length };
}

function OrderPlaced({ order }: { order: Order }) {
  const user = useUser();

  return (
    <Stack gap="lg">
      <Title order={1} size="h2">
        Order placed
      </Title>

      <Paper withBorder p="md" radius="sm" maw={420}>
        <Stack gap="xs">
          <Group justify="space-between">
            <Text size="sm">Order</Text>
            <Text size="sm" ff="monospace">
              {order.order_number}
            </Text>
          </Group>
          <Group justify="space-between">
            <Text size="sm">Paid</Text>
            <Text size="sm" ff="monospace">
              {order.total_coins_used.toLocaleString()} coins
            </Text>
          </Group>
          <Group justify="space-between">
            <Text size="sm" c="dimmed">
              Balance
            </Text>
            <Text size="sm" c="dimmed" ff="monospace">
              {user?.coins.toLocaleString()} coins
            </Text>
          </Group>
        </Stack>
      </Paper>

      <Group>
        <Button component={Link} to="/products">
          Keep shopping
        </Button>
        <Button component={Link} to={`/orders/${order.id}`} variant="default">
          View order
        </Button>
      </Group>
    </Stack>
  );
}
