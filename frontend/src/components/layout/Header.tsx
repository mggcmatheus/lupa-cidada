import { Link, useLocation } from 'react-router-dom';
import { Search, Scale, BarChart3, Menu, X } from 'lucide-react';
import { useState } from 'react';
import { cn } from '../../lib/utils';
import { useComparacaoStore } from '../../stores/useComparacaoStore';

const navItems = [
  { path: '/', label: 'Início', icon: BarChart3 },
  { path: '/politicos', label: 'Políticos', icon: Search },
  { path: '/comparar', label: 'Comparar', icon: Scale },
];

export function Header() {
  const location = useLocation();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const { politicosSelecionados } = useComparacaoStore();

  return (
    <header className="fixed top-0 left-0 right-0 z-50 glass border-b border-border/50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-3 group">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-accent-primary to-accent-secondary flex items-center justify-center shadow-glow-cyan group-hover:scale-105 transition-transform">
              <Search className="w-5 h-5 text-background" />
            </div>
            <div>
              <h1 className="text-xl font-bold text-content-primary">
                Lupa <span className="text-gradient">Cidadã</span>
              </h1>
              <p className="text-xs text-content-muted hidden sm:block">
                Transparência Política
              </p>
            </div>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;
              const isComparar = item.path === '/comparar';

              return (
                <Link
                  key={item.path}
                  to={item.path}
                  className={cn(
                    'relative flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all duration-200',
                    isActive
                      ? 'text-accent-primary bg-accent-primary/10'
                      : 'text-content-secondary hover:text-content-primary hover:bg-background-secondary'
                  )}
                >
                  <Icon className="w-4 h-4" />
                  {item.label}
                  {isComparar && politicosSelecionados.length > 0 && (
                    <span className="absolute -top-1 -right-1 w-5 h-5 bg-accent-primary text-background text-xs font-bold rounded-full flex items-center justify-center">
                      {politicosSelecionados.length}
                    </span>
                  )}
                </Link>
              );
            })}
          </nav>

          {/* Mobile menu button */}
          <button
            className="md:hidden p-2 rounded-lg text-content-secondary hover:text-content-primary hover:bg-background-secondary"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          >
            {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <nav className="md:hidden py-4 border-t border-border animate-slide-up">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.path;
              const isComparar = item.path === '/comparar';

              return (
                <Link
                  key={item.path}
                  to={item.path}
                  onClick={() => setMobileMenuOpen(false)}
                  className={cn(
                    'flex items-center gap-3 px-4 py-3 rounded-lg font-medium transition-all duration-200',
                    isActive
                      ? 'text-accent-primary bg-accent-primary/10'
                      : 'text-content-secondary hover:text-content-primary hover:bg-background-secondary'
                  )}
                >
                  <Icon className="w-5 h-5" />
                  {item.label}
                  {isComparar && politicosSelecionados.length > 0 && (
                    <span className="ml-auto w-6 h-6 bg-accent-primary text-background text-sm font-bold rounded-full flex items-center justify-center">
                      {politicosSelecionados.length}
                    </span>
                  )}
                </Link>
              );
            })}
          </nav>
        )}
      </div>
    </header>
  );
}

