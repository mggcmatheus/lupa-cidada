// Script de inicialização do MongoDB
// Executado automaticamente quando o container sobe pela primeira vez

db = db.getSiblingDB('lupa_cidada');

// Criar coleções
db.createCollection('politicos');
db.createCollection('votacoes');
db.createCollection('proposicoes');
db.createCollection('despesas');
db.createCollection('presencas');
db.createCollection('partidos');

// Índices para políticos
db.politicos.createIndex({ "nome": "text", "nome_civil": "text" });
db.politicos.createIndex({ "partido.sigla": 1 });
db.politicos.createIndex({ "cargo_atual.tipo": 1 });
db.politicos.createIndex({ "cargo_atual.esfera": 1 });
db.politicos.createIndex({ "cargo_atual.estado": 1 });
db.politicos.createIndex({ "cargo_atual.em_exercicio": 1 });
db.politicos.createIndex({ "genero": 1 });
db.politicos.createIndex({ "created_at": -1 });

// Índices para votações
db.votacoes.createIndex({ "politico_id": 1 });
db.votacoes.createIndex({ "proposicao_id": 1 });
db.votacoes.createIndex({ "data": -1 });
db.votacoes.createIndex({ "voto": 1 });
db.votacoes.createIndex({ "politico_id": 1, "data": -1 });

// Índices para proposições
db.proposicoes.createIndex({ "tipo": 1 });
db.proposicoes.createIndex({ "ano": 1 });
db.proposicoes.createIndex({ "autor_id": 1 });
db.proposicoes.createIndex({ "situacao": 1 });
db.proposicoes.createIndex({ "tema": 1 });
db.proposicoes.createIndex({ "ementa": "text" });

// Índices para despesas
db.despesas.createIndex({ "politico_id": 1 });
db.despesas.createIndex({ "tipo": 1 });
db.despesas.createIndex({ "data": -1 });
db.despesas.createIndex({ "ano_referencia": 1, "mes_referencia": 1 });
db.despesas.createIndex({ "politico_id": 1, "ano_referencia": 1 });
db.despesas.createIndex({ "valor": -1 });

// Índices para presenças
db.presencas.createIndex({ "politico_id": 1 });
db.presencas.createIndex({ "data": -1 });
db.presencas.createIndex({ "politico_id": 1, "data": -1 });

// Inserir partidos iniciais
db.partidos.insertMany([
  { sigla: "PT", nome: "Partido dos Trabalhadores", cor: "#CC0000" },
  { sigla: "PL", nome: "Partido Liberal", cor: "#003366" },
  { sigla: "UNIÃO", nome: "União Brasil", cor: "#2E3092" },
  { sigla: "PP", nome: "Progressistas", cor: "#0066CC" },
  { sigla: "MDB", nome: "Movimento Democrático Brasileiro", cor: "#00AA00" },
  { sigla: "PSD", nome: "Partido Social Democrático", cor: "#FF6600" },
  { sigla: "REPUBLICANOS", nome: "Republicanos", cor: "#0033CC" },
  { sigla: "PDT", nome: "Partido Democrático Trabalhista", cor: "#FF0000" },
  { sigla: "PSDB", nome: "Partido da Social Democracia Brasileira", cor: "#003399" },
  { sigla: "PSOL", nome: "Partido Socialismo e Liberdade", cor: "#FFD700" },
  { sigla: "PSB", nome: "Partido Socialista Brasileiro", cor: "#FF6347" },
  { sigla: "PODE", nome: "Podemos", cor: "#00CED1" },
  { sigla: "CIDADANIA", nome: "Cidadania", cor: "#9932CC" },
  { sigla: "AVANTE", nome: "Avante", cor: "#FF8C00" },
  { sigla: "SOLIDARIEDADE", nome: "Solidariedade", cor: "#FF4500" },
  { sigla: "PCdoB", nome: "Partido Comunista do Brasil", cor: "#8B0000" },
  { sigla: "PV", nome: "Partido Verde", cor: "#228B22" },
  { sigla: "NOVO", nome: "Partido Novo", cor: "#FF6600" },
  { sigla: "REDE", nome: "Rede Sustentabilidade", cor: "#00AA66" }
]);

print('✅ Banco de dados lupa_cidada inicializado com sucesso!');

