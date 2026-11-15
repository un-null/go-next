import { Route, Routes } from "react-router";

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<p style={{ padding: 16 }}>go-next</p>} />
    </Routes>
  );
}
