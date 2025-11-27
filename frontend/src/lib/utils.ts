import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatCurrency(value: number): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(value);
}

export function formatNumber(value: number): string {
  return new Intl.NumberFormat('pt-BR').format(value);
}

export function formatPercentage(value: number): string {
  return `${value.toFixed(1)}%`;
}

export function formatDate(date: string): string {
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(new Date(date));
}

export function formatDateLong(date: string): string {
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
  }).format(new Date(date));
}

export function getCargoLabel(cargo: string): string {
  const labels: Record<string, string> = {
    DEPUTADO_FEDERAL: 'Deputado Federal',
    SENADOR: 'Senador',
    DEPUTADO_ESTADUAL: 'Deputado Estadual',
    DEPUTADO_DISTRITAL: 'Deputado Distrital',
    VEREADOR: 'Vereador',
    PREFEITO: 'Prefeito',
    GOVERNADOR: 'Governador',
    PRESIDENTE: 'Presidente',
  };
  return labels[cargo] || cargo;
}

export function getEsferaLabel(esfera: string): string {
  const labels: Record<string, string> = {
    FEDERAL: 'Federal',
    ESTADUAL: 'Estadual',
    MUNICIPAL: 'Municipal',
  };
  return labels[esfera] || esfera;
}

export function getVotoLabel(voto: string): string {
  const labels: Record<string, string> = {
    SIM: 'Sim',
    NAO: 'Não',
    ABSTENCAO: 'Abstenção',
    AUSENTE: 'Ausente',
    OBSTRUCAO: 'Obstrução',
  };
  return labels[voto] || voto;
}

export function getVotoColor(voto: string): string {
  const colors: Record<string, string> = {
    SIM: 'text-accent-success',
    NAO: 'text-accent-danger',
    ABSTENCAO: 'text-accent-warning',
    AUSENTE: 'text-content-muted',
    OBSTRUCAO: 'text-accent-secondary',
  };
  return colors[voto] || 'text-content-secondary';
}

export function getSituacaoLabel(situacao: string): string {
  const labels: Record<string, string> = {
    EM_TRAMITACAO: 'Em Tramitação',
    APROVADA: 'Aprovada',
    REJEITADA: 'Rejeitada',
    ARQUIVADA: 'Arquivada',
    RETIRADA: 'Retirada',
  };
  return labels[situacao] || situacao;
}

export function getRegiaoByEstado(estado: string): string {
  const regioes: Record<string, string> = {
    AC: 'NORTE', AM: 'NORTE', AP: 'NORTE', PA: 'NORTE', RO: 'NORTE', RR: 'NORTE', TO: 'NORTE',
    AL: 'NORDESTE', BA: 'NORDESTE', CE: 'NORDESTE', MA: 'NORDESTE', PB: 'NORDESTE',
    PE: 'NORDESTE', PI: 'NORDESTE', RN: 'NORDESTE', SE: 'NORDESTE',
    DF: 'CENTRO_OESTE', GO: 'CENTRO_OESTE', MT: 'CENTRO_OESTE', MS: 'CENTRO_OESTE',
    ES: 'SUDESTE', MG: 'SUDESTE', RJ: 'SUDESTE', SP: 'SUDESTE',
    PR: 'SUL', RS: 'SUL', SC: 'SUL',
  };
  return regioes[estado] || 'DESCONHECIDA';
}

export function calculateAge(birthDate: string): number {
  const today = new Date();
  const birth = new Date(birthDate);
  let age = today.getFullYear() - birth.getFullYear();
  const monthDiff = today.getMonth() - birth.getMonth();
  
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
    age--;
  }
  
  return age;
}

export function truncateText(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength).trim() + '...';
}

export function debounce<T extends (...args: unknown[]) => unknown>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: ReturnType<typeof setTimeout>;
  
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}

