import { useLoaderData } from "react-router";

type User = {
  id: number;
  name: string;
};

// Loader to fetch user data
// export async function loader() {
//   const res = await fetch("http://localhost:8080/users");
//   if (!res.ok) {
//     throw new Response("Failed to fetch users", { status: res.status });
//   }
//   return res.json();
// }

export default function UserPage() {
  // const users = useLoaderData() as User[];

  return (
    <main className="p-4">
      <h1 className="font-bold text-2xl">Users</h1>
    </main>
  );
}
