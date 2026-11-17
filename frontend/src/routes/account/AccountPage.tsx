import {
  Alert,
  Button,
  Flex,
  Paper,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";

import { transactionsQuery } from "../../api/queries";
import type { CoinTransaction } from "../../api/types";
import { formatDate } from "../../lib/format";
import { type Session, useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";
import { useUser } from "../../lib/useUser";

export function AccountPage() {
  const session = useSession();

  return session ? <Account session={session} /> : null;
}

function Account({ session }: { session: Session }) {
  const user = useUser();
  const activity = useQuery(transactionsQuery(session.token));
  const isSlow = useSlowHint(activity.isPending);

  return (
    <Stack gap="lg">
      <Title order={1} size="h2">
        Account
      </Title>

      <Paper withBorder p="md" radius="sm">
        <Stack gap={4}>
          <Text size="xs" tt="uppercase" ff="monospace" c="dimmed">
            Balance
          </Text>
          <Text size="28px" fw={600} ff="monospace">
            {user ? user.coins.toLocaleString() : "—"} coins
          </Text>
          <Text size="xs" c="dimmed">
            Coins are demo currency, with no monetary value.
          </Text>
        </Stack>
      </Paper>

      {user && (
        <Paper withBorder p="md" radius="sm">
          <Stack gap="xs">
            <Row label="Name" value={user.name} />
            <Row label="Email" value={user.email} />
            <Row label="Member since" value={formatDate(user.created_at)} />
          </Stack>
        </Paper>
      )}

      <div>
        <Title order={2} size="h4" mb="xs">
          Coin activity
        </Title>

        {isSlow && activity.isPending && (
          <Alert color="yellow" variant="light" title="Waking the API up">
            The free instance (Render) sleeps after fifteen minutes. The first
            request can take about 20 sec. It is working, not broken.
          </Alert>
        )}

        {activity.isError && (
          <Alert
            color="red"
            variant="light"
            title="Could not load your activity"
          >
            <Text size="sm">{activity.error.message}</Text>
            <Button
              mt="sm"
              size="xs"
              variant="light"
              color="red"
              onClick={() => activity.refetch()}
              loading={activity.isFetching}
            >
              Try again
            </Button>
          </Alert>
        )}

        {activity.isPending && (
          <Stack gap="xs">
            <Skeleton height={58} radius="sm" />
            <Skeleton height={58} radius="sm" />
            <Skeleton height={58} radius="sm" />
          </Stack>
        )}

        {activity.data &&
          (activity.data.transactions.length === 0 ? (
            <Text size="sm" c="dimmed">
              No coin activity yet.
            </Text>
          ) : (
            <Stack gap="xs">
              {activity.data.transactions.map((transaction) => (
                <TransactionRow
                  key={transaction.id}
                  transaction={transaction}
                />
              ))}
            </Stack>
          ))}
      </div>
    </Stack>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <Flex justify="space-between" gap="md">
      <Text size="sm" c="dimmed">
        {label}
      </Text>
      <Text size="sm" ta="right" style={{ overflowWrap: "anywhere" }}>
        {value}
      </Text>
    </Flex>
  );
}

function TransactionRow({ transaction }: { transaction: CoinTransaction }) {
  const isCredit = transaction.amount > 0;

  return (
    <Paper withBorder p="sm" radius="sm">
      <Flex justify="space-between" gap="sm" wrap="wrap">
        <Stack gap={2} style={{ minWidth: 0 }}>
          <Text size="sm">{transaction.description}</Text>
          <Text size="xs" c="dimmed">
            {formatDate(transaction.created_at)}
          </Text>
        </Stack>
        <Stack gap={2} align="flex-end" ml="auto">
          <Text size="sm" ff="monospace" c={isCredit ? "teal.7" : "red.8"}>
            {isCredit ? "+" : "−"}
            {Math.abs(transaction.amount).toLocaleString()}
          </Text>
          <Text size="xs" c="dimmed" ff="monospace">
            {transaction.balance_after.toLocaleString()} coins
          </Text>
        </Stack>
      </Flex>
    </Paper>
  );
}
