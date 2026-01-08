import React, { useState, useEffect } from 'react';
import api from '../../api/api';
import { SalesOutlet, Product, StockItem } from '../../types';
import OutletOptions from '../OutletOptions';
import styles from './AdminWarehouse.module.css'

const AdminWarehouse: React.FC = () => {
  const [outlets, setOutlets] = useState<SalesOutlet[]>([]);
  const [selectedOutlet, setSelectedOutlet] = useState<string>('');
  const [stockItems, setStockItems] = useState<StockItem[]>([]);
  const [allProducts, setAllProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [outletsLoading, setOutletsLoading] = useState<boolean>(true);
  const [productsLoading, setProductsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [showAddForm, setShowAddForm] = useState<boolean>(false);
  const [newStockItem, setNewStockItem] = useState({
    product_id: '',
    size: '',
    amount: 0,
  });
  const [updating, setUpdating] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<string | null>(null);

  useEffect(() => {
    fetchOutlets();
    fetchAllProducts();
  }, []);

  useEffect(() => {
    if (selectedOutlet) {
      fetchStockItems(selectedOutlet);
    } else {
      setStockItems([]);
    }
  }, [selectedOutlet]);

  const fetchOutlets = async () => {
    try {
      setOutletsLoading(true);
      const response = await api.getSalesOutlets();
      setOutlets(response.data || []);
    } catch (err: any) {
      console.error('Ошибка загрузки точек продаж:', err);
    } finally {
      setOutletsLoading(false);
    }
  };

  const fetchAllProducts = async () => {
    try {
      setProductsLoading(true);
      const response = await api.getAllProducts();
      setAllProducts(response.data || []);
    } catch (err: any) {
      console.error('Ошибка загрузки товаров:', err);
    } finally {
      setProductsLoading(false);
    }
  };

  const fetchStockItems = async (outletId: string) => {
    try {
      setLoading(true);
      setError(null);
      const response = await api.getProductsByOutlet(outletId);
      
      const data = response.data === null ? [] : response.data;
      
      if (Array.isArray(data)) {
        const transformedData = data.map((item: any) => ({
          product_id: item.product?.id || '',
          name: item.product?.name || 'Неизвестный товар',
          size: item.size || 0,
          amount: item.amount || 0,
          sales_outlet_id: item.sales_outlet_id || '',
          product: item.product || {}
        }));
        setStockItems(transformedData);
      } else {
        console.error('Ожидался массив товаров на складе, но получено:', data);
        setStockItems([]);
      }
    } catch (err: any) {
      console.error('Ошибка загрузки товаров на складе:', err);
      setError('Не удалось загрузить товары на складе');
    } finally {
      setLoading(false);
    }
  };

  const handleAddStockItem = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!selectedOutlet) {
      alert('Пожалуйста, выберите точку продаж');
      return;
    }
    
    if (!newStockItem.product_id || !newStockItem.size || newStockItem.amount < 0) {
      alert('Пожалуйста, заполните все поля корректно');
      return;
    }

    try {
      const stockItemData = {
        product: { id: newStockItem.product_id },
        sales_outlet_id: selectedOutlet,
        size: parseInt(newStockItem.size) || newStockItem.size,
        amount: newStockItem.amount,
      };
      await api.addStockItem(stockItemData);
      
      alert('Товар успешно добавлен на склад!');
      setShowAddForm(false);
      setNewStockItem({ product_id: '', size: '', amount: 0 });
      fetchStockItems(selectedOutlet);
    } catch (err: any) {
      console.error('Ошибка добавления товара на склад:', err);
      alert('Не удалось добавить товар на склад. Возможно, такой товар уже существует на этой точке.');
    }
  };

  const handleUpdateStockAmount = async (productId: string, size: string | number, newAmount: number) => {
    if (!selectedOutlet) return;
    
    try {
      setUpdating(`${productId}-${size}`);
      await api.updateStockItem(selectedOutlet, productId, newAmount, size);
      
      setStockItems(prevItems =>
        prevItems.map(item =>
          item.product.id === productId && item.size === size
            ? { ...item, amount: newAmount }
            : item
        )
      );
      
      alert('Количество товара успешно обновлено!');
    } catch (err: any) {
      console.error('Ошибка обновления количества:', err);
      alert('Не удалось обновить количество товара');
    } finally {
      setUpdating(null);
    }
  };

  const handleDeleteStockItem = async (productId: string, size: string | number, productName: string) => {
    if (!selectedOutlet) return;
    
    if (!window.confirm(`Вы уверены, что хотите удалить товар "${productName}" размера ${size} со склада?`)) {
      return;
    }

    try {
      setDeleting(`${productId}-${size}`);
      await api.deleteStockItem(selectedOutlet, productId, size);
    
      setStockItems(prevItems =>
        prevItems.filter(item => !(item.product.id === productId && item.size === size))
      );
      
      alert(`Товар "${productName}" размера ${size} успешно удален со склада!`);
    } catch (err: any) {
      console.error('Ошибка удаления товара:', err);
      alert('Не удалось удалить товар со склада. Возможно, есть связанные заказы.');
    } finally {
      setDeleting(null);
    }
  };

  const getProductName = (productId: string): string => {
    const stockItem = stockItems.find(item => item.product.id === productId);
    if (stockItem && stockItem.product && stockItem.product.name) {
      return stockItem.product.name;
    }
    
    const product = allProducts.find(p => p.id === productId);
    return product ? product.name : `Товар #${productId}`;
  };

  const getOutletName = (outletId: string): string => {
    const outlet = outlets.find(o => o.id === outletId);
    return outlet ? `${outlet.address}` : `Точка #${outletId}`;
  };

  const handleAmountChange = (productId: string, size: string | number, currentAmount: number) => {
    const newAmount = prompt('Введите новое количество:', currentAmount.toString());
    if (newAmount !== null) {
      const amountNum = parseInt(newAmount);
      if (!isNaN(amountNum) && amountNum >= 0) {
        handleUpdateStockAmount(productId, size, amountNum);
      } else {
        alert('Пожалуйста, введите корректное число');
      }
    }
  };

  const getUniqueProducts = (items: StockItem[]): StockItem[] => {
    const uniqueProducts = new Map<string, StockItem>();
    items.forEach(item => {
      if (item.product.id && !uniqueProducts.has(item.product.id)) {
        uniqueProducts.set(item.product.id, item);
      }
    });
    return Array.from(uniqueProducts.values());
  };

  if (outletsLoading || productsLoading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <p>Загрузка данных...</p>
      </div>
    );
  }

  if (error && !selectedOutlet) {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.errorText}>{error}</p>
        <button onClick={() => selectedOutlet && fetchStockItems(selectedOutlet)} className={styles.retryButton}>
          Попробовать снова
        </button>
      </div>
    );
  }

  const displayStockItems = Array.isArray(stockItems) ? stockItems : [];
  const uniqueProducts = getUniqueProducts(displayStockItems);

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h2 className={styles.title}>Управление складом</h2>
      </header>

      <div className={styles.controls}>
        <div className={styles.outletSelector}>
          <label className={styles.label}>Выберите точку продаж:</label>
          <select
            value={selectedOutlet}
            onChange={(e) => setSelectedOutlet(e.target.value)}
            className={styles.select}
          >
            <OutletOptions outlets={outlets} />
          </select>
        </div>

        {selectedOutlet && (
          <div className={styles.outletInfo}>
            <span className={styles.outletInfoText}>
              Точка: <strong>{getOutletName(selectedOutlet)}</strong>
            </span>
            <button
              onClick={() => setShowAddForm(true)}
              className={styles.addButton}
            >
              + Добавить товар на склад
            </button>
          </div>
        )}
      </div>

      {showAddForm && (
        <div className={styles.addFormOverlay}>
          <div className={styles.addForm}>
            <div className={styles.addFormHeader}>
              <h3 className={styles.addFormTitle}>Добавить товар на склад</h3>
              <button 
                onClick={() => setShowAddForm(false)} 
                className={styles.closeFormButton}
              >
                ✕
              </button>
            </div>
            <form onSubmit={handleAddStockItem} className={styles.form}>
              <div className={styles.formGroup}>
                <label className={styles.label}>Товар *</label>
                <select
                  value={newStockItem.product_id}
                  onChange={(e) => setNewStockItem({...newStockItem, product_id: e.target.value})}
                  className={styles.select}
                  required
                >
                  <option value="">Выберите товар</option>
                  {allProducts.map((product) => (
                    <option key={product.id} value={product.id}>
                      {product.name} (ID: {product.id})
                    </option>
                  ))}
                </select>
              </div>
              
              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label className={styles.label}>Размер *</label>
                  <input
                    type="text"
                    value={newStockItem.size}
                    onChange={(e) => setNewStockItem({...newStockItem, size: e.target.value})}
                    className={styles.input}
                    placeholder="Например: 42"
                    required
                  />
                </div>
                <div className={styles.formGroup}>
                  <label className={styles.label}>Количество *</label>
                  <input
                    type="number"
                    value={newStockItem.amount}
                    onChange={(e) => setNewStockItem({...newStockItem, amount: parseInt(e.target.value) || 0})}
                    className={styles.input}
                    min="0"
                    required
                  />
                </div>
              </div>
              
              <div className={styles.formButtons}>
                <button type="button" onClick={() => setShowAddForm(false)} className={styles.cancelButton}>
                  Отмена
                </button>
                <button type="submit" className={styles.submitButton}>
                  Добавить на склад
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {selectedOutlet ? (
        loading ? (
          <div className={styles.loadingContainer}>
            <div className={styles.spinner}></div>
            <p>Загрузка товаров на складе...</p>
          </div>
        ) : displayStockItems.length === 0 ? (
          <div className={styles.emptyState}>
            <p>На этой точке нет товаров на складе</p>
            <button onClick={() => setShowAddForm(true)} className={styles.addButton}>
              Добавить первый товар
            </button>
          </div>
        ) : (
          <div className={styles.tableContainer}>
            <table className={styles.table}>
              <thead>
                <tr className={styles.tableHeader}>
                  <th className={styles.th}>ID товара</th>
                  <th className={styles.th}>Название товара</th>
                  <th className={styles.th}>Размер</th>
                  <th className={styles.th}>Количество</th>
                  <th className={styles.th}>Действия</th>
                </tr>
              </thead>
              <tbody>
                {displayStockItems.map((item, index) => (
                  <tr key={`${item.product.id}-${item.size}-${index}`} className={styles.tableRow}>
                    <td className={styles.td}>{item.product.id || 'N/A'}</td>
                    <td className={styles.td}>
                      <strong>{getProductName(item.product.id)}</strong>
                    </td>
                    <td className={styles.td}>{item.size}</td>
                    <td className={styles.td}>
                      <span className={styles.amountCell}>
                        {item.amount}
                        {updating === `${item.product.id}-${item.size}` && (
                          <span className={styles.updatingText}>...</span>
                        )}
                      </span>
                    </td>
                    <td className={styles.td}>
                      <div className={styles.actions}>
                        <button
                          onClick={() => handleAmountChange(item.product.id, item.size, item.amount)}
                          disabled={updating === `${item.product.id}-${item.size}`}
                          className={styles.updateButton}
                          title="Изменить количество"
                        >
                          <img src='/icons/edit.svg' />
                        </button>
                        <button
                          onClick={() => handleDeleteStockItem(item.product.id, item.size, getProductName(item.product.id))}
                          disabled={deleting === `${item.product.id}-${item.size}`}
                          className={deleting === `${item.product.id}-${item.size}` ? styles.deleteButtonDisabled : styles.deleteButton}
                          title="Удалить этот размер товара"
                        >
                          {deleting === `${item.product.id}-${item.size}` ? '...' : <img src='/icons/delete.png' />}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            
            <div className={styles.summary}>
              <p>Всего товаров на точке: {uniqueProducts.length}</p>
              <p>Общее количество единиц: {displayStockItems.reduce((sum, item) => sum + item.amount, 0)}</p>
            </div>
          </div>
        )
      ) : (
        <div className={styles.emptyState}>
          <p>Выберите точку продаж для просмотра товаров на складе</p>
        </div>
      )}
    </div>
  );
};

export default AdminWarehouse;