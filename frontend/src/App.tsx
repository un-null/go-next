import { Route, Routes } from "react-router";

import { AppLayout } from "./components/AppLayout";

export default function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route path="/" element={<p>go-next</p>} />
      </Route>
    </Routes>
  );
}
