import { Route, Routes } from "react-router";

import { AppLayout } from "./components/AppLayout";
import { RequireAuth } from "./components/RequireAuth";
import { LoginPage } from "./routes/auth/LoginPage";
import { SignupPage } from "./routes/auth/SignupPage";
import { CartPage } from "./routes/cart/CartPage";
import { HomePage } from "./routes/home/HomePage";
import { OrderDetailPage } from "./routes/orders/OrderDetailPage";
import { OrdersPage } from "./routes/orders/OrdersPage";
import { ProductDetailPage } from "./routes/products/ProductDetailPage";
import { ProductsPage } from "./routes/products/ProductsPage";

export default function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/products" element={<ProductsPage />} />
        <Route path="/products/:id" element={<ProductDetailPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignupPage />} />

        <Route element={<RequireAuth />}>
          <Route path="/cart" element={<CartPage />} />
          <Route path="/orders" element={<OrdersPage />} />
          <Route path="/orders/:id" element={<OrderDetailPage />} />
        </Route>
      </Route>
    </Routes>
  );
}
