import { create } from 'zustand';
import type { FiltrosPoliticos } from '../types';

interface FiltrosStore {
  filtros: FiltrosPoliticos;
  
  // Actions
  setFiltro: <K extends keyof FiltrosPoliticos>(key: K, value: FiltrosPoliticos[K]) => void;
  setFiltros: (filtros: Partial<FiltrosPoliticos>) => void;
  limparFiltros: () => void;
  resetarPaginacao: () => void;
}

const filtrosIniciais: FiltrosPoliticos = {
  pagina: 1,
  porPagina: 12,
  ordenarPor: 'nome',
  ordem: 'asc',
};

export const useFiltrosStore = create<FiltrosStore>((set) => ({
  filtros: filtrosIniciais,

  setFiltro: (key, value) => {
    set((state) => ({
      filtros: {
        ...state.filtros,
        [key]: value,
        // Resetar paginação quando filtros mudam (exceto paginação)
        ...(key !== 'pagina' && key !== 'porPagina' ? { pagina: 1 } : {}),
      },
    }));
  },

  setFiltros: (novosFiltros) => {
    set((state) => ({
      filtros: {
        ...state.filtros,
        ...novosFiltros,
        pagina: 1, // Resetar paginação
      },
    }));
  },

  limparFiltros: () => {
    set({ filtros: filtrosIniciais });
  },

  resetarPaginacao: () => {
    set((state) => ({
      filtros: {
        ...state.filtros,
        pagina: 1,
      },
    }));
  },
}));

