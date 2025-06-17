import { useEffect, useState } from 'react';
import { Cloud, CloudLightning } from 'lucide-react';
import { useConfig } from './hooks';
import {
  TabNavigation,
  SetupPhase,
  InfrastructurePhase,
  ClusterPhase,
  AppsPhase,
  Advanced,
  ErrorBoundary
} from './components';
import { ThemeToggle } from './components/ThemeToggle';
import type { Phase, Tab } from './components/TabNavigation';

function App() {
  const [currentTab, setCurrentTab] = useState<Tab>('setup');
  const [completedPhases, setCompletedPhases] = useState<Phase[]>([]);

  const { config } = useConfig();

  // Update phase state from config when it changes
  useEffect(() => {
    console.log('Config changed:', config);
    console.log('config?.wildcloud:', config?.wildcloud);
    if (config?.wildcloud?.currentPhase) {
      console.log('Setting currentTab to:', config.wildcloud.currentPhase);
      setCurrentTab(config.wildcloud.currentPhase as Phase);
    }
    if (config?.wildcloud?.completedPhases) {
      console.log('Setting completedPhases to:', config.wildcloud.completedPhases);
      setCompletedPhases(config.wildcloud.completedPhases as Phase[]);
    }
  }, [config]);

  const handlePhaseComplete = (phase: Phase) => {
    if (!completedPhases.includes(phase)) {
      setCompletedPhases(prev => [...prev, phase]);
    }
    
    // Auto-advance to next phase (excluding advanced)
    const phases: Phase[] = ['setup', 'infrastructure', 'cluster', 'apps'];
    const currentIndex = phases.indexOf(phase);
    if (currentIndex < phases.length - 1) {
      setCurrentTab(phases[currentIndex + 1]);
    }
  };

  const renderCurrentTab = () => {
    switch (currentTab) {
      case 'setup':
        return (
          <ErrorBoundary>
            <SetupPhase onComplete={() => handlePhaseComplete('setup')} />
          </ErrorBoundary>
        );
      case 'infrastructure':
        return (
          <ErrorBoundary>
            <InfrastructurePhase onComplete={() => handlePhaseComplete('infrastructure')} />
          </ErrorBoundary>
        );
      case 'cluster':
        return (
          <ErrorBoundary>
            <ClusterPhase onComplete={() => handlePhaseComplete('cluster')} />
          </ErrorBoundary>
        );
      case 'apps':
        return (
          <ErrorBoundary>
            <AppsPhase onComplete={() => handlePhaseComplete('apps')} />
          </ErrorBoundary>
        );
      case 'advanced':
        return (
          <ErrorBoundary>
            <Advanced />
          </ErrorBoundary>
        );
      default:
        return (
          <ErrorBoundary>
            <SetupPhase onComplete={() => handlePhaseComplete('setup')} />
          </ErrorBoundary>
        );
    }
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-8 max-w-7xl">
        <header className="mb-12">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-4">
              <div className="p-2 bg-primary/10 rounded-lg">
                <CloudLightning className="h-8 w-8 text-primary" />
              </div>
              <div>
                <h1 className="text-3xl font-bold tracking-tight text-foreground">
                  Wild Cloud
                </h1>
                <p className="text-muted-foreground text-lg">
                  Central
                </p>
              </div>
            </div>
            <ThemeToggle />
          </div>
        </header>
        
        <TabNavigation
          currentTab={currentTab}
          onTabChange={setCurrentTab}
          completedPhases={completedPhases}
        />
        
        {renderCurrentTab()}
      </div>
    </div>
  );
}

export default App;