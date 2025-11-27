import { Link } from 'react-router-dom';
import { Scale, X, ArrowRight, Users, Vote, Receipt, FileText, TrendingUp } from 'lucide-react';
import { useComparacaoStore } from '../stores/useComparacaoStore';
import { Button } from '../components/ui/Button';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';
import { Avatar } from '../components/ui/Avatar';
import { cn, formatCurrency, formatPercentage, getCargoLabel } from '../lib/utils';
import type { Politico, EstatisticasPolitico } from '../types';

// Mock data para demonstração
const MOCK_POLITICOS_DETALHES: Record<string, { politico: Politico; estatisticas: EstatisticasPolitico }> = {
  '1': {
    politico: {
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
    estatisticas: {
      totalVotacoes: 245,
      votosSim: 180,
      votosNao: 45,
      abstencoes: 10,
      ausencias: 10,
      percentualPresenca: 95.8,
      totalProposicoes: 32,
      proposicoesAprovadas: 8,
      totalDespesas: 156000.00,
      mediaGastoMensal: 13000.00,
    },
  },
  '2': {
    politico: {
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
    estatisticas: {
      totalVotacoes: 189,
      votosSim: 120,
      votosNao: 50,
      abstencoes: 15,
      ausencias: 4,
      percentualPresenca: 97.8,
      totalProposicoes: 45,
      proposicoesAprovadas: 12,
      totalDespesas: 198000.00,
      mediaGastoMensal: 16500.00,
    },
  },
};

export function Comparar() {
  const { politicosSelecionados, removerPolitico, limparSelecao } = useComparacaoStore();

  const politicosComDados = politicosSelecionados
    .map((id) => MOCK_POLITICOS_DETALHES[id])
    .filter(Boolean);

  if (politicosSelecionados.length === 0) {
    return (
      <div className="min-h-screen py-16">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <div className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-background-secondary mb-6">
              <Scale className="w-10 h-10 text-content-muted" />
            </div>
            <h1 className="text-3xl font-bold text-content-primary mb-4">
              Comparar Políticos
            </h1>
            <p className="text-content-secondary mb-8 max-w-md mx-auto">
              Selecione até 4 políticos na página de listagem para comparar suas atuações lado a lado.
            </p>
            <Link to="/politicos">
              <Button>
                <Users className="w-5 h-5" />
                Selecionar Políticos
              </Button>
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
          <div>
            <h1 className="text-3xl font-bold text-content-primary mb-2">
              Comparar Políticos
            </h1>
            <p className="text-content-secondary">
              {politicosSelecionados.length} político{politicosSelecionados.length !== 1 ? 's' : ''} selecionado{politicosSelecionados.length !== 1 ? 's' : ''}
            </p>
          </div>
          <div className="flex gap-3">
            <Link to="/politicos">
              <Button variant="secondary">
                <Users className="w-4 h-4" />
                Adicionar mais
              </Button>
            </Link>
            <Button variant="ghost" onClick={limparSelecao}>
              Limpar seleção
            </Button>
          </div>
        </div>

        {/* Cards dos políticos */}
        <div className={cn(
          'grid gap-6 mb-8',
          politicosComDados.length === 1 && 'grid-cols-1 max-w-md',
          politicosComDados.length === 2 && 'grid-cols-1 md:grid-cols-2',
          politicosComDados.length === 3 && 'grid-cols-1 md:grid-cols-3',
          politicosComDados.length >= 4 && 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4'
        )}>
          {politicosComDados.map(({ politico, estatisticas }) => (
            <Card key={politico.id} className="relative">
              {/* Botão remover */}
              <button
                onClick={() => removerPolitico(politico.id)}
                className="absolute top-3 right-3 p-1.5 rounded-lg bg-background-secondary hover:bg-accent-danger/20 text-content-muted hover:text-accent-danger transition-colors"
              >
                <X className="w-4 h-4" />
              </button>

              <CardContent className="pt-6">
                {/* Foto e nome */}
                <div className="text-center mb-6">
                  <Avatar
                    src={politico.fotoUrl}
                    alt={politico.nome}
                    size="lg"
                    className="mx-auto mb-3"
                  />
                  <h3 className="font-semibold text-content-primary">
                    {politico.nome}
                  </h3>
                  <p className="text-sm text-content-secondary">
                    {politico.partido.sigla} - {politico.cargoAtual.estado}
                  </p>
                  <Badge variant="info" className="mt-2">
                    {getCargoLabel(politico.cargoAtual.tipo)}
                  </Badge>
                </div>

                {/* Estatísticas */}
                <div className="space-y-4">
                  <StatItem
                    icon={TrendingUp}
                    label="Presença"
                    value={formatPercentage(estatisticas.percentualPresenca)}
                    color="text-accent-success"
                  />
                  <StatItem
                    icon={Vote}
                    label="Votações"
                    value={estatisticas.totalVotacoes.toString()}
                    color="text-accent-primary"
                  />
                  <StatItem
                    icon={FileText}
                    label="Proposições"
                    value={`${estatisticas.proposicoesAprovadas}/${estatisticas.totalProposicoes}`}
                    color="text-accent-secondary"
                  />
                  <StatItem
                    icon={Receipt}
                    label="Gasto médio/mês"
                    value={formatCurrency(estatisticas.mediaGastoMensal)}
                    color="text-accent-warning"
                  />
                </div>

                {/* Link para detalhes */}
                <div className="mt-6 pt-4 border-t border-border">
                  <Link
                    to={`/politicos/${politico.id}`}
                    className="flex items-center justify-center gap-2 text-sm text-accent-primary hover:text-accent-primary/80 transition-colors"
                  >
                    Ver perfil completo
                    <ArrowRight className="w-4 h-4" />
                  </Link>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Tabela comparativa detalhada */}
        {politicosComDados.length >= 2 && (
          <Card>
            <CardHeader>
              <CardTitle>Comparativo Detalhado</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-border">
                      <th className="text-left py-3 px-4 text-sm font-medium text-content-secondary">
                        Métrica
                      </th>
                      {politicosComDados.map(({ politico }) => (
                        <th key={politico.id} className="text-center py-3 px-4 text-sm font-medium text-content-primary">
                          {politico.nome}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    <CompareRow
                      label="Partido"
                      values={politicosComDados.map(({ politico }) => politico.partido.sigla)}
                    />
                    <CompareRow
                      label="Cargo"
                      values={politicosComDados.map(({ politico }) => getCargoLabel(politico.cargoAtual.tipo))}
                    />
                    <CompareRow
                      label="Presença"
                      values={politicosComDados.map(({ estatisticas }) => formatPercentage(estatisticas.percentualPresenca))}
                      highlight="max"
                    />
                    <CompareRow
                      label="Total Votações"
                      values={politicosComDados.map(({ estatisticas }) => estatisticas.totalVotacoes.toString())}
                      highlight="max"
                    />
                    <CompareRow
                      label="Votos Sim"
                      values={politicosComDados.map(({ estatisticas }) => estatisticas.votosSim.toString())}
                    />
                    <CompareRow
                      label="Votos Não"
                      values={politicosComDados.map(({ estatisticas }) => estatisticas.votosNao.toString())}
                    />
                    <CompareRow
                      label="Abstenções"
                      values={politicosComDados.map(({ estatisticas }) => estatisticas.abstencoes.toString())}
                    />
                    <CompareRow
                      label="Proposições"
                      values={politicosComDados.map(({ estatisticas }) => estatisticas.totalProposicoes.toString())}
                      highlight="max"
                    />
                    <CompareRow
                      label="Proposições Aprovadas"
                      values={politicosComDados.map(({ estatisticas }) => estatisticas.proposicoesAprovadas.toString())}
                      highlight="max"
                    />
                    <CompareRow
                      label="Salário Bruto"
                      values={politicosComDados.map(({ politico }) => formatCurrency(politico.salarioBruto))}
                    />
                    <CompareRow
                      label="Gasto Médio Mensal"
                      values={politicosComDados.map(({ estatisticas }) => formatCurrency(estatisticas.mediaGastoMensal))}
                      highlight="min"
                    />
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}

interface StatItemProps {
  icon: React.ElementType;
  label: string;
  value: string;
  color: string;
}

function StatItem({ icon: Icon, label, value, color }: StatItemProps) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-2">
        <Icon className={cn('w-4 h-4', color)} />
        <span className="text-sm text-content-secondary">{label}</span>
      </div>
      <span className={cn('text-sm font-semibold font-mono', color)}>
        {value}
      </span>
    </div>
  );
}

interface CompareRowProps {
  label: string;
  values: string[];
  highlight?: 'max' | 'min';
}

function CompareRow({ label, values, highlight }: CompareRowProps) {
  // Encontrar índice do melhor valor se houver highlight
  let bestIndex = -1;
  if (highlight) {
    const numericValues = values.map((v) => parseFloat(v.replace(/[^0-9.,]/g, '').replace(',', '.')));
    if (highlight === 'max') {
      bestIndex = numericValues.indexOf(Math.max(...numericValues));
    } else {
      bestIndex = numericValues.indexOf(Math.min(...numericValues));
    }
  }

  return (
    <tr className="border-b border-border last:border-0">
      <td className="py-3 px-4 text-sm text-content-secondary">{label}</td>
      {values.map((value, index) => (
        <td
          key={index}
          className={cn(
            'py-3 px-4 text-center text-sm font-mono',
            index === bestIndex ? 'text-accent-success font-semibold' : 'text-content-primary'
          )}
        >
          {value}
        </td>
      ))}
    </tr>
  );
}

