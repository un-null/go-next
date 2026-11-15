import { Route, Routes } from "react-router";

import { AppLayout } from "./components/AppLayout";
import { ProductDetailPage } from "./routes/products/ProductDetailPage";
import { ProductsPage } from "./routes/products/ProductsPage";

export default function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route path="/" element={<p>go-next</p>} />
        <Route path="/products" element={<ProductsPage />} />
        <Route path="/products/:id" element={<ProductDetailPage />} />
      </Route>
    </Routes>
  );
}
