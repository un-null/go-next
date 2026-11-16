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
import { useMutation } from "@tanstack/react-query";
import { Link, Navigate, useNavigate, useSearchParams } from "react-router";

import { login } from "../../api/auth";
import { safeNext, withNext } from "../../lib/redirect";
import { setSession, useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";

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

  const signIn = useMutation({
    mutationFn: login,
    onSuccess: (data) => {
      setSession({ token: data.token, user: data.user });
      navigate(next, { replace: true });
    },
  });

  const isSlow = useSlowHint(signIn.isPending);

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

      {signIn.isError && (
        <Alert color="red" variant="light" title="Could not sign in">
          <Text size="sm">{signIn.error.message}</Text>
        </Alert>
      )}

      <form onSubmit={form.onSubmit((values) => signIn.mutate(values))}>
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

          <Button type="submit" fullWidth loading={signIn.isPending}>
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

      <Text size="sm" c="dimmed">
        No account yet?{" "}
        <Anchor component={Link} to={withNext("/signup", next)}>
          Create one
        </Anchor>
      </Text>
    </Stack>
  );
}
