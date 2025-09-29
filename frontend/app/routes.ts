import type { RouteConfig } from "@react-router/dev/routes";
import { index, route } from "@react-router/dev/routes";

export default [
  index("./routes/home.tsx"),
  route("products", "./routes/product.tsx"),
] satisfies RouteConfig;
