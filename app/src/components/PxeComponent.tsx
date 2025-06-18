import { Card } from './ui/card';
import { Button } from './ui/button';
import { HardDrive } from 'lucide-react';

export function PxeComponent() {
  return (
    <div className="space-y-6">
      <Card className="p-6">
        <div className="flex items-center gap-4 mb-6">
          <div className="p-2 bg-primary/10 rounded-lg">
            <HardDrive className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h2 className="text-2xl font-semibold">PXE Configuration</h2>
            <p className="text-muted-foreground">
              Manage PXE boot assets and network boot configuration
            </p>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center gap-2 mb-4">
            <span className="text-sm font-medium">Status:</span>
            <span className="text-sm text-green-600">Active</span>
          </div>

          <div>
            <h4 className="font-medium mb-2">Boot Assets</h4>
            <p className="text-sm text-muted-foreground mb-4">
              Manage Talos Linux boot images and iPXE configurations for network booting.
            </p>
          </div>

          <div className="flex gap-2 justify-end mt-4">
            <Button variant="outline" onClick={() => console.log('View assets')}>
              View Assets
            </Button>
            <Button onClick={() => console.log('Download PXE assets')}>
              Download Assets
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
}