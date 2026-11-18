import {
  Alert,
  Anchor,
  AspectRatio,
  Badge,
  Button,
  Divider,
  Group,
  Image,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link, useParams } from "react-router";

import { setCartQuantity } from "../../api/cart";
import { ApiError } from "../../api/client";
import { cartKeys, cartQuery, productQuery } from "../../api/queries";
import type { Product } from "../../api/types";
import { formatPrice, thumbnailUrl } from "../../lib/format";
import { placeholder } from "../../lib/placeholder";
import { withNext } from "../../lib/redirect";
import { clearSession, type Session, useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";
import classes from "./ProductDetailPage.module.css";

const PLACEHOLDER = placeholder(900, 600);

export function ProductDetailPage() {
  const { id } = useParams();
  const productId = Number(id);

  if (!Number.isInteger(productId) || productId <= 0) {
    return <NotFound />;
  }

  return <ProductDetail productId={productId} />;
}

function ProductDetail({ productId }: { productId: number }) {
  const { data, isPending, isError, error, refetch, isFetching } = useQuery(
    productQuery(productId),
  );
  const isSlow = useSlowHint(isPending);

  if (isError && error instanceof ApiError && error.status === 404) {
    return <NotFound />;
  }

  return (
    <Stack gap="lg">
      <Anchor component={Link} to="/products" size="sm" c="dimmed">
        All products
      </Anchor>

      {isSlow && (
        <Alert color="yellow" variant="light" title="Waking the API up">
          The free instance sleeps after fifteen minutes. This can take about
          twenty seconds.
        </Alert>
      )}

      {isError && (
        <Alert color="red" variant="light" title="Could not load this product">
          <Text size="sm">{error.message}</Text>
          <Button
            mt="sm"
            size="xs"
            variant="light"
            color="red"
            onClick={() => refetch()}
            loading={isFetching}
          >
            Try again
          </Button>
        </Alert>
      )}

      {isPending && (
        <SimpleGrid cols={{ base: 1, md: 2 }} spacing="xl">
          <Skeleton height={360} radius="md" />
          <Skeleton height={280} radius="md" />
        </SimpleGrid>
      )}

      {data && (
        <SimpleGrid
          cols={{ base: 1, md: 2 }}
          spacing="xl"
          className={classes.layout}
        >
          <AspectRatio ratio={3 / 2} bg="gray.1" className={classes.media}>
            <Image
              src={thumbnailUrl(data.image_url, 900)}
              fallbackSrc={PLACEHOLDER}
              alt=""
              w="100%"
              h="100%"
              fit="cover"
            />
          </AspectRatio>

          <Stack gap="md">
            <div>
              <Title order={1} size="h2">
                {data.name}
              </Title>
              <Text c="dimmed" mt={4}>
                {data.description}
              </Text>
            </div>

            <Group gap="xs">
              <Text fw={700} fz="h3" ff="monospace">
                {formatPrice(data.price)}
              </Text>
              <Text size="sm" c="dimmed">
                from {Math.ceil(data.price)} coins
              </Text>
            </Group>

            <Group gap="xs">
              <Badge
                variant="light"
                color={data.stock_quantity > 0 ? "teal" : "gray"}
              >
                {data.stock_quantity > 0
                  ? `${data.stock_quantity} in stock`
                  : "Out of stock"}
              </Badge>
              <Text size="sm" c="dimmed">
                {data.average_rating.toFixed(1)} stars
              </Text>
            </Group>

            <Divider />

            <AddToCart product={data} />
          </Stack>
        </SimpleGrid>
      )}
    </Stack>
  );
}

function AddToCart({ product }: { product: Product }) {
  const session = useSession();

  if (!session) {
    return (
      <Group>
        <Button
          component={Link}
          to={withNext("/login", `/products/${product.id}`)}
        >
          Sign in to add to cart
        </Button>
      </Group>
    );
  }

  return <AddToCartButton product={product} session={session} />;
}

function AddToCartButton({
  product,
  session,
}: {
  product: Product;
  session: Session;
}) {
  const queryClient = useQueryClient();
  const cart = useQuery(cartQuery(session.token));

  const line = cart.data?.items.find((item) => item.product_id === product.id);
  const quantity = line?.quantity ?? 0;

  const add = useMutation({
    mutationFn: () =>
      setCartQuantity(session.token, {
        productId: product.id,
        quantity: quantity + 1,
      }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: cartKeys.all }),
    onError: (error) => {
      if (error instanceof ApiError && error.status === 401) {
        clearSession();
      }
    },
  });

  return (
    <Stack gap="xs">
      <Group gap="md">
        <Button
          onClick={() => add.mutate()}
          loading={add.isPending}
          disabled={cart.isPending || cart.isError}
        >
          Add to cart
        </Button>
        {quantity > 0 && (
          <Anchor component={Link} to="/cart" size="sm">
            View cart
          </Anchor>
        )}
      </Group>

      {quantity > 0 && (
        <Text size="sm" c="dimmed">
          {quantity} in your cart.
        </Text>
      )}

      {add.isError && (
        <Text size="sm" c="red.8">
          {add.error.message}
        </Text>
      )}

      {cart.isError && (
        <Text size="sm" c="dimmed">
          Could not read your cart, so the quantity is unknown. Reload the page
          to try again.
        </Text>
      )}

      <Text size="xs" c="dimmed">
        Coins are demo currency, with no monetary value.
      </Text>
    </Stack>
  );
}

function NotFound() {
  return (
    <Stack gap="md">
      <Title order={1} size="h2">
        Product not found
      </Title>
      <Text c="dimmed">This item does not exist, or it was removed.</Text>
      <Group>
        <Button component={Link} to="/products" variant="light">
          Back to the catalogue
        </Button>
      </Group>
    </Stack>
  );
}
