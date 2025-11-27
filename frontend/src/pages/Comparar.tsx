import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { Scale, X, ArrowRight, Users, Vote, Receipt, FileText, TrendingUp, Loader2 } from 'lucide-react';
import { politicosApi } from '../services/api';
import { useComparacaoStore } from '../stores/useComparacaoStore';
import { Button } from '../components/ui/Button';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';
import { Avatar } from '../components/ui/Avatar';
import { cn, formatCurrency, formatPercentage, getCargoLabel } from '../lib/utils';
import type { EstatisticasPolitico } from '../types';

export function Comparar() {
  const { politicosSelecionados, removerPolitico, limparSelecao } = useComparacaoStore();

  const { data, isLoading } = useQuery({
    queryKey: ['comparar', politicosSelecionados],
    queryFn: () => politicosApi.comparar(politicosSelecionados),
    enabled: politicosSelecionados.length >= 1,
  });

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

  if (isLoading) {
    return (
      <div className="min-h-screen py-16 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="w-8 h-8 animate-spin text-accent-primary mx-auto mb-4" />
          <p className="text-content-secondary">Carregando comparação...</p>
        </div>
      </div>
    );
  }

  const politicos = data?.politicos || [];
  const estatisticas = (data?.estatisticas || {}) as Record<string, EstatisticasPolitico>;

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
          politicos.length === 1 && 'grid-cols-1 max-w-md',
          politicos.length === 2 && 'grid-cols-1 md:grid-cols-2',
          politicos.length === 3 && 'grid-cols-1 md:grid-cols-3',
          politicos.length >= 4 && 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4'
        )}>
          {politicos.map((politico) => {
            const stats = estatisticas[politico.id];
            return (
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
                  {stats && (
                    <div className="space-y-4">
                      <StatItem
                        icon={TrendingUp}
                        label="Presença"
                        value={formatPercentage(stats.percentualPresenca)}
                        color="text-accent-success"
                      />
                      <StatItem
                        icon={Vote}
                        label="Votações"
                        value={stats.totalVotacoes.toString()}
                        color="text-accent-primary"
                      />
                      <StatItem
                        icon={FileText}
                        label="Proposições"
                        value={`${stats.proposicoesAprovadas}/${stats.totalProposicoes}`}
                        color="text-accent-secondary"
                      />
                      <StatItem
                        icon={Receipt}
                        label="Gasto médio/mês"
                        value={formatCurrency(stats.mediaGastoMensal)}
                        color="text-accent-warning"
                      />
                    </div>
                  )}

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
            );
          })}
        </div>

        {/* Tabela comparativa detalhada */}
        {politicos.length >= 2 && (
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
                      {politicos.map((politico) => (
                        <th key={politico.id} className="text-center py-3 px-4 text-sm font-medium text-content-primary">
                          {politico.nome}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    <CompareRow
                      label="Partido"
                      values={politicos.map((p) => p.partido.sigla)}
                    />
                    <CompareRow
                      label="Cargo"
                      values={politicos.map((p) => getCargoLabel(p.cargoAtual.tipo))}
                    />
                    <CompareRow
                      label="Presença"
                      values={politicos.map((p) => formatPercentage(estatisticas[p.id]?.percentualPresenca || 0))}
                      highlight="max"
                    />
                    <CompareRow
                      label="Total Votações"
                      values={politicos.map((p) => (estatisticas[p.id]?.totalVotacoes || 0).toString())}
                      highlight="max"
                    />
                    <CompareRow
                      label="Votos Sim"
                      values={politicos.map((p) => (estatisticas[p.id]?.votosSim || 0).toString())}
                    />
                    <CompareRow
                      label="Votos Não"
                      values={politicos.map((p) => (estatisticas[p.id]?.votosNao || 0).toString())}
                    />
                    <CompareRow
                      label="Abstenções"
                      values={politicos.map((p) => (estatisticas[p.id]?.abstencoes || 0).toString())}
                    />
                    <CompareRow
                      label="Proposições"
                      values={politicos.map((p) => (estatisticas[p.id]?.totalProposicoes || 0).toString())}
                      highlight="max"
                    />
                    <CompareRow
                      label="Proposições Aprovadas"
                      values={politicos.map((p) => (estatisticas[p.id]?.proposicoesAprovadas || 0).toString())}
                      highlight="max"
                    />
                    <CompareRow
                      label="Salário Bruto"
                      values={politicos.map((p) => formatCurrency(p.salarioBruto))}
                    />
                    <CompareRow
                      label="Gasto Médio Mensal"
                      values={politicos.map((p) => formatCurrency(estatisticas[p.id]?.mediaGastoMensal || 0))}
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
