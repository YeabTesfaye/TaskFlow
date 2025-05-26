"use client";

import { TaskForm } from "@/components/task-form";
import { Header } from "@/components/layout/header";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { motion } from "framer-motion";

export default function NewTaskPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="flex-1">
        <div className="container mx-auto max-w-3xl px-4 py-8">
          <div className="mb-6 flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              asChild
              className="h-8 w-8"
            >
              <Link href="/">
                <ArrowLeft className="h-4 w-4" />
                <span className="sr-only">Go back</span>
              </Link>
            </Button>
            <motion.h1 
              className="text-2xl font-bold tracking-tight"
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3 }}
            >
              Create New Task
            </motion.h1>
          </div>
          
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
          >
            <div className="rounded-lg border bg-card p-6 shadow-sm">
              <TaskForm mode="create" />
            </div>
          </motion.div>
        </div>
      </main>
    </div>
  );
}