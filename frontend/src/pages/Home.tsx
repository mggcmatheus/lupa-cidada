import { Link } from 'react-router-dom';
import { Search, Scale, TrendingUp, Users, Vote, Receipt, FileText, ArrowRight } from 'lucide-react';
import { Button } from '../components/ui/Button';
import { Card, CardContent } from '../components/ui/Card';

const stats = [
  { icon: Users, value: '594', label: 'Políticos', color: 'text-accent-primary' },
  { icon: Vote, value: '12.847', label: 'Votações', color: 'text-accent-success' },
  { icon: Receipt, value: 'R$ 2.3B', label: 'Em Despesas', color: 'text-accent-warning' },
  { icon: FileText, value: '8.234', label: 'Proposições', color: 'text-accent-secondary' },
];

const features = [
  {
    icon: Search,
    title: 'Busca Avançada',
    description: 'Encontre políticos por nome, partido, cargo, estado e muito mais.',
  },
  {
    icon: TrendingUp,
    title: 'Estatísticas',
    description: 'Acompanhe votações, presenças, proposições e gastos detalhados.',
  },
  {
    icon: Scale,
    title: 'Compare',
    description: 'Compare a atuação de até 4 políticos lado a lado.',
  },
];

export function Home() {
  return (
    <div className="relative">
      {/* Hero Section */}
      <section className="relative py-20 lg:py-32 overflow-hidden">
        {/* Background effects */}
        <div className="absolute inset-0 bg-grid-pattern opacity-30" />
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-accent-primary/10 rounded-full blur-3xl" />
        <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-accent-secondary/10 rounded-full blur-3xl" />

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-4xl mx-auto">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-accent-primary/10 border border-accent-primary/30 text-accent-primary text-sm font-medium mb-8 animate-fade-in">
              <span className="w-2 h-2 rounded-full bg-accent-primary animate-pulse" />
              Dados atualizados diariamente
            </div>

            {/* Heading */}
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-content-primary mb-6 animate-slide-up">
              Acompanhe seus{' '}
              <span className="text-gradient">representantes</span>
            </h1>

            <p className="text-lg sm:text-xl text-content-secondary mb-10 max-w-2xl mx-auto animate-slide-up animate-delay-100">
              Uma plataforma transparente para consultar votações, despesas e a atuação
              dos políticos brasileiros. Porque a democracia começa com informação.
            </p>

            {/* CTA Buttons */}
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 animate-slide-up animate-delay-200">
              <Link to="/politicos">
                <Button size="lg" className="w-full sm:w-auto">
                  <Search className="w-5 h-5" />
                  Explorar Políticos
                </Button>
              </Link>
              <Link to="/comparar">
                <Button variant="secondary" size="lg" className="w-full sm:w-auto">
                  <Scale className="w-5 h-5" />
                  Comparar Políticos
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-16 border-y border-border bg-background-secondary/30">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-6">
            {stats.map((stat, index) => {
              const Icon = stat.icon;
              return (
                <div
                  key={stat.label}
                  className="text-center animate-slide-up"
                  style={{ animationDelay: `${index * 100}ms` }}
                >
                  <div className="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-background-card border border-border mb-4">
                    <Icon className={`w-6 h-6 ${stat.color}`} />
                  </div>
                  <p className="text-3xl font-bold font-mono text-content-primary mb-1">
                    {stat.value}
                  </p>
                  <p className="text-sm text-content-secondary">{stat.label}</p>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-content-primary mb-4">
              Como funciona
            </h2>
            <p className="text-content-secondary max-w-2xl mx-auto">
              Acesse dados públicos de forma simples e organizada
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            {features.map((feature, index) => {
              const Icon = feature.icon;
              return (
                <Card
                  key={feature.title}
                  variant="hover"
                  className="text-center animate-slide-up"
                  style={{ animationDelay: `${index * 100}ms` }}
                >
                  <CardContent className="pt-6">
                    <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-gradient-to-br from-accent-primary to-accent-secondary mb-6">
                      <Icon className="w-7 h-7 text-background" />
                    </div>
                    <h3 className="text-xl font-semibold text-content-primary mb-3">
                      {feature.title}
                    </h3>
                    <p className="text-content-secondary">
                      {feature.description}
                    </p>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-r from-accent-primary/10 via-transparent to-accent-secondary/10" />
        
        <div className="relative max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl font-bold text-content-primary mb-6">
            Pronto para fiscalizar?
          </h2>
          <p className="text-content-secondary mb-8 max-w-xl mx-auto">
            Comece agora a acompanhar a atuação dos políticos que representam você.
            Informação é poder.
          </p>
          <Link to="/politicos">
            <Button size="lg">
              Começar agora
              <ArrowRight className="w-5 h-5" />
            </Button>
          </Link>
        </div>
      </section>
    </div>
  );
}

