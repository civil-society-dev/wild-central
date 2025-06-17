import { useState } from 'react';
import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Cloud, Server, Network, Settings, CheckCircle, AlertCircle, Clock, HelpCircle, Check } from 'lucide-react';
import { Input, Label } from '../ui';

interface SetupPhaseProps {
  onComplete?: () => void;
}

interface SetupStep {
  id: string;
  title: string;
  description: string;
  status: 'pending' | 'in-progress' | 'completed' | 'error';
  icon: React.ComponentType<{ className?: string }>;
}

export function SetupPhase({ onComplete }: SetupPhaseProps) {
  const [setupSteps, setSetupSteps] = useState<SetupStep[]>([
    {
      id: 'server-config',
      title: 'Server Configuration',
      description: 'Configure wild-cloud central server settings',
      status: 'pending',
      icon: Server,
    },
    {
      id: 'dnsmasq-setup',
      title: 'DNS & DHCP Setup',
      description: 'Configure dnsmasq for DNS and DHCP services',
      status: 'pending',
      icon: Network,
    },
    {
      id: 'pxe-assets',
      title: 'PXE Boot Assets',
      description: 'Download Talos Linux boot images and iPXE configurations',
      status: 'pending',
      icon: Cloud,
    },
    {
      id: 'service-validation',
      title: 'Service Validation',
      description: 'Verify all services are running correctly',
      status: 'pending',
      icon: Settings,
    },
  ]);

  const getStatusIcon = (status: SetupStep['status']) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="h-5 w-5 text-green-500" />;
      case 'error':
        return <AlertCircle className="h-5 w-5 text-red-500" />;
      case 'in-progress':
        return <Clock className="h-5 w-5 text-blue-500 animate-spin" />;
      default:
        return null;
    }
  };

  const getStatusBadge = (status: SetupStep['status']) => {
    const variants = {
      pending: 'secondary',
      'in-progress': 'default',
      completed: 'success',
      error: 'destructive',
    } as const;

    const labels = {
      pending: 'Pending',
      'in-progress': 'In Progress',
      completed: 'Completed',
      error: 'Error',
    };

    return (
      <Badge variant={variants[status] as any}>
        {labels[status]}
      </Badge>
    );
  };

  const handleStepAction = (stepId: string) => {
    console.log(`Starting step: ${stepId}`);
  };

  const completedSteps = setupSteps.filter(step => step.status === 'completed').length;
  const totalSteps = setupSteps.length;
  const isComplete = completedSteps === totalSteps;

  return (
    <div className="space-y-6">

      <Card className="p-6">
        
        <div className="flex items-center gap-4 mb-6">
          <div className="p-2 bg-primary/10 rounded-lg">
            <Cloud className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h2 className="text-2xl font-semibold">Central Setup</h2>
            <p className="text-muted-foreground">
              Configure your central server to manage the wild-cloud infrastructure
            </p>
          </div>
        </div>

        <div>
          <Label htmlFor="upstream">Upstream</Label>
          <div className="flex w-full items-center mt-1">
            <Input id="internalDomain" value="https://mywildcloud.org"/>
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

      </Card>
      <Card className="p-6">
        <div>
          <h3 className="text-lg font-medium mb-4">Central Service</h3>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 mb-6">
            <div className="flex items-center gap-2">
              <Server className="h-5 w-5 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">IP Address: 192.168.8.50</span>
            </div>
            <div className="flex items-center gap-2">
              <Network className="h-5 w-5 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Network: 192.168.8.0/24</span>
            </div>
            <div className="flex items-center gap-2">
              <Settings className="h-5 w-5 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Version: 1.0.0 (update available)</span>
            </div>
            <div className="flex items-center gap-2">
              <Clock className="h-5 w-5 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Age: 12s</span>
            </div>
            <div className="flex items-center gap-2">
              <HelpCircle className="h-5 w-5 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Platform: ARM</span>
            </div>
            <div className="flex items-center gap-2">
              <CheckCircle className="h-5 w-5 text-green-500" />
              <span className="text-sm text-green-500">File permissions: Good</span>
            </div>
          </div>
  
          <div>
            <Label htmlFor="ip">IP</Label>
            <div className="flex w-full items-center mt-1">
              <Input id="ip" value="192.168.5.80"/>
              <Button variant="ghost">
                <HelpCircle/>
              </Button>
            </div>
          </div>
          <div>
            <Label htmlFor="interface">Interface</Label>
            <div className="flex w-full items-center mt-1">
              <Input id="interface" value="eth0"/>
              <Button variant="ghost">
                <HelpCircle/>
              </Button>
            </div>
          </div>
          <div className="flex gap-2 justify-end mt-4">
            <Button onClick={() => console.log('Update service')}>
              Update
            </Button>
            <Button onClick={() => console.log('Restart service')}>
              Restart
            </Button>
            <Button onClick={() => console.log('Save Settings')}>
              View log
            </Button>
          </div>
        </div>
      </Card>


      <Card className="p-6">
        <div>
          <h3 className="text-lg font-medium mb-4">DNS</h3>
          <div>Local resolution: ✅</div>
        </div>
      </Card>

      <Card className="p-6">
        <div>
          <h3 className="text-lg font-medium mb-4">DHCP</h3>
          <div>Status: On</div>
          <Label htmlFor="dhcpRange">Range</Label>
          <div className="flex w-full items-center mt-1">
            <Input id="dhcpRange" value="192.168.8.100,192.168.8.239"/>
            <Button variant="ghost">
              <HelpCircle/>
            </Button>
          </div>
          <Button onClick={() => console.log('View DHCP clients.')} className="mt-2">
            Clients
          </Button>
        </div>
      </Card>

      <Card className="p-6">
        <div>
          <h3 className="text-lg font-medium mb-4">PXE</h3>
          <div>Status: On</div>
          <Button onClick={() => console.log('Download PXE assets.')} className="mt-2">
            Download PXE Assets
          </Button>
        </div>

        {isComplete && (
          <div className="mt-6 p-4 bg-green-50 dark:bg-green-950 rounded-lg border border-green-200 dark:border-green-800">
            <div className="flex items-center gap-2 mb-2">
              <CheckCircle className="h-5 w-5 text-green-600" />
              <h3 className="font-medium text-green-800 dark:text-green-200">
                Setup Complete!
              </h3>
            </div>
            <p className="text-sm text-green-700 dark:text-green-300 mb-3">
              Your wild-cloud central server is now configured and ready to manage infrastructure.
            </p>
            <Button onClick={onComplete} className="bg-green-600 hover:bg-green-700">
              Continue to Infrastructure Setup
            </Button>
          </div>
        )}
      </Card>
    </div>
  );
}