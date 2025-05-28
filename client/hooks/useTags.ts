import { initialTags } from '@/lib/data';
import { useEffect, useState } from 'react';

export interface Tag {
  id: string;
  name: string;
  color: string;
}

export function useTags() {
  const [tagList, setTagList] = useState<Tag[]>(initialTags);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTags = async () => {
      setLoading(true);
      try {
        // const data = await tags.getAll();
        setTagList(tagList);
      } catch (error) {
        setError('Failed to load tags');
      } finally {
        setLoading(false);
      }
    };
    fetchTags();
  }, []);
  return { tagList, loading, error };
}
