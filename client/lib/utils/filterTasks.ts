import { Task } from '@/types';
import { FilterState } from '@/components/task-filters';

export function filterTasks(
  tasks: Task[],
  filters: FilterState,
  searchQuery: string,
): Task[] {
  let filtered = [...tasks];

  if (searchQuery) {
    const q = searchQuery.toLowerCase();
    filtered = filtered.filter(
      (task) =>
        task.title.toLowerCase().includes(q) ||
        task.description.toLowerCase().includes(q),
    );
  }

  filtered = filtered.filter((task) => filters.statuses.includes(task.status));
  filtered = filtered.filter((task) =>
    filters.priorities.includes(task.priority),
  );

  if (filters.tagIds.length > 0) {
    filtered = filtered.filter((task) =>
      task.tags.some((tagId) => filters.tagIds.includes(tagId)),
    );
  }

  filtered.sort((a, b) => {
    if (filters.sortBy === 'dueDate') {
      const aDate = a.dueDate ? new Date(a.dueDate).getTime() : 0;
      const bDate = b.dueDate ? new Date(b.dueDate).getTime() : 0;
      return filters.sortDirection === 'asc' ? aDate - bDate : bDate - aDate;
    }

    if (filters.sortBy === 'priority') {
      const priorityOrder = { Low: 0, Medium: 1, High: 2, Urgent: 3 };
      return filters.sortDirection === 'asc'
        ? priorityOrder[a.priority] - priorityOrder[b.priority]
        : priorityOrder[b.priority] - priorityOrder[a.priority];
    }

    const aCreated = new Date(a.createdAt).getTime();
    const bCreated = new Date(b.createdAt).getTime();
    return filters.sortDirection === 'asc'
      ? aCreated - bCreated
      : bCreated - aCreated;
  });

  return filtered;
}
