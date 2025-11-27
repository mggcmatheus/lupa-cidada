import { Github, Heart } from 'lucide-react';

export function Footer() {
  return (
    <footer className="border-t border-border bg-background-secondary/50 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Sobre */}
          <div>
            <h3 className="text-lg font-semibold text-content-primary mb-3">
              Lupa Cidadã
            </h3>
            <p className="text-content-secondary text-sm leading-relaxed">
              Uma plataforma de transparência política para ajudar os brasileiros
              a acompanhar a atuação de seus representantes eleitos.
            </p>
          </div>

          {/* Links */}
          <div>
            <h3 className="text-lg font-semibold text-content-primary mb-3">
              Fontes de Dados
            </h3>
            <ul className="space-y-2 text-sm">
              <li>
                <a
                  href="https://dadosabertos.camara.leg.br"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-content-secondary hover:text-accent-primary transition-colors"
                >
                  Portal da Câmara dos Deputados
                </a>
              </li>
              <li>
                <a
                  href="https://www12.senado.leg.br/dados-abertos"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-content-secondary hover:text-accent-primary transition-colors"
                >
                  Portal do Senado Federal
                </a>
              </li>
              <li>
                <a
                  href="https://portaldatransparencia.gov.br"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-content-secondary hover:text-accent-primary transition-colors"
                >
                  Portal da Transparência
                </a>
              </li>
            </ul>
          </div>

          {/* Contribuir */}
          <div>
            <h3 className="text-lg font-semibold text-content-primary mb-3">
              Contribua
            </h3>
            <p className="text-content-secondary text-sm mb-3">
              Este é um projeto open source. Contribuições são bem-vindas!
            </p>
            <a
              href="https://github.com/seu-usuario/lupa-cidada"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 text-sm text-accent-primary hover:text-accent-primary/80 transition-colors"
            >
              <Github className="w-4 h-4" />
              Ver no GitHub
            </a>
          </div>
        </div>

        <div className="mt-8 pt-8 border-t border-border flex flex-col sm:flex-row items-center justify-between gap-4">
          <p className="text-content-muted text-sm">
            © {new Date().getFullYear()} Lupa Cidadã. Todos os direitos reservados.
          </p>
          <p className="text-content-muted text-sm flex items-center gap-1">
            Feito com <Heart className="w-4 h-4 text-accent-danger" /> para o Brasil
          </p>
        </div>
      </div>
    </footer>
  );
}

