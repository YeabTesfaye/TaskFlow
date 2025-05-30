export default function Loading() {
  return (
    <div className="fixed inset-0 flex items-center justify-center bg-white bg-opacity-75 dark:bg-gray-900 dark:bg-opacity-75 z-50">
      <div className="flex flex-col items-center">
        <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-indigo-500 dark:border-indigo-400 mb-4"></div>
        <p className="text-lg font-medium text-gray-700 dark:text-gray-200">
          Loading...
        </p>
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
          Please wait while we prepare your content
        </p>
      </div>
    </div>
  );
}
