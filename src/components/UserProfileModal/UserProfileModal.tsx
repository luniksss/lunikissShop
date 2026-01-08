import React, { useState, useEffect } from 'react';
import api from '../../api/api';
import { User } from '../../types';
import { useOutlets } from '../../hooks/useOutlets';
import OutletOptions from '../OutletOptions';
import styles from './UserProfileModal.module.css';

interface UserProfileModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const UserProfileModal: React.FC<UserProfileModalProps> = ({ isOpen, onClose }) => {
  const [userData, setUserData] = useState<User>({
    id: '',
    email: '',
    name: '',
    surname: '',
    phone: '',
    role: '',
    default_outlet_id: '',
  });
  
  const [originalData, setOriginalData] = useState<User | null>(null);
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const [saving, setSaving] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const { outlets, loading: outletsLoading, getOutletById } = useOutlets();

  useEffect(() => {
    if (isOpen) {
      loadUserData();
    }
  }, [isOpen]);

  const loadUserData = async (): Promise<void> => {
    try {
      setLoading(true);
      setError('');
      
      const user = localStorage.getItem('user');
      if (!user) {
        throw new Error('Пользователь не найден');
      }

      const parsedUser = JSON.parse(user);
      if (!parsedUser.id) {
        throw new Error('ID пользователя не найден');
      }

      const response = await api.getUserById(parsedUser.id);
      
      const userDataFromApi = {
        id: response.data.id,
        email: response.data.email,
        name: response.data.name,
        surname: response.data.surname,
        phone: response.data.phone || '',
        role: response.data.role,
        default_outlet_id: response.data.default_outlet_id || '',
      };

      setUserData(userDataFromApi);
      setOriginalData(userDataFromApi);
    } catch (err: any) {
      console.error('Ошибка загрузки данных пользователя:', err);
      setError('Не удалось загрузить данные пользователя');
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>): void => {
    const { name, value } = e.target;
    setUserData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    setSaving(true);
    setError('');
    setSuccess('');

    try {
      const updateData = {
        id: userData.id,
        name: userData.name,
        surname: userData.surname,
        email: userData.email,
        phone: userData.phone,
        role: userData.role,
        default_outlet_id: userData.default_outlet_id || null,
      };

      await api.updateUser(updateData);
      
      const currentUser = JSON.parse(localStorage.getItem('user') || '{}');
      const updatedUser = { 
        ...currentUser, 
        name: userData.name,
        surname: userData.surname,
        phone: userData.phone,
        default_outlet_id: userData.default_outlet_id,
      };
      localStorage.setItem('user', JSON.stringify(updatedUser));
      
      setSuccess('Данные успешно сохранены!');
      setIsEditing(false);
      setOriginalData(userData);
      
      setTimeout(() => {
        setSuccess('');
      }, 3000);
    } catch (err: any) {
      console.error('Ошибка сохранения данных:', err);
      setError(err.response?.data?.message || 'Ошибка при сохранении данных');
    } finally {
      setSaving(false);
    }
  };

  const handleCancel = (): void => {
    if (originalData) {
      setUserData(originalData);
    }
    setIsEditing(false);
    setError('');
  };

  const handleOverlayClick = (e: React.MouseEvent): void => {
    if (e.target === e.currentTarget) {
      if (isEditing) {
        if (window.confirm('У вас есть несохраненные изменения. Закрыть без сохранения?')) {
          handleCancel();
          onClose();
        }
      } else {
        onClose();
      }
    }
  };

  const getOutletName = (outletId: string): string => {
    if (!outletId) return 'Не выбрано';
    const outlet = getOutletById(outletId);
    return outlet ? `${outlet.address}` : `ID: ${outletId}`;
  };

  if (!isOpen) return null;

  return (
    <div className={styles.overlay} onClick={handleOverlayClick}>
      <div className={styles.modal}>
        <div className={styles.header}>
          <h2 className={styles.title}>Мой профиль</h2>
          
          <div className={styles.headerButtons}>
            {!isEditing && !loading && (
              <button
                type="button"
                onClick={() => setIsEditing(true)}
                className={styles.editButton}
                title="Редактировать"
              >
                <img className={styles.iconImg} src='/icons/edit.svg'/>
              </button>
            )}
            
            <button onClick={onClose} className={styles.closeButton}>
              ×
            </button>
          </div>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          {loading || outletsLoading ? (
            <div className={styles.loading}>
              <div className={styles.spinner}></div>
              <p>Загрузка данных...</p>
            </div>
          ) : (
            <>
              <div className={styles.formGroup}>
                <label htmlFor="email" className={styles.label}>
                  Email
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  value={userData.email}
                  onChange={handleInputChange}
                  className={styles.input}
                  disabled={!isEditing}
                />
              </div>

              <div className={styles.twoColumns}>
                <div className={styles.column}>
                  <label htmlFor="name" className={styles.label}>
                    Имя
                  </label>
                  <input
                    id="name"
                    name="name"
                    type="text"
                    value={userData.name}
                    onChange={handleInputChange}
                    placeholder="Имя"
                    className={styles.input}
                    disabled={!isEditing}
                  />
                </div>
                
                <div className={styles.column}>
                  <label htmlFor="surname" className={styles.label}>
                    Фамилия
                  </label>
                  <input
                    id="surname"
                    name="surname"
                    type="text"
                    value={userData.surname}
                    onChange={handleInputChange}
                    placeholder="Фамилия"
                    className={styles.input}
                    disabled={!isEditing}
                  />
                </div>
              </div>

              <div className={styles.formGroup}>
                <label htmlFor="phone" className={styles.label}>
                  Телефон
                </label>
                <input
                  id="phone"
                  name="phone"
                  type="tel"
                  value={userData.phone}
                  onChange={handleInputChange}
                  placeholder="+7 (999) 999-99-99"
                  className={styles.input}
                  disabled={!isEditing}
                />
              </div>

              <div className={styles.formGroup}>
                <label htmlFor="default_outlet_id" className={styles.label}>
                  Точка продаж по умолчанию
                </label>
                {isEditing ? (
                  <select
                    id="default_outlet_id"
                    name="default_outlet_id"
                    value={userData.default_outlet_id || ''}
                    onChange={handleInputChange}
                    className={styles.select}
                    disabled={!isEditing}
                  >
                    <OutletOptions outlets={outlets} />
                  </select>
                ) : (
                  <div className={styles.readOnlyField}>
                    {getOutletName(userData.default_outlet_id)}
                  </div>
                )}
              </div>

              {error && (
                <div className={styles.error}>
                  {error}
                </div>
              )}

              {success && (
                <div className={styles.success}>
                  {success}
                </div>
              )}

              <div className={styles.buttons}>
                {isEditing ? (
                  <>
                    <button
                      type="submit"
                      disabled={saving}
                      className={styles.saveButton}
                    >
                      {saving ? 'Сохранение...' : 'Сохранить'}
                    </button>
                    <button
                      type="button"
                      onClick={handleCancel}
                      className={styles.cancelButton}
                      disabled={saving}
                    >
                      Отмена
                    </button>
                  </>
                ) : null}
              </div>
            </>
          )}
        </form>
      </div>
    </div>
  );
};

export default UserProfileModal;