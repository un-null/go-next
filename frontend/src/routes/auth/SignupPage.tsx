import {
  Alert,
  Anchor,
  Button,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useActionState, useTransition } from "react";
import { Link, Navigate, useNavigate, useSearchParams } from "react-router";

import { login, type SignupInput, signup } from "../../api/auth";
import { safeNext, withNext } from "../../lib/redirect";
import { setSession, useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";

type SignupState = { status: "idle" } | { status: "error"; message: string };

const IDLE: SignupState = { status: "idle" };

export function SignupPage() {
  const session = useSession();
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const next = safeNext(params.get("next"));

  const form = useForm({
    initialValues: { name: "", email: "", password: "" },
    validate: {
      name: (value) => (value.trim() ? null : "Enter your name"),
      email: (value) =>
        /^\S+@\S+\.\S+$/.test(value) ? null : "Enter an email address",
      password: (value) => (value.length >= 8 ? null : "At least 8 characters"),
    },
  });

  const createAccount = async (
    _previous: SignupState,
    values: SignupInput,
  ): Promise<SignupState> => {
    try {
      await signup(values);
      const data = await login({
        email: values.email,
        password: values.password,
      });
      setSession({ token: data.token, user: data.user });
      navigate(next, { replace: true });
      return IDLE;
    } catch (error) {
      return {
        status: "error",
        message:
          error instanceof Error
            ? error.message
            : "Could not create the account",
      };
    }
  };

  const [state, dispatch, isPending] = useActionState(createAccount, IDLE);
  const [, startTransition] = useTransition();

  const submit = form.onSubmit((values) => {
    if (isPending) return;
    startTransition(() => dispatch(values));
  });

  const isSlow = useSlowHint(isPending);

  if (session) {
    return <Navigate to={next} replace />;
  }

  return (
    <Stack gap="lg" maw={400} mx="auto" py="md">
      <div>
        <Title order={1} size="h2">
          Create an account
        </Title>
        <Text c="dimmed" mt={4}>
          You start with 1,000 coins. No card, no confirmation email.
        </Text>
      </div>

      {state.status === "error" && (
        <Alert color="red" variant="light" title="Could not create the account">
          <Text size="sm">{state.message}</Text>
        </Alert>
      )}

      <form onSubmit={submit}>
        <Stack gap="sm">
          <TextInput
            label="Name"
            autoComplete="name"
            withAsterisk
            {...form.getInputProps("name")}
          />
          <TextInput
            label="Email"
            type="email"
            autoComplete="email"
            withAsterisk
            {...form.getInputProps("email")}
          />
          <PasswordInput
            label="Password"
            autoComplete="new-password"
            description="At least 8 characters."
            withAsterisk
            {...form.getInputProps("password")}
          />

          <Button type="submit" fullWidth loading={isPending}>
            Create account
          </Button>

          {isSlow && (
            <Text size="xs" c="dimmed">
              Waking the API up. The first request after a break can take about
              twenty seconds.
            </Text>
          )}
        </Stack>
      </form>

      <Text size="sm" c="dimmed">
        Already have an account?{" "}
        <Anchor component={Link} to={withNext("/login", next)}>
          Sign in
        </Anchor>
      </Text>
    </Stack>
  );
}
