"use client";

import { Task, Priority, TaskStatus, Tag } from "@/types/task";
import { create } from "zustand";
import { persist } from "zustand/middleware";

// Demo tags with colors
export const initialTags: Tag[] = [
  { id: "1", name: "Work", color: "#3B82F6" },
  { id: "2", name: "Personal", color: "#10B981" },
  { id: "3", name: "Urgent", color: "#EF4444" },
  { id: "4", name: "Learning", color: "#8B5CF6" },
  { id: "5", name: "Health", color: "#EC4899" },
];

// Initial tasks for demo
const initialTasks: Task[] = [
  {
    id: "1",
    title: "Complete project proposal",
    description: "Draft the initial project proposal for client review",
    priority: "high",
    status: "in-progress",
    dueDate: new Date(Date.now() + 86400000 * 2), // 2 days from now
    createdAt: new Date(),
    tags: ["1", "3"],
  },
  {
    id: "2",
    title: "Schedule dentist appointment",
    description: "Call dentist office to schedule annual checkup",
    priority: "medium",
    status: "pending",
    dueDate: new Date(Date.now() + 86400000 * 7), // 7 days from now
    createdAt: new Date(),
    tags: ["2", "5"],
  },
  {
    id: "3",
    title: "Learn Next.js 14",
    description: "Complete the tutorial on the official website",
    priority: "low",
    status: "pending",
    dueDate: new Date(Date.now() + 86400000 * 14), // 14 days from now
    createdAt: new Date(),
    tags: ["4"],
  },
];

interface TaskStore {
  tasks: Task[];
  tags: Tag[];
  addTask: (task: Omit<Task, "id" | "createdAt">) => void;
  updateTask: (id: string, task: Partial<Task>) => void;
  deleteTask: (id: string) => void;
  addTag: (tag: Omit<Tag, "id">) => void;
  updateTag: (id: string, tag: Partial<Tag>) => void;
  deleteTag: (id: string) => void;
}

export const useTaskStore = create<TaskStore>()(
  persist(
    (set) => ({
      tasks: initialTasks,
      tags: initialTags,
      addTask: (task) =>
        set((state) => ({
          tasks: [
            ...state.tasks,
            {
              ...task,
              id: crypto.randomUUID(),
              createdAt: new Date(),
            },
          ],
        })),
      updateTask: (id, updatedTask) =>
        set((state) => ({
          tasks: state.tasks.map((task) =>
            task.id === id ? { ...task, ...updatedTask } : task
          ),
        })),
      deleteTask: (id) =>
        set((state) => ({
          tasks: state.tasks.filter((task) => task.id !== id),
        })),
      addTag: (tag) =>
        set((state) => ({
          tags: [...state.tags, { ...tag, id: crypto.randomUUID() }],
        })),
      updateTag: (id, updatedTag) =>
        set((state) => ({
          tags: state.tags.map((tag) =>
            tag.id === id ? { ...tag, ...updatedTag } : tag
          ),
        })),
      deleteTag: (id) =>
        set((state) => ({
          tags: state.tags.filter((tag) => tag.id !== id),
        })),
    }),
    {
      name: "task-storage",
    }
  )
);