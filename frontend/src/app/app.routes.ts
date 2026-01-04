import { Routes } from '@angular/router';
import { LandingPage } from './pages/landing/landing';
import { Workspace } from './components/workspace/workspace';

export const routes: Routes = [
  {
    path: '',
    component: LandingPage,
  },
  {
    path: 'app',
    component: Workspace,
  },
];
