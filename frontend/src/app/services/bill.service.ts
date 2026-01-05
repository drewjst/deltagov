import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { AppConfigService } from '../config/app-config.service';

// API Response Types
export interface BillVersion {
  id: number;
  versionCode: string;
  date: string;
  contentHash: string;
  label: string;
}

export interface Bill {
  id: number;
  congress: number;
  billNumber: number;
  billType: string;
  title: string;
  sponsor: string;
  originChamber: string;
  currentStatus: string;
  updateDate: string;
  versions?: BillVersion[];
}

export interface DiffLine {
  lineNumber: number;
  type: 'insertion' | 'deletion' | 'unchanged';
  text: string;
}

export interface DiffSegment {
  type: 'insertion' | 'deletion' | 'unchanged';
  text: string;
}

export interface DiffResponse {
  fromVersion: string;
  toVersion: string;
  insertions: number;
  deletions: number;
  lines: DiffLine[];
  segments: DiffSegment[];
}

export interface BillsListResponse {
  bills: Bill[];
  total: number;
}

@Injectable({
  providedIn: 'root',
})
export class BillService {
  private readonly http = inject(HttpClient);
  private readonly configService = inject(AppConfigService);

  private get apiUrl(): string {
    return this.configService.apiUrl;
  }

  /**
   * Get H.R. 1 - The One Big Beautiful Bill
   * This endpoint auto-fetches from Congress.gov if not cached.
   */
  getHR1(): Observable<Bill> {
    return this.http.get<Bill>(`${this.apiUrl}/bills/hr1`);
  }

  /**
   * Trigger fetching H.R. 1 from Congress.gov
   */
  fetchHR1(): Observable<Bill> {
    return this.http.post<Bill>(`${this.apiUrl}/bills/hr1/fetch`, {});
  }

  /**
   * Get all bills
   */
  getBills(): Observable<Bill[]> {
    return this.http.get<BillsListResponse>(`${this.apiUrl}/bills`).pipe(
      map((response) => response.bills),
    );
  }

  /**
   * Get a specific bill by ID
   */
  getBill(id: number): Observable<Bill> {
    return this.http.get<Bill>(`${this.apiUrl}/bills/${id}`);
  }

  /**
   * Get versions for a specific bill
   */
  getBillVersions(billId: number): Observable<BillVersion[]> {
    return this.http.get<{ billId: number; versions: BillVersion[] }>(`${this.apiUrl}/bills/${billId}/versions`).pipe(
      map((response) => response.versions),
    );
  }

  /**
   * Compute diff between two versions
   */
  computeDiff(billId: number, fromVersionId: number, toVersionId: number): Observable<DiffResponse> {
    return this.http.get<DiffResponse>(`${this.apiUrl}/bills/${billId}/diff/${fromVersionId}/${toVersionId}`);
  }
}
