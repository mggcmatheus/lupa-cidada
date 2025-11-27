import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { SlidersHorizontal, ChevronLeft, ChevronRight } from 'lucide-react';
import { politicosApi } from '../services/api';
import { PoliticosList } from '../components/politicos/PoliticosList';
import { FiltrosPanel } from '../components/filtros/FiltrosPanel';
import { Button } from '../components/ui/Button';
import { useFiltrosStore } from '../stores/useFiltrosStore';
import { useComparacaoStore } from '../stores/useComparacaoStore';
import { formatNumber } from '../lib/utils';

export function Politicos() {
  const [showFilters, setShowFilters] = useState(false);
  const { filtros, setFiltro } = useFiltrosStore();
  const { politicosSelecionados } = useComparacaoStore();

  // Buscar dados da API
  const { data, isLoading } = useQuery({
    queryKey: ['politicos', filtros],
    queryFn: () => politicosApi.listar(filtros),
  });

  const totalPaginas = data?.totalPaginas || 1;
  const paginaAtual = filtros.pagina || 1;

  return (
    <div className="min-h-screen py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-content-primary mb-2">
            Políticos
          </h1>
          <p className="text-content-secondary">
            Encontre e acompanhe a atuação dos políticos brasileiros
          </p>
        </div>

        <div className="flex flex-col lg:flex-row gap-8">
          {/* Sidebar com filtros (desktop) */}
          <aside className="hidden lg:block w-80 flex-shrink-0">
            <div className="sticky top-24">
              <FiltrosPanel />
            </div>
          </aside>

          {/* Conteúdo principal */}
          <div className="flex-1">
            {/* Toolbar */}
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center gap-4">
                {/* Botão filtros mobile */}
                <Button
                  variant="secondary"
                  className="lg:hidden"
                  onClick={() => setShowFilters(true)}
                >
                  <SlidersHorizontal className="w-4 h-4" />
                  Filtros
                </Button>

                {/* Contador */}
                <p className="text-sm text-content-secondary">
                  {isLoading ? (
                    'Carregando...'
                  ) : (
                    <>
                      <span className="font-semibold text-content-primary">
                        {formatNumber(data?.total || 0)}
                      </span>{' '}
                      políticos encontrados
                    </>
                  )}
                </p>
              </div>

              {/* Indicador de selecionados */}
              {politicosSelecionados.length > 0 && (
                <div className="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-lg bg-accent-primary/10 border border-accent-primary/30">
                  <span className="text-sm text-accent-primary">
                    {politicosSelecionados.length} selecionados para comparar
                  </span>
                </div>
              )}
            </div>

            {/* Lista de políticos */}
            <PoliticosList 
              politicos={data?.data || []} 
              isLoading={isLoading} 
            />

            {/* Paginação */}
            {data && data.totalPaginas > 1 && (
              <div className="mt-8 flex items-center justify-center gap-2">
                <Button
                  variant="secondary"
                  size="sm"
                  disabled={paginaAtual === 1}
                  onClick={() => setFiltro('pagina', paginaAtual - 1)}
                >
                  <ChevronLeft className="w-4 h-4" />
                  Anterior
                </Button>
                
                <div className="flex items-center gap-1 px-4">
                  <span className="text-content-secondary text-sm">
                    Página{' '}
                    <span className="font-semibold text-content-primary">
                      {paginaAtual}
                    </span>{' '}
                    de{' '}
                    <span className="font-semibold text-content-primary">
                      {totalPaginas}
                    </span>
                  </span>
                </div>

                <Button
                  variant="secondary"
                  size="sm"
                  disabled={paginaAtual === totalPaginas}
                  onClick={() => setFiltro('pagina', paginaAtual + 1)}
                >
                  Próxima
                  <ChevronRight className="w-4 h-4" />
                </Button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Modal de filtros mobile */}
      {showFilters && (
        <div className="fixed inset-0 z-50 lg:hidden">
          <div 
            className="absolute inset-0 bg-black/50 backdrop-blur-sm"
            onClick={() => setShowFilters(false)}
          />
          <div className="absolute right-0 top-0 bottom-0 w-full max-w-sm bg-background overflow-y-auto animate-slide-in-right">
            <FiltrosPanel 
              isMobile 
              onClose={() => setShowFilters(false)} 
            />
          </div>
        </div>
      )}
    </div>
  );
}
