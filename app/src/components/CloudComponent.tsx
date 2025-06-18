import { Card } from './ui/card';
import { Button } from './ui/button';
import { Cloud, HelpCircle } from 'lucide-react';
import { Input, Label } from './ui';

export function CloudComponent() {
  return (
    <div className="space-y-6">
      <Card className="p-6">
        <div className="flex items-center gap-4 mb-6">
          <div className="p-2 bg-primary/10 rounded-lg">
            <Cloud className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h2 className="text-2xl font-semibold">Cloud Configuration</h2>
            <p className="text-muted-foreground">
              Configure top-level cloud settings and domains
            </p>
          </div>
        </div>

        <div className="space-y-4">
          <div>
            <Label htmlFor="upstream">Upstream</Label>
            <div className="flex w-full items-center mt-1">
              <Input id="upstream" value="https://mywildcloud.org"/>
              <Button variant="ghost">
                <HelpCircle/>
              </Button>
            </div>
          </div>
          <div>
            <Label htmlFor="domain">Domain</Label>
            <div className="flex w-full items-center mt-1">
              <Input id="domain" value="cloud.payne.io"/>
              <Button variant="ghost">
                <HelpCircle/>
              </Button>
            </div>
          </div>
          <div>
            <Label htmlFor="internalDomain">Internal Domain</Label>
            <div className="flex w-full items-center mt-1">
              <Input id="internalDomain" value="internal.cloud.payne.io"/>
              <Button variant="ghost">
                <HelpCircle/>
              </Button>
            </div>
          </div>
        </div>

        <div className="flex gap-2 justify-end mt-6">
          <Button variant="outline" onClick={() => console.log('Reset settings')}>
            Reset
          </Button>
          <Button onClick={() => console.log('Save settings')}>
            Save Configuration
          </Button>
        </div>
      </Card>
    </div>
  );
}