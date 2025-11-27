import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { 
  ArrowLeft, 
  Mail, 
  Phone, 
  Building2, 
  MapPin, 
  Calendar,
  Vote,
  Receipt,
  FileText,
  TrendingUp,
  ExternalLink,
  Twitter,
  Instagram,
  Loader2
} from 'lucide-react';
import { politicosApi } from '../services/api';
import { Button } from '../components/ui/Button';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';
import { Avatar } from '../components/ui/Avatar';
import { 
  formatCurrency, 
  formatPercentage, 
  formatDate,
  getCargoLabel,
  getEsferaLabel,
  calculateAge
} from '../lib/utils';

export function PoliticoDetalhe() {
  const { id } = useParams<{ id: string }>();

  const { data: politico, isLoading: loadingPolitico } = useQuery({
    queryKey: ['politico', id],
    queryFn: () => politicosApi.buscar(id!),
    enabled: !!id,
  });

  const { data: estatisticas, isLoading: loadingStats } = useQuery({
    queryKey: ['politico-estatisticas', id],
    queryFn: () => politicosApi.buscarEstatisticas(id!),
    enabled: !!id,
  });

  const isLoading = loadingPolitico || loadingStats;

  if (isLoading) {
    return (
      <div className="min-h-screen py-16 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="w-8 h-8 animate-spin text-accent-primary mx-auto mb-4" />
          <p className="text-content-secondary">Carregando dados do político...</p>
        </div>
      </div>
    );
  }

  if (!politico) {
    return (
      <div className="min-h-screen py-16">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-3xl font-bold text-content-primary mb-4">
            Político não encontrado
          </h1>
          <p className="text-content-secondary mb-8">
            O político que você está procurando não existe ou foi removido.
          </p>
          <Link to="/politicos">
            <Button>
              <ArrowLeft className="w-5 h-5" />
              Voltar para a lista
            </Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen py-8">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Voltar */}
        <Link
          to="/politicos"
          className="inline-flex items-center gap-2 text-content-secondary hover:text-content-primary mb-6 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          Voltar para a lista
        </Link>

        {/* Header do político */}
        <div className="bg-background-card border border-border rounded-xl p-6 md:p-8 mb-8">
          <div className="flex flex-col md:flex-row gap-6">
            {/* Avatar */}
            <div className="flex-shrink-0">
              <Avatar
                src={politico.fotoUrl}
                alt={politico.nome}
                size="xl"
              />
            </div>

            {/* Info */}
            <div className="flex-1">
              <div className="flex flex-wrap items-start justify-between gap-4">
                <div>
                  <h1 className="text-2xl md:text-3xl font-bold text-content-primary mb-1">
                    {politico.nome}
                  </h1>
                  <p className="text-content-secondary mb-3">
                    {politico.nomeCivil}
                  </p>
                  
                  <div className="flex flex-wrap items-center gap-2 mb-4">
                    <Badge 
                      variant="info"
                      className="text-sm"
                      style={{ 
                        backgroundColor: `${politico.partido.cor}20`,
                        color: politico.partido.cor 
                      }}
                    >
                      {politico.partido.sigla}
                    </Badge>
                    <Badge variant="info">
                      <Building2 className="w-3 h-3" />
                      {getCargoLabel(politico.cargoAtual.tipo)}
                    </Badge>
                    <Badge variant={politico.cargoAtual.emExercicio ? 'success' : 'warning'}>
                      {politico.cargoAtual.emExercicio ? 'Em exercício' : 'Fora de exercício'}
                    </Badge>
                  </div>

                  <div className="flex flex-wrap items-center gap-4 text-sm text-content-secondary">
                    <span className="flex items-center gap-1">
                      <MapPin className="w-4 h-4" />
                      {politico.cargoAtual.municipio 
                        ? `${politico.cargoAtual.municipio}, ${politico.cargoAtual.estado}`
                        : politico.cargoAtual.estado
                      }
                    </span>
                    {politico.dataNascimento && (
                      <span className="flex items-center gap-1">
                        <Calendar className="w-4 h-4" />
                        {calculateAge(politico.dataNascimento)} anos
                      </span>
                    )}
                  </div>
                </div>

                {/* Salário */}
                <div className="bg-background-secondary rounded-xl p-4 text-center min-w-[140px]">
                  <p className="text-xs text-content-muted mb-1">Salário Bruto</p>
                  <p className="text-xl font-bold font-mono text-accent-primary">
                    {formatCurrency(politico.salarioBruto)}
                  </p>
                </div>
              </div>

              {/* Contato e Redes Sociais */}
              <div className="flex flex-wrap items-center gap-3 mt-4 pt-4 border-t border-border">
                {politico.contato?.email && (
                  <a
                    href={`mailto:${politico.contato.email}`}
                    className="flex items-center gap-1.5 text-sm text-content-secondary hover:text-accent-primary transition-colors"
                  >
                    <Mail className="w-4 h-4" />
                    {politico.contato.email}
                  </a>
                )}
                {politico.contato?.telefone && (
                  <span className="flex items-center gap-1.5 text-sm text-content-secondary">
                    <Phone className="w-4 h-4" />
                    {politico.contato.telefone}
                  </span>
                )}
                {politico.redesSociais?.twitter && (
                  <a
                    href={`https://twitter.com/${politico.redesSociais.twitter.replace('@', '')}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1.5 text-sm text-content-secondary hover:text-accent-primary transition-colors"
                  >
                    <Twitter className="w-4 h-4" />
                    {politico.redesSociais.twitter}
                  </a>
                )}
                {politico.redesSociais?.instagram && (
                  <a
                    href={`https://instagram.com/${politico.redesSociais.instagram.replace('@', '')}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1.5 text-sm text-content-secondary hover:text-accent-primary transition-colors"
                  >
                    <Instagram className="w-4 h-4" />
                    {politico.redesSociais.instagram}
                  </a>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Grid de estatísticas */}
        {estatisticas && (
          <>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
              <StatCard
                icon={TrendingUp}
                label="Presença"
                value={formatPercentage(estatisticas.percentualPresenca)}
                color="text-accent-success"
              />
              <StatCard
                icon={Vote}
                label="Votações"
                value={estatisticas.totalVotacoes.toString()}
                color="text-accent-primary"
              />
              <StatCard
                icon={FileText}
                label="Proposições"
                value={estatisticas.totalProposicoes.toString()}
                color="text-accent-secondary"
              />
              <StatCard
                icon={Receipt}
                label="Gasto Médio/Mês"
                value={formatCurrency(estatisticas.mediaGastoMensal)}
                color="text-accent-warning"
              />
            </div>

            {/* Detalhes em cards */}
            <div className="grid md:grid-cols-2 gap-6">
              {/* Votações */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Vote className="w-5 h-5 text-accent-primary" />
                    Votações
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Votos Sim</span>
                      <span className="font-mono font-semibold text-accent-success">
                        {estatisticas.votosSim}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Votos Não</span>
                      <span className="font-mono font-semibold text-accent-danger">
                        {estatisticas.votosNao}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Abstenções</span>
                      <span className="font-mono font-semibold text-accent-warning">
                        {estatisticas.abstencoes}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Ausências</span>
                      <span className="font-mono font-semibold text-content-muted">
                        {estatisticas.ausencias}
                      </span>
                    </div>
                    <div className="pt-3 border-t border-border">
                      <div className="flex justify-between items-center">
                        <span className="text-content-primary font-medium">Total</span>
                        <span className="font-mono font-bold text-content-primary">
                          {estatisticas.totalVotacoes}
                        </span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Proposições */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <FileText className="w-5 h-5 text-accent-secondary" />
                    Proposições
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Total apresentadas</span>
                      <span className="font-mono font-semibold text-content-primary">
                        {estatisticas.totalProposicoes}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Aprovadas</span>
                      <span className="font-mono font-semibold text-accent-success">
                        {estatisticas.proposicoesAprovadas}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Taxa de aprovação</span>
                      <span className="font-mono font-semibold text-accent-primary">
                        {estatisticas.totalProposicoes > 0 
                          ? formatPercentage((estatisticas.proposicoesAprovadas / estatisticas.totalProposicoes) * 100)
                          : '0%'
                        }
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Despesas */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Receipt className="w-5 h-5 text-accent-warning" />
                    Despesas
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Total acumulado</span>
                      <span className="font-mono font-semibold text-content-primary">
                        {formatCurrency(estatisticas.totalDespesas)}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Média mensal</span>
                      <span className="font-mono font-semibold text-accent-warning">
                        {formatCurrency(estatisticas.mediaGastoMensal)}
                      </span>
                    </div>
                  </div>
                  <div className="mt-4 pt-4 border-t border-border">
                    <Button variant="secondary" size="sm" className="w-full">
                      <ExternalLink className="w-4 h-4" />
                      Ver detalhes das despesas
                    </Button>
                  </div>
                </CardContent>
              </Card>

              {/* Informações do mandato */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Building2 className="w-5 h-5 text-accent-primary" />
                    Mandato Atual
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Cargo</span>
                      <span className="text-content-primary">
                        {getCargoLabel(politico.cargoAtual.tipo)}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Esfera</span>
                      <span className="text-content-primary">
                        {getEsferaLabel(politico.cargoAtual.esfera)}
                      </span>
                    </div>
                    <div className="flex justify-between items-center">
                      <span className="text-content-secondary">Início do mandato</span>
                      <span className="text-content-primary">
                        {formatDate(politico.cargoAtual.dataInicio)}
                      </span>
                    </div>
                    {politico.contato?.gabinete && (
                      <div className="flex justify-between items-center">
                        <span className="text-content-secondary">Gabinete</span>
                        <span className="text-content-primary">
                          {politico.contato.gabinete}
                        </span>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

interface StatCardProps {
  icon: React.ElementType;
  label: string;
  value: string;
  color: string;
}

function StatCard({ icon: Icon, label, value, color }: StatCardProps) {
  return (
    <div className="bg-background-card border border-border rounded-xl p-4 text-center">
      <Icon className={`w-6 h-6 mx-auto mb-2 ${color}`} />
      <p className={`text-xl font-bold font-mono ${color}`}>{value}</p>
      <p className="text-xs text-content-muted mt-1">{label}</p>
    </div>
  );
}
