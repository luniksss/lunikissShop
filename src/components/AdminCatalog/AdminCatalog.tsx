import React, { useState, useEffect, useRef } from 'react';
import api from '../../api/api';
import { ProductInfo } from '../../types';
import styles from './AdminCatalog.module.css'

interface ProductFormData {
  id: string;
  name: string;
  description: string;
  price: number;
  image: {
    image_path: string;
    product_id?: string;
  };
}

const AdminCatalog: React.FC = () => {
  const [products, setProducts] = useState<ProductInfo[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [showAddForm, setShowAddForm] = useState<boolean>(false);
  const [showEditForm, setShowEditForm] = useState<boolean>(false);
  const [editingProduct, setEditingProduct] = useState<ProductInfo | null>(null);
  const [deleting, setDeleting] = useState<string | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [imagePreviewUrl, setImagePreviewUrl] = useState<string>('');
  const [editImagePreviewUrl, setEditImagePreviewUrl] = useState<string>('');
  const [imageLoading, setImageLoading] = useState<boolean>(false);
  const [editImageLoading, setEditImageLoading] = useState<boolean>(false);
  const [imageError, setImageError] = useState<boolean>(false);
  const [editImageError, setEditImageError] = useState<boolean>(false);

  const imageTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const editImageTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const [newProduct, setNewProduct] = useState<ProductFormData>({
    id: '',
    name: '',
    description: '',
    price: 0,
    image: {
      image_path: '',
    },
  });

  const [editForm, setEditForm] = useState<ProductFormData>({
    id: '',
    name: '',
    description: '',
    price: 0,
    image: {
      image_path: '',
    },
  });

  useEffect(() => {
    fetchAllProducts();
    return () => {
      if (imageTimeoutRef.current) clearTimeout(imageTimeoutRef.current);
      if (editImageTimeoutRef.current) clearTimeout(editImageTimeoutRef.current);
    };
  }, []);

  useEffect(() => {
    if (newProduct.image.image_path) {
      setImageLoading(true);
      setImageError(false);
      if (imageTimeoutRef.current) clearTimeout(imageTimeoutRef.current);
      imageTimeoutRef.current = setTimeout(() => {
        setImagePreviewUrl(newProduct.image.image_path);
        setImageLoading(false);
      }, 500);
    } else {
      setImagePreviewUrl('');
      setImageError(false);
    }
  }, [newProduct.image.image_path]);

  useEffect(() => {
    if (editForm.image.image_path) {
      setEditImageLoading(true);
      setEditImageError(false);
      if (editImageTimeoutRef.current) clearTimeout(editImageTimeoutRef.current);
      editImageTimeoutRef.current = setTimeout(() => {
        setEditImagePreviewUrl(editForm.image.image_path);
        setEditImageLoading(false);
      }, 500);
    } else {
      setEditImagePreviewUrl('');
      setEditImageError(false);
    }
  }, [editForm.image.image_path]);

  const fetchAllProducts = async () => {
    try {
      setLoading(true);
      setError(null);
      setDeleteError(null);
      const response = await api.getAllProducts();

      const data = response.data === null ? [] : response.data;
    
      if (Array.isArray(data)) {
        setProducts(data);
      } else {
        console.error('Ожидался массив товаров, но получено:', data);
        setProducts([]);
      }
    } catch (err: any) {
      console.error('Ошибка загрузки товаров:', err);
      
      if (err.response?.status === 404) {
        setError('Эндпоинт для получения товаров не найден.');
      } else if (err.response?.status === 401) {
        setError('Необходима авторизация. Пожалуйста, войдите снова.');
      } else if (err.response?.status === 403) {
        setError('У вас нет прав для просмотра товаров.');
      } else {
        setError('Не удалось загрузить товары. Пожалуйста, попробуйте позже.');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleAddProduct = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newProduct.name.trim() || !newProduct.description.trim() || newProduct.price <= 0 || !newProduct.image.image_path.trim()) {
      alert('Пожалуйста, заполните все поля корректно. Цена должна быть больше 0.');
      return;
    }

    try {
      const productData = {
        ...newProduct,
      };

      await api.createProduct(productData);
      
      alert('Товар успешно добавлен!');
      setShowAddForm(false);
      resetNewProductForm();
      setImagePreviewUrl('');
      setImageError(false);
      fetchAllProducts();
    } catch (err: any) {
      console.error('Ошибка добавления товара:', err);
      alert('Не удалось добавить товар. Возможно, товар с таким названием уже существует.');
    }
  };

  const handleStartEdit = async (product: ProductInfo) => {
    try {
      const response = await api.getProductById(product.id);
      const fullProduct = response.data;
      
      setEditingProduct(fullProduct);
      setEditForm({
        id: fullProduct.id,
        name: fullProduct.name || '',
        description: fullProduct.description || '',
        price: fullProduct.price || 0,
        image: {
          image_path: fullProduct.image?.image_path || '',
          product_id: fullProduct.id,
        },
      });
      setEditImagePreviewUrl(fullProduct.image?.image_path || '');
      setEditImageError(false);
      setShowEditForm(true);
    } catch (err: any) {
      console.error('Ошибка загрузки данных товара:', err);
      alert('Не удалось загрузить данные товара для редактирования.');
    }
  };

  const handleCancelEdit = () => {
    setShowEditForm(false);
    setEditingProduct(null);
    setEditForm({
      id: '',
      name: '',
      description: '',
      price: 0,
      image: {
        image_path: '',
      },
    });
    setEditImagePreviewUrl('');
    setEditImageError(false);
  };

  const handleSaveEdit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!editForm.name.trim() || !editForm.description.trim() || editForm.price <= 0 || !editForm.image.image_path.trim()) {
      alert('Пожалуйста, заполните все поля корректно. Цена должна быть больше 0.');
      return;
    }

    try {
      const productData = {
        ...editForm,
        image: {
          ...editForm.image,
          id: editingProduct?.image.id || '',
          product_id: editingProduct?.id,
        },
      };

      await api.updateProduct(editForm.id, productData);
      
      alert('Товар успешно обновлен!');
      setShowEditForm(false);
      setEditingProduct(null);
      resetEditForm();
      setEditImagePreviewUrl('');
      setEditImageError(false);
      fetchAllProducts();
    } catch (err: any) {
      console.error('Ошибка обновления товара:', err);
      alert('Не удалось обновить товар.');
    }
  };

  const handleDeleteProduct = async (productId: string, productName: string) => {
    if (!window.confirm(`Вы уверены, что хотите удалить товар "${productName}"? Это действие нельзя отменить.`)) {
      return;
    }

    setDeleting(productId);
    setDeleteError(null);

    try {
      await api.deleteProduct(productId);
      
      setProducts(prevProducts => 
        prevProducts.filter(product => product.id !== productId)
      );
      
      alert(`Товар "${productName}" успешно удален!`);
    } catch (err: any) {
      console.error('Ошибка удаления товара:', err);
      setDeleteError('Не удалось удалить товар. Возможно, есть связанные заказы или товары на складе.');
    } finally {
      setDeleting(null);
    }
  };

  const resetNewProductForm = () => {
    setNewProduct({
      id: '',
      name: '',
      description: '',
      price: 0,
      image: {
        image_path: '',
      },
    });
  };

  const resetEditForm = () => {
    setEditForm({
      id: '',
      name: '',
      description: '',
      price: 0,
      image: {
        image_path: '',
      },
    });
  };

  const formatPrice = (price: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RUB',
      minimumFractionDigits: 0,
    }).format(price);
  };

  const isValidImageUrl = (url: string): boolean => {
    return url.startsWith('http') || url.startsWith('/');
  };

  const ProductImage: React.FC<{ src: string; alt: string; className: string;  }> = ({ src, alt, className }) => {
    const [hasError, setHasError] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
      setHasError(false);
      setIsLoading(true);
    }, [src]);

    if (!src) {
      return <div className={styles.noImage}>Нет изображения</div>;
    }

    return (
      <div style={{ position: 'relative', width: '60px', height: '60px' }}>
        {isLoading && (
          <div style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            backgroundColor: '#f8f9fa',
            borderRadius: '4px',
            fontFamily: 'Rubik',
            fontSize: '10px',
            color: '#6c757d',
          }}>
            Загрузка...
          </div>
        )}
        <img 
          src={src} 
          alt={alt}
          className={className}
          style={{
            display: hasError ? 'none' : 'block',
            opacity: isLoading ? 0 : 1,
            transition: 'opacity 0.3s',
          }}
          onError={() => {
            setHasError(true);
            setIsLoading(false);
          }}
          onLoad={() => setIsLoading(false)}
        />
        {hasError && (
          <div style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            backgroundColor: '#f8f9fa',
            borderRadius: '4px',
            fontSize: '10px',
            color: '#6c757d',
            textAlign: 'center' as 'center',
          }}>
            Ошибка загрузки
          </div>
        )}
      </div>
    );
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <p>Загрузка товаров...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.errorText}>{error}</p>
        <button onClick={fetchAllProducts} className={styles.retryButton}>
          Попробовать снова
        </button>
      </div>
    );
  }

  const displayProducts = Array.isArray(products) ? products : [];

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h2 className={styles.title}>Управление каталогом</h2>
          <p className={styles.subtitle}>Всего товаров: {displayProducts.length}</p>
        </div>
        <div className={styles.headerButtons}>
          <button onClick={() => setShowAddForm(true)} className={styles.addButton}>
            + Добавить товар
          </button>
          <button onClick={fetchAllProducts} className={styles.refreshButton}>
            Обновить
          </button>
        </div>
      </header>

      {deleteError && (
        <div className={styles.deleteError}>
          <p className={styles.deleteErrorText}>{deleteError}</p>
          <button 
            onClick={() => setDeleteError(null)}
            className={styles.dismissButton}
          >
            ✕
          </button>
        </div>
      )}

      {showAddForm && (
        <div className={styles.formOverlay}>
          <div className={styles.form}>
            <div className={styles.formHeader}>
              <h3 className={styles.formTitle}>Добавить новый товар</h3>
              <button 
                onClick={() => {
                  setShowAddForm(false);
                  resetNewProductForm();
                  setImagePreviewUrl('');
                  setImageError(false);
                }} 
                className={styles.closeFormButton}
              >
                ✕
              </button>
            </div>
            <form onSubmit={handleAddProduct} className={styles.formContent}>
              <div className={styles.formGroup}>
                <label className={styles.label}>Название товара *</label>
                <input
                  type="text"
                  value={newProduct.name}
                  onChange={(e) => setNewProduct({...newProduct, name: e.target.value})}
                  className={styles.input}
                  placeholder="Например: Футболка Black Swan"
                  required
                />
              </div>
              
              <div className={styles.formGroup}>
                <label className={styles.label}>Описание *</label>
                <textarea
                  value={newProduct.description}
                  onChange={(e) => setNewProduct({...newProduct, description: e.target.value})}
                  className={styles.textarea}
                  placeholder="Подробное описание товара"
                  rows={4}
                  required
                />
              </div>
              
              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label className={styles.label}>Цена (руб) *</label>
                  <input
                    type="number"
                    value={newProduct.price}
                    onChange={(e) => setNewProduct({...newProduct, price: parseFloat(e.target.value) || 0})}
                    className={styles.input}
                    min="1"
                    step="0.01"
                    required
                  />
                </div>
              </div>
              
              <div className={styles.formGroup}>
                <label className={styles.label}>Путь к изображению *</label>
                <input
                  type="text"
                  value={newProduct.image.image_path}
                  onChange={(e) => setNewProduct({
                    ...newProduct, 
                    image: { ...newProduct.image, image_path: e.target.value }
                  })}
                  className={styles.input}
                  placeholder="/products/image.jpg"
                  required
                />
                {newProduct.image.image_path && isValidImageUrl(newProduct.image.image_path) && (
                  <div className={styles.imagePreviewContainer}>
                    <p className={styles.imagePreviewLabel}>
                      Предпросмотр {imageLoading && '(загрузка...)'}:
                    </p>
                    {!imageLoading && imagePreviewUrl && (
                      <ProductImage
                        src={imagePreviewUrl}
                        alt="Preview"
                        className={styles.imagePreview}
                      />
                    )}
                  </div>
                )}
                {newProduct.image.image_path && !isValidImageUrl(newProduct.image.image_path) && (
                  <p className={styles.imageError}>Некорректный путь к изображению</p>
                )}
              </div>
              
              <div className={styles.formButtons}>
                <button 
                  type="button" 
                  onClick={() => {
                    setShowAddForm(false);
                    resetNewProductForm();
                    setImagePreviewUrl('');
                    setImageError(false);
                  }} 
                  className={styles.cancelButton}
                >
                  Отмена
                </button>
                <button type="submit" className={styles.submitButton}>
                  Добавить товар
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showEditForm && editingProduct && (
        <div className={styles.formOverlay}>
          <div className={styles.form}>
            <div className={styles.formHeader}>
              <h3 className={styles.formTitle}>Редактировать товар</h3>
              <button 
                onClick={handleCancelEdit} 
                className={styles.closeFormButton}
              >
                ✕
              </button>
            </div>
            <form onSubmit={handleSaveEdit} className={styles.formContent}>
              <div className={styles.formGroup}>
                <label className={styles.label}>Название товара *</label>
                <input
                  type="text"
                  value={editForm.name}
                  onChange={(e) => setEditForm({...editForm, name: e.target.value})}
                  className={styles.input}
                  required
                />
              </div>
              
              <div className={styles.formGroup}>
                <label className={styles.label}>Описание *</label>
                <textarea
                  value={editForm.description}
                  onChange={(e) => setEditForm({...editForm, description: e.target.value})}
                  className={styles.textarea}
                  rows={4}
                  required
                />
              </div>
              
              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label className={styles.label}>Цена (руб) *</label>
                  <input
                    type="number"
                    value={editForm.price}
                    onChange={(e) => setEditForm({...editForm, price: parseFloat(e.target.value) || 0})}
                    className={styles.input}
                    min="1"
                    step="0.01"
                    required
                  />
                </div>
              </div>
              
              <div className={styles.formGroup}>
                <label className={styles.label}>Путь к изображению *</label>
                <input
                  type="text"
                  value={editForm.image.image_path}
                  onChange={(e) => setEditForm({
                    ...editForm, 
                    image: { ...editForm.image, image_path: e.target.value }
                  })}
                  className={styles.input}
                  required
                />
                {editForm.image.image_path && isValidImageUrl(editForm.image.image_path) && (
                  <div className={styles.imagePreviewContainer}>
                    <p className={styles.imagePreviewLabel}>
                      Текущее изображение {editImageLoading && '(загрузка...)'}:
                    </p>
                    {!editImageLoading && editImagePreviewUrl && (
                      <ProductImage
                        src={editImagePreviewUrl}
                        alt="Preview"
                        className={styles.imagePreview}
                      />
                    )}
                  </div>
                )}
                {editForm.image.image_path && !isValidImageUrl(editForm.image.image_path) && (
                  <p className={styles.imageError}>Некорректный путь к изображению</p>
                )}
              </div>
              
              <div className={styles.formButtons}>
                <button 
                  type="button" 
                  onClick={handleCancelEdit} 
                  className={styles.cancelButton}
                >
                  Отмена
                </button>
                <button type="submit" className={styles.submitButton}>
                  Сохранить изменения
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className={styles.tableContainer}>
        {displayProducts.length === 0 ? (
          <div className={styles.emptyState}>
            <p>Товары не найдены</p>
            <button onClick={() => setShowAddForm(true)} className={styles.addButton}>
              Добавить первый товар
            </button>
          </div>
        ) : (
          <table className={styles.table}>
            <thead>
              <tr className={styles.tableHeader}>
                <th className={styles.th}>Изображение</th>
                <th className={styles.th}>ID</th>
                <th className={styles.th}>Название</th>
                <th className={styles.th}>Описание</th>
                <th className={styles.th}>Цена</th>
                <th className={styles.th}>Действия</th>
              </tr>
            </thead>
            <tbody>
              {displayProducts.map((product) => (
                <tr key={product.id} className={styles.tableRow}>
                  <td className={styles.td}>
                    <ProductImage
                      src={product.image?.image_path || ''}
                      alt={product.name}
                      className={styles.productImage}
                    />
                  </td>
                  <td className={styles.td}>
                    <span className={styles.productId}>{product.id}</span>
                  </td>
                  <td className={styles.td}>
                    <strong>{product.name}</strong>
                  </td>
                  <td className={styles.td}>
                    <div className={styles.descriptionCell}>
                      {product.description.length > 100 
                        ? `${product.description.substring(0, 100)}...` 
                        : product.description
                      }
                    </div>
                  </td>
                  <td className={styles.td}>
                    <span className={styles.price}>{formatPrice(product.price)}</span>
                  </td>
                  <td className={styles.td}>
                    <div className={styles.actions}>
                      <button
                        onClick={() => handleStartEdit(product)}
                        className={styles.editButton}
                        title="Редактировать"
                      >
                        <img src='/icons/edit.svg' />
                      </button>
                      <button
                        onClick={() => handleDeleteProduct(product.id, product.name)}
                        disabled={deleting === product.id}
                        className={deleting === product.id ? styles.deleteButtonDisabled : styles.deleteButton}
                        title="Удалить товар"
                      >
                        {deleting === product.id ? '...' : <img src='/icons/delete.png' />}
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default AdminCatalog;