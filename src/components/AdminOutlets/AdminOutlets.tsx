import React, { useState, useEffect } from 'react';
import api from '../../api/api';
import { SalesOutlet } from '../../types';
import styles from './AdminOutlets.module.css'

const AdminOutlets: React.FC = () => {
  const [outlets, setOutlets] = useState<SalesOutlet[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [showAddForm, setShowAddForm] = useState<boolean>(false);
  const [newOutlet, setNewOutlet] = useState({
    address: '',
  });
  const [editForm, setEditForm] = useState({
    address: '',
  });

  useEffect(() => {
    fetchAllOutlets();
  }, []);

  const fetchAllOutlets = async () => {
    try {
      setLoading(true);
      setError(null);
      setDeleteError(null);
      const response = await api.getSalesOutlets();

      const data = response.data === null ? [] : response.data;
    
      if (Array.isArray(data)) {
        setOutlets(data);
      } else {
        console.error('Ожидался массив точек продаж, но получено:', data);
        setOutlets([]);
      }
    } catch (err: any) {
      console.error('Ошибка загрузки точек продаж:', err);
      setError('Не удалось загрузить точки продаж. Пожалуйста, попробуйте позже.');
    } finally {
      setLoading(false);
    }
  };

  const handleStartEdit = (outlet: SalesOutlet) => {
    setEditingId(outlet.id);
    setEditForm({
      address: outlet.address || '',
    });
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setEditForm({ address: '' });
  };

  const handleSaveEdit = async (outletId: string) => {
    if (!editForm.address.trim()) {
      alert('Пожалуйста, заполните поле');
      return;
    }

    try {
      await api.updateOutlet(outletId, editForm.address);
      
      setOutlets(prevOutlets =>
        prevOutlets.map(outlet =>
          outlet.id === outletId 
            ? { ...outlet, address: editForm.address }
            : outlet
        )
      );
      
      setEditingId(null);
      setEditForm({ address: '' });
      alert('Точка продаж успешно обновлена!');
    } catch (err: any) {
      console.error('Ошибка обновления точки:', err);
      alert('Не удалось обновить точку продаж');
    }
  };

  const handleDeleteOutlet = async (outletId: string, outletAddress: string) => {
    if (!window.confirm(`Вы уверены, что хотите удалить точку "${outletAddress}"? Это действие нельзя отменить.`)) {
      return;
    }

    setDeletingId(outletId);
    setDeleteError(null);

    try {
      await api.deleteOutlet(outletId);
      
      setOutlets(prevOutlets => 
        prevOutlets.filter(outlet => outlet.id !== outletId)
      );
      
      alert(`Точка "${outletAddress}" успешно удалена!`);
    } catch (err: any) {
      console.error('Ошибка удаления точки:', err);
      setDeleteError('Не удалось удалить точку продаж. Пожалуйста, попробуйте позже.');
    } finally {
      setDeletingId(null);
    }
  };

  const handleAddOutlet = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!newOutlet.address.trim()) {
      alert('Пожалуйста, заполните поле');
      return;
    }

    try {
      const response = await api.createOutlet(newOutlet.address);
      
      alert('Точка продаж успешно добавлена!');
      setShowAddForm(false);
      setNewOutlet({ address: '' });
      fetchAllOutlets();
    } catch (err: any) {
      console.error('Ошибка добавления точки:', err);
      alert('Не удалось добавить точку продаж');
    }
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <p>Загрузка точек продаж...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.errorText}>{error}</p>
        <button onClick={fetchAllOutlets} className={styles.retryButton}>
          Попробовать снова
        </button>
      </div>
    );
  }

  const displayOutlets = Array.isArray(outlets) ? outlets : [];

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h2 className={styles.title}>Управление точками продаж</h2>
          <p className={styles.subtitle}>Всего точек: {displayOutlets.length}</p>
        </div>
        <div className={styles.headerButtons}>
          <button onClick={() => setShowAddForm(true)} className={styles.addButton}>
            + Добавить точку
          </button>
          <button onClick={fetchAllOutlets} className={styles.refreshButton}>
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
        <div className={styles.addFormOverlay}>
          <div className={styles.addForm}>
            <div className={styles.addFormHeader}>
              <h3 className={styles.addFormTitle}>Добавить новую точку продаж</h3>
              <button 
                onClick={() => setShowAddForm(false)} 
                className={styles.closeFormButton}
              >
                ✕
              </button>
            </div>
            <form onSubmit={handleAddOutlet} className={styles.form}>
              <div className={styles.formGroup}>
                <label className={styles.label}>Адрес *</label>
                <input
                  type="text"
                  value={newOutlet.address}
                  onChange={(e) => setNewOutlet({...newOutlet, address: e.target.value})}
                  className={styles.input}
                  placeholder="Например: ул. Ленина, д. 1"
                  required
                />
              </div>
              <div className={styles.formButtons}>
                <button type="button" onClick={() => setShowAddForm(false)} className={styles.cancelButton}>
                  Отмена
                </button>
                <button type="submit" className={styles.submitButton}>
                  Добавить точку
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className={styles.tableContainer}>
        {displayOutlets.length === 0 ? (
          <div className={styles.emptyState}>
            <p>Точки продаж не найдены</p>
            <button onClick={() => setShowAddForm(true)} className={styles.addButton}>
              Добавить первую точку
            </button>
          </div>
        ) : (
          <table className={styles.table}>
            <thead>
              <tr className={styles.tableHeader}>
                <th className={styles.th}>ID</th>
                <th className={styles.th}>Адрес</th>
                <th className={styles.th}>Действия</th>
              </tr>
            </thead>
            <tbody>
              {displayOutlets.map((outlet) => (
                <tr key={outlet.id} className={styles.tableRow}>
                  <td className={styles.td}>
                    <span className={styles.outletId}>{outlet.id}</span>
                  </td>
                  <td className={styles.td}>
                    {editingId === outlet.id ? (
                      <input
                        type="text"
                        value={editForm.address}
                        onChange={(e) => setEditForm({...editForm, address: e.target.value})}
                        className={styles.editInput}
                      />
                    ) : (
                      outlet.address || 'Адрес не указан'
                    )}
                  </td>
                  <td className={styles.td}>
                    <div className={styles.actions}>
                      {editingId === outlet.id ? (
                        <>
                          <button
                            onClick={() => handleSaveEdit(outlet.id)}
                            className={styles.saveButton}
                            disabled={!editForm.address.trim()}
                          >
                            Сохранить
                          </button>
                          <button
                            onClick={handleCancelEdit}
                            className={styles.cancelEditButton}
                          >
                            Отмена
                          </button>
                        </>
                      ) : (
                        <>
                          <button
                            onClick={() => handleStartEdit(outlet)}
                            className={styles.editButton}
                            title="Редактировать"
                          >
                            <img src='/icons/edit.svg' />
                          </button>
                          <button
                            onClick={() => handleDeleteOutlet(outlet.id, outlet.address || 'точка')}
                            disabled={deletingId === outlet.id}
                            className={deletingId === outlet.id ? styles.deleteButtonDisabled : styles.deleteButton}
                            title="Удалить точку"
                          >
                            {deletingId === outlet.id ? '...' : <img src='/icons/delete.png' />}
                          </button>
                        </>
                      )}
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

export default AdminOutlets;