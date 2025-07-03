import { useLoaderData } from "react-router";

type User = {
  id: number;
  name: string;
};

// Loader to fetch user data
export async function loader() {
  const res = await fetch("http://localhost:8080/users");
  if (!res.ok) {
    throw new Response("Failed to fetch users", { status: res.status });
  }
  return res.json();
}

export default function UserPage() {
  const users = useLoaderData() as User[];

  return (
    <main className="p-4">
      <h1 className="text-2xl font-bold">Users</h1>
      <ul className="mt-2 list-disc pl-5">
        {users.map((user) => (
          <li key={user.id}>
            {user.id}. {user.name}
          </li>
        ))}
      </ul>
    </main>
  );
}
