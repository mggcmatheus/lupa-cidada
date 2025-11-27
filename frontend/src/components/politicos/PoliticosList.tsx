import { PoliticoCard } from './PoliticoCard';
import { SkeletonCard } from '../ui/Skeleton';
import { Users } from 'lucide-react';
import type { Politico } from '../../types';

interface PoliticosListProps {
  politicos: Politico[];
  isLoading?: boolean;
  showCompareButton?: boolean;
}

export function PoliticosList({ politicos, isLoading, showCompareButton = true }: PoliticosListProps) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {Array.from({ length: 6 }).map((_, i) => (
          <SkeletonCard key={i} />
        ))}
      </div>
    );
  }

  if (politicos.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-center">
        <div className="w-16 h-16 rounded-full bg-background-secondary flex items-center justify-center mb-4">
          <Users className="w-8 h-8 text-content-muted" />
        </div>
        <h3 className="text-lg font-semibold text-content-primary mb-2">
          Nenhum político encontrado
        </h3>
        <p className="text-content-secondary max-w-md">
          Tente ajustar os filtros para encontrar políticos que correspondam aos seus critérios.
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {politicos.map((politico, index) => (
        <div
          key={politico.id}
          className="animate-fade-in"
          style={{ animationDelay: `${index * 50}ms` }}
        >
          <PoliticoCard politico={politico} showCompareButton={showCompareButton} />
        </div>
      ))}
    </div>
  );
}

