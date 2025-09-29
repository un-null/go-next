import { useLoaderData } from "react-router";
import type { Route } from "./+types/product";

type Product = {
  id: number;
  category_id: number;
  name: string;
  description: string;
  price: number;
  stock_quantity: number;
  image_url: string;
  average_rating: number;
  total_comments: number;
  created_at: string;
  updated_at: string;
};

export async function loader({ context }: Route.LoaderArgs) {
  const res = await fetch(`${context.cloudflare.env.API_URL}/products`);
  if (!res.ok) {
    throw new Response("Failed to fetch product", { status: res.status });
  }
  return res.json();
}

export default function UserPage() {
  const products: Product[] = useLoaderData();

  console.log(products);

  return (
    <main className="space-y-8 p-4">
      <h1 className="font-bold text-2xl">Products</h1>
      <ul className="grid grid-cols-3 gap-4">
        {products.map((p) => (
          <li key={p.id} className="border p-4">
            <img
              src={p.image_url}
              alt={p.name}
              className="aspect-video w-full object-contain"
            />
            <p className="px-4">{p.name}</p>
          </li>
        ))}
      </ul>
    </main>
  );
}
