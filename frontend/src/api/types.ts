export interface Product {
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
}

export interface Category {
  id: number;
  name: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  coins: number;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface SignupResponse {
  message: string;
  user: User;
}

export interface CartItem {
  id: number;
  user_id: string;
  product_id: number;
  quantity: number;
  created_at: string;
  updated_at: string;
}

export interface CartResponse {
  items: CartItem[];
  total_items: number;
  success: boolean;
  message?: string;
}

export interface OrderItem {
  id: number;
  order_id: number;
  product_id: number;
  product_name: string;
  product_price: number;
  quantity: number;
  subtotal: number;
}

export type OrderStatus = "pending" | "completed" | "cancelled";

export interface Order {
  id: number;
  user_id: string;
  order_number: string;
  total_amount: number;
  total_coins_used: number;
  status: OrderStatus;
  created_at: string;
  updated_at: string;
  items?: OrderItem[];
}

export interface CreateOrderInput {
  items: { product_id: number; quantity: number }[];
}
