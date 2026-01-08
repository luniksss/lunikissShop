import { useState, useEffect } from 'react';
import api from '../api/api';
import { SalesOutlet } from '../types';

export const useOutlets = () => {
  const [outlets, setOutlets] = useState<SalesOutlet[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchOutlets = async () => {
      try {
        setLoading(true);
        const response = await api.getSalesOutlets();
        setOutlets(response.data);
      } catch (err) {
        console.error('Ошибка загрузки точек продаж:', err);
        setError('Не удалось загрузить точки продаж');
      } finally {
        setLoading(false);
      }
    };

    if (outlets.length === 0) {
      fetchOutlets();
    }
  }, []);

  const getOutletById = (id: string): SalesOutlet | undefined => {
    return outlets.find(outlet => outlet.id === id);
  };

  const getOutletAddress = (id: string): string => {
    const outlet = getOutletById(id);
    return outlet ? `${outlet.address}` : 'Неизвестная точка';
  };

  return {
    outlets,
    loading,
    error,
    getOutletById,
    getOutletAddress
  };
};