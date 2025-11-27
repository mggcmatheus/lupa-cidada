import { useState } from 'react';
import { Link } from 'react-router-dom';
import { MapPin, Building2, Check, Plus, ExternalLink, User } from 'lucide-react';
import { cn, formatCurrency, getCargoLabel } from '../../lib/utils';
import { Badge } from '../ui/Badge';
import { useComparacaoStore } from '../../stores/useComparacaoStore';
import type { Politico } from '../../types';

interface PoliticoCardProps {
  politico: Politico;
  showCompareButton?: boolean;
}

export function PoliticoCard({ politico, showCompareButton = true }: PoliticoCardProps) {
  const { togglePolitico, estaSelecionado, podeAdicionar } = useComparacaoStore();
  const selecionado = estaSelecionado(politico.id);
  const podeSelecionarMais = podeAdicionar();
  const [imageError, setImageError] = useState(false);

  const handleCompareClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    togglePolitico(politico.id);
  };

  return (
    <Link
      to={`/politicos/${politico.id}`}
      className={cn(
        'group block bg-background-card border rounded-xl p-5 transition-all duration-300',
        'hover:shadow-glow-cyan/10 hover:border-border-light',
        selecionado ? 'border-accent-primary shadow-glow-cyan/20' : 'border-border'
      )}
    >
      <div className="flex items-start gap-4">
        {/* Foto */}
        <div className="relative flex-shrink-0">
          {!imageError && politico.fotoUrl ? (
            <img
              src={politico.fotoUrl}
              alt={politico.nome}
              className="w-16 h-16 rounded-full object-cover border-2 border-border group-hover:border-accent-primary/50 transition-colors"
              onError={() => setImageError(true)}
            />
          ) : (
            <div className="w-16 h-16 rounded-full bg-background-secondary border-2 border-border group-hover:border-accent-primary/50 transition-colors flex items-center justify-center">
              <User className="w-8 h-8 text-content-muted" />
            </div>
          )}
          {/* Indicador de partido */}
          <div
            className="absolute -bottom-1 -right-1 w-6 h-6 rounded-full border-2 border-background-card flex items-center justify-center text-[10px] font-bold text-white"
            style={{ backgroundColor: politico.partido.cor }}
            title={politico.partido.nome}
          >
            {politico.partido.sigla.slice(0, 2)}
          </div>
        </div>

        {/* Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2">
            <div>
              <h3 className="font-semibold text-content-primary truncate group-hover:text-accent-primary transition-colors">
                {politico.nome}
              </h3>
              <p className="text-sm text-content-secondary">
                {politico.partido.sigla}
              </p>
            </div>
            
            {/* Botão de comparar */}
            {showCompareButton && (
              <button
                onClick={handleCompareClick}
                disabled={!selecionado && !podeSelecionarMais}
                className={cn(
                  'flex-shrink-0 p-2 rounded-lg transition-all duration-200',
                  selecionado
                    ? 'bg-accent-primary text-background'
                    : podeSelecionarMais
                    ? 'bg-background-secondary text-content-secondary hover:bg-accent-primary/20 hover:text-accent-primary'
                    : 'bg-background-secondary text-content-muted cursor-not-allowed opacity-50'
                )}
                title={selecionado ? 'Remover da comparação' : 'Adicionar à comparação'}
              >
                {selecionado ? (
                  <Check className="w-4 h-4" />
                ) : (
                  <Plus className="w-4 h-4" />
                )}
              </button>
            )}
          </div>

          {/* Tags */}
          <div className="flex flex-wrap items-center gap-2 mt-3">
            <Badge variant="info">
              <Building2 className="w-3 h-3" />
              {getCargoLabel(politico.cargoAtual.tipo)}
            </Badge>
            <Badge variant={politico.cargoAtual.emExercicio ? 'success' : 'secondary'}>
              {politico.cargoAtual.emExercicio ? 'Em exercício' : 'Fora de exercício'}
            </Badge>
          </div>

          {/* Local */}
          <div className="flex items-center gap-1 mt-3 text-sm text-content-muted">
            <MapPin className="w-3.5 h-3.5" />
            <span>
              {politico.cargoAtual.municipio
                ? `${politico.cargoAtual.municipio}, ${politico.cargoAtual.estado}`
                : politico.cargoAtual.estado}
            </span>
          </div>
        </div>
      </div>

      {/* Footer com salário */}
      <div className="mt-4 pt-4 border-t border-border flex items-center justify-between">
        <div>
          <p className="text-xs text-content-muted">Salário bruto</p>
          <p className="text-sm font-mono font-semibold text-accent-primary">
            {formatCurrency(politico.salarioBruto)}
          </p>
        </div>
        <div className="flex items-center gap-1 text-content-muted group-hover:text-accent-primary transition-colors">
          <span className="text-xs">Ver detalhes</span>
          <ExternalLink className="w-3.5 h-3.5" />
        </div>
      </div>
    </Link>
  );
}

