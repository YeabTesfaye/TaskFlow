'use client';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { useAuth } from '@/hooks/use-auth';
import { Header } from '@/components/layout/header';
import { PriorityChart } from '@/components/ui/priority-chart';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { useTaskStats } from '@/hooks/use-task-stats';
import Loading from '@/components/ui/loading';

const DashboardPage = () => {
  const { isAuthenticated, loading: authLoading } = useAuth();
  const router = useRouter();

  const { stats, loading: statsLoading } = useTaskStats(
    !authLoading && isAuthenticated,
  );

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [authLoading, isAuthenticated, router]);

  if (authLoading) return <Loading />;

  if (authLoading || (!authLoading && !isAuthenticated)) {
    return (
      <div className="flex items-center justify-center h-screen text-muted-foreground">
        Checking authentication...
      </div>
    );
  }

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

          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
            <StatCard
              title="Total Tasks"
              value={stats?.total_tasks}
              loading={statsLoading}
              icon="📋"
            />
            <StatCard
              title="Completed"
              value={stats?.completed_tasks}
              loading={statsLoading}
              icon="✅"
            />
            <StatCard
              title="Completion Rate"
              value={`${stats?.completion_rate?.toFixed(1) ?? 0}%`}
              loading={statsLoading}
              icon="📈"
            />
            <StatCard
              title="Overdue"
              value={stats?.overdue_tasks}
              loading={statsLoading}
              icon="⚠️"
            />
          </div>

          <Card className="mb-8">
            <CardHeader>
              <CardTitle className="text-xl">Tasks by Priority</CardTitle>
            </CardHeader>
            <CardContent className="h-[400px]">
              {statsLoading ? (
                <LoadingState message="Loading priority data..." />
              ) : stats?.by_priority ? (
                <PriorityChart data={stats.by_priority} />
              ) : (
                <LoadingState message="No priority data available" />
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
};

const StatCard = ({
  title,
  value,
  loading,
  icon,
}: {
  title: string;
  value: string | number | undefined;
  loading: boolean;
  icon: string;
}) => (
  <Card className="bg-card text-card-foreground border border-border hover:shadow-lg transition-shadow">
    <CardHeader className="flex flex-row items-center justify-between pb-2">
      <CardTitle className="text-sm font-medium">{title}</CardTitle>
      <span className="text-muted-foreground dark:text-gray-400">{icon}</span>
    </CardHeader>
    <CardContent>
      <div className="text-2xl font-bold">{loading ? '...' : value}</div>
    </CardContent>
  </Card>
);

const LoadingState = ({ message }: { message: string }) => (
  <p className="text-muted-foreground dark:text-gray-300 text-center w-full">
    {message}
  </p>
);

export default DashboardPage;
