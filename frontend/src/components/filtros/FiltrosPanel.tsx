import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Filter, X, ChevronDown, ChevronUp, RotateCcw } from 'lucide-react';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Select } from '../ui/Select';
import { Badge } from '../ui/Badge';
import { cn } from '../../lib/utils';
import { useFiltrosStore } from '../../stores/useFiltrosStore';
import { politicosApi } from '../../services/api';
import type { Cargo, Esfera, Genero } from '../../types';

const CARGOS: { value: Cargo; label: string }[] = [
  { value: 'DEPUTADO_FEDERAL', label: 'Deputado Federal' },
  { value: 'SENADOR', label: 'Senador' },
  { value: 'DEPUTADO_ESTADUAL', label: 'Deputado Estadual' },
  { value: 'DEPUTADO_DISTRITAL', label: 'Deputado Distrital' },
  { value: 'VEREADOR', label: 'Vereador' },
  { value: 'PREFEITO', label: 'Prefeito' },
  { value: 'GOVERNADOR', label: 'Governador' },
  { value: 'PRESIDENTE', label: 'Presidente' },
];

const ESFERAS: { value: Esfera; label: string }[] = [
  { value: 'FEDERAL', label: 'Federal' },
  { value: 'ESTADUAL', label: 'Estadual' },
  { value: 'MUNICIPAL', label: 'Municipal' },
];

const ESTADOS = [
  'AC', 'AL', 'AP', 'AM', 'BA', 'CE', 'DF', 'ES', 'GO', 'MA', 'MT', 'MS',
  'MG', 'PA', 'PB', 'PR', 'PE', 'PI', 'RJ', 'RN', 'RS', 'RO', 'RR', 'SC',
  'SP', 'SE', 'TO',
];

const GENEROS: { value: Genero; label: string }[] = [
  { value: 'M', label: 'Masculino' },
  { value: 'F', label: 'Feminino' },
  { value: 'OUTRO', label: 'Outro' },
];

interface FiltrosPanelProps {
  onClose?: () => void;
  isMobile?: boolean;
}

