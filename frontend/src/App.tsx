import { Route, Routes } from "react-router";

import { AppLayout } from "./components/AppLayout";
import { LoginPage } from "./routes/auth/LoginPage";
import { SignupPage } from "./routes/auth/SignupPage";
import { HomePage } from "./routes/home/HomePage";
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
      </Route>
    </Routes>
  );
}
