// Enums
export type Cargo =
  | 'DEPUTADO_FEDERAL'
  | 'SENADOR'
  | 'DEPUTADO_ESTADUAL'
  | 'DEPUTADO_DISTRITAL'
  | 'VEREADOR'
  | 'PREFEITO'
  | 'GOVERNADOR'
  | 'PRESIDENTE';

export type Esfera = 'FEDERAL' | 'ESTADUAL' | 'MUNICIPAL';

export type Genero = 'M' | 'F' | 'OUTRO';

export type TipoVoto = 'SIM' | 'NAO' | 'ABSTENCAO' | 'AUSENTE' | 'OBSTRUCAO';

export type SituacaoProposicao =
  | 'EM_TRAMITACAO'
  | 'APROVADA'
  | 'REJEITADA'
  | 'ARQUIVADA'
  | 'RETIRADA';

export type Regiao = 'NORTE' | 'NORDESTE' | 'CENTRO_OESTE' | 'SUDESTE' | 'SUL';

// Interfaces principais
export interface Partido {
  sigla: string;
  nome: string;
  cor: string;
}

export interface CargoAtual {
  tipo: Cargo;
  esfera: Esfera;
  estado: string;
  municipio?: string;
  dataInicio: string;
  dataFim?: string;
  emExercicio: boolean;
}

export interface Politico {
  id: string;
  nome: string;
  nomeCivil: string;
  fotoUrl: string;
  dataNascimento: string;
  genero: Genero;
  partido: Partido;
  cargoAtual: CargoAtual;
  historicoCargos: CargoAtual[];
  contato: {
    email?: string;
    telefone?: string;
    gabinete?: string;
  };
  redesSociais: {
    twitter?: string;
    instagram?: string;
    facebook?: string;
  };
  salarioBruto: number;
  salarioLiquido: number;
  createdAt: string;
  updatedAt: string;
}

export interface Votacao {
  id: string;
  politicoId: string;
  proposicaoId: string;
  proposicao?: Proposicao;
  voto: TipoVoto;
  data: string;
  sessao: string;
}

export interface Proposicao {
  id: string;
  tipo: string;
  numero: string;
  ano: number;
  ementa: string;
  autorId: string;
  autor?: Politico;
  coautoresIds: string[];
  situacao: SituacaoProposicao;
  tema: string[];
  tramitacao: TramitacaoItem[];
}

export interface TramitacaoItem {
  data: string;
  descricao: string;
  orgao: string;
}

export interface Despesa {
  id: string;
  politicoId: string;
  tipo: string;
  descricao: string;
  fornecedor: string;
  cnpjFornecedor: string;
  valor: number;
  data: string;
  mesReferencia: number;
  anoReferencia: number;
  documentoUrl?: string;
}

export interface Presenca {
  id: string;
  politicoId: string;
  data: string;
  tipoSessao: string;
  presente: boolean;
}

// Estatísticas agregadas
export interface EstatisticasPolitico {
  totalVotacoes: number;
  votosSim: number;
  votosNao: number;
  abstencoes: number;
  ausencias: number;
  percentualPresenca: number;
  totalProposicoes: number;
  proposicoesAprovadas: number;
  totalDespesas: number;
  mediaGastoMensal: number;
}

// Filtros
export interface FiltrosPoliticos {
  nome?: string;
  partido?: string[];
  cargo?: Cargo[];
  esfera?: Esfera[];
  estado?: string[];
  municipio?: string[];
  regiao?: Regiao[];
  emExercicio?: boolean;
  primeiroMandato?: boolean;
  reeleito?: boolean;
  genero?: Genero[];
  idadeMinima?: number;
  idadeMaxima?: number;
  presencaMinima?: number;
  proposicoesMinima?: number;
  gastoMensalMaximo?: number;
  ordenarPor?: 'nome' | 'partido' | 'presenca' | 'proposicoes' | 'gastos';
  ordem?: 'asc' | 'desc';
  pagina?: number;
  porPagina?: number;
}

// Resposta paginada
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  pagina: number;
  porPagina: number;
  totalPaginas: number;
}

// Estado da comparação
export interface ComparacaoState {
  politicosSelecionados: string[];
  maxPoliticos: number;
}

