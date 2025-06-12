import { useEffect, useState } from 'react';
import { Cloud } from 'lucide-react';
import { useConfig } from './hooks';
import {
  PhaseNavigation,
  SetupPhase,
  InfrastructurePhase,
  ClusterPhase,
  AppsPhase,
  ErrorBoundary
} from './components';
import { ThemeToggle } from './components/ThemeToggle';
import type { Phase } from './components/PhaseNavigation';

function App() {
  const [currentPhase, setCurrentPhase] = useState<Phase>('setup');
  const [completedPhases, setCompletedPhases] = useState<Phase[]>([]);

  const { config } = useConfig();

  // Update phase state from config when it changes
  useEffect(() => {
    console.log('Config changed:', config);
    console.log('config?.wildcloud:', config?.wildcloud);
    if (config?.wildcloud?.currentPhase) {
      console.log('Setting currentPhase to:', config.wildcloud.currentPhase);
      setCurrentPhase(config.wildcloud.currentPhase as Phase);
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
    
    // Auto-advance to next phase
    const phases: Phase[] = ['setup', 'infrastructure', 'cluster', 'apps'];
    const currentIndex = phases.indexOf(phase);
    if (currentIndex < phases.length - 1) {
      setCurrentPhase(phases[currentIndex + 1]);
    }
  };

  const renderCurrentPhase = () => {
    switch (currentPhase) {
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
                <Cloud className="h-8 w-8 text-primary" />
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
        
        <PhaseNavigation
          currentPhase={currentPhase}
          onPhaseChange={setCurrentPhase}
          completedPhases={completedPhases}
        />
        
        {renderCurrentPhase()}
      </div>
    </div>
  );
}

export default App;