
export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex  flex-col">
      <main className="flex-1 flex items-center justify-center">
        {children}
      </main>
    </div>
  );
}