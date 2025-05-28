'use client';

import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { PlusCircle } from 'lucide-react';

import { Header } from '@/components/layout/header';
import { TaskGrid } from '@/components/task-grid';
import { TaskFilters, FilterState } from '@/components/task-filters';
import { Button } from '@/components/ui/button';
import { Task } from '@/types/task';
import { getTasks } from '@/lib/actions/task.action';
import { initialTags } from '@/lib/data';

export default function Home() {
  const searchParams = useSearchParams();
  const searchQuery = searchParams?.get('search') || '';

  const [tasks, setTasks] = useState<Task[]>([]);
  const [filteredTasks, setFilteredTasks] = useState<Task[]>([]);
  const [filters, setFilters] = useState<FilterState>({
    statuses: ['Pending', 'In Progress', 'Completed'],
    priorities: ['Low', 'Medium', 'High', 'Urgent'],
    tagIds: [],
    sortBy: 'createdAt',
    sortDirection: 'desc',
    view: 'grid',
  });

  useEffect(() => {
    const fetchTasks = async () => {
      const fetchedTasks = await getTasks();
      setTasks(fetchedTasks);
    };
    fetchTasks();
  }, []);

  useEffect(() => {
    let filtered = [...tasks];

    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (task) =>
          task.title.toLowerCase().includes(q) ||
          task.description.toLowerCase().includes(q),
      );
    }

    filtered = filtered.filter((task) =>
      filters.statuses.includes(task.status),
    );

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

    setFilteredTasks(filtered);
  }, [tasks, filters, searchQuery]);

  const handleFilterChange = (newFilters: FilterState) => {
    setFilters(newFilters);
  };

  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="flex-1">
        <div className="container mx-auto px-4 py-8 md:px-6 lg:px-8">
          <div className="mb-6 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
            <h1 className="text-3xl font-bold tracking-tight">
              {searchQuery ? `Search: "${searchQuery}"` : 'My Tasks'}
            </h1>
            <Button asChild size="sm" className="w-full sm:w-auto">
              <Link href="/new">
                <PlusCircle className="mr-2 h-4 w-4" />
                New Task
              </Link>
            </Button>
          </div>

          <TaskFilters tags={initialTags} onFilterChange={handleFilterChange} />

          <TaskGrid tasks={filteredTasks} />

          {filteredTasks.length === 0 && tasks.length > 0 && (
            <div className="mt-4 text-center text-sm text-muted-foreground">
              No tasks match your current filters.
            </div>
          )}

          {tasks.length === 0 && (
            <div className="mt-8 flex flex-col items-center justify-center rounded-lg border-2 border-dashed p-12 text-center">
              <h3 className="mb-2 text-xl font-medium">No tasks yet</h3>
              <p className="mb-6 text-muted-foreground">
                Get started by creating your first task
              </p>
              <Button asChild>
                <Link href="/new">
                  <PlusCircle className="mr-2 h-5 w-5" />
                  Create Your First Task
                </Link>
              </Button>
            </div>
          )}
        </div>
      </main>
      <footer className="border-t py-6">
        <div className="container mx-auto flex items-center justify-center px-4 text-sm text-muted-foreground">
          <p>TaskFlow &copy; {new Date().getFullYear()}</p>
        </div>
      </footer>
    </div>
  );
}
