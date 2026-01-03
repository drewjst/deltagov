import { ChangeDetectionStrategy, Component } from '@angular/core';
import { Header } from './components/header/header';
import { Workspace } from './components/workspace/workspace';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [Header, Workspace],
  templateUrl: './app.html',
  styleUrl: './app.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class App {}
