'use client';

import { useEffect, useState } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { PriorityChart } from '@/components/ui/priority-chart';
import { CategoryTable } from '@/components/category-table';
import { statistics } from '@/api';
import { Header } from '@/components/layout/header';

export default function DashboardPage() {
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const [categoryStats, setCategoryStats] = useState([]);
  const [categoryLoading, setCategoryLoading] = useState(false);

  useEffect(() => {
    const fetchStats = async () => {
      setLoading(true);
      try {
        const data = await statistics.getTaskStatistics();
        setStats(data);
      } catch (err) {
        console.error('Error fetching task statistics:', err);
      } finally {
        setLoading(false);
      }
    };

    const fetchCategoryStats = async () => {
      setCategoryLoading(true);
      try {
        const data = await statistics.getTaskStatisticsByCategory();
        setCategoryStats(data.categories || []);
      } catch (err) {
        console.error('Error fetching category stats:', err);
      } finally {
        setCategoryLoading(false);
      }
    };

    fetchStats();
    fetchCategoryStats();
  }, []);

  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <main className="flex-1">
        <div className="container mx-auto px-4 py-8 md:px-6 lg:px-8">
          <div className="mb-6 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
            <h1 className="text-3xl font-bold tracking-tight">
              Task Flow Dashboard
            </h1>
          </div>

          {/* Overview Cards */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
            <Card className="bg-card text-card-foreground border border-border hover:shadow-lg transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Total Tasks</CardTitle>
                <span className="text-muted-foreground dark:text-gray-400">📋</span>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {loading ? '...' : stats?.total_tasks}
                </div>
              </CardContent>
            </Card>

            <Card className="bg-card text-card-foreground border border-border hover:shadow-lg transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Completed</CardTitle>
                <span className="text-muted-foreground dark:text-green-400">
                  ✅
                </span>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {loading ? '...' : stats?.completed_tasks}
                </div>
              </CardContent>
            </Card>

            <Card className="bg-card text-card-foreground border border-border hover:shadow-lg transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">
                  Completion Rate
                </CardTitle>
                <span className="text-muted-foreground dark:text-blue-400">📈</span>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {loading ? '...' : `${stats?.completion_rate?.toFixed(1) ?? 0}%`}
                </div>
              </CardContent>
            </Card>

            <Card className="bg-card text-card-foreground border border-border hover:shadow-lg transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium">Overdue</CardTitle>
                <span className="text-muted-foreground dark:text-red-400">⚠️</span>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {loading ? '...' : stats?.overdue_tasks}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Priority Breakdown */}
          <Card className="mb-8">
            <CardHeader>
              <CardTitle className="text-xl">Tasks by Priority</CardTitle>
            </CardHeader>
            <CardContent className="h-[400px]">
              {loading ? (
                <div className="flex h-full items-center justify-center text-muted-foreground dark:text-gray-300">
                  <p>Loading priority data...</p>
                </div>
              ) : stats?.by_priority ? (
                <PriorityChart data={stats.by_priority} />
              ) : (
                <p className="text-muted-foreground dark:text-gray-300">
                  No priority data available
                </p>
              )}
            </CardContent>
          </Card>
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
