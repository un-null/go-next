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
import { useActionState, useState, useTransition } from "react";
import { Link, Navigate, useNavigate, useSearchParams } from "react-router";

import { login } from "../../api/auth";
import { safeNext, withNext } from "../../lib/redirect";
import { setSession, useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";

export const DEMO = { email: "demo@go-next.dev", password: "go-next-demo" };

interface SignInValues {
  email: string;
  password: string;
}

type SignInState = { status: "idle" } | { status: "error"; message: string };

const IDLE: SignInState = { status: "idle" };

export function LoginPage() {
  const session = useSession();
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const next = safeNext(params.get("next"));

  const form = useForm({
    initialValues: { email: "", password: "" },
    validate: {
      email: (value) =>
        /^\S+@\S+\.\S+$/.test(value) ? null : "Enter an email address",
      password: (value) => (value ? null : "Enter your password"),
    },
  });

  const signIn = async (
    _previous: SignInState,
    values: SignInValues,
  ): Promise<SignInState> => {
    try {
      const data = await login(values);
      setSession({ token: data.token, user: data.user });
      navigate(next, { replace: true });
      return IDLE;
    } catch (error) {
      return {
        status: "error",
        message: error instanceof Error ? error.message : "Could not sign in",
      };
    }
  };

  const [state, dispatch, isPending] = useActionState(signIn, IDLE);
  const [, startTransition] = useTransition();

  const [isDemo, setIsDemo] = useState(false);

  const submit = form.onSubmit((values) => {
    if (isPending) return;
    setIsDemo(false);
    startTransition(() => dispatch(values));
  });

  const signInAsDemo = () => {
    if (isPending) return;
    setIsDemo(true);
    startTransition(() => dispatch(DEMO));
  };

  const isSlow = useSlowHint(isPending);

  if (session) {
    return <Navigate to={next} replace />;
  }

  return (
    <Stack gap="lg" maw={400} mx="auto" py="md">
      <div>
        <Title order={1} size="h2">
          Sign in
        </Title>
        <Text c="dimmed" mt={4}>
          Your cart and coins come back with you.
        </Text>
      </div>

      {state.status === "error" && (
        <Alert color="red" variant="light" title="Could not sign in">
          <Text size="sm">{state.message}</Text>
        </Alert>
      )}

      <form onSubmit={submit}>
        <Stack gap="sm">
          <TextInput
            label="Email"
            type="email"
            autoComplete="email"
            withAsterisk
            {...form.getInputProps("email")}
          />
          <PasswordInput
            label="Password"
            autoComplete="current-password"
            withAsterisk
            {...form.getInputProps("password")}
          />

          <Button type="submit" fullWidth loading={isPending && !isDemo}>
            Sign in
          </Button>

          {isSlow && (
            <Text size="xs" c="dimmed">
              Waking the API up. The first request after a break can take about
              twenty seconds.
            </Text>
          )}
        </Stack>
      </form>

      <Stack gap={6}>
        <Button
          variant="default"
          fullWidth
          onClick={signInAsDemo}
          loading={isPending && isDemo}
        >
          Use the demo account
        </Button>
        <Text size="xs" c="dimmed">
          The demo account is shared. Its cart and coins are whatever the last
          visitor left.
        </Text>
      </Stack>

      <Text size="sm" c="dimmed">
        No account yet?{" "}
        <Anchor component={Link} to={withNext("/signup", next)}>
          Create one
        </Anchor>
      </Text>
    </Stack>
  );
}
