import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface ComparacaoStore {
  politicosSelecionados: string[];
  maxPoliticos: number;
  
  // Actions
  adicionarPolitico: (id: string) => void;
  removerPolitico: (id: string) => void;
  togglePolitico: (id: string) => void;
  limparSelecao: () => void;
  estaSelecionado: (id: string) => boolean;
  podeAdicionar: () => boolean;
}

export const useComparacaoStore = create<ComparacaoStore>()(
  persist(
    (set, get) => ({
      politicosSelecionados: [],
      maxPoliticos: 4,

      adicionarPolitico: (id: string) => {
        const { politicosSelecionados, maxPoliticos } = get();
        if (politicosSelecionados.length < maxPoliticos && !politicosSelecionados.includes(id)) {
          set({ politicosSelecionados: [...politicosSelecionados, id] });
        }
      },

      removerPolitico: (id: string) => {
        const { politicosSelecionados } = get();
        set({ politicosSelecionados: politicosSelecionados.filter((pid) => pid !== id) });
      },

      togglePolitico: (id: string) => {
        const { politicosSelecionados, adicionarPolitico, removerPolitico } = get();
        if (politicosSelecionados.includes(id)) {
          removerPolitico(id);
        } else {
          adicionarPolitico(id);
        }
      },

      limparSelecao: () => {
        set({ politicosSelecionados: [] });
      },

      estaSelecionado: (id: string) => {
        return get().politicosSelecionados.includes(id);
      },

      podeAdicionar: () => {
        const { politicosSelecionados, maxPoliticos } = get();
        return politicosSelecionados.length < maxPoliticos;
      },
    }),
    {
      name: 'lupa-comparacao-storage',
    }
  )
);

