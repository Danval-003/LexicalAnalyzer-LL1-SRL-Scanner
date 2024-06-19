/* eslint-disable @typescript-eslint/no-explicit-any */
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SimulateText } from '../../../apiCalls';

@Component({
  selector: 'app-simulater',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './simulater.component.html',
  styleUrls: ['./simulater.component.scss'],
})
export class SimulaterComponent {
  success: boolean = false;
  message: string = '';
  resultScanner: string = '';

  constructor() {}

  onSimulate(content: string, scanner: string, slr: string) {
    console.log(content, scanner, slr);
    SimulateText(content, scanner, slr)
      .then((response: any) => {
        console.log(response);
        this.success = response.sim?.accept || false;
        this.message = response.message || 'No message returned';
        this.resultScanner = response.scannerResult || 'No scanner result';
        console.log(this.success);
      })
      .catch((error: any) => {
        console.log(error);
        this.success = false;
        this.message = error.message || 'An error occurred';
        this.resultScanner = '';
      });
  }
}
