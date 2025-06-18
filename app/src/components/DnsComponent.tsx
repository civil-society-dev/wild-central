import { Card } from './ui/card';
import { Button } from './ui/button';
import { Globe, CheckCircle } from 'lucide-react';

export function DnsComponent() {
  return (
    <div className="space-y-6">
      <Card className="p-6">
        <div className="flex items-center gap-4 mb-6">
          <div className="p-2 bg-primary/10 rounded-lg">
            <Globe className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h2 className="text-2xl font-semibold">DNS Configuration</h2>
            <p className="text-muted-foreground">
              Manage DNS settings and domain resolution
            </p>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <CheckCircle className="h-5 w-5 text-green-500" />
            <span className="text-sm">Local resolution: Active</span>
          </div>
          
          <div className="mt-4">
            <h4 className="font-medium mb-2">DNS Status</h4>
            <p className="text-sm text-muted-foreground">
              DNS service is running and resolving domains correctly.
            </p>
          </div>

          <div className="flex gap-2 justify-end mt-4">
            <Button variant="outline" onClick={() => console.log('Test DNS')}>
              Test DNS
            </Button>
            <Button onClick={() => console.log('Configure DNS')}>
              Configure
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}