import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule} from '@angular/material/dialog';

@Component({
  selector: 'app-image-modal',
  standalone: true,
  imports: [MatDialogModule],
  templateUrl: './image-modal.component.html',
  styleUrl: './image-modal.component.scss'
})
export class ImageModalComponent {
  image2: string = '';
  constructor(
    public dialogRef: MatDialogRef<ImageModalComponent>,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {
    this.image2 = data.image;
    console.log("Hola",this.image2)
  }

  onClose():void {
    this.dialogRef.close();
  }
}
