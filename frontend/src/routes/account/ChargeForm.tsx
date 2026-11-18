import {
  Alert,
  Button,
  Group,
  NumberInput,
  Paper,
  Stack,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useQueryClient } from "@tanstack/react-query";
import { useActionState, useTransition } from "react";

import { ApiError } from "../../api/client";
import { chargeCoins } from "../../api/coins";
import { transactionKeys, userKeys } from "../../api/queries";
import { clearSession, type Session, setSession } from "../../lib/session";

// The handler rejects anything outside this range, so the form says so first.
const MAX_CHARGE = 10000;

interface ChargeValues {
  amount: number | string;
  description: string;
}

type ChargeState =
  | { status: "idle" }
  | { status: "ok"; message: string }
  | { status: "error"; message: string };

const IDLE: ChargeState = { status: "idle" };

export function ChargeForm({ session }: { session: Session }) {
  const queryClient = useQueryClient();
  const [, startTransition] = useTransition();

  const charge = async (
    _previous: ChargeState,
    values: ChargeValues,
  ): Promise<ChargeState> => {
    try {
      const response = await chargeCoins(session.token, {
        amount: Number(values.amount),
        description: values.description.trim(),
      });

      setSession({ ...session, user: response.user });
      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: userKeys.detail(session.user.id),
        }),
        queryClient.invalidateQueries({ queryKey: transactionKeys.all }),
      ]);

      return {
        status: "ok",
        message: `Charged ${response.transaction.amount.toLocaleString()} coins.`,
      };
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) {
        clearSession();
      }

      return {
        status: "error",
        message:
          error instanceof Error ? error.message : "Could not charge coins",
      };
    }
  };

  const [state, dispatch, isPending] = useActionState(charge, IDLE);

  const form = useForm({
    initialValues: { amount: 500, description: "Demo top-up" },
    validate: {
      amount: (value) => {
        const amount = Number(value);
        return Number.isInteger(amount) && amount >= 1 && amount <= MAX_CHARGE
          ? null
          : "Enter a whole number from 1 to 10,000";
      },
      description: (value) =>
        value.trim() ? null : "Say what the coins are for",
    },
  });

  const submit = form.onSubmit((values) => {
    startTransition(() => dispatch(values));
  });

  return (
    <Paper withBorder p="md" radius="sm">
      <Stack gap="sm">
        <div>
          <Title order={2} size="h4">
            Add coins
          </Title>
          <Text size="xs" c="dimmed">
            This is a demo top-up. No payment is involved.
          </Text>
        </div>

        {state.status === "ok" && (
          <Alert color="teal" variant="light">
            <Text size="sm">{state.message}</Text>
          </Alert>
        )}

        {state.status === "error" && (
          <Alert color="red" variant="light" title="Could not charge coins">
            <Text size="sm">{state.message}</Text>
          </Alert>
        )}

        <form onSubmit={submit}>
          <Stack gap="sm">
            <NumberInput
              label="Coins"
              description={`1 to ${MAX_CHARGE.toLocaleString()}`}
              min={1}
              max={MAX_CHARGE}
              allowDecimal={false}
              allowNegative={false}
              withAsterisk
              {...form.getInputProps("amount")}
            />
            <TextInput
              label="Description"
              withAsterisk
              {...form.getInputProps("description")}
            />

            <Group justify="flex-end">
              <Button type="submit" loading={isPending}>
                Add coins
              </Button>
            </Group>
          </Stack>
        </form>
      </Stack>
    </Paper>
  );
}
