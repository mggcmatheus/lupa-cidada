# 📋 Implementação de Sincronização para Outros Cargos

Este documento descreve o que é necessário para implementar a sincronização de dados para os cargos que ainda não estão implementados.

## 🏗️ Estrutura Atual

O sistema já possui sincronização para:
- ✅ **Deputados Federais** - API da Câmara dos Deputados
- ✅ **Senadores** - API do Senado Federal

## 📦 Cargos Pendentes

### 1. Deputados Estaduais

**Fonte de Dados:**
- Cada Assembleia Legislativa tem sua própria API/portal
- Não há uma API unificada nacional
- Algumas assembleias têm APIs REST, outras apenas dados em CSV/PDF

**O que é necessário:**
- [ ] Criar pacote `backend/internal/sync/assembleias/`
- [ ] Implementar sincronizador para cada estado (ou estados prioritários)
- [ ] Mapear estruturas de dados diferentes de cada assembleia
- [ ] Criar tipos Go para cada API de assembleia
- [ ] Implementar rate limiting específico (cada assembleia tem limites diferentes)

**Exemplos de APIs disponíveis:**
- ALESP (SP): https://www.al.sp.gov.br/dados-abertos/
- ALERJ (RJ): https://www.alerj.rj.gov.br/
- ALMG (MG): https://www.almg.gov.br/
- ALCE (CE): https://www.al.ce.gov.br/

**Complexidade:** 🔴 Alta (27 estados = 27 APIs diferentes)

---

### 2. Deputados Distritais

**Fonte de Dados:**
- Câmara Legislativa do Distrito Federal
- API: https://www.cl.df.gov.br/

**O que é necessário:**
- [ ] Criar pacote `backend/internal/sync/distrital/`
- [ ] Implementar sincronizador similar ao da Câmara/Senado
- [ ] Mapear estrutura de dados da CLDF
- [ ] Criar tipos Go para a API

**Complexidade:** 🟢 Baixa (1 fonte única)

---

### 3. Vereadores

**Fonte de Dados:**
- Cada Câmara Municipal tem sua própria estrutura
- TSE fornece dados eleitorais (candidatos eleitos)
- Não há API unificada

**O que é necessário:**
- [ ] Decidir estratégia:
  - Opção A: Sincronizar apenas cidades grandes (prioritárias)
  - Opção B: Usar dados do TSE (eleições) + scraping de câmaras
- [ ] Criar pacote `backend/internal/sync/vereadores/`
- [ ] Integrar com API do TSE para dados eleitorais
- [ ] Implementar scraping ou integração com APIs municipais

**Fontes possíveis:**
- TSE: https://dadosabertos.tse.jus.br/ (dados eleitorais)
- APIs municipais (varia por cidade)

**Complexidade:** 🔴 Muito Alta (5570+ municípios)

---

### 4. Prefeitos

**Fonte de Dados:**
- TSE (dados eleitorais)
- Portal da Transparência
- Sites das prefeituras

**O que é necessário:**
- [ ] Criar pacote `backend/internal/sync/prefeitos/`
- [ ] Integrar com API do TSE para dados eleitorais
- [ ] Buscar dados do Portal da Transparência
- [ ] Implementar scraping de sites de prefeituras (se necessário)

**Fontes:**
- TSE: https://dadosabertos.tse.jus.br/
- Portal da Transparência: https://portaldatransparencia.gov.br/

**Complexidade:** 🟡 Média-Alta (5570+ municípios, mas dados mais centralizados)

---

### 5. Governadores

**Fonte de Dados:**
- TSE (dados eleitorais)
- Portal da Transparência
- Sites dos governos estaduais

**O que é necessário:**
- [ ] Criar pacote `backend/internal/sync/governadores/`
- [ ] Integrar com API do TSE
- [ ] Buscar dados do Portal da Transparência
- [ ] Implementar sincronização de dados dos governos estaduais

**Complexidade:** 🟢 Baixa (27 estados, dados centralizados)

---

### 6. Presidente

**Fonte de Dados:**
- TSE (dados eleitorais)
- Portal da Transparência
- Site da Presidência

**O que é necessário:**
- [ ] Criar pacote `backend/internal/sync/presidente/`
- [ ] Integrar com API do TSE
- [ ] Buscar dados do Portal da Transparência
- [ ] Implementar sincronização manual ou via scraping

**Complexidade:** 🟢 Muito Baixa (1 pessoa, dados centralizados)

---

## 🛠️ Padrão de Implementação

Cada sincronizador deve seguir o padrão existente:

```go
// Estrutura básica
type XxxSync struct {
    client *sync.HTTPClient
    db     *mongo.Database
}

func NewXxxSync(db *mongo.Database) *XxxSync {
    return &XxxSync{
        client: sync.NewHTTPClient(requestsPerSecond),
        db:     db,
    }
}

func (s *XxxSync) SyncXxx(ctx context.Context) error {
    // 1. Buscar lista de políticos
    // 2. Para cada político, buscar detalhes
    // 3. Verificar se já existe no banco
    // 4. Criar/atualizar registro
    // 5. Gerenciar histórico de cargos
}
```

## 📝 Checklist de Implementação

Para cada novo cargo, é necessário:

1. **Análise de Dados:**
   - [ ] Identificar fonte de dados oficial
   - [ ] Verificar disponibilidade de API
   - [ ] Analisar estrutura de dados
   - [ ] Verificar rate limits e políticas de uso

2. **Implementação:**
   - [ ] Criar pacote de sincronização
   - [ ] Definir tipos Go para estruturas de dados
   - [ ] Implementar função de sincronização principal
   - [ ] Implementar busca de político existente
   - [ ] Implementar mapeamento de dados
   - [ ] Adicionar ao `cmd/sync/main.go`

3. **Testes:**
   - [ ] Testar sincronização com dados reais
   - [ ] Validar mapeamento de campos
   - [ ] Verificar histórico de cargos
   - [ ] Testar atualização de políticos existentes

4. **Documentação:**
   - [ ] Documentar fonte de dados
   - [ ] Documentar limitações conhecidas
   - [ ] Atualizar README

## 🎯 Priorização Sugerida

1. **Fase 1 (Mais fácil):**
   - Presidente
   - Governadores
   - Deputados Distritais

2. **Fase 2 (Média complexidade):**
   - Prefeitos (cidades grandes primeiro)

3. **Fase 3 (Alta complexidade):**
   - Deputados Estaduais (estados prioritários)
   - Vereadores (cidades grandes primeiro)

## 🔗 Links Úteis

- **TSE Dados Abertos:** https://dadosabertos.tse.jus.br/
- **Portal da Transparência:** https://portaldatransparencia.gov.br/
- **Câmara dos Deputados:** https://dadosabertos.camara.leg.br/
- **Senado Federal:** https://legis.senado.leg.br/dadosabertos
- **IBGE:** https://www.ibge.gov.br/ (dados demográficos)

## 📌 Notas Importantes

- **LGPD:** Garantir conformidade com Lei Geral de Proteção de Dados
- **Rate Limiting:** Respeitar limites de cada API
- **Atualização:** Implementar rotinas de atualização periódica
- **Fallback:** Ter estratégias de fallback quando APIs estiverem indisponíveis

