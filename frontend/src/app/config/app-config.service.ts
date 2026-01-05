import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { AppConfig } from './app-config';

@Injectable({
  providedIn: 'root',
})
export class AppConfigService {
  private config: AppConfig = {
    apiUrl: '',
    congressApiKey: '',
  };

  constructor(private http: HttpClient) {}

  async loadConfig(): Promise<void> {
    const config = await firstValueFrom(
      this.http.get<AppConfig>('/assets/config.json')
    );
    this.config = config;
  }

  get apiUrl(): string {
    return this.config.apiUrl;
  }

  get congressApiKey(): string {
    return this.config.congressApiKey;
  }
}
