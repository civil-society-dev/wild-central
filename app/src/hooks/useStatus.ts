import { useQuery } from '@tanstack/react-query';
import { apiService } from '../services/api';
import { Status } from '../types';

export const useStatus = () => {
  return useQuery<Status>({
    queryKey: ['status'],
    queryFn: apiService.getStatus,
    refetchInterval: 30000, // Refetch every 30 seconds
  });
};