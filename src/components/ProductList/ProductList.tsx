import React, { useState, useEffect } from 'react';
import api from '../../api/api';
import { Product, StockItem, ProductSize, User } from '../../types';
import ProductModal from '../ProductModal/ProductModal';
import styles from "./ProductList.module.css"

interface ProductListProps {
  outletId: string;
}

const ProductList: React.FC<ProductListProps> = ({ outletId }) => {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [bookingLoading, setBookingLoading] = useState<boolean>(false);

  useEffect(() => {
    if (outletId) {
      loadProducts(outletId);
    } else {
      setProducts([]);
    }
  }, [outletId]);

  const transformBackendData = (backendData: StockItem[]): Product[] => {
    const productsMap = new Map<string, Product>();
    
    backendData.forEach(item => {
      const productId = item.product.id;
      
      if (!productsMap.has(productId)) {
        productsMap.set(productId, {
          id: productId,
          name: item.product.name,
          description: item.product.description,
          price: item.product.price,
          stock: 0,
          image_url: item.product.image?.image_path,
          sizes: []
        });
      }
      
      const product = productsMap.get(productId)!;
      
      const sizeItem: ProductSize = {
        size: item.size,
        amount: item.amount,
        available: item.amount > 0
      };
      
      product.sizes.push(sizeItem);
      product.stock += item.amount;
    });
    
    return Array.from(productsMap.values());
  };

  const loadProducts = async (outletId: string): Promise<void> => {
    setLoading(true);
    setError('');
    
    try {
      const response = await api.getProductsByOutlet(outletId);
      const backendData: StockItem[] = response.data || [];
      
      const transformedProducts = transformBackendData(backendData);
      
      setProducts(transformedProducts);
    } catch (error: any) {
      console.error('Ошибка загрузки товаров:', error);
      setError('Не удалось загрузить товары. Попробуйте снова.');
    } finally {
      setLoading(false);
    }
  };

  const handleProductClick = (product: Product): void => {
    setSelectedProduct(product);
    setIsModalOpen(true);
  };

  const handleCloseModal = (): void => {
    setIsModalOpen(false);
    setSelectedProduct(null);
  };

 const handleBookProduct = async (productId: string, size?: number): Promise<void> => {
    const token = localStorage.getItem('token');
    if (!token) {
      alert('Для бронирования необходимо войти в систему');
      return;
    }

    if (!outletId) {
      alert('Для бронирования необходимо выбрать точку продаж');
      return;
    }

    const userStr = localStorage.getItem('user');
    if (!userStr) {
      alert('Ошибка: информация о пользователе не найдена');
      return;
    }

    const productToBook = products.find(p => p.id === productId);
    if (!productToBook) {
      alert('Ошибка: товар не найден');
      return;
    }

    let selectedSizeItem: ProductSize | undefined;
    if (size !== undefined) {
      selectedSizeItem = productToBook.sizes.find(s => s.size === size);
      if (!selectedSizeItem || selectedSizeItem.amount === 0) {
        alert('Выбранный размер недоступен');
        return;
      }
    } else if (productToBook.sizes.length > 0) {
      alert('Пожалуйста, выберите размер');
      return;
    }

    setBookingLoading(true);

    try {
      const user: User = JSON.parse(userStr);

      const orderData = {
        UserID: user.id,
        SalesOutletID: outletId,
        OrderItems: [
          {
            ProductID: productId,
            Amount: 1,
            Price: productToBook.price,
            Size: size || 0,
          }
        ]
      };

      const response = await api.createOrder(orderData);
    
      if (response.status === 201) {
        alert('Товар успешно забронирован!');
        await loadProducts(outletId);
        setIsModalOpen(false);
        setSelectedProduct(null);
      } else {
        throw new Error(`Ошибка сервера: ${response.status}`);
      }
    } catch (error: any) {
      console.error('Полная ошибка бронирования:', error);   
      alert('Не удалось забронировать товар');
    } finally {
      setBookingLoading(false);
    }
  };

  const formatPrice = (price: number): string => {
    return price?.toLocaleString('ru-RU') || '0';
  };

  const formatSizes = (sizes: ProductSize[]): string => {
    const availableSizes = sizes
      .filter(size => size.available)
      .map(size => size.size)
      .sort((a, b) => a - b);
    
    if (availableSizes.length === 0) return 'Нет размеров';
    if (availableSizes.length <= 3) return availableSizes.join(', ');
    
    return `${availableSizes[0]}-${availableSizes[availableSizes.length - 1]}`;
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <p className={styles.loadingText}>Загрузка товаров...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.errorText}>{error}</p>
        <button 
          onClick={() => outletId && loadProducts(outletId)}
          className={styles.retryButton}
        >
          Попробовать снова
        </button>
      </div>
    );
  }

  if (!outletId) {
    return (
      <div className={styles.emptyState}>
        <p className={styles.emptyText}>Выберите точку продаж, чтобы увидеть товары</p>
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className={styles.emptyState}>
        <p className={styles.emptyText}>В выбранной точке нет товаров</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2 className={styles.title}>Товары в точке</h2>
        <p className={styles.countText}>Найдено товаров: {products.length}</p>
      </div>
      
      <div className={styles.productsGrid}>
        {products.map((product: Product) => (
          <div 
            key={product.id} 
            className={product.stock > 0 ? styles.inStockCard : styles.outOfStockCard}
            onClick={() => handleProductClick(product)}
            title="Нажмите для подробной информации"
          >
            <div className={styles.productImage}>
              {product.image_url ? (
                <img 
                  src={product.image_url} 
                  alt={product.name}
                  className={styles.image}
                  onError={(e) => {
                    const target = e.target as HTMLImageElement;
                    target.style.display = 'none';
                  }}
                />
              ) : (
                <div className={styles.placeholderImage}>
                  <span className={styles.placeholderText}>Нет фото</span>
                </div>
              )}
              {product.stock > 0 && (
                <div className={styles.inStockBadge}>
                  В наличии
                </div>
              )}
            </div>
            
            <div className={styles.productInfo}>
              <h3 className={styles.productName}>{product.name}</h3>
              <div className={styles.stockInfo}>
                {product.sizes.length > 0 && (
                  <span className={styles.sizesText}>
                    Размеры: {formatSizes(product.sizes)}
                  </span>
                )}
              </div>
              <div className={styles.detailsRow}>
                <span className={styles.price}>{formatPrice(product.price)} ₽</span>
                <span className={styles.inStockText}>
                  Доступно: {product.stock} шт.
                </span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {isModalOpen && selectedProduct && (
        <ProductModal
          product={selectedProduct}
          onClose={handleCloseModal}
          onBook={handleBookProduct}
        />
      )}
    </div>
  );
};

export default ProductList;