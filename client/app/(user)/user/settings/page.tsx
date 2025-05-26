'use client';

import { useState, useEffect } from 'react';
import { useToast } from '@/hooks/use-toast';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardContent, CardTitle, CardDescription } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { user } from '@/api';
// import { user } from '@/lib/api';

export default function SettingsPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [settings, setSettings] = useState({
    emailNotifications: false,
    pushNotifications: false,
    dailyDigest: false,
  });
  const { toast } = useToast();

  useEffect(() => {
    loadSettings();
  }, []);

  const loadSettings = async () => {
    try {
      const data = await user.getProfile();
      setSettings({
        emailNotifications: data.preferences?.emailNotifications ?? false,
        pushNotifications: data.preferences?.pushNotifications ?? false,
        dailyDigest: data.preferences?.dailyDigest ?? false,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to load settings',
        variant: 'destructive',
      });
    }
  };

  const handleSaveSettings = async () => {
    setIsLoading(true);
    try {
      await user.updateProfile({ preferences: settings });
      toast({
        title: 'Success',
        description: 'Settings updated successfully',
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.error || 'Failed to update settings',
        variant: 'destructive',
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
          <CardDescription>Manage how you receive notifications</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <div className="font-medium">Email Notifications</div>
              <div className="text-sm text-muted-foreground">
                Receive notifications about your tasks via email
              </div>
            </div>
            <Switch
              checked={settings.emailNotifications}
              onCheckedChange={(checked) =>
                setSettings({ ...settings, emailNotifications: checked })
              }
              disabled={isLoading}
            />
          </div>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <div className="font-medium">Push Notifications</div>
              <div className="text-sm text-muted-foreground">
                Receive notifications in your browser
              </div>
            </div>
            <Switch
              checked={settings.pushNotifications}
              onCheckedChange={(checked) =>
                setSettings({ ...settings, pushNotifications: checked })
              }
              disabled={isLoading}
            />
          </div>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <div className="font-medium">Daily Digest</div>
              <div className="text-sm text-muted-foreground">
                Receive a daily summary of your tasks
              </div>
            </div>
            <Switch
              checked={settings.dailyDigest}
              onCheckedChange={(checked) =>
                setSettings({ ...settings, dailyDigest: checked })
              }
              disabled={isLoading}
            />
          </div>

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