import { useState } from 'react';
import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Container, Shield, Network, Database, CheckCircle, AlertCircle, Clock, Terminal, FileText } from 'lucide-react';

interface ClusterPhaseProps {
  onComplete?: () => void;
}

interface ClusterComponent {
  id: string;
  name: string;
  description: string;
  status: 'pending' | 'installing' | 'ready' | 'error';
  version?: string;
  logs?: string[];
}

export function ClusterPhase({ onComplete }: ClusterPhaseProps) {
  const [components, setComponents] = useState<ClusterComponent[]>([
    {
      id: 'talos-config',
      name: 'Talos Configuration',
      description: 'Generate and apply Talos cluster configuration',
      status: 'pending',
    },
    {
      id: 'kubernetes-bootstrap',
      name: 'Kubernetes Bootstrap',
      description: 'Initialize Kubernetes control plane',
      status: 'pending',
      version: 'v1.29.0',
    },
    {
      id: 'cni-plugin',
      name: 'Container Network Interface',
      description: 'Install and configure Cilium CNI',
      status: 'pending',
      version: 'v1.14.5',
    },
    {
      id: 'storage-class',
      name: 'Storage Classes',
      description: 'Configure persistent volume storage',
      status: 'pending',
    },
    {
      id: 'ingress-controller',
      name: 'Ingress Controller',
      description: 'Install Traefik ingress controller',
      status: 'pending',
      version: 'v3.0.0',
    },
    {
      id: 'monitoring',
      name: 'Cluster Monitoring',
      description: 'Deploy Prometheus and Grafana stack',
      status: 'pending',
    },
  ]);

  const [showLogs, setShowLogs] = useState<string | null>(null);

  const getStatusIcon = (status: ClusterComponent['status']) => {
    switch (status) {
      case 'ready':
        return <CheckCircle className="h-5 w-5 text-green-500" />;
      case 'error':
        return <AlertCircle className="h-5 w-5 text-red-500" />;
      case 'installing':
        return <Clock className="h-5 w-5 text-blue-500 animate-spin" />;
      default:
        return null;
    }
  };

  const getStatusBadge = (status: ClusterComponent['status']) => {
    const variants = {
      pending: 'secondary',
      installing: 'default',
      ready: 'success',
      error: 'destructive',
    } as const;

    const labels = {
      pending: 'Pending',
      installing: 'Installing',
      ready: 'Ready',
      error: 'Error',
    };

    return (
      <Badge variant={variants[status] as any}>
        {labels[status]}
      </Badge>
    );
  };

  const getComponentIcon = (id: string) => {
    switch (id) {
      case 'talos-config':
        return <FileText className="h-5 w-5" />;
      case 'kubernetes-bootstrap':
        return <Container className="h-5 w-5" />;
      case 'cni-plugin':
        return <Network className="h-5 w-5" />;
      case 'storage-class':
        return <Database className="h-5 w-5" />;
      case 'ingress-controller':
        return <Shield className="h-5 w-5" />;
      case 'monitoring':
        return <Terminal className="h-5 w-5" />;
      default:
        return <Container className="h-5 w-5" />;
    }
  };

  const handleComponentAction = (componentId: string, action: 'install' | 'retry') => {
    console.log(`${action} component: ${componentId}`);
  };

  const readyComponents = components.filter(component => component.status === 'ready').length;
  const totalComponents = components.length;
  const isComplete = readyComponents === totalComponents;

  return (
    <div className="space-y-6">
      <Card className="p-6">
        <div className="flex items-center gap-4 mb-6">
          <div className="p-2 bg-primary/10 rounded-lg">
            <Container className="h-6 w-6 text-primary" />
          </div>
          <div>
            <h2 className="text-2xl font-semibold">Kubernetes Installation</h2>
            <p className="text-muted-foreground">
              Install and configure Kubernetes on your cluster nodes
            </p>
          </div>
        </div>

        <div className="mb-6">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium">Installation Progress</span>
            <span className="text-sm text-muted-foreground">
              {readyComponents} of {totalComponents} components ready
            </span>
          </div>
          <div className="w-full bg-muted rounded-full h-2">
            <div
              className="bg-primary h-2 rounded-full transition-all duration-300"
              style={{ width: `${(readyComponents / totalComponents) * 100}%` }}
            />
          </div>
        </div>

        <div className="space-y-4">
          {components.map((component) => (
            <div key={component.id}>
              <div className="flex items-center gap-4 p-4 rounded-lg border bg-card">
                <div className="p-2 bg-muted rounded-lg">
                  {getComponentIcon(component.id)}
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="font-medium">{component.name}</h3>
                    {component.version && (
                      <Badge variant="outline" className="text-xs">
                        {component.version}
                      </Badge>
                    )}
                    {getStatusIcon(component.status)}
                  </div>
                  <p className="text-sm text-muted-foreground">{component.description}</p>
                </div>
                <div className="flex items-center gap-3">
                  {getStatusBadge(component.status)}
                  {(component.status === 'installing' || component.status === 'error') && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setShowLogs(showLogs === component.id ? null : component.id)}
                    >
                      <Terminal className="h-4 w-4 mr-1" />
                      Logs
                    </Button>
                  )}
                  {component.status === 'pending' && (
                    <Button
                      size="sm"
                      onClick={() => handleComponentAction(component.id, 'install')}
                    >
                      Install
                    </Button>
                  )}
                  {component.status === 'error' && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleComponentAction(component.id, 'retry')}
                    >
                      Retry
                    </Button>
                  )}
                </div>
              </div>
              
              {showLogs === component.id && (
                <Card className="mt-2 p-4 bg-black text-green-400 font-mono text-sm">
                  <div className="max-h-40 overflow-y-auto">
                    <div>Installing {component.name}...</div>
                    <div>✓ Checking prerequisites</div>
                    <div>✓ Downloading manifests</div>
                    {component.status === 'installing' && (
                      <div className="animate-pulse">⏳ Applying configuration...</div>
                    )}
                    {component.status === 'error' && (
                      <div className="text-red-400">✗ Installation failed: timeout waiting for pods</div>
                    )}
                  </div>
                </Card>
              )}
            </div>
          ))}
        </div>

        {isComplete && (
          <div className="mt-6 p-4 bg-green-50 dark:bg-green-950 rounded-lg border border-green-200 dark:border-green-800">
            <div className="flex items-center gap-2 mb-2">
              <CheckCircle className="h-5 w-5 text-green-600" />
              <h3 className="font-medium text-green-800 dark:text-green-200">
                Kubernetes Cluster Ready!
              </h3>
            </div>
            <p className="text-sm text-green-700 dark:text-green-300 mb-3">
              Your Kubernetes cluster is fully configured and ready for application deployment.
            </p>
            <Button onClick={onComplete} className="bg-green-600 hover:bg-green-700">
              Continue to App Management
            </Button>
          </div>
        )}
      </Card>

      <Card className="p-6">
        <h3 className="text-lg font-medium mb-4">Cluster Information</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div>
            <div className="font-medium mb-2">Control Plane</div>
            <div className="space-y-1 text-muted-foreground">
              <div>• API Server: https://cluster.wildcloud.local:6443</div>
              <div>• Nodes: 1 controller, 2 workers</div>
              <div>• Version: Kubernetes v1.29.0</div>
            </div>
          </div>
          <div>
            <div className="font-medium mb-2">Network Configuration</div>
            <div className="space-y-1 text-muted-foreground">
              <div>• Pod CIDR: 10.244.0.0/16</div>
              <div>• Service CIDR: 10.96.0.0/12</div>
              <div>• CNI: Cilium v1.14.5</div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}