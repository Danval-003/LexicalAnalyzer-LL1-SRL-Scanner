import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SimulateText } from '../../../apiCalls';

@Component({
  selector: 'app-simulater',
  standalone: true,
  imports: [
    CommonModule
  ],
  templateUrl: './simulater.component.html',
  styleUrls: ['./simulater.component.scss']
})
export class SimulaterComponent {

  success: boolean = false;
  message: string = '';
  resultScanner: string = '';

  constructor() { }

  onSimulate(content: string, scanner: string, slr: string) {
    console.log(content, scanner, slr);
    SimulateText(content, scanner, slr).then((response) => {
      console.log(response);
      this.success = response.sim.accept;
      this.message = response.message;
      this.resultScanner = response.scannerResult;
      console.log(this.success);
    }).catch((error) => {
      console.log(error);
      this.success = false;
      this.message = error.message; // Mostrar el mensaje de error específico
      this.resultScanner = '';
    });
  }
}
