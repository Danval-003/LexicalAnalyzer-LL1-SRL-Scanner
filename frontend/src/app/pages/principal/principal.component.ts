import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CreateScanner, CreateSLR } from '../../../apiCalls';
import { ScannerComponent } from '../../components/scanner/scanner.component';
import { TableSlrComponent } from '../../components/table-slr/table-slr.component';

@Component({
  selector: 'app-principal',
  standalone: true,
  imports: [CommonModule, ScannerComponent, TableSlrComponent],
  templateUrl: './principal.component.html',
  styleUrls: ['./principal.component.scss']
})
export class PrincipalComponent {
  title = 'frontend';
  yalex: string = 'Yalex';
  scanners: {
    name: string,
    image: string
  }[] = [];

  constructor() {}

  slr: {
    name: string,
    image: string
  } = {
    name: '',
    image: ''
  };

  yapar: string = 'Yapar';
  errorYalex: boolean = false;
  errorYapar: boolean = false;

  onContentYalexChange(value: string) {
    this.yalex = value;
    console.log(this.yalex);

    CreateScanner(value).then((response) => {
      console.log(response);
      for (let i = 0; i < response.names.length; i++) {
        const scanner = {
          name: response.names[i],
          image: "http://localhost:8000/image/"+response.filesId[i]
        };

        this.scanners.push(scanner);
      }
    }).catch((error) => {
      console.log(error);
    });
  }

  onContentYaparChange(value: string) {
    this.yapar = value;
    console.log(this.yapar);

    CreateSLR(value).then((response) => {
      console.log(response);
      this.slr = {
        name: response.name,
        image: response.imageURL
      };
    }).catch((error) => {
      console.log(error);
    });
  }
}
