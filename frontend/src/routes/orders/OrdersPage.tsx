import {
  Alert,
  Anchor,
  Badge,
  Button,
  Flex,
  Group,
  Paper,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router";

import { ordersQuery } from "../../api/queries";
import type { Order, OrderStatus } from "../../api/types";
import { formatDate, formatPrice } from "../../lib/format";
import { type Session, useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";

const SKELETONS = [0, 1, 2];

const STATUS_COLOR: Record<OrderStatus, string> = {
  completed: "teal",
  pending: "yellow",
  cancelled: "gray",
};

export function OrdersPage() {
  const session = useSession();
  return session ? <Orders session={session} /> : null;
}

function Orders({ session }: { session: Session }) {
  const orders = useQuery(ordersQuery(session.token, session.user.id));
  const isSlow = useSlowHint(orders.isPending);

  return (
    <Stack gap="lg">
      <div>
        <Title order={1} size="h2">
          Orders
        </Title>
      </div>

      {isSlow && (
        <Alert color="yellow" variant="light" title="Waking the API up">
          The free instance (Render) sleeps after fifteen minutes. The first
          request can take about 20 sec. It is working, not broken.
        </Alert>
      )}

      {orders.isError && (
        <Alert color="red" variant="light" title="Could not load your orders">
          <Text size="sm">{orders.error.message}</Text>
          <Button
            mt="sm"
            size="xs"
            variant="light"
            color="red"
            onClick={() => orders.refetch()}
            loading={orders.isFetching}
          >
            Try again
          </Button>
        </Alert>
      )}

      {orders.isPending && (
        <Stack gap="sm">
          {SKELETONS.map((skeleton) => (
            <Skeleton key={skeleton} height={92} radius="sm" />
          ))}
        </Stack>
      )}

      {orders.data?.length === 0 && (
        <Text c="dimmed">
          No orders yet.{" "}
          <Anchor component={Link} to="/products">
            Browse products
          </Anchor>
        </Text>
      )}

      {orders.data && orders.data.length > 0 && (
        <Stack gap="sm">
          {orders.data.map((order) => (
            <OrderRow key={order.id} order={order} />
          ))}
        </Stack>
      )}
    </Stack>
  );
}

function OrderRow({ order }: { order: Order }) {
  return (
    <Paper withBorder p="md" radius="sm">
      <Flex
        direction={{ base: "column", sm: "row" }}
        gap="sm"
        align={{ base: "stretch", sm: "center" }}
        justify="space-between"
      >
        <div>
          <Text fw={600} ff="monospace">
            {order.order_number}
          </Text>
          <Text size="xs" c="dimmed">
            {formatDate(order.created_at)}
          </Text>
        </div>

        <Group gap="md" w={{ base: "100%", sm: "auto" }}>
          <Group gap="xs">
            <Badge variant="light" color={STATUS_COLOR[order.status]}>
              {order.status}
            </Badge>
            <Text size="sm">{formatPrice(order.total_amount)}</Text>
            <Text size="sm" c="dimmed">
              {order.total_coins_used.toLocaleString()} coins
            </Text>
          </Group>

          <Anchor
            component={Link}
            to={`/orders/${order.id}`}
            size="sm"
            ml="auto"
          >
            View order
          </Anchor>
        </Group>
      </Flex>
    </Paper>
  );
}
