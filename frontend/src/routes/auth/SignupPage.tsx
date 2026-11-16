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
import { Link, Navigate, useNavigate } from "react-router";

import { login, type SignupInput, signup } from "../../api/auth";
import { setSession, useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";

export function SignupPage() {
  const session = useSession();
  const navigate = useNavigate();

  const form = useForm({
    initialValues: { name: "", email: "", password: "" },
    validate: {
      name: (value) => (value.trim() ? null : "Enter your name"),
      email: (value) =>
        /^\S+@\S+\.\S+$/.test(value) ? null : "Enter an email address",
      password: (value) => (value.length >= 8 ? null : "At least 8 characters"),
    },
  });

  const createAccount = useMutation({
    mutationFn: async (values: SignupInput) => {
      await signup(values);
      return login({ email: values.email, password: values.password });
    },
    onSuccess: (data) => {
      setSession({ token: data.token, user: data.user });
      navigate("/products");
    },
  });

  const isSlow = useSlowHint(createAccount.isPending);

  if (session) {
    return <Navigate to="/products" replace />;
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

      {createAccount.isError && (
        <Alert color="red" variant="light" title="Could not create the account">
          <Text size="sm">{createAccount.error.message}</Text>
        </Alert>
      )}

      <form onSubmit={form.onSubmit((values) => createAccount.mutate(values))}>
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

          <Button type="submit" fullWidth loading={createAccount.isPending}>
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
        <Anchor component={Link} to="/login">
          Sign in
        </Anchor>
      </Text>
    </Stack>
  );
}
