import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';

const API_URL = 'http://localhost:8080';

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors(): void {
    this.api.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        const token = localStorage.getItem('token');
        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    this.api.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          if (window.location.pathname !== '/login') {
            window.location.href = '/login';
          }
        }
        return Promise.reject(error);
      }
    );
  }

  async getSalesOutlets(): Promise<AxiosResponse> {
    return this.api.get('/outlet/list');
  }

  async getAllProducts(): Promise<AxiosResponse> {
    return this.api.get('/product/list');
  }

  async getProductsByOutlet(
    outletId: string,
  ): Promise<AxiosResponse> {
    return this.api.get(`/products/outlet/${outletId}`);
  }

  async getProductById(
    id: string,
  ): Promise<AxiosResponse> {
    return this.api.get(`/product/${id}`);
  }

  async createProduct(
    productData: any,
  ): Promise<AxiosResponse> {
    return this.api.post(`/api/v1/product/add`, productData)
  }

  async updateProduct(
    productId: string, 
    productData: any,
  ): Promise<AxiosResponse> {
    return this.api.post(`/api/v1/product/update/${productId}`, productData)
  }

  async deleteProduct(
    productId: string,
  ): Promise<AxiosResponse> {
    return this.api.delete(`/api/v1/product/delete/${productId}`)
  }

  async login(
    email: string, 
    password: string,
  ): Promise<AxiosResponse> {
    return this.api.post('/api/v1/auth/login', { email, password });
  }

  async register(
    userData: any,
  ): Promise<AxiosResponse> {
    return this.api.post('/api/v1/auth/register', userData);
  }

  async getUserOrders(
    userID: string,
  ): Promise<AxiosResponse> {
    return this.api.get(`api/v1/users/${userID}/orders`);
  }

  async createOrder(
    orderData: any,
  ): Promise<AxiosResponse> {
    return this.api.post('/api/v1/order', orderData);
  }

  async getOrderInfo(
    orderID: string,
  ): Promise<AxiosResponse> {
    return this.api.get(`/api/v1/orders/${orderID}`);
  }

  async deleteOrder(
    orderID: string,
  ): Promise<AxiosResponse> {
    return this.api.delete(`/api/v1/order/${orderID}`);
  }

  async deleteOrderItem(
    orderItemID: string,
  ): Promise<AxiosResponse> {
    return this.api.delete(`/api/v1/order-items/${orderItemID}`);
  }

  async getAllOrders(): Promise<AxiosResponse> {
    return this.api.get('/api/v1/orders/list');
  }

  async updateOrderStatus(
    orderId: string, 
    status: string,
  ): Promise<AxiosResponse> {
    return this.api.patch(`/api/v1/order/${orderId}/status`, { status });
  }

  async getAllUsers(): Promise<AxiosResponse> {
    return this.api.get('/api/v1/users');
  }

  async getUserById(
    userId: string,
  ): Promise<AxiosResponse> {
    return this.api.get(`/api/v1/users/${userId}`);
  }

  async updateUser(
    userData: any,
  ): Promise<AxiosResponse> {
    return this.api.put('/api/v1/users', userData);
  }

  async updateUserRole(
    userId: string, 
    role: string,
  ): Promise<AxiosResponse> {
    return this.api.patch(`/api/v1/users/${userId}/role`, { role });
  }

  async deleteUser(
    userId: string,
  ): Promise<AxiosResponse> {
    return this.api.delete(`/api/v1/users/${userId}`);
  }

  async createOutlet(
    newOutletAddress: string,
  ): Promise<AxiosResponse> {
    return this.api.post('/api/v1/outlet/add', newOutletAddress);
  }

  async updateOutlet(
    outletId: string, 
    newAddress: string,
  ): Promise<AxiosResponse> {
    return this.api.post(`/api/v1/outlet/update/${outletId}`, newAddress);
  }
  
  async deleteOutlet(
    outletId: string,
  ): Promise<AxiosResponse> {
    return this.api.delete(`/api/v1/outlet/delete/${outletId}`);
  }

  async addStockItem(
    stockItemData: any,
  ): Promise<AxiosResponse> {
    return this.api.post('/api/v1/stock/add', stockItemData)
  }

  async updateStockItem(
    selectedOutlet: string, 
    productId: string, 
    newAmount: number, 
    size: any,
  ): Promise<AxiosResponse> {
    return this.api.post(`/api/v1/stock/update/${selectedOutlet}/${productId}/${newAmount}/${size}`)
  }

  async deleteStockItem(
    selectedOutlet: string, 
    productId: string, 
    size: any,
  ): Promise<AxiosResponse> {
    return this.api.delete(`/api/v1/stock/delete/${selectedOutlet}/${productId}/${size}`)
  }
}

const apiService = new ApiService();
export default apiService;