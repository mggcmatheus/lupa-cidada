import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { SlidersHorizontal, ChevronLeft, ChevronRight } from 'lucide-react';
import { PoliticosList } from '../components/politicos/PoliticosList';
import { FiltrosPanel } from '../components/filtros/FiltrosPanel';
import { Button } from '../components/ui/Button';
import { useFiltrosStore } from '../stores/useFiltrosStore';
import { useComparacaoStore } from '../stores/useComparacaoStore';
import { formatNumber } from '../lib/utils';
import type { Politico } from '../types';

// Dados mockados para demonstração
const MOCK_POLITICOS: Politico[] = [
  {
    id: '1',
    nome: 'João Silva',
    nomeCivil: 'João Pedro da Silva Santos',
    fotoUrl: 'https://www.camara.leg.br/internet/deputado/bandep/placeholder.png',
    dataNascimento: '1965-03-15',
    genero: 'M',
    partido: { sigla: 'PT', nome: 'Partido dos Trabalhadores', cor: '#CC0000' },
    cargoAtual: {
      tipo: 'DEPUTADO_FEDERAL',
      esfera: 'FEDERAL',
      estado: 'SP',
      dataInicio: '2023-02-01',
      emExercicio: true,
    },
    historicoCargos: [],
    contato: { email: 'joao.silva@camara.leg.br' },
    redesSociais: {},
    salarioBruto: 33763.00,
    salarioLiquido: 25000.00,
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01',
  },
  {
    id: '2',
    nome: 'Maria Santos',
    nomeCivil: 'Maria das Graças Santos',
    fotoUrl: 'https://www.camara.leg.br/internet/deputado/bandep/placeholder.png',
    dataNascimento: '1970-07-22',
    genero: 'F',
    partido: { sigla: 'PL', nome: 'Partido Liberal', cor: '#003366' },
    cargoAtual: {
      tipo: 'SENADOR',
      esfera: 'FEDERAL',
      estado: 'RJ',
      dataInicio: '2023-02-01',
      emExercicio: true,
    },
    historicoCargos: [],
    contato: { email: 'maria.santos@senado.leg.br' },
    redesSociais: {},
    salarioBruto: 41650.92,
    salarioLiquido: 30000.00,
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01',
  },
  {
    id: '3',
    nome: 'Pedro Costa',
    nomeCivil: 'Pedro Henrique Costa Lima',
    fotoUrl: 'https://www.camara.leg.br/internet/deputado/bandep/placeholder.png',
    dataNascimento: '1980-11-08',
    genero: 'M',
    partido: { sigla: 'MDB', nome: 'Movimento Democrático Brasileiro', cor: '#00AA00' },
    cargoAtual: {
      tipo: 'DEPUTADO_ESTADUAL',
      esfera: 'ESTADUAL',
      estado: 'MG',
      dataInicio: '2023-02-01',
      emExercicio: true,
    },
    historicoCargos: [],
    contato: {},
    redesSociais: {},
    salarioBruto: 25322.25,
    salarioLiquido: 18000.00,
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01',
  },
  {
    id: '4',
    nome: 'Ana Oliveira',
    nomeCivil: 'Ana Paula Oliveira Souza',
    fotoUrl: 'https://www.camara.leg.br/internet/deputado/bandep/placeholder.png',
    dataNascimento: '1975-04-30',
    genero: 'F',
    partido: { sigla: 'PSOL', nome: 'Partido Socialismo e Liberdade', cor: '#FFD700' },
    cargoAtual: {
      tipo: 'VEREADOR',
      esfera: 'MUNICIPAL',
      estado: 'SP',
      municipio: 'São Paulo',
      dataInicio: '2021-01-01',
      emExercicio: true,
    },
    historicoCargos: [],
    contato: {},
    redesSociais: { twitter: '@anaoliveira' },
    salarioBruto: 18991.68,
    salarioLiquido: 14000.00,
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01',
  },
  {
    id: '5',
    nome: 'Carlos Ferreira',
    nomeCivil: 'Carlos Alberto Ferreira',
    fotoUrl: 'https://www.camara.leg.br/internet/deputado/bandep/placeholder.png',
    dataNascimento: '1960-09-12',
    genero: 'M',
    partido: { sigla: 'UNIÃO', nome: 'União Brasil', cor: '#2E3092' },
    cargoAtual: {
      tipo: 'GOVERNADOR',
      esfera: 'ESTADUAL',
      estado: 'BA',
      dataInicio: '2023-01-01',
      emExercicio: true,
    },
    historicoCargos: [],
    contato: {},
    redesSociais: {},
    salarioBruto: 35462.22,
    salarioLiquido: 26000.00,
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01',
  },
  {
    id: '6',
    nome: 'Fernanda Lima',
    nomeCivil: 'Fernanda Cristina Lima',
    fotoUrl: 'https://www.camara.leg.br/internet/deputado/bandep/placeholder.png',
    dataNascimento: '1985-02-28',
    genero: 'F',
    partido: { sigla: 'NOVO', nome: 'Partido Novo', cor: '#FF6600' },
    cargoAtual: {
      tipo: 'DEPUTADO_FEDERAL',
      esfera: 'FEDERAL',
      estado: 'RS',
      dataInicio: '2023-02-01',
      emExercicio: true,
    },
    historicoCargos: [],
    contato: { email: 'fernanda.lima@camara.leg.br' },
    redesSociais: { instagram: '@fernanda.lima' },
    salarioBruto: 33763.00,
    salarioLiquido: 25000.00,
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01',
  },
];

export function Politicos() {
  const [showFilters, setShowFilters] = useState(false);
  const { filtros, setFiltro } = useFiltrosStore();
  const { politicosSelecionados } = useComparacaoStore();

  // Por enquanto usamos dados mockados, depois conectamos com a API real
  const { data, isLoading } = useQuery({
    queryKey: ['politicos', filtros],
    queryFn: async () => {
      // Simular delay de API
      await new Promise((resolve) => setTimeout(resolve, 500));
      
      // Filtrar dados mockados
      let resultado = MOCK_POLITICOS;
      
      if (filtros.nome) {
        const termo = filtros.nome.toLowerCase();
        resultado = resultado.filter((p) => 
          p.nome.toLowerCase().includes(termo) ||
          p.nomeCivil.toLowerCase().includes(termo)
        );
      }
      
      if (filtros.cargo?.length) {
        resultado = resultado.filter((p) => 
          filtros.cargo!.includes(p.cargoAtual.tipo)
        );
      }
      
      if (filtros.esfera?.length) {
        resultado = resultado.filter((p) => 
          filtros.esfera!.includes(p.cargoAtual.esfera)
        );
      }
      
      if (filtros.estado?.length) {
        resultado = resultado.filter((p) => 
          filtros.estado!.includes(p.cargoAtual.estado)
        );
      }

      if (filtros.emExercicio !== undefined) {
        resultado = resultado.filter((p) => 
          p.cargoAtual.emExercicio === filtros.emExercicio
        );
      }
      
      return {
        data: resultado,
        total: resultado.length,
        pagina: 1,
        porPagina: 12,
        totalPaginas: 1,
      };
    },
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

