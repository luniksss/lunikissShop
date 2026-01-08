import React, { useState, useEffect } from 'react';
import api from '../../api/api';
import { User } from '../../types';
import styles from './AdminEmployees.module.css'

const AdminEmployees: React.FC = () => {
  const [employees, setEmployees] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [roleChanging, setRoleChanging] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<string | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  useEffect(() => {
    fetchAllEmployees();
  }, []);

  const fetchAllEmployees = async () => {
    try {
      setLoading(true);
      setError(null);
      setDeleteError(null);
      const response = await api.getAllUsers();

      const data = response.data === null ? [] : response.data;
    
      if (Array.isArray(data)) {
        setEmployees(data);
      } else {
        console.error('Ожидался массив пользователей, но получено:', data);
        setEmployees([]);
      }
    } catch (err: any) {
      console.error('Ошибка загрузки пользователей:', err);
      
      if (err.response?.status === 404) {
        setError('Эндпоинт для получения пользователей не найден. Убедитесь, что бэкенд поддерживает этот запрос.');
      } else if (err.response?.status === 401) {
        setError('Необходима авторизация. Пожалуйста, войдите снова.');
      } else if (err.response?.status === 403) {
        setError('У вас нет прав для просмотра пользователей.');
      } else {
        setError('Не удалось загрузить пользователей. Пожалуйста, попробуйте позже.');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleRoleChange = async (userId: string, newRole: string) => {
    try {
      setRoleChanging(userId);
      await api.updateUserRole(userId, newRole);
      
      setEmployees(prevEmployees =>
        prevEmployees.map(employee =>
          employee.id === userId ? { ...employee, role: newRole } : employee
        )
      );
      
      alert('Роль пользователя обновлена!');
    } catch (err: any) {
      console.error('Ошибка обновления роли:', err);
      alert('Не удалось обновить роль пользователя');
    } finally {
      setRoleChanging(null);
    }
  };

  const handleDeleteUser = async (userId: string, email: string) => {
    if (!window.confirm(`Вы уверены, что хотите удалить пользователя ${email}? Это действие нельзя отменить.`)) {
      return;
    }

    setDeleting(userId);
    setDeleteError(null);

    try {
      await api.deleteUser(userId);
      
      setEmployees(prevEmployees => 
        prevEmployees.filter(employee => employee.id !== userId)
      );
      
      alert(`Пользователь ${email} успешно удален!`);
    } catch (err: any) {
      console.error('Ошибка удаления пользователя:', err);
      setDeleteError('Не удалось удалить пользователя. Пожалуйста, попробуйте позже.');
    } finally {
      setDeleting(null);
    }
  };

  const getRoleOptions = () => {
    return [
      { value: 'user', label: 'Пользователь' },
      { value: 'seller', label: 'Продавец' },
      { value: 'admin', label: 'Администратор' },
    ];
  };

  const getRoleClass = (role: string) => {
    switch(role) {
      case 'admin': return styles.roleAdmin;
      case 'seller': return styles.roleSeller;
      case 'customer': return styles.roleCustomer;
      default: return styles.roleCustomer;
    }
  };

  const getRoleText = (role: string): string => {
    const roleMap: Record<string, string> = {
      'admin': 'Администратор',
      'seller': 'Продавец',
      'user': 'Пользователь',
    };
    return roleMap[role.toLowerCase()] || role;
  };

  const getFullName = (user: User): string => {
    return `${user.name || ''} ${user.surname || ''}`.trim() || 'Не указано';
  };

  const getCurrentUser = () => {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  };

  const canDeleteUser = (userId: string): boolean => {
    const currentUser = getCurrentUser();
    return currentUser && currentUser.id !== userId;
  };

  const canEditRole = (userId: string): boolean => {
    const currentUser = getCurrentUser();
    return currentUser && currentUser.id !== userId;
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <p>Загрузка пользователей...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.errorText}>{error}</p>
        <button onClick={fetchAllEmployees} className={styles.retryButton}>
          Попробовать снова
        </button>
      </div>
    );
  }

  const displayEmployees = Array.isArray(employees) ? employees : [];

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h2 className={styles.title}>Управление сотрудниками</h2>
          <p className={styles.subtitle}>Всего пользователей: {displayEmployees.length}</p>
        </div>
        <div className={styles.headerButtons}>
          <button onClick={fetchAllEmployees} className={styles.refreshButton}>
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

      <div className={styles.tableContainer}>
        {displayEmployees.length === 0 ? (
          <div className={styles.emptyState}>
            <p>Пользователи не найдены</p>
          </div>
        ) : (
          <table className={styles.table}>
            <thead>
              <tr className={styles.tableHeader}>
                <th className={styles.th}>ID</th>
                <th className={styles.th}>Имя</th>
                <th className={styles.th}>Email</th>
                <th className={styles.th}>Телефон</th>
                <th className={styles.th}>Роль</th>
                <th className={styles.th}>Точка</th>
                <th className={styles.th}>Действия</th>
              </tr>
            </thead>
            <tbody>
              {displayEmployees.map((employee) => {
                const isCurrentUser = getCurrentUser()?.id === employee.id;
                
                return (
                  <tr key={employee.id} className={styles.tableRow}>
                    <td className={styles.td}>{employee.id}</td>
                    <td className={styles.td}>
                      {getFullName(employee)}
                      {isCurrentUser && (
                        <span className={styles.currentUserBadge}> (Вы)</span>
                      )}
                    </td>
                    <td className={styles.td}>{employee.email}</td>
                    <td className={styles.td}>{employee.phone || 'Не указан'}</td>
                    <td className={styles.td}>
                      <span className={`${styles.roleBadge} ${getRoleClass(employee.role)}`}>
                        {getRoleText(employee.role)}
                      </span>
                    </td>
                    <td className={styles.td}>{employee.default_outlet_id || 'Не указана'}</td>
                    <td className={styles.td}>
                      <div className={styles.actions}>
                        <select
                          value={employee.role ? employee.role.toLowerCase() : 'user'}
                          onChange={(e) => handleRoleChange(employee.id, e.target.value)}
                          disabled={roleChanging === employee.id || !canEditRole(employee.id)}
                          className={`${styles.roleSelect} ${!canEditRole(employee.id) ? styles.disabledSelect : ''}`}
                          title={!canEditRole(employee.id) ? "Нельзя изменить свою роль" : ""}
                        >
                          {getRoleOptions().map(option => (
                            <option key={option.value} value={option.value}>
                              {option.label}
                            </option>
                          ))}
                        </select>
                        
                        {canDeleteUser(employee.id) && (
                          <button
                            onClick={() => handleDeleteUser(employee.id, employee.email)}
                            disabled={deleting === employee.id}
                            className={deleting === employee.id ? styles.deleteButtonDisabled : styles.deleteButton}
                            title="Удалить пользователя"
                          >
                            {deleting === employee.id ? 'Удаление...' : <img src='/icons/delete.png' />}
                          </button>
                        )}
                        
                        {(roleChanging === employee.id || deleting === employee.id) && (
                          <span className={styles.loadingText}>...</span>
                        )}
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default AdminEmployees;