import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./pages/landing/landing').then((m) => m.LandingPage),
  },
  {
    path: 'diffs',
    loadComponent: () =>
      import('./pages/bills/bills').then((m) => m.BillsPage),
  },
  {
    path: 'bills',
    redirectTo: 'diffs',
    pathMatch: 'full',
  },
  {
    path: 'history',
    loadComponent: () =>
      import('./pages/history/history').then((m) => m.HistoryPage),
  },
  {
    path: 'influence',
    loadComponent: () =>
      import('./pages/influence/influence').then((m) => m.InfluencePage),
  },
  {
    path: 'money-trail',
    loadComponent: () =>
      import('./pages/money-trail/money-trail').then((m) => m.MoneyTrailPage),
  },
];
