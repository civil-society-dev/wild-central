import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { ConfigEditor } from './ConfigEditor';

export function Advanced() {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Advanced Configuration</CardTitle>
          <CardDescription>
            Advanced settings and system configuration options
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div>
            <h3 className="text-sm font-medium mb-2">Configuration Management</h3>
            <p className="text-sm text-muted-foreground mb-4">
              Edit the raw YAML configuration file directly. This provides full access to all configuration options.
            </p>
            <ConfigEditor />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}