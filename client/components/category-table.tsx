import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface CategoryStat {
  category_id: string;
  category_name: string;
  total_tasks: number;
  completed_tasks: number;
  completion_rate: number;
  overdue_tasks: number;
}

export function CategoryTable({ data }: { data: CategoryStat[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Category</TableHead>
          <TableHead>Total</TableHead>
          <TableHead>Completed</TableHead>
          <TableHead>Completion Rate</TableHead>
          <TableHead>Overdue</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.map((stat) => (
          <TableRow key={stat.category_id}>
            <TableCell>{stat.category_name}</TableCell>
            <TableCell>{stat.total_tasks}</TableCell>
            <TableCell>{stat.completed_tasks}</TableCell>
            <TableCell>{stat.completion_rate.toFixed(1)}%</TableCell>
            <TableCell>{stat.overdue_tasks}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}