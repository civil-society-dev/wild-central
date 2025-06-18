import { Card } from './ui/card';
import { Button } from './ui/button';
import { Wifi, HelpCircle } from 'lucide-react';
import { Input, Label } from './ui';

export function DhcpComponent() {
  return (
    <div className="space-y-6">
      <Card className="p-6">
        <div className="flex items-center gap-4 mb-6">
          <div className="p-2 bg-primary/10 rounded-lg">
            <Wifi className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h2 className="text-2xl font-semibold">DHCP Configuration</h2>
            <p className="text-muted-foreground">
              Manage DHCP settings and IP address allocation
            </p>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center gap-2 mb-4">
            <span className="text-sm font-medium">Status:</span>
            <span className="text-sm text-green-600">Active</span>
          </div>

          <div>
            <Label htmlFor="dhcpRange">IP Range</Label>
            <div className="flex w-full items-center mt-1">
              <Input id="dhcpRange" value="192.168.8.100,192.168.8.239"/>
              <Button variant="ghost">
                <HelpCircle/>
              </Button>
            </div>
          </div>

          <div className="flex gap-2 justify-end mt-4">
            <Button variant="outline" onClick={() => console.log('View DHCP clients')}>
              View Clients
            </Button>
            <Button onClick={() => console.log('Configure DHCP')}>
              Configure
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}