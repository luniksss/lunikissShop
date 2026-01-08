export interface SalesOutlet {
  id: string;
  name: string;
  address: string;
  phone?: string;
  email?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  stock: number;
  image_id?: string;
  image_url?: string;
  sizes: ProductSize[];
  created_at?: string;
}

export interface ProductInfo {
  id: string;
  name: string;
  description: string;
  price: number;
  image: {
    id: string;
    product_id: string;
    image_path: string;
  };
}

export interface ProductSize {
  size: number;
  amount: number;
  available: boolean;
}

export interface StockItem {
  size: number;
  amount: number;
  sales_outlet_id: string;
  product: ProductInfo;
}

export interface User {
  id: string;
  name: string;
  surname: string;
  email: string;
  role: string;
  phone?: string;
  default_outlet_id: string;
}

export interface UserInfo extends User {
  password: string;
}

export interface Order {
  id: string;
  user_id: string;
  sales_outlet_id: string;
  created_at: string;
  status_name: string;
}

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  product_name?: string;
  product_image?: string;
  amount: number;
  price: number;
  size: number;
}

export interface OrderDetails {
  order: Order;
  items: OrderItem[];
  total: number;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
  error?: string;
}