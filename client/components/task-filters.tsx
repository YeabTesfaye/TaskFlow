"use client";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { TaskStatus, Priority } from "@/types/task";
import { useTaskStore } from "@/lib/data";
import { useState } from "react";
import { Filter, SortAsc, SortDesc, CalendarDays } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { motion } from "framer-motion";

interface TaskFiltersProps {
  onFilterChange: (filters: FilterState) => void;
}

export interface FilterState {
  statuses: TaskStatus[];
  priorities: Priority[];
  tagIds: string[];
  sortBy: "dueDate" | "priority" | "createdAt";
  sortDirection: "asc" | "desc";
  view: "grid" | "list";
}

const initialFilters: FilterState = {
  statuses: ["pending", "in-progress", "completed"],
  priorities: ["low", "medium", "high", "urgent"],
  tagIds: [],
  sortBy: "createdAt",
  sortDirection: "desc",
  view: "grid",
};

export function TaskFilters({ onFilterChange }: TaskFiltersProps) {
  const [filters, setFilters] = useState<FilterState>(initialFilters);
  const { tags } = useTaskStore();

  const updateFilters = (newFilters: Partial<FilterState>) => {
    const updatedFilters = { ...filters, ...newFilters };
    setFilters(updatedFilters);
    onFilterChange(updatedFilters);
  };

  const toggleStatus = (status: TaskStatus) => {
    const statuses = filters.statuses.includes(status)
      ? filters.statuses.filter((s) => s !== status)
      : [...filters.statuses, status];
    updateFilters({ statuses });
  };

  const togglePriority = (priority: Priority) => {
    const priorities = filters.priorities.includes(priority)
      ? filters.priorities.filter((p) => p !== priority)
      : [...filters.priorities, priority];
    updateFilters({ priorities });
  };

  const toggleTag = (tagId: string) => {
    const tagIds = filters.tagIds.includes(tagId)
      ? filters.tagIds.filter((id) => id !== tagId)
      : [...filters.tagIds, tagId];
    updateFilters({ tagIds });
  };

  const sortOptions = [
    { value: "createdAt", label: "Creation Date" },
    { value: "dueDate", label: "Due Date" },
    { value: "priority", label: "Priority" },
  ] as const;

  return (
    <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex flex-wrap gap-2">
        <motion.div 
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-1">
                <Filter className="h-4 w-4" />
                <span>Status</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48">
              <DropdownMenuLabel>Filter by Status</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuCheckboxItem
                checked={filters.statuses.includes("pending")}
                onCheckedChange={() => toggleStatus("pending")}
              >
                Pending
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.statuses.includes("in-progress")}
                onCheckedChange={() => toggleStatus("in-progress")}
              >
                In Progress
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.statuses.includes("completed")}
                onCheckedChange={() => toggleStatus("completed")}
              >
                Completed
              </DropdownMenuCheckboxItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </motion.div>
        
        <motion.div 
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.1 }}
        >
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-1">
                <Filter className="h-4 w-4" />
                <span>Priority</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48">
              <DropdownMenuLabel>Filter by Priority</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuCheckboxItem
                checked={filters.priorities.includes("low")}
                onCheckedChange={() => togglePriority("low")}
              >
                Low
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.priorities.includes("medium")}
                onCheckedChange={() => togglePriority("medium")}
              >
                Medium
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.priorities.includes("high")}
                onCheckedChange={() => togglePriority("high")}
              >
                High
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.priorities.includes("urgent")}
                onCheckedChange={() => togglePriority("urgent")}
              >
                Urgent
              </DropdownMenuCheckboxItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </motion.div>
        
        <motion.div 
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.2 }}
        >
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-1">
                <Filter className="h-4 w-4" />
                <span>Tags</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48">
              <DropdownMenuLabel>Filter by Tags</DropdownMenuLabel>
              <DropdownMenuSeparator />
              {tags.length > 0 ? (
                tags.map((tag) => (
                  <DropdownMenuCheckboxItem
                    key={tag.id}
                    checked={filters.tagIds.includes(tag.id)}
                    onCheckedChange={() => toggleTag(tag.id)}
                  >
                    <span className="inline-flex items-center gap-1.5">
                      <span
                        className="h-2 w-2 rounded-full"
                        style={{ backgroundColor: tag.color }}
                      ></span>
                      {tag.name}
                    </span>
                  </DropdownMenuCheckboxItem>
                ))
              ) : (
                <div className="px-2 py-1.5 text-sm text-muted-foreground">
                  No tags available
                </div>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </motion.div>
      </div>
      
      <div className="flex items-center gap-2">
        <motion.div 
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.3 }}
        >
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-1">
                {filters.sortDirection === "asc" ? (
                  <SortAsc className="h-4 w-4" />
                ) : (
                  <SortDesc className="h-4 w-4" />
                )}
                <span>Sort</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              <DropdownMenuLabel>Sort by</DropdownMenuLabel>
              <DropdownMenuSeparator />
              {sortOptions.map((option) => (
                <DropdownMenuCheckboxItem
                  key={option.value}
                  checked={filters.sortBy === option.value}
                  onCheckedChange={() =>
                    updateFilters({ sortBy: option.value })
                  }
                >
                  {option.label}
                </DropdownMenuCheckboxItem>
              ))}
              <DropdownMenuSeparator />
              <DropdownMenuCheckboxItem
                checked={filters.sortDirection === "asc"}
                onCheckedChange={() =>
                  updateFilters({
                    sortDirection: filters.sortDirection === "asc" ? "desc" : "asc",
                  })
                }
              >
                Ascending Order
              </DropdownMenuCheckboxItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </motion.div>
      </div>
    </div>
  );
}