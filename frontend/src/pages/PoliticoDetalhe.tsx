import { useParams, Link } from 'react-router-dom';
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
  Instagram
} from 'lucide-react';
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
import type { Politico, EstatisticasPolitico } from '../types';

// Mock data para demonstração
const MOCK_POLITICOS: Record<string, { politico: Politico; estatisticas: EstatisticasPolitico }> = {
  '1': {
    politico: {
      id: '1',
      nome: 'João Silva',
      nomeCivil: 'João Pedro da Silva Santos',
      fotoUrl: '',
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
      contato: { 
        email: 'joao.silva@camara.leg.br',
        telefone: '(61) 3215-5000',
        gabinete: 'Anexo IV, Gabinete 300'
      },
      redesSociais: {
        twitter: '@joaosilva',
        instagram: '@dep.joaosilva'
      },
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
      fotoUrl: '',
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
      contato: { 
        email: 'maria.santos@senado.leg.br',
        telefone: '(61) 3303-1000',
        gabinete: 'Anexo II, Gabinete 15'
      },
      redesSociais: {
        twitter: '@mariasantos',
        instagram: '@sen.mariasantos'
      },
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
  '3': {
    politico: {
      id: '3',
      nome: 'Pedro Costa',
      nomeCivil: 'Pedro Henrique Costa Lima',
      fotoUrl: '',
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
      contato: { email: 'pedro.costa@almg.gov.br' },
      redesSociais: {},
      salarioBruto: 25322.25,
      salarioLiquido: 18000.00,
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
    },
    estatisticas: {
      totalVotacoes: 156,
      votosSim: 100,
      votosNao: 40,
      abstencoes: 8,
      ausencias: 8,
      percentualPresenca: 94.9,
      totalProposicoes: 18,
      proposicoesAprovadas: 5,
      totalDespesas: 89000.00,
      mediaGastoMensal: 7416.67,
    },
  },
  '4': {
    politico: {
      id: '4',
      nome: 'Ana Oliveira',
      nomeCivil: 'Ana Paula Oliveira Souza',
      fotoUrl: '',
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
      contato: { email: 'ana.oliveira@saopaulo.sp.leg.br' },
      redesSociais: { twitter: '@anaoliveira' },
      salarioBruto: 18991.68,
      salarioLiquido: 14000.00,
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
    },
    estatisticas: {
      totalVotacoes: 312,
      votosSim: 200,
      votosNao: 80,
      abstencoes: 20,
      ausencias: 12,
      percentualPresenca: 96.2,
      totalProposicoes: 67,
      proposicoesAprovadas: 23,
      totalDespesas: 45000.00,
      mediaGastoMensal: 3750.00,
    },
  },
  '5': {
    politico: {
      id: '5',
      nome: 'Carlos Ferreira',
      nomeCivil: 'Carlos Alberto Ferreira',
      fotoUrl: '',
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
      contato: { email: 'gabinete@ba.gov.br' },
      redesSociais: {},
      salarioBruto: 35462.22,
      salarioLiquido: 26000.00,
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
    },
    estatisticas: {
      totalVotacoes: 0,
      votosSim: 0,
      votosNao: 0,
      abstencoes: 0,
      ausencias: 0,
      percentualPresenca: 100,
      totalProposicoes: 156,
      proposicoesAprovadas: 89,
      totalDespesas: 0,
      mediaGastoMensal: 0,
    },
  },
  '6': {
    politico: {
      id: '6',
      nome: 'Fernanda Lima',
      nomeCivil: 'Fernanda Cristina Lima',
      fotoUrl: '',
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
      contato: { 
        email: 'fernanda.lima@camara.leg.br',
        telefone: '(61) 3215-5001'
      },
      redesSociais: { instagram: '@fernanda.lima' },
      salarioBruto: 33763.00,
      salarioLiquido: 25000.00,
      createdAt: '2024-01-01',
      updatedAt: '2024-01-01',
    },
    estatisticas: {
      totalVotacoes: 245,
      votosSim: 150,
      votosNao: 70,
      abstencoes: 15,
      ausencias: 10,
      percentualPresenca: 95.9,
      totalProposicoes: 28,
      proposicoesAprovadas: 6,
      totalDespesas: 120000.00,
      mediaGastoMensal: 10000.00,
    },
  },
};

export function PoliticoDetalhe() {
  const { id } = useParams<{ id: string }>();
  
  const data = id ? MOCK_POLITICOS[id] : null;

  if (!data) {
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

  const { politico, estatisticas } = data;

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
                    <span className="flex items-center gap-1">
                      <Calendar className="w-4 h-4" />
                      {calculateAge(politico.dataNascimento)} anos
                    </span>
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
                {politico.contato.email && (
                  <a
                    href={`mailto:${politico.contato.email}`}
                    className="flex items-center gap-1.5 text-sm text-content-secondary hover:text-accent-primary transition-colors"
                  >
                    <Mail className="w-4 h-4" />
                    {politico.contato.email}
                  </a>
                )}
                {politico.contato.telefone && (
                  <span className="flex items-center gap-1.5 text-sm text-content-secondary">
                    <Phone className="w-4 h-4" />
                    {politico.contato.telefone}
                  </span>
                )}
                {politico.redesSociais.twitter && (
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
                {politico.redesSociais.instagram && (
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
                {politico.contato.gabinete && (
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

