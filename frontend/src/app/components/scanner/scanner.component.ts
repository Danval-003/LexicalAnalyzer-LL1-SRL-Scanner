import { Component, Input } from '@angular/core';
import { MatDialogModule, MatDialog } from '@angular/material/dialog';
import { ImageModalComponent } from '../image-modal/image-modal.component';

@Component({
  selector: 'app-scanner',
  standalone: true,
  imports: [MatDialogModule, ImageModalComponent],
  templateUrl: './scanner.component.html',
  styleUrl: './scanner.component.scss'
})
export class ScannerComponent {
  @Input() name: string = '';
  @Input() image: string = '';

  constructor(public dialog: MatDialog) {}

  openDialog(): void {
    const dialogRef = this.dialog.open(ImageModalComponent, {
      maxWidth: '100vw',
      width: 'fit-content',
      data: { name: 'Angular', image: this.image } // Puedes pasar datos al diálogo si es necesario
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log(result);
      console.log('The dialog was closed');
    });
  }

  // Function to copy name
  copyName(): void {
    navigator.clipboard.writeText(this.name).then(
      () => {
        console.log('Name copied to clipboard');
      },
      (err) => {
        console.error('Could not copy name: ', err);
      }
    );
  }
}
