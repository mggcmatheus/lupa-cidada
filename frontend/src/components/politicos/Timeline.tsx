import { Building2, MapPin, Calendar, CheckCircle2 } from 'lucide-react';
import type { CargoAtual, Cargo } from '../../types';
import { formatDate, getCargoLabel, getEsferaLabel } from '../../lib/utils';
import { Badge } from '../ui/Badge';
import { cn } from '../../lib/utils';

interface TimelineProps {
  cargoAtual: CargoAtual;
  historicoCargos: CargoAtual[];
}

export function Timeline({ cargoAtual, historicoCargos }: TimelineProps) {
  // Combinar cargo atual e histórico, ordenando por data (mais recente primeiro)
  const todosCargos: (CargoAtual & { isAtual: boolean })[] = [];
  
  // Adicionar cargo atual se estiver em exercício
  if (cargoAtual.emExercicio) {
    todosCargos.push({ ...cargoAtual, isAtual: true });
  }
  
  // Adicionar histórico
  historicoCargos.forEach(cargo => {
    todosCargos.push({ ...cargo, isAtual: false });
  });
  
  // Ordenar por data de início (mais recente primeiro)
  todosCargos.sort((a, b) => {
    const dateA = new Date(a.dataInicio).getTime();
    const dateB = new Date(b.dataInicio).getTime();
    return dateB - dateA;
  });

  if (todosCargos.length === 0) {
    return (
      <div className="text-center py-8 text-content-secondary">
        <p>Nenhum histórico de cargo disponível.</p>
      </div>
    );
  }

  return (
    <div className="relative">
      {/* Linha vertical da timeline */}
      <div className="absolute left-6 top-0 bottom-0 w-0.5 bg-border" />
      
      <div className="space-y-6">
        {todosCargos.map((cargo, index) => (
          <TimelineItem
            key={`${cargo.tipo}-${cargo.dataInicio}-${index}`}
            cargo={cargo}
            isFirst={index === 0}
            isLast={index === todosCargos.length - 1}
          />
        ))}
      </div>
    </div>
  );
}

interface TimelineItemProps {
  cargo: CargoAtual & { isAtual: boolean };
  isFirst: boolean;
  isLast: boolean;
}

function TimelineItem({ cargo, isFirst, isLast }: TimelineItemProps) {
  const getCargoIcon = (tipo: Cargo) => {
    switch (tipo) {
      case 'DEPUTADO_FEDERAL':
      case 'DEPUTADO_ESTADUAL':
      case 'DEPUTADO_DISTRITAL':
        return '🏛️';
      case 'SENADOR':
        return '⚖️';
      case 'PREFEITO':
        return '🏢';
      case 'VEREADOR':
        return '🏛️';
      case 'GOVERNADOR':
        return '🏛️';
      case 'PRESIDENTE':
        return '🇧🇷';
      default:
        return '👤';
    }
  };

  const getCargoColor = (tipo: Cargo) => {
    switch (tipo) {
      case 'DEPUTADO_FEDERAL':
      case 'DEPUTADO_ESTADUAL':
      case 'DEPUTADO_DISTRITAL':
        return 'bg-blue-500';
      case 'SENADOR':
        return 'bg-purple-500';
      case 'PREFEITO':
        return 'bg-green-500';
      case 'VEREADOR':
        return 'bg-cyan-500';
      case 'GOVERNADOR':
        return 'bg-orange-500';
      case 'PRESIDENTE':
        return 'bg-yellow-500';
      default:
        return 'bg-gray-500';
    }
  };

  return (
    <div className="relative flex gap-4">
      {/* Ícone do ponto na timeline */}
      <div className="relative z-10 flex-shrink-0">
        <div
          className={cn(
            'w-12 h-12 rounded-full flex items-center justify-center text-xl border-2 border-background',
            cargo.isAtual
              ? 'bg-accent-primary border-accent-primary shadow-lg shadow-accent-primary/30'
              : getCargoColor(cargo.tipo)
          )}
        >
          {cargo.isAtual ? (
            <CheckCircle2 className="w-6 h-6 text-background" />
          ) : (
            <span>{getCargoIcon(cargo.tipo)}</span>
          )}
        </div>
      </div>

      {/* Conteúdo do item */}
      <div className="flex-1 pb-6">
        <div className="bg-background-card border border-border rounded-xl p-4 hover:border-accent-primary/50 transition-colors">
          <div className="flex flex-wrap items-start justify-between gap-3 mb-3">
            <div className="flex-1">
              <div className="flex flex-wrap items-center gap-2 mb-2">
                <Badge
                  variant={cargo.isAtual ? 'success' : 'info'}
                  className="text-sm font-semibold"
                >
                  {getCargoLabel(cargo.tipo)}
                </Badge>
                {cargo.isAtual && (
                  <Badge variant="success" className="text-xs">
                    Em exercício
                  </Badge>
                )}
                {!cargo.emExercicio && !cargo.isAtual && (
                  <Badge variant="secondary" className="text-xs">
                    Finalizado
                  </Badge>
                )}
              </div>
              
              <div className="flex flex-wrap items-center gap-3 text-sm text-content-secondary">
                <span className="flex items-center gap-1">
                  <Building2 className="w-3.5 h-3.5" />
                  {getEsferaLabel(cargo.esfera)}
                </span>
                <span className="flex items-center gap-1">
                  <MapPin className="w-3.5 h-3.5" />
                  {cargo.municipio 
                    ? `${cargo.municipio}, ${cargo.estado}`
                    : cargo.estado
                  }
                </span>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-4 text-sm text-content-secondary pt-3 border-t border-border">
            <span className="flex items-center gap-1.5">
              <Calendar className="w-4 h-4" />
              <span>
                <strong>Início:</strong> {formatDate(cargo.dataInicio)}
              </span>
            </span>
            {cargo.dataFim && (
              <span className="flex items-center gap-1.5">
                <Calendar className="w-4 h-4" />
                <span>
                  <strong>Fim:</strong> {formatDate(cargo.dataFim)}
                </span>
              </span>
            )}
            {!cargo.dataFim && !cargo.isAtual && (
              <span className="text-content-muted italic">
                Data de fim não informada
              </span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

