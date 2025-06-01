'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useToast } from '@/hooks/use-toast';
import { useAuth } from '@/hooks/use-auth';
import { user } from '@/api';

import { Button } from '@/components/ui/button';
import {
  Card,
  CardHeader,
  CardContent,
  CardTitle,
  CardDescription,
} from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';

const defaultSettings = {
  emailNotifications: false,
  pushNotifications: false,
  dailyDigest: false,
};

export default function SettingsPage() {
  const [settings, setSettings] = useState(defaultSettings);
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();
  const router = useRouter();
  const { isAuthenticated, loading } = useAuth();

  const loadSettings = useCallback(async () => {
    try {
      const data = await user.getProfile();
      const prefs = data.preferences ?? {};
      setSettings({
        emailNotifications: prefs.emailNotifications ?? false,
        pushNotifications: prefs.pushNotifications ?? false,
        dailyDigest: prefs.dailyDigest ?? false,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to load settings',
        variant: 'destructive',
        duration: 1000,
      });
    }
  }, [toast]);

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push('/login');
    }
  }, [loading, isAuthenticated, router]);

  useEffect(() => {
    loadSettings();
  }, [loadSettings]);

  const handleToggleChange =
    (key: keyof typeof defaultSettings) => (checked: boolean) => {
      setSettings((prev) => ({ ...prev, [key]: checked }));
    };

  const handleSaveSettings = async () => {
    setIsLoading(true);
    try {
      await user.updateProfile({ preferences: settings });
      toast({
        title: 'Success',
        description: 'Settings updated successfully',
        duration: 1000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to update settings',
        variant: 'destructive',
        duration: 1000,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-8">
      <Card>
        <CardHeader>
          <CardTitle>Notification Preferences</CardTitle>
          <CardDescription>
            Manage how you receive notifications
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {[
            {
              label: 'Email Notifications',
              description: 'Receive notifications about your tasks via email',
              key: 'emailNotifications',
            },
            {
              label: 'Push Notifications',
              description: 'Receive notifications in your browser',
              key: 'pushNotifications',
            },
            {
              label: 'Daily Digest',
              description: 'Receive a daily summary of your tasks',
              key: 'dailyDigest',
            },
          ].map((item) => (
            <div className="flex items-center justify-between" key={item.key}>
              <div className="space-y-0.5">
                <div className="font-medium">{item.label}</div>
                <div className="text-sm text-muted-foreground">
                  {item.description}
                </div>
              </div>
              <Switch
                checked={settings[item.key as keyof typeof settings]}
                onCheckedChange={handleToggleChange(
                  item.key as keyof typeof settings,
                )}
                disabled={isLoading}
              />
            </div>
          ))}

          <Button
            onClick={handleSaveSettings}
            disabled={isLoading}
            className="mt-6"
          >
            {isLoading ? 'Saving...' : 'Save Settings'}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
