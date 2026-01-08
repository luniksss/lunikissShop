import React, { useState } from 'react';
import { Product, ProductSize } from '../../types';
import styles from "./ProductModal.module.css"

interface ProductModalProps {
  product: Product;
  onClose: () => void;
  onBook: (productId: string, size?: number) => Promise<void>;
  bookingLoading?: boolean;
}

const ProductModal: React.FC<ProductModalProps> = ({ 
  product, 
  onClose, 
  onBook,
  bookingLoading = false 
}) => {
  const [selectedSize, setSelectedSize] = useState<number | null>(null);
  const [hoveredSize, setHoveredSize] = useState<number | null>(null);
  const [isBooking, setIsBooking] = useState<boolean>(false);
  
  const formatPrice = (price: number): string => {
    return price?.toLocaleString('ru-RU') || '0';
  };

   const handleBookClick = async () => {
    if (bookingLoading || isBooking) return;
    
    setIsBooking(true);
    try {
      await onBook(product.id, selectedSize || undefined);
    } catch (error) {
    } finally {
      setIsBooking(false);
    }
  };

  const getAvailableSizes = (): ProductSize[] => {
    return product.sizes.filter(size => size.available);
  };

  const getTotalStock = (): number => {
    return product.sizes.reduce((total, size) => total + size.amount, 0);
  };

  const getSizeButtonClass = (sizeItem: ProductSize): string => {
    const isSelected = selectedSize === sizeItem.size;
    const isHovered = hoveredSize === sizeItem.size;
    const isDisabled = sizeItem.amount === 0;
    
    let className = `${styles.sizeButton} ${styles.sizeButtonDefault}`;
    
    if (isDisabled) {
      className = `${styles.sizeButton} ${styles.sizeButtonDisabled}`;
    } else if (isSelected) {
      className = `${styles.sizeButton} ${styles.sizeButtonSelected}`;
    } else if (isHovered) {
      className = `${styles.sizeButton} ${styles.sizeButtonHover}`;
    }
    
    return className;
  };

  const getBookButtonClass = (): string => {
    const isDisabled = getTotalStock() === 0 || 
      (product.sizes.length > 0 && !selectedSize) || 
      isBooking || 
      bookingLoading;
    
    return `${styles.bookButton} ${isDisabled ? styles.bookButtonDisabled : styles.bookButtonActive}`;
  };

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
        <button className={styles.closeButton} onClick={onClose} aria-label="Закрыть">
          ×
        </button>
        
        <div className={styles.modalBody}>
          <div className={styles.imageSection}>
            {product.image_url ? (
              <img 
                src={product.image_url} 
                alt={product.name}
                className={styles.mainImage}
                onError={(e) => {
                  const target = e.target as HTMLImageElement;
                  target.style.display = 'none';
                  target.parentElement!.innerHTML = '<div style="width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #999; font-size: 16px;">Нет фото</div>';
                }}
              />
            ) : (
              <div className={styles.placeholderImage}>
                <span>Нет фото</span>
              </div>
            )}
          </div>
          
          <div className={styles.infoSection}>
            <h2 className={styles.productName}>{product.name}</h2>
            
            <div className={styles.priceSection}>
              <span className={styles.priceLabel}>Цена:</span>
              <span className={styles.price}>{formatPrice(product.price)} ₽</span>
            </div>
            <div className={styles.delimiterImage}></div>
            
            <div className={styles.descriptionSection}>
              <h3 className={styles.sectionTitle}>Описание</h3>
              <p className={styles.description}>
                {product.description || 'Описание отсутствует'}
              </p>
            </div>
            
            <div className={styles.sizesSection}>
              <h3 className={styles.sectionTitle}>Доступные размеры</h3>
              {getAvailableSizes().length === 0 ? (
                <p className={styles.noSizesText}>Нет доступных размеров</p>
              ) : (
                <div className={styles.sizesContainer}>
                  {getAvailableSizes().map((sizeItem) => (
                    <button
                      key={sizeItem.size}
                      className={getSizeButtonClass(sizeItem)}
                      onClick={() => sizeItem.amount > 0 && setSelectedSize(sizeItem.size)}
                      onMouseEnter={() => !sizeItem.amount && setHoveredSize(sizeItem.size)}
                      onMouseLeave={() => setHoveredSize(null)}
                      disabled={sizeItem.amount === 0}
                      title={`${sizeItem.amount} шт.`}
                    >
                      <span className={styles.sizeText}>{sizeItem.size}</span>
                      <span className={styles.sizeAmount}>{sizeItem.amount} шт.</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
            
            <div className={styles.stockSection}>
              {selectedSize && (
                <div className={styles.selectedSizeInfo}>
                  <span className={styles.selectedSizeLabel}>Выбран размер:</span>
                  <span className={styles.selectedSizeValue}>{selectedSize}</span>
                </div>
              )}
            </div>
            
            <div className={styles.actionSection}>
              <button
                onClick={handleBookClick}
                disabled={getTotalStock() === 0 || (product.sizes.length > 0 && !selectedSize) || isBooking || bookingLoading}
                 className={getBookButtonClass()}
              >
                {isBooking || bookingLoading 
                  ? 'Бронирование...' 
                  : getTotalStock() === 0 
                    ? 'Нет в наличии' 
                    : product.sizes.length > 0 && !selectedSize
                      ? 'Выберите размер'
                      : `Забронировать${selectedSize ? ` (размер ${selectedSize})` : ''}`
                }
              </button>
              
              <p className={styles.bookingInfo}>
                Товар будет зарезервирован на 24 часа. Оплата при получении в выбранной точке.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductModal;