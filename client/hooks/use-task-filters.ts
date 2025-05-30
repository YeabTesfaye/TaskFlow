import { useState } from 'react';
import { FilterState } from '@/components/task-filters';

export const useTaskFilters = () => {
  const [filters, setFilters] = useState<FilterState>({
    statuses: ['Pending', 'In Progress', 'Completed'],
    priorities: ['Low', 'Medium', 'High', 'Urgent'],
    tagIds: [],
    sortBy: 'createdAt',
    sortDirection: 'desc',
    view: 'grid',
  });

  return { filters, setFilters };
};
