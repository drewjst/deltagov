import { Injectable } from '@angular/core';
import { Observable, of, delay } from 'rxjs';

export interface SponsorProfile {
  name: string;
  role: string;
  imageUrl: string;
  contributions: {
    industry: string;
    amount: string;
  }[];
}

export type StakeholderPosition = 'Support' | 'Oppose' | 'Amend';

export interface Stakeholder {
  organization: string;
  position: StakeholderPosition;
  reportedSpend: string;
  textConnection: {
    section: string;
    description: string;
  };
}

export interface InfluenceData {
  sponsor: SponsorProfile;
  stakeholders: Stakeholder[];
}

@Injectable({
  providedIn: 'root',
})
export class InfluenceService {
  getInfluenceData(): Observable<InfluenceData> {
    // Mock data
    const data: InfluenceData = {
      sponsor: {
        name: 'Sen. Elena Rodriguez',
        role: 'Chair, Senate Committee on Commerce',
        imageUrl: 'https://i.pravatar.cc/150?u=senator',
        contributions: [
          { industry: 'Technology', amount: '$1.2M' },
          { industry: 'Green Energy', amount: '$850k' },
          { industry: 'Telecommunications', amount: '$500k' },
        ],
      },
      stakeholders: [
        {
          organization: 'Amazon.com',
          position: 'Support',
          reportedSpend: '$1.2M',
          textConnection: {
            section: 'Sec. 402',
            description: 'Tax Credits for Cloud Infrastructure',
          },
        },
        {
          organization: 'Sierra Club',
          position: 'Amend',
          reportedSpend: '$450k',
          textConnection: {
            section: 'Sec. 201',
            description: 'Environmental Impact Assessments',
          },
        },
        {
          organization: 'National Coal Association',
          position: 'Oppose',
          reportedSpend: '$2.1M',
          textConnection: {
            section: 'Sec. 105',
            description: 'Carbon Emission Caps',
          },
        },
        {
          organization: 'Google',
          position: 'Support',
          reportedSpend: '$900k',
          textConnection: {
            section: 'Sec. 405',
            description: 'AI Research Grants',
          },
        },
        {
          organization: 'Local Manufacturers Union',
          position: 'Oppose',
          reportedSpend: '$120k',
          textConnection: {
            section: 'Sec. 301',
            description: 'Labor Standards',
          },
        },
      ],
    };

    return of(data).pipe(delay(500)); // Simulate API latency
  }
}