export function FiltrosPanel({ onClose, isMobile }: FiltrosPanelProps) {
  const { filtros, setFiltro, limparFiltros } = useFiltrosStore();
  const [expandedSections, setExpandedSections] = useState<string[]>(['cargo', 'status']);

  // Buscar contagens por cargo usando useQueries
  const contagensQueries = useQuery({
    queryKey: ['politicos-contagens-cargo'],
    queryFn: async () => {
      const contagens: Record<Cargo, number> = {} as Record<Cargo, number>;
      
      // Buscar contagem para cada cargo
      await Promise.all(
        CARGOS.map(async (cargo) => {
          try {
            const result = await politicosApi.listar({ 
              cargo: [cargo.value], 
              porPagina: 1 
            });
            contagens[cargo.value] = result.total || 0;
          } catch (error) {
            contagens[cargo.value] = 0;
          }
        })
      );
      
      return contagens;
    },
    staleTime: 5 * 60 * 1000, // Cache por 5 minutos
  });

  const contagensPorCargo = contagensQueries.data || ({} as Record<Cargo, number>);

  const toggleSection = (section: string) => {
    setExpandedSections((prev) =>
      prev.includes(section)
        ? prev.filter((s) => s !== section)
        : [...prev, section]
    );
  };

  const toggleArrayFilter = <T extends string>(
    key: 'cargo' | 'esfera' | 'estado' | 'genero',
    value: T
  ) => {
    const current = (filtros[key] as T[] | undefined) || [];
    const updated = current.includes(value)
      ? current.filter((v) => v !== value)
      : [...current, value];
    setFiltro(key, updated.length > 0 ? updated : undefined);
  };

  const getActiveFiltersCount = () => {
    let count = 0;
    if (filtros.nome) count++;
    if (filtros.cargo?.length) count += filtros.cargo.length;
    if (filtros.esfera?.length) count += filtros.esfera.length;
    if (filtros.estado?.length) count += filtros.estado.length;
    if (filtros.genero?.length) count += filtros.genero.length;
    if (filtros.emExercicio !== undefined) count++;
    return count;
  };

  const activeCount = getActiveFiltersCount();

  return (
    <div className={cn(
      'bg-background-card border border-border rounded-xl',
      isMobile ? 'p-4' : 'p-6'
    )}>
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-accent-primary" />
          <h2 className="text-lg font-semibold text-content-primary">Filtros</h2>
          {activeCount > 0 && (
            <Badge variant="info">{activeCount}</Badge>
          )}
        </div>
        <div className="flex items-center gap-2">
          {activeCount > 0 && (
            <Button
              variant="ghost"
              size="sm"
              onClick={limparFiltros}
              className="text-content-muted hover:text-accent-danger"
            >
              <RotateCcw className="w-4 h-4" />
              Limpar
            </Button>
          )}
          {isMobile && onClose && (
            <button
              onClick={onClose}
              className="p-2 rounded-lg hover:bg-background-secondary text-content-muted"
            >
              <X className="w-5 h-5" />
            </button>
          )}
        </div>
      </div>

      <div className="space-y-4">
        {/* Busca por nome */}
        <Input
          placeholder="Buscar por nome..."
          value={filtros.nome || ''}
          onChange={(e) => setFiltro('nome', e.target.value || undefined)}
        />

        {/* Cargo */}
        <FilterSection
          title="Cargo"
          isExpanded={expandedSections.includes('cargo')}
          onToggle={() => toggleSection('cargo')}
        >
          <div className="flex flex-wrap gap-2">
            {CARGOS.map((cargo) => {
              const count = contagensPorCargo[cargo.value] || 0;
              return (
                <button
                  key={cargo.value}
                  onClick={() => toggleArrayFilter('cargo', cargo.value)}
                  className={cn(
                    'px-3 py-1.5 rounded-lg text-sm font-medium transition-all flex items-center gap-2',
                    filtros.cargo?.includes(cargo.value)
                      ? 'bg-accent-primary text-background'
                      : 'bg-background-secondary text-content-secondary hover:bg-background-hover'
                  )}
                >
                  <span>{cargo.label}</span>
                  {count > 0 && (
                    <Badge 
                      variant="secondary" 
                      className={cn(
                        'text-xs',
                        filtros.cargo?.includes(cargo.value)
                          ? 'bg-background/30 text-background'
                          : ''
                      )}
                    >
                      {count}
                    </Badge>
                  )}
                </button>
              );
            })}
          </div>
        </FilterSection>

        {/* Esfera */}
        <FilterSection
          title="Esfera"
          isExpanded={expandedSections.includes('esfera')}
          onToggle={() => toggleSection('esfera')}
        >
          <div className="flex flex-wrap gap-2">
            {ESFERAS.map((esfera) => (
              <button
                key={esfera.value}
                onClick={() => toggleArrayFilter('esfera', esfera.value)}
                className={cn(
                  'px-3 py-1.5 rounded-lg text-sm font-medium transition-all',
                  filtros.esfera?.includes(esfera.value)
                    ? 'bg-accent-secondary text-background'
                    : 'bg-background-secondary text-content-secondary hover:bg-background-hover'
                )}
              >
                {esfera.label}
              </button>
            ))}
          </div>
        </FilterSection>

        {/* Estado */}
        <FilterSection
          title="Estado"
          isExpanded={expandedSections.includes('estado')}
          onToggle={() => toggleSection('estado')}
        >
          <div className="flex flex-wrap gap-1.5">
            {ESTADOS.map((estado) => (
              <button
                key={estado}
                onClick={() => toggleArrayFilter('estado', estado)}
                className={cn(
                  'px-2 py-1 rounded text-xs font-medium transition-all',
                  filtros.estado?.includes(estado)
                    ? 'bg-accent-primary text-background'
                    : 'bg-background-secondary text-content-secondary hover:bg-background-hover'
                )}
              >
                {estado}
              </button>
            ))}
          </div>
        </FilterSection>

        {/* Status */}
        <FilterSection
          title="Status"
          isExpanded={expandedSections.includes('status')}
          onToggle={() => toggleSection('status')}
        >
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => setFiltro('emExercicio', filtros.emExercicio === true ? undefined : true)}
              className={cn(
                'px-3 py-1.5 rounded-lg text-sm font-medium transition-all',
                filtros.emExercicio === true
                  ? 'bg-accent-success text-background'
                  : 'bg-background-secondary text-content-secondary hover:bg-background-hover'
              )}
            >
              Em exercício
            </button>
            <button
              onClick={() => setFiltro('emExercicio', filtros.emExercicio === false ? undefined : false)}
              className={cn(
                'px-3 py-1.5 rounded-lg text-sm font-medium transition-all',
                filtros.emExercicio === false
                  ? 'bg-accent-warning text-background'
                  : 'bg-background-secondary text-content-secondary hover:bg-background-hover'
              )}
            >
              Fora de exercício
            </button>
          </div>
        </FilterSection>

        {/* Gênero */}
        <FilterSection
          title="Gênero"
          isExpanded={expandedSections.includes('genero')}
          onToggle={() => toggleSection('genero')}
        >
          <div className="flex flex-wrap gap-2">
            {GENEROS.map((genero) => (
              <button
                key={genero.value}
                onClick={() => toggleArrayFilter('genero', genero.value)}
                className={cn(
                  'px-3 py-1.5 rounded-lg text-sm font-medium transition-all',
                  filtros.genero?.includes(genero.value)
                    ? 'bg-accent-secondary text-background'
                    : 'bg-background-secondary text-content-secondary hover:bg-background-hover'
                )}
              >
                {genero.label}
              </button>
            ))}
          </div>
        </FilterSection>

        {/* Ordenação */}
        <div className="pt-4 border-t border-border">
          <Select
            label="Ordenar por"
            value={filtros.ordenarPor || 'nome'}
            onChange={(e) => setFiltro('ordenarPor', e.target.value as typeof filtros.ordenarPor)}
            options={[
              { value: 'nome', label: 'Nome' },
              { value: 'partido', label: 'Partido' },
              { value: 'presenca', label: 'Presença' },
              { value: 'proposicoes', label: 'Proposições' },
              { value: 'gastos', label: 'Gastos' },
            ]}
          />
        </div>
      </div>
    </div>
  );
}

interface FilterSectionProps {
  title: string;
  isExpanded: boolean;
  onToggle: () => void;
  children: React.ReactNode;
}

function FilterSection({ title, isExpanded, onToggle, children }: FilterSectionProps) {
  return (
    <div className="border-t border-border pt-4">
      <button
        onClick={onToggle}
        className="w-full flex items-center justify-between text-left mb-3"
      >
        <span className="text-sm font-medium text-content-primary">{title}</span>
        {isExpanded ? (
          <ChevronUp className="w-4 h-4 text-content-muted" />
        ) : (
          <ChevronDown className="w-4 h-4 text-content-muted" />
        )}
      </button>
      {isExpanded && children}
    </div>
  );
}

