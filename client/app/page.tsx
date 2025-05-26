"use client";

import { useTaskStore } from "@/lib/data";
import { TaskGrid } from "@/components/task-grid";
import { Header } from "@/components/layout/header";
import { TaskFilters, FilterState } from "@/components/task-filters";
import { useState, useEffect } from "react";
import { Task } from "@/types/task";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { PlusCircle } from "lucide-react";
import Link from "next/link";

export default function Home() {
  const { tasks, tags } = useTaskStore();
  const searchParams = useSearchParams();
  const searchQuery = searchParams?.get("search") || "";

  const [filteredTasks, setFilteredTasks] = useState<Task[]>(tasks);
  const [filters, setFilters] = useState<FilterState>({
    statuses: ["pending", "in-progress", "completed"],
    priorities: ["low", "medium", "high", "urgent"],
    tagIds: [],
    sortBy: "createdAt",
    sortDirection: "desc",
    view: "grid",
  });

  useEffect(() => {
    let filtered = [...tasks];

    // Apply search filter if present
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (task) =>
          task.title.toLowerCase().includes(query) ||
          task.description.toLowerCase().includes(query)
      );
    }

    // Apply status filter
    filtered = filtered.filter((task) =>
      filters.statuses.includes(task.status)
    );

    // Apply priority filter
    filtered = filtered.filter((task) =>
      filters.priorities.includes(task.priority)
    );

    // Apply tag filter if any tags are selected
    if (filters.tagIds.length > 0) {
      filtered = filtered.filter((task) =>
        task.tags.some((tagId) => filters.tagIds.includes(tagId))
      );
    }

    // Apply sorting
    filtered.sort((a, b) => {
      if (filters.sortBy === "dueDate") {
        if (!a.dueDate) return 1;
        if (!b.dueDate) return -1;
        return filters.sortDirection === "asc"
          ? new Date(a.dueDate).getTime() - new Date(b.dueDate).getTime()
          : new Date(b.dueDate).getTime() - new Date(a.dueDate).getTime();
      }
      
      if (filters.sortBy === "priority") {
        const priorityOrder = { low: 0, medium: 1, high: 2, urgent: 3 };
        return filters.sortDirection === "asc"
          ? priorityOrder[a.priority] - priorityOrder[b.priority]
          : priorityOrder[b.priority] - priorityOrder[a.priority];
      }
      
      // Default sort by createdAt
      return filters.sortDirection === "asc"
        ? new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
        : new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
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
              {searchQuery ? `Search: "${searchQuery}"` : "My Tasks"}
            </h1>
            <Button asChild size="sm" className="w-full sm:w-auto">
              <Link href="/new">
                <PlusCircle className="mr-2 h-4 w-4" />
                New Task
              </Link>
            </Button>
          </div>

          <TaskFilters onFilterChange={handleFilterChange} />
          
          <TaskGrid tasks={filteredTasks} />
          
          {filteredTasks.length === 0 && tasks.length > 0 && (
            <div className="mt-4 text-center text-sm text-muted-foreground">
              No tasks match your current filters. Try adjusting your search or filter criteria.
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