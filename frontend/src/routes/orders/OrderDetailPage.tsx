import {
  Alert,
  Anchor,
  Badge,
  Button,
  Divider,
  Flex,
  Group,
  Paper,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router";

import { ApiError } from "../../api/client";
import { orderQuery } from "../../api/queries";
import type { OrderItem } from "../../api/types";
import { formatDate, formatPrice } from "../../lib/format";
import { type Session, useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";

const STATUS_COLOR = {
  completed: "teal",
  pending: "yellow",
  cancelled: "gray",
} as const;

export function OrderDetailPage() {
  const { id } = useParams();
  const orderId = Number(id);
  const session = useSession();

  if (!Number.isInteger(orderId) || orderId <= 0) {
    return <OrderNotFound />;
  }

  return session ? <OrderDetail orderId={orderId} session={session} /> : null;
}

function OrderDetail({
  orderId,
  session,
}: {
  orderId: number;
  session: Session;
}) {
  const { data, isPending, isError, error, refetch, isFetching } = useQuery(
    orderQuery(session.token, orderId),
  );
  const isSlow = useSlowHint(isPending);

  // A 404 means the order is not this shopper's, or it never existed. Both
  // read the same, and neither is worth a retry button.
  if (isError && error instanceof ApiError && error.status === 404) {
    return <OrderNotFound />;
  }

  return (
    <Stack gap="lg">
      <Anchor component={Link} to="/orders" size="sm" c="dimmed">
        All orders
      </Anchor>

      {isSlow && (
        <Alert color="yellow" variant="light" title="Waking the API up">
          The free instance sleeps after fifteen minutes. This can take about
          twenty seconds.
        </Alert>
      )}

      {isError && (
        <Alert color="red" variant="light" title="Could not load this order">
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
        <>
          <Skeleton height={64} radius="sm" />
          <Skeleton height={240} radius="sm" />
        </>
      )}

      {data && (
        <>
          <div>
            <Group gap="sm">
              <Title order={1} size="h2" ff="monospace">
                {data.order_number}
              </Title>
              <Badge variant="light" color={STATUS_COLOR[data.status]}>
                {data.status}
              </Badge>
            </Group>
            <Text c="dimmed" mt={4}>
              Placed {formatDate(data.created_at)}
            </Text>
          </div>

          <Paper withBorder p="md" radius="sm">
            <Stack gap="sm">
              {(data.items ?? []).length === 0 ? (
                <Text size="sm" c="dimmed">
                  This order has no lines.
                </Text>
              ) : (
                (data.items ?? []).map((item) => (
                  <ItemRow key={item.id} item={item} />
                ))
              )}
            </Stack>

            <Divider my="md" />

            <Stack gap="xs">
              <Group justify="space-between">
                <Text size="sm">Total</Text>
                <Text size="sm" ff="monospace">
                  {formatPrice(data.total_amount)}
                </Text>
              </Group>
              <Group justify="space-between">
                <Text size="sm">Paid</Text>
                <Text size="sm" ff="monospace">
                  {data.total_coins_used.toLocaleString()} coins
                </Text>
              </Group>
            </Stack>
          </Paper>
        </>
      )}
    </Stack>
  );
}

function ItemRow({ item }: { item: OrderItem }) {
  return (
    <Flex gap="sm" justify="space-between" wrap="wrap">
      <Stack gap={2} style={{ minWidth: 0 }}>
        <Anchor
          component={Link}
          to={`/products/${item.product_id}`}
          size="sm"
          lineClamp={1}
        >
          {item.product_name}
        </Anchor>
        <Text size="xs" c="dimmed">
          {item.quantity} × {formatPrice(item.product_price)}
        </Text>
      </Stack>
      <Text size="sm" ff="monospace" ml="auto">
        {formatPrice(item.subtotal)}
      </Text>
    </Flex>
  );
}

function OrderNotFound() {
  return (
    <Stack gap="md">
      <Title order={1} size="h2">
        Order not found
      </Title>
      <Text c="dimmed">
        This order does not exist, or it belongs to another account.
      </Text>
      <Group>
        <Button component={Link} to="/orders" variant="light">
          Back to your orders
        </Button>
      </Group>
    </Stack>
  );
}
