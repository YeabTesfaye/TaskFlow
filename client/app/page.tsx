'use client';

import { useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { PlusCircle } from 'lucide-react';

import { Header } from '@/components/layout/header';
import { TaskFilters } from '@/components/task-filters';
import { TaskGrid } from '@/components/task-grid';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/hooks/use-auth';
import { useTasks } from '@/hooks/use-tasks';
import { useTaskFilters } from '@/hooks/use-task-filters';
import { filterTasks } from '@/lib/utils/filterTasks';
import { Task } from '@/types';
import { useEffect, useState } from 'react';
import Loading from '@/components/ui/loading';
import { useTags } from '@/hooks/useTags';

export default function Home() {
  const searchParams = useSearchParams();
  const searchQuery = searchParams?.get('search') || '';
  const { isAuthenticated } = useAuth();
  const { filters, setFilters } = useTaskFilters();
  const { tasks, loading, refetch } = useTasks();
  const [filteredTasks, setFilteredTasks] = useState<Task[]>([]);
  const { tagList } = useTags();


  useEffect(() => {
    const result = filterTasks(tasks, filters, searchQuery);
    setFilteredTasks(result);
  }, [tasks, filters, searchQuery]);

  const handleStatusChange = (taskId: string, newStatus: Task['status']) => {
    refetch();
  };

  const handleDelete = (taskId: string) => {
    refetch();
  };

  if (loading) return <Loading />;
  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="flex-1">
        <div className="container mx-auto px-4 py-8 md:px-6 lg:px-8">
          <div className="mb-6 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
            <h1 className="text-3xl font-bold tracking-tight">
              {searchQuery ? `Search: "${searchQuery}"` : 'My Tasks'}
            </h1>
            <Button
              asChild
              size="sm"
              className="w-full sm:w-auto"
              disabled={!isAuthenticated}
            >
              <Link href="/new">
                <PlusCircle className="mr-2 h-4 w-4" />
                New Task
              </Link>
            </Button>
          </div>

          <TaskFilters tags={tagList} onFilterChange={setFilters} />

          {loading ? (
            <p className="text-muted-foreground">Loading tasks...</p>
          ) : tasks.length === 0 ? (
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
          ) : filteredTasks.length === 0 ? (
            <div className="mt-4 text-center text-sm text-muted-foreground">
              No tasks match your current filters.
            </div>
          ) : (
            <TaskGrid
              tasks={filteredTasks}
              onDelete={handleDelete}
              onStatusChange={handleStatusChange}
            />
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
