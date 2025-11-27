import axios from 'axios';
import type {
  Politico,
  Votacao,
  Proposicao,
  Despesa,
  Presenca,
  EstatisticasPolitico,
  FiltrosPoliticos,
  PaginatedResponse,
  Partido,
} from '../types';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para tratamento de erros
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// Políticos
export const politicosApi = {
  listar: async (filtros?: FiltrosPoliticos): Promise<PaginatedResponse<Politico>> => {
    const { data } = await api.get('/politicos', { params: filtros });
    return data;
  },

  buscar: async (id: string): Promise<Politico> => {
    const { data } = await api.get(`/politicos/${id}`);
    return data;
  },

  buscarEstatisticas: async (id: string): Promise<EstatisticasPolitico> => {
    const { data } = await api.get(`/politicos/${id}/estatisticas`);
    return data;
  },

  buscarVotacoes: async (
    id: string,
    pagina = 1,
    porPagina = 20
  ): Promise<PaginatedResponse<Votacao>> => {
    const { data } = await api.get(`/politicos/${id}/votacoes`, {
      params: { pagina, porPagina },
    });
    return data;
  },

  buscarDespesas: async (
    id: string,
    ano?: number,
    mes?: number,
    pagina = 1,
    porPagina = 20
  ): Promise<PaginatedResponse<Despesa>> => {
    const { data } = await api.get(`/politicos/${id}/despesas`, {
      params: { ano, mes, pagina, porPagina },
    });
    return data;
  },

  buscarProposicoes: async (
    id: string,
    pagina = 1,
    porPagina = 20
  ): Promise<PaginatedResponse<Proposicao>> => {
    const { data } = await api.get(`/politicos/${id}/proposicoes`, {
      params: { pagina, porPagina },
    });
    return data;
  },

  buscarPresencas: async (
    id: string,
    ano?: number,
    mes?: number,
    pagina = 1,
    porPagina = 50
  ): Promise<PaginatedResponse<Presenca>> => {
    const { data } = await api.get(`/politicos/${id}/presencas`, {
      params: { ano, mes, pagina, porPagina },
    });
    return data;
  },

  comparar: async (ids: string[]): Promise<{
    politicos: Politico[];
    estatisticas: Record<string, EstatisticasPolitico>;
  }> => {
    const { data } = await api.get('/politicos/comparar', {
      params: { ids: ids.join(',') },
    });
    return data;
  },
};

// Filtros
export const filtrosApi = {
  partidos: async (): Promise<Partido[]> => {
    const { data } = await api.get('/filtros/partidos');
    return data;
  },

  estados: async (): Promise<{ sigla: string; nome: string }[]> => {
    const { data } = await api.get('/filtros/estados');
    return data;
  },

  cargos: async (): Promise<{ valor: string; label: string }[]> => {
    const { data } = await api.get('/filtros/cargos');
    return data;
  },
};

// Estatísticas gerais
export const estatisticasApi = {
  geral: async (): Promise<{
    totalPoliticos: number;
    totalVotacoes: number;
    totalProposicoes: number;
    totalDespesas: number;
  }> => {
    const { data } = await api.get('/estatisticas/geral');
    return data;
  },

  ranking: async (
    tipo: 'presenca' | 'proposicoes' | 'gastos',
    limite = 10
  ): Promise<{ politico: Politico; valor: number }[]> => {
    const { data } = await api.get('/estatisticas/ranking', {
      params: { tipo, limite },
    });
    return data;
  },
};

// Busca
export const buscaApi = {
  buscar: async (query: string, limite = 10): Promise<Politico[]> => {
    const { data } = await api.get('/busca', {
      params: { q: query, limite },
    });
    return data;
  },
};

export default api;

