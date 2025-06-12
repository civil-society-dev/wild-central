import { useState } from 'react';
import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Cloud, Server, Network, Settings, CheckCircle, AlertCircle, Clock } from 'lucide-react';

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
            <h2 className="text-2xl font-semibold">Wild-cloud Central Setup</h2>
            <p className="text-muted-foreground">
              Configure your central server to manage the wild-cloud infrastructure
            </p>
          </div>
        </div>

        <div className="mb-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium">Setup Progress</span>
            <span className="text-sm text-muted-foreground">
              {completedSteps} of {totalSteps} steps completed
            </span>
          </div>
          <div className="w-full bg-muted rounded-full h-2">
            <div
              className="bg-primary h-2 rounded-full transition-all duration-300"
              style={{ width: `${(completedSteps / totalSteps) * 100}%` }}
            />
          </div>
        </div>

        <div className="space-y-4">
          {setupSteps.map((step) => {
            const Icon = step.icon;
            return (
              <div
                key={step.id}
                className="flex items-center gap-4 p-4 rounded-lg border bg-card"
              >
                <div className="p-2 bg-muted rounded-lg">
                  <Icon className="h-5 w-5" />
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="font-medium">{step.title}</h3>
                    {getStatusIcon(step.status)}
                  </div>
                  <p className="text-sm text-muted-foreground">{step.description}</p>
                </div>
                <div className="flex items-center gap-3">
                  {getStatusBadge(step.status)}
                  {step.status === 'pending' && (
                    <Button
                      size="sm"
                      onClick={() => handleStepAction(step.id)}
                    >
                      Start
                    </Button>
                  )}
                  {step.status === 'error' && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleStepAction(step.id)}
                    >
                      Retry
                    </Button>
                  )}
                </div>
              </div>
            );
          })}
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